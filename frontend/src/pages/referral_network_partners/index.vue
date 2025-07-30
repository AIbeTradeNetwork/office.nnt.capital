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

  const loading = ref(false)
  const partnersCount = ref(0)
  const friends = ref([])
  const partners = ref([])
  const availableSlots = ref(0)
  const maxPartners = ref(1) // По умолчанию Researcher (1 партнёр)

  // Загружаем данные
  async function loadData() {
    loading.value = true
    try {
      console.log('Current user uid:', $me.data.uid)
      console.log('Current user nickname:', $me.data.nickname)
      
      // Загружаем количество партнёров
      partnersCount.value = await $requests.partners.count($me.data.uid)
      
      // Загружаем список друзей (клиентов)
      friends.value = await $requests.friends.list($me.data.uid, 100, 0)
      
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

  // Добавить партнёра
  async function addPartner(friendUid: string) {
    try {
      const success = await $requests.partners.addPartner($me.data.uid, friendUid)
      if (success) {
        await loadData() // Перезагружаем данные
      }
    } catch (error) {
      console.error('Error adding partner:', error)
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

  // Проверяем, является ли друг уже партнёром
  const isAlreadyPartner = (friendUid: string) => {
    return partners.value.some(partner => partner.uid === friendUid)
  }

  onMounted(() => {
    loadData()
  })
</script>

<template>
  <div class="p-4">
    <!-- Заголовок и статистика -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold mb-4">{{ $t('Partners') }}</h1>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
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
      </div>
    </div>

    <!-- Список назначенных партнёров -->
    <div class="mb-8">
      <h2 class="text-xl font-semibold mb-4">{{ $t('Assigned Partners') }}</h2>
      
      <div v-if="loading" class="flex justify-center">
        <span class="loading loading-spinner loading-lg"></span>
      </div>
      
      <div v-else-if="partners.length === 0" class="text-center py-8 text-gray-500">
        {{ $t('No partners assigned yet') }}
      </div>
      
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div 
          v-for="partner in partners" 
          :key="partner.uid"
          class="card bg-base-100 shadow-xl"
        >
          <div class="card-body">
            <div class="flex items-center gap-3 mb-3">
              <div class="avatar">
                <div class="w-12 h-12 rounded-full">
                  <Avatar :name="partner.uid" variant="marble" />
                </div>
              </div>
              <div>
                <h3 class="font-semibold">{{ partner.nickname }}</h3>
                <p class="text-sm text-gray-500">{{ partner.email }}</p>
              </div>
            </div>
            
            <div class="text-sm text-gray-600 mb-3">
              {{ $t('Joined') }}: {{ $formatDate(partner.created_at) }}
            </div>
            
            <div class="card-actions justify-end">
              <button 
                @click="removePartner(partner.uid)"
                class="btn btn-sm btn-error"
              >
                {{ $t('Remove') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Список доступных клиентов для назначения -->
    <div v-if="availableSlots > 0">
      <h2 class="text-xl font-semibold mb-4">{{ $t('Available Clients') }}</h2>
      
      <div v-if="friends.length === 0" class="text-center py-8 text-gray-500">
        {{ $t('No clients available') }}
      </div>
      
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div 
          v-for="friend in friends" 
          :key="friend.uid"
          class="card bg-base-100 shadow-xl"
        >
          <div class="card-body">
            <div class="flex items-center gap-3 mb-3">
              <div class="avatar">
                <div class="w-12 h-12 rounded-full">
                  <Avatar :name="friend.uid" variant="marble" />
                </div>
              </div>
              <div>
                <h3 class="font-semibold">{{ friend.nickname }}</h3>
                <p class="text-sm text-gray-500">{{ friend.email }}</p>
              </div>
            </div>
            
            <div class="text-sm text-gray-600 mb-3">
              {{ $t('Joined') }}: {{ $formatDate(friend.created_at) }}
            </div>
            
            <div class="card-actions justify-end">
              <button 
                v-if="!isAlreadyPartner(friend.uid)"
                @click="addPartner(friend.uid)"
                class="btn btn-sm btn-primary"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
                {{ $t('Add Partner') }}
              </button>
              
              <span v-else class="badge badge-success">
                {{ $t('Already Partner') }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>