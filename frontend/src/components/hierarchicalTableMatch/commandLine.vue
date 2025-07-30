<script lang="ts">
  export default {
    name: 'CommandLine',
  }
</script>

<script setup lang="ts">
  import { ref } from 'vue'
  import CommandLine from './commandLine.vue'
  import { toggle } from './tableToggler'
  import { $formatDate } from 'utils/date'
  import Avatar from 'vue-boring-avatars'
  import { $requests } from 'queries/index'

  const props = defineProps<{
    user: TeamUser
    hiddenClass?: string
    lineid: number
  }>()

  const loading = ref(false)
  const isHasChildrens = props.user?.team_count > 0
  const isShow = ref(false)
  const childrens = ref<TeamUser[]>([])

  async function loadData() {
    loading.value = true

    try {
      childrens.value = await $requests.teams.getMatch({
        user_uid: props.user.uid,
        // $me.data.uid
      })
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  async function doToggle(event) {
    if (!isHasChildrens) return

    if (!childrens.value.length) {
      await loadData()
    }

    toggle(event)
  }
</script>

<template>
  <tr :id="lineid + '' || '0'" :class="props.hiddenClass">
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
          <div class="bg-neutral text-neutral-content rounded-full w-12">
            <Avatar class="w-full h-full" variant="marble" :name="user.uid" />
          </div>
        </div>
        <div>
          <p>{{ user.nickname }}</p>
          <p>{{ user.email }}</p>
        </div>
      </div>
    </td>
    <td>
      <div v-if="user?.plan">
        <b>{{ user.plan.code }}</b>
        {{ $t('before') }}: {{ $formatDate(user.plan.end_at) }}
      </div>
      <div v-else>-</div>
      <!-- <span class="badge badge-outline text-brilliant">Brilliant</span> -->
    </td>
    <td>
      <div v-if="user?.rank">
        <b>{{ user.rank.code }}</b>
        {{ $t('before') }}: {{ $formatDate(user.rank.end_at) }}
      </div>
      <div v-else>-</div>
    </td>
    <td>
      <div v-if="user?.activity">{{ $t('before') }}: {{ $formatDate(user.activity.end_at) }}</div>
      <div v-else>-</div>
    </td>
  </tr>

  <CommandLine
    v-for="item in childrens"
    :key="item.uid"
    :lineid="lineid + 1"
    :hidden-class="isShow ? '' : 'hidden'"
    :user="item"
  />
</template>
