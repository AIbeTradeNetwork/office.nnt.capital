<script lang="ts">
  export default {
    name: 'BranchContent',
  }
</script>

<script setup lang="ts">
  import { $formatDate } from 'utils/date'
  import { $screen } from 'utils/screen'
  import { inject } from 'vue'
  import Avatar from 'vue-boring-avatars'
  const { user } = defineProps<{
    user: TeamUser
  }>()

  const loadData = inject('hierarchicalTreeLoadData')

  function show() {
    if (typeof loadData === 'function') {
      loadData(user)
    }
  }
</script>

<template>
  <div
    class="branch-content cursor-pointer"
    :class="`${!user ? 'opacity-30 pointer-events-none' : ''}`"
    @click="$screen.isPc.value ? show() : () => {}"
  >
    <div
      v-if="user"
      class="info absolute -top-1 bottom-0 -left-1 rounded-tl-full rounded-bl-full rounded-br-[999px] rounded-tr-[999px] bg-base-100 hidden pl-[110px] pr-4 text-left z-[2] text-xs"
    >
      <div class="py-2 whitespace-nowrap flex flex-col justify-center h-full">
        <div>
          <div class="font-bold">{{ user?.nickname }}</div>
          <div>
            <span class="opacity-50">{{ $t('Plan') }}:</span>
            {{ user?.plan?.code || '-' }}
            <template v-if="user?.plan">
              {{ $t('before') }} {{ $formatDate(user.plan.end_at) }}
            </template>
          </div>
          <div>
            <span class="opacity-50">{{ $t('Rank') }}:</span>
            <span>{{ user?.rank?.code || '-' }}</span>
            <template v-if="user?.rank">
              {{ $t('before') }} {{ $formatDate(user?.rank?.end_at) }}
            </template>
          </div>
          <div>
            <span class="opacity-50">{{ $t('Activity') }}:&nbsp;</span>
            <template v-if="user?.activity">
              {{ $t('before') }} {{ $formatDate(user.activity.end_at) }}
            </template>
            <template v-else>-</template>
          </div>
          <div
            v-if="$screen.isMobile.value"
            class="cursor-pointer mt-1 font-bold text-success"
            @click="show"
          >
            {{ $t('Go') }}
          </div>
        </div>
      </div>
    </div>

    <div class="avatar w-[6em] h-[6em] rounded-full">
      <Avatar class="relative w-full h-full" variant="marble" :name="user ? user?.uid : '3'" />
    </div>

    <div
      v-if="user"
      class="absolute bottom-0 left-0 w-full badge badge-neutral text-xs z-[4] overflow-hidden pl-1 whitespace-nowrap"
    >
      {{ user.nickname }}
    </div>
  </div>
</template>
