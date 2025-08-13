<script lang="ts">
  export default {
    name: 'PagePartners',
  }
</script>

<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { $requests } from 'queries/index'
  import { $me } from 'utils/me'
  import { $t } from 'i18n/index'
  import Avatar from 'vue-boring-avatars'
  import { $formatDate } from 'utils/date'
  import { $notify } from 'utils/notify'
  import { 
    getPartnerApplications, 
    getPartnerApplicationsCount,
    processPartnerApplication,
    type PartnerApplication,
    type PartnerApplicationResponseReq
  } from 'queries/partnerApplications'

  const loading = ref(false)
  const partnersCount = ref(0)
  const partners = ref([])
  const availableSlots = ref(0)
  const maxPartners = ref(1) // По умолчанию Researcher (1 партнёр)

  // Данные для заявок
  const applications = ref<PartnerApplication[]>([])
  const applicationsCount = ref(0)
  const applicationsLoading = ref(false)

  // Загрузка при обработке заявки
  const processLoading = ref(false)

  // Загружаем данные
  async function loadData() {
    loading.value = true
    try {
      console.log('Current user uid:', $me.data.uid)
      console.log('Current user nickname:', $me.data.nickname)
      
      // Загружаем количество партнёров
      partnersCount.value = await $requests.partners.count($me.data.uid)
      
      // Загружаем список партнёров
      partners.value = await $requests.teams.getMatch({ user_uid: $me.data.uid })
      
      // Получаем лимит партнёров из абонемента
      // Используем premium_invites как лимит партнёров (временно)
      maxPartners.value = $me.data.premium_invites || 1
      
      // Вычисляем доступные слоты
      availableSlots.value = maxPartners.value - partnersCount.value
    } catch (error) {
      console.error('Error loading partners data:', error)
    } finally {
      loading.value = false
    }
  }

  // Загружаем заявки
  async function loadApplications() {
    console.log('Loading applications...')
    applicationsLoading.value = true
    try {
      const [applicationsData, count] = await Promise.all([
        getPartnerApplications(20, 0),
        getPartnerApplicationsCount()
      ])
      console.log('Applications loaded:', applicationsData)
      console.log('Applications count:', count)
      applications.value = applicationsData
      applicationsCount.value = count
    } catch (error) {
      console.error('Error loading applications:', error)
      $notify.show({
        type: 'error',
        text: $t('Error loading applications')
      })
    } finally {
      applicationsLoading.value = false
    }
  }

  // Удалить партнёра
  async function removePartner(partnerUid: string) {
    try {
      const success = await $requests.partners.removePartner($me.data.uid, partnerUid)
      if (success) {
        await loadData() // Перезагружаем данные
      }
    } catch (error) {
      console.error('Error removing partner:', error)
    }
  }

  // Принять заявку напрямую
  async function acceptApplication(application: PartnerApplication) {
    processLoading.value = true
    try {
      await processPartnerApplication({
        application_uid: application.uid,
        status: 'approved',
        response: ''
      })
      
      $notify.show({
        type: 'success',
        text: $t('Application accepted successfully')
      })
      
      // Перезагружаем данные
      await Promise.all([loadApplications(), loadData()])
    } catch (error) {
      console.error('Error accepting application:', error)
      $notify.show({
        type: 'error',
        text: $t('Error accepting application')
      })
    } finally {
      processLoading.value = false
    }
  }

  // Получить оставшееся время жизни заявки (4 часа)
  function getApplicationTimeLeft(createdAt: number) {
    const fourHours = 4 * 60 * 60 * 1000 // 4 часа в миллисекундах
    
    // Если createdAt в секундах, конвертируем в миллисекунды
    const createdAtMs = createdAt < 10000000000 ? createdAt * 1000 : createdAt
    
    const expiresAt = createdAtMs + fourHours
    const now = Date.now()
    const timeLeft = expiresAt - now
    
    console.log('Time left calculation:', { 
      createdAt, 
      createdAtMs, 
      expiresAt, 
      now, 
      timeLeft 
    })
    
    if (timeLeft <= 0) {
      return $t('Expired')
    }
    
    const hours = Math.floor(timeLeft / (60 * 60 * 1000))
    const minutes = Math.floor((timeLeft % (60 * 60 * 1000)) / (60 * 1000))
    
    if (hours > 0) {
      return `${hours}ч ${minutes}м`
    } else {
      return `${minutes}м`
    }
  }

  // Проверить, не истекла ли заявка
  function isApplicationExpired(createdAt: number) {
    const fourHours = 4 * 60 * 60 * 1000
    
    // Если createdAt в секундах, конвертируем в миллисекунды
    const createdAtMs = createdAt < 10000000000 ? createdAt * 1000 : createdAt
    
    const expiresAt = createdAtMs + fourHours
    const now = Date.now()
    
    console.log('Checking expiration:', { 
      createdAt, 
      createdAtMs, 
      expiresAt, 
      now, 
      diff: now - createdAtMs,
      expired: now > expiresAt 
    })
    
    return now > expiresAt
  }

  // Получить класс для статуса
  function getStatusBadge(status: string) {
    switch (status) {
      case 'pending':
        return 'badge-warning'
      case 'approved':
        return 'badge-success'
      case 'rejected':
        return 'badge-error'
      default:
        return 'badge-neutral'
    }
  }

  // Получить текст статуса
  function getStatusText(status: string) {
    switch (status) {
      case 'pending':
        return $t('Pending')
      case 'approved':
        return $t('Approved')
      case 'rejected':
        return $t('Rejected')
      default:
        return status
    }
  }

  // Фильтрованные заявки - только активные и не истекшие
  const activeApplications = computed(() => {
    console.log('Filtering applications:', applications.value)
    const filtered = applications.value.filter(app => {
      console.log('App:', app.uid, 'Status:', app.status, 'Created:', app.created_at, 'Expired:', isApplicationExpired(app.created_at))
      return app.status === 'pending' && !isApplicationExpired(app.created_at)
    })
    console.log('Active applications after filter:', filtered)
    return filtered
  })

  onMounted(async () => {
    await Promise.all([loadData(), loadApplications()])
  })
</script>

<template>
  <div class="p-4">
    <!-- Заголовок и статистика -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold mb-4">{{ $t('Partners') }}</h1>
      
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div class="stat bg-base-200 rounded-lg">
          <div class="stat-title">{{ $t('Assigned Partners') }}</div>
          <div class="stat-value text-primary">{{ partnersCount }}</div>
        </div>
        
        <div class="stat bg-base-200 rounded-lg">
          <div class="stat-title">{{ $t('Available Slots') }}</div>
          <div class="stat-value" :class="availableSlots > 0 ? 'text-success' : 'text-error'">
            {{ availableSlots }}
          </div>
        </div>
        
        <div class="stat bg-base-200 rounded-lg">
          <div class="stat-title">{{ $t('Total Slots') }}</div>
          <div class="stat-value">{{ maxPartners }}</div>
        </div>

        <div class="stat bg-base-200 rounded-lg">
          <div class="stat-title">{{ $t('Pending Applications') }}</div>
          <div class="stat-value text-warning">{{ activeApplications.length }}</div>
        </div>
      </div>
    </div>

    <!-- Список назначенных партнёров -->
    <div class="mb-8">
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-3">
          <h2 class="text-xl font-semibold">{{ $t('Assigned Partners') }}</h2>
          <div class="badge badge-neutral">{{ partnersCount }}</div>
        </div>
        <button class="btn btn-ghost btn-sm">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>
      
      <div v-if="loading" class="flex justify-center">
        <span class="loading loading-spinner loading-lg"></span>
      </div>
      
      <div v-else-if="partners.length === 0" class="text-center py-8 text-gray-500">
        {{ $t('No partners assigned yet') }}
      </div>
      
      <div v-else class="space-y-3">
        <div 
          v-for="partner in partners" 
          :key="partner.uid"
          class="flex items-center justify-between p-4 bg-base-200 rounded-lg hover:bg-base-300 transition-colors"
        >
          <div class="flex items-center gap-4">
            <div class="avatar">
              <div class="w-12 h-12 rounded-full">
                <Avatar :name="partner.uid" variant="marble" />
              </div>
            </div>
            <div class="flex-1">
              <h3 class="font-semibold text-base">{{ partner.nickname }}</h3>
              <p class="text-sm text-gray-500">{{ partner.email }}</p>
            </div>
            <div class="text-sm text-gray-500">
              {{ $formatDate(partner.created_at) }}
            </div>
          </div>
          
          <div class="flex items-center gap-2">
            <div class="dropdown dropdown-end">
              <div tabindex="0" role="button" class="btn btn-ghost btn-sm btn-circle">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01" />
                </svg>
              </div>
              <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                <li>
                  <a @click="removePartner(partner.uid)" class="text-error">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    {{ $t('Remove') }}
                  </a>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Список заявок на партнёрство -->
    <div v-if="availableSlots > 0" class="mb-8">
      <h2 class="text-xl font-semibold mb-4">{{ $t('Partner Applications') }}</h2>
      
      <div v-if="applicationsLoading" class="flex justify-center">
        <span class="loading loading-spinner loading-lg"></span>
      </div>
      
      <div v-else-if="activeApplications.length === 0" class="text-center py-8 text-gray-500">
        {{ $t('No active applications') }}
      </div>
      
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div 
          v-for="application in activeApplications" 
          :key="application.uid"
          class="card bg-base-100 shadow-xl"
        >
          <div class="card-body">
            <div class="flex items-center gap-3 mb-3">
              <div class="avatar">
                <div class="w-12 h-12 rounded-full">
                  <Avatar :name="application.applicant_uid" variant="marble" />
                </div>
              </div>
              <div>
                <h3 class="font-semibold">{{ application.applicant_uid }}</h3>
                <p class="text-sm text-gray-500">{{ $t('Applied') }}</p>
              </div>
            </div>
            
            <div class="text-sm text-gray-600 mb-3">
              <div>{{ $t('Applied') }}: {{ $formatDate(application.created_at) }}</div>
              <div class="flex items-center gap-2 mt-1">
                <span>{{ $t('Time left') }}:</span>
                <span 
                  :class="isApplicationExpired(application.created_at) ? 'text-error font-semibold' : 'text-warning font-semibold'"
                >
                  {{ getApplicationTimeLeft(application.created_at) }}
                </span>
              </div>
            </div>
            
            <div class="mb-3">
              <p class="text-sm text-gray-600 mb-2">{{ $t('Message') }}:</p>
              <p class="bg-base-200 p-3 rounded-lg text-sm">{{ application.message }}</p>
            </div>
            
            <div class="card-actions justify-end">
              <button 
                @click="acceptApplication(application)"
                :disabled="isApplicationExpired(application.created_at) || processLoading"
                :class="[
                  'btn btn-sm',
                  isApplicationExpired(application.created_at) ? 'btn-disabled' : 'btn-success'
                ]"
              >
                <span v-if="processLoading" class="loading loading-spinner loading-sm mr-1"></span>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                {{ isApplicationExpired(application.created_at) ? $t('Expired') : $t('Accept') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>


  </div>
</template>