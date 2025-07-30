<script lang="ts">
  export default {
    name: 'Notifications',
  }
</script>

<script setup lang="ts">
  import Dropdown from 'components/dropdown/index.vue'
  import { $i18nGlobal } from 'i18n/index'
  import { $formatDate } from 'utils/date'
  import { $screen } from 'utils/screen'
  import { $store } from 'utils/store'
  import { computed, onMounted } from 'vue'

  const locale = computed(() => {
    const l = $i18nGlobal.locale.value
    return l.split('-')[0]
  })

  const notificationsList = computed(() => {
    return $store.get('notifications')
  })

  const notificationsCount = computed(() => {
    return notificationsList.value?.length || 0
  })

  function removeScripts(html: string) {
    return html.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
  }

  onMounted(() => {
    $store.updateNotifications()
  })
</script>

<template>
  <Dropdown :close-bnt="$screen.isMobile.value" :class="`rounded-none`">
    <template #trigger>
      <button class="btn btn-circle btn-ghost relative mobile:gap-0">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-6 h-6 relative top-[1px]"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M14.857 17.082a23.848 23.848 0 0 0 5.454-1.31A8.967 8.967 0 0 1 18 9.75V9A6 6 0 0 0 6 9v.75a8.967 8.967 0 0 1-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 0 1-5.714 0m5.714 0a3 3 0 1 1-5.714 0"
          />
        </svg>
        <div
          v-if="notificationsCount > 0"
          class="badge badge-sm badge-error absolute top-0 right-0"
        >
          {{ notificationsCount }}
        </div>
        <div class="pc:hidden text-xs font-normal truncate">{{ $t('Events') }}</div>
      </button>
    </template>
    <template #drop>
      <div class="relative">
        <ul tabindex="0" class="p-3 max-w-[500px] max-h-[100vh] overflow-hidden overflow-y-auto">
          <template v-if="notificationsCount > 0">
            <li
              v-for="item in notificationsList"
              :key="item.uid"
              class="py-4 [&:not(:last-child)]:border-b border-b-gray-500"
            >
              <div class="text-right italic text-sm">{{ $formatDate(item.created_at) }}</div>
              <div
                style="white-space: break-spaces"
                v-html="removeScripts(item.texts[locale])"
              ></div>
            </li>
          </template>

          <template v-else>
            {{ $t('notifications.empty') }}
          </template>
        </ul>
      </div>
    </template>
  </Dropdown>
</template>
