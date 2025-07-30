<script lang="ts">
  export default {
    name: 'CommandLine',
  }
</script>

<script setup lang="ts">
  import { ref } from 'vue'
  import Avatar from 'vue-boring-avatars'
  import { $formatDate } from 'utils/date'

  const props = defineProps<{
    user: any
    hiddenClass?: string
    lineid: number
    showRefCount?: boolean
  }>()

  const loading = ref(false)
  const colspan = 3
  const isHasChildrens = props.user?.team_count > 0
  const isShow = ref(false)
  const childrens = ref<any[]>([])
  const limit = 15
  let page = 0
  const isEnd = ref(true)

  async function getPartners() {
    if (!page || page < 0) page = 0
    loading.value = true
    try {
      // TODO: заменить на реальный запрос к API партнёров
      childrens.value.push(...[])
      isEnd.value = true
      page++
    } catch (error) {
      console.error(error)
    }
    loading.value = false
  }

  async function doToggle(event) {
    if (!isHasChildrens) return
    if (!childrens.value.length) {
      await getPartners()
    }
    // TODO: добавить toggle функциональность
  }
</script>

<template>
  <tr :id="lineid + '' || '0'">
    <td class="relative" :style="`padding-left: calc(1rem + (${lineid / 10}rem * 30))`">
      <div class="flex items-center gap-2">
        <button :class="!isHasChildrens ? `invisible pointer-events-none` : ''" class="btn btn-square btn-xs hi" @click.stop="doToggle">
          <span v-if="loading" class="loading loading-spinner"></span>
          <template v-else>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 inline [.open_&]:hidden">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 hidden [.open_&]:inline">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 12h14" />
            </svg>
          </template>
        </button>
        <div class="avatar placeholder">
          <div class="bg-neutral text-neutral-content rounded-full w-12 relative overflow-visible">
            <Avatar class="w-full h-full" variant="marble" :name="user.uid" />
          </div>
        </div>
        <div class="relative overflow-visible">
          <p>{{ user.nickname }}</p>
          <p>{{ user.email }}</p>
        </div>
      </div>
    </td>
    <td>{{ $formatDate(user.created_at) }}</td>
    <td>
      <div v-if="props.showRefCount && user.team_count > 0" class="badge badge-secondary">
        {{ user.team_count }}
      </div>
    </td>
  </tr>
</template>