<script lang="ts">
  export default {
    name: 'PagePartnerApplications',
  }
</script>

<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { $t } from 'i18n/index'
  import { $me } from 'utils/me'
  import { $formatDate } from 'utils/date'
  import { $notify } from 'utils/notify'
  import Avatar from 'vue-boring-avatars'
  import { 
    getPartnerApplications, 
    getMyApplications, 
    getPartnerApplicationsCount, 
    getMyApplicationsCount,
    createPartnerApplication,
    processPartnerApplication,
    type PartnerApplication,
    type PartnerApplicationReq,
    type PartnerApplicationResponseReq
  } from 'queries/partnerApplications'

  const loading = ref(false)
  const activeTab = ref('incoming') // 'incoming' | 'outgoing'
  
  // Данные для входящих заявок
  const incomingApplications = ref<PartnerApplication[]>([])
  const incomingCount = ref(0)
  const incomingLimit = 20
  const incomingSkip = ref(0)
  
  // Данные для исходящих заявок
  const outgoingApplications = ref<PartnerApplication[]>([])
  const outgoingCount = ref(0)
  const outgoingLimit = 20
  const outgoingSkip = ref(0)
  
  // Модальное окно для создания заявки
  const showCreateModal = ref(false)
  const createForm = ref<PartnerApplicationReq>({
    partner_uid: '',
    message: ''
  })
  const createLoading = ref(false)
  
  // Чекбоксы согласия
  const agreements = ref({
    privacyPolicy: false,
    termsOfService: false,
    riskWarning: false
  })

  // На этой странице никто не должен видеть входящие заявки
  // Входящие заявки обрабатываются на странице партнерской сети
  const canViewIncomingApplications = computed(() => {
    return false
  })
  
  // Модальное окно для обработки заявки
  const showProcessModal = ref(false)
  const processForm = ref<PartnerApplicationResponseReq>({
    application_uid: '',
    status: 'approved',
    response: ''
  })
  const processLoading = ref(false)
  const currentApplication = ref<PartnerApplication | null>(null)

  // Загружаем входящие заявки
  async function loadIncomingApplications() {
    loading.value = true
    try {
      const [applications, count] = await Promise.all([
        getPartnerApplications(incomingLimit, incomingSkip.value),
        getPartnerApplicationsCount()
      ])
      incomingApplications.value = applications
      incomingCount.value = count
    } catch (error) {
      console.error('Error loading incoming applications:', error)
      $notify.show({
        type: 'error',
        text: $t('Error loading applications')
      })
    } finally {
      loading.value = false
    }
  }

  // Загружаем исходящие заявки
  async function loadOutgoingApplications() {
    loading.value = true
    try {
      const [applications, count] = await Promise.all([
        getMyApplications(outgoingLimit, outgoingSkip.value),
        getMyApplicationsCount()
      ])
      outgoingApplications.value = applications
      outgoingCount.value = count
    } catch (error) {
      console.error('Error loading outgoing applications:', error)
      $notify.show({
        type: 'error',
        text: $t('Error loading applications')
      })
    } finally {
      loading.value = false
    }
  }

  // Создать заявку на партнёрство
  // Функция открытия модального окна создания заявки
  function openCreateModal() {
    // Проверяем, есть ли у пользователя спонсор
    if (!$me.data.ref_uid) {
      $notify.show({
        type: 'error',
        text: $t('You are not assigned to any partner')
      })
      return
    }

    // Автоматически заполняем partner_uid из ref_uid пользователя
    createForm.value.partner_uid = $me.data.ref_uid
    createForm.value.message = ''
    
    // Сбрасываем чекбоксы
    agreements.value = {
      privacyPolicy: false,
      termsOfService: false,
      riskWarning: false
    }
    
    showCreateModal.value = true
  }

  // Проверка всех соглашений
  const allAgreementsChecked = computed(() => {
    return agreements.value.privacyPolicy && 
           agreements.value.termsOfService && 
           agreements.value.riskWarning
  })

  async function submitCreateApplication() {
    if (!createForm.value.partner_uid || !createForm.value.message) {
      $notify.show({
        type: 'error',
        text: $t('Please fill all fields')
      })
      return
    }

    if (!allAgreementsChecked.value) {
      $notify.show({
        type: 'error',
        text: $t('Please accept all agreements')
      })
      return
    }

    createLoading.value = true
    try {
      await createPartnerApplication(createForm.value)
      $notify.show({
        type: 'success',
        text: $t('Application created successfully')
      })
      showCreateModal.value = false
      createForm.value = { partner_uid: '', message: '' }
      await loadOutgoingApplications()
    } catch (error) {
      console.error('Error creating application:', error)
      $notify.show({
        type: 'error',
        text: $t('Error creating application')
      })
    } finally {
      createLoading.value = false
    }
  }

  // Обработать заявку
  async function submitProcessApplication() {
    if (!processForm.value.response && processForm.value.status === 'rejected') {
      $notify.show({
        type: 'error',
        text: $t('Please provide a reason for rejection')
      })
      return
    }

    processLoading.value = true
    try {
      await processPartnerApplication(processForm.value)
      $notify.show({
        type: 'success',
        text: processForm.value.status === 'approved' 
          ? $t('Application approved successfully')
          : $t('Application rejected successfully')
      })
      showProcessModal.value = false
      processForm.value = { application_uid: '', status: 'approved', response: '' }
      currentApplication.value = null
      await loadIncomingApplications()
    } catch (error) {
      console.error('Error processing application:', error)
      $notify.show({
        type: 'error',
        text: $t('Error processing application')
      })
    } finally {
      processLoading.value = false
    }
  }

  // Открыть модальное окно для обработки заявки
  function openProcessModal(application: PartnerApplication) {
    currentApplication.value = application
    processForm.value.application_uid = application.uid
    processForm.value.status = 'approved'
    processForm.value.response = ''
    showProcessModal.value = true
  }

  // Получить статус заявки
  const getStatusBadge = (status: string) => {
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
  const getStatusText = (status: string) => {
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

  // Переключение вкладок
  const switchTab = async (tab: string) => {
    // Клиенты не могут переключиться на входящие заявки
    if (tab === 'incoming' && !canViewIncomingApplications.value) {
      return
    }
    
    activeTab.value = tab
    if (tab === 'incoming') {
      await loadIncomingApplications()
    } else {
      await loadOutgoingApplications()
    }
  }

  onMounted(async () => {
    // На этой странице показываем только исходящие заявки для всех пользователей
    activeTab.value = 'outgoing'
    await loadOutgoingApplications()
  })
</script>

<template>
  <div class="p-4">
    <!-- Заголовок -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold mb-4">{{ $t('Partner Applications') }}</h1>
      
      <!-- Счетчик заявок -->
      <div class="stats shadow">
        <div class="stat">
          <div class="stat-title">{{ $t('My Applications') }}</div>
          <div class="stat-value">{{ outgoingCount }}</div>
          <div class="stat-desc">{{ $t('Total applications sent') }}</div>
        </div>
      </div>
    </div>

    <!-- Кнопка создания заявки -->
    <div class="mb-6">
      <button 
        @click="openCreateModal"
        class="btn btn-primary"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {{ $t('Create Application') }}
      </button>
    </div>

    <!-- Список заявок -->
    <div v-if="loading" class="flex justify-center">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
    
    <div v-else>
      <!-- Входящие заявки удалены - обрабатываются на странице партнерской сети -->
      <div v-if="false">
        <div v-if="incomingApplications.length === 0" class="text-center py-8 text-gray-500">
          {{ $t('No incoming applications') }}
        </div>
        
        <div v-else class="space-y-4">
          <div 
            v-for="application in incomingApplications" 
            :key="application.uid"
            class="card bg-base-100 shadow-xl"
          >
            <div class="card-body">
              <div class="flex items-center justify-between mb-4">
                <div class="flex items-center gap-3">
                  <div class="avatar">
                    <div class="w-12 h-12 rounded-full">
                      <Avatar :name="application.applicant?.uid || application.applicant_uid" variant="marble" />
                    </div>
                  </div>
                  <div>
                    <h3 class="font-semibold">{{ application.applicant?.nickname || application.applicant_uid }}</h3>
                    <p class="text-sm text-gray-500">{{ application.applicant?.email || 'N/A' }}</p>
                  </div>
                </div>
                <div class="badge" :class="getStatusBadge(application.status)">
                  {{ getStatusText(application.status) }}
                </div>
              </div>
              
              <div class="mb-4">
                <p class="text-sm text-gray-600 mb-2">{{ $t('Message') }}:</p>
                <p class="bg-base-200 p-3 rounded-lg">{{ application.message }}</p>
              </div>
              
              <div v-if="application.response" class="mb-4">
                <p class="text-sm text-gray-600 mb-2">{{ $t('Response') }}:</p>
                <p class="bg-base-200 p-3 rounded-lg">{{ application.response }}</p>
              </div>
              
              <div class="text-sm text-gray-500 mb-4">
                {{ $t('Created') }}: {{ $formatDate(application.created_at) }}
                <span v-if="application.processed_at">
                  | {{ $t('Processed') }}: {{ $formatDate(application.processed_at) }}
                </span>
              </div>
              
              <div v-if="application.status === 'pending'" class="card-actions justify-end">
                <button 
                  @click="openProcessModal(application)"
                  class="btn btn-primary btn-sm"
                >
                  {{ $t('Process Application') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Мои заявки -->
      <div>
        <div v-if="outgoingApplications.length === 0" class="text-center py-8 text-gray-500">
          {{ $t('No outgoing applications') }}
        </div>
        
        <div v-else class="space-y-4">
          <div 
            v-for="application in outgoingApplications" 
            :key="application.uid"
            class="card bg-base-100 shadow-xl"
          >
            <div class="card-body">
              <div class="flex items-center justify-between mb-4">
                <div class="flex items-center gap-3">
                  <div class="avatar">
                    <div class="w-12 h-12 rounded-full">
                      <Avatar :name="application.partner?.uid || application.partner_uid" variant="marble" />
                    </div>
                  </div>
                  <div>
                    <h3 class="font-semibold">{{ application.partner?.nickname || application.partner_uid }}</h3>
                    <p class="text-sm text-gray-500">{{ application.partner?.email || 'N/A' }}</p>
                  </div>
                </div>
                <div class="badge" :class="getStatusBadge(application.status)">
                  {{ getStatusText(application.status) }}
                </div>
              </div>
              
              <div class="mb-4">
                <p class="text-sm text-gray-600 mb-2">{{ $t('Message') }}:</p>
                <p class="bg-base-200 p-3 rounded-lg">{{ application.message }}</p>
              </div>
              
              <div v-if="application.response" class="mb-4">
                <p class="text-sm text-gray-600 mb-2">{{ $t('Response') }}:</p>
                <p class="bg-base-200 p-3 rounded-lg">{{ application.response }}</p>
              </div>
              
              <div class="text-sm text-gray-500">
                {{ $t('Created') }}: {{ $formatDate(application.created_at) }}
                <span v-if="application.processed_at">
                  | {{ $t('Processed') }}: {{ $formatDate(application.processed_at) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Модальное окно создания заявки -->
    <div v-if="showCreateModal" class="modal modal-open">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">{{ $t('Create Partner Application') }}</h3>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">{{ $t('Partner UID') }}</span>
          </label>
          <input 
            v-model="createForm.partner_uid"
            type="text" 
            class="input input-bordered bg-base-200" 
            :placeholder="$t('Your assigned partner')"
            readonly
            disabled
          />
          <label class="label">
            <span class="label-text-alt text-info">{{ $t('Automatically assigned to your sponsor') }}</span>
          </label>
        </div>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">{{ $t('Message') }}</span>
          </label>
          <textarea 
            v-model="createForm.message"
            class="textarea textarea-bordered" 
            :placeholder="$t('Enter your message')"
            rows="4"
          ></textarea>
        </div>

        <!-- Чекбоксы согласия -->
        <div class="divider">{{ $t('Agreements') }}</div>
        
        <div class="form-control">
          <label class="label cursor-pointer">
            <span class="label-text">{{ $t('I agree with Privacy Policy') }}</span>
            <input 
              v-model="agreements.privacyPolicy"
              type="checkbox" 
              class="checkbox checkbox-primary" 
            />
          </label>
        </div>

        <div class="form-control">
          <label class="label cursor-pointer">
            <span class="label-text">{{ $t('I agree with Terms of Service') }}</span>
            <input 
              v-model="agreements.termsOfService"
              type="checkbox" 
              class="checkbox checkbox-primary" 
            />
          </label>
        </div>

        <div class="form-control">
          <label class="label cursor-pointer">
            <span class="label-text">{{ $t('I agree with Risk Warning') }}</span>
            <input 
              v-model="agreements.riskWarning"
              type="checkbox" 
              class="checkbox checkbox-primary" 
            />
          </label>
        </div>
        
        <div class="modal-action">
          <button 
            @click="showCreateModal = false"
            class="btn"
            :disabled="createLoading"
          >
            {{ $t('Cancel') }}
          </button>
          <button 
            @click="submitCreateApplication"
            class="btn btn-primary"
            :disabled="createLoading || !allAgreementsChecked"
          >
            <span v-if="createLoading" class="loading loading-spinner loading-sm"></span>
            {{ $t('Create') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Модальное окно обработки заявки -->
    <div v-if="showProcessModal" class="modal modal-open">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">{{ $t('Process Application') }}</h3>
        
        <div v-if="currentApplication" class="mb-4">
          <p class="text-sm text-gray-600">{{ $t('From') }}: {{ currentApplication.applicant?.nickname || currentApplication.applicant_uid }}</p>
          <p class="text-sm text-gray-600 mb-2">{{ $t('Message') }}: {{ currentApplication.message }}</p>
        </div>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">{{ $t('Status') }}</span>
          </label>
          <select v-model="processForm.status" class="select select-bordered">
            <option value="approved">{{ $t('Approve') }}</option>
            <option value="rejected">{{ $t('Reject') }}</option>
          </select>
        </div>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">{{ $t('Response') }}</span>
          </label>
          <textarea 
            v-model="processForm.response"
            class="textarea textarea-bordered" 
            :placeholder="processForm.status === 'rejected' ? $t('Enter rejection reason') : $t('Enter response (optional)')"
            rows="4"
          ></textarea>
        </div>
        
        <div class="modal-action">
          <button 
            @click="showProcessModal = false"
            class="btn"
            :disabled="processLoading"
          >
            {{ $t('Cancel') }}
          </button>
          <button 
            @click="submitProcessApplication"
            class="btn btn-primary"
            :disabled="processLoading"
          >
            <span v-if="processLoading" class="loading loading-spinner loading-sm"></span>
            {{ $t('Process') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template> 