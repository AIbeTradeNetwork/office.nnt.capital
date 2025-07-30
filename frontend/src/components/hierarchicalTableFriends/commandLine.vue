<script lang="ts">
  export default {
    name: 'CommandLine',
  }
</script>

<script setup lang="ts">
  import { ref } from 'vue'
  import CommandLine from './commandLine.vue'
  import { toggle } from './tableToggler'
  import Avatar from 'vue-boring-avatars'
  import { $requests } from 'queries/index'
  import { $formatDate } from 'utils/date'

  const props = defineProps<{
    user: FriendUser
    hiddenClass?: string
    lineid: number
    showRefCount?: boolean
  }>()

  const loading = ref(false)
  const colspan = 3
  const isHasChildrens = props.user?.team_count > 0
  const isShow = ref(false)
  const childrens = ref<FriendUser[]>([])
  const limit = 15
  let page = 0
  const isEnd = ref(true)

  async function getRefs() {
    if (!page || page < 0) page = 0

    loading.value = true

    try {
      childrens.value.push(...(await $requests.friends.list(props.user.uid, limit, limit * page)))

      if (childrens.value.length >= props.user.team_count) {
        isEnd.value = true
      } else {
        isEnd.value = false
      }

      page++
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  async function doToggle(event) {
    if (!isHasChildrens) return

    if (!childrens.value.length) {
      await getRefs()
    }

    toggle(event)
  }
</script>

<template>
  <tr :id="lineid + '' || '0'">
    <td class="relative" :style="`padding-left: calc(1rem + (${lineid / 10}rem * 30))`">
      <div class="flex items-center gap-2">
        <button
          :class="!isHasChildrens ? `invisible pointer-events-none` : ''"
          class="btn btn-square btn-xs hi"
          @click.stop="doToggle"
        >
          <span v-if="loading" class="loading loading-spinner"></span>
          <template v-else>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-4 h-4 inline [.open_&]:hidden"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-6 h-6 hidden [.open_&]:inline"
            >
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

        <div v-if="props.showRefCount && user.team_count > 0" class="badge badge-secondary ml-6">
          {{ user.team_count }}
        </div>
      </div>
    </td>
    <td>{{ $formatDate(user.created_at) }}</td>
    <td :colspan="colspan - 2"></td>
  </tr>

  <CommandLine
    v-for="item in childrens"
    :key="item.uid"
    :lineid="lineid + 1"
    :hidden-class="isShow ? '' : 'hidden'"
    :user="item"
    :show-ref-count="true"
  />

  <tr v-if="!isEnd" :id="lineid + 1 + '' || '0'">
    <td :colspan="colspan" :style="`padding-left: calc(1rem + (${(lineid + 1) / 10}rem * 30))`">
      <div class="">
        <div v-if="loading" class="font-bold p-4">{{ $t('Loading') }}</div>
        <button v-else class="btn btn-primary btn-sm" @click="getRefs">
          {{ $t('ShowMore') }}
        </button>
      </div>
    </td>
  </tr>
</template>
