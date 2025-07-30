<script lang="ts">
  export default {
    name: 'InvitedTotal',
  }
</script>

<script setup lang="ts">
  import { $isAUnlimitedInvite } from 'utils/checks'
  import { $me } from 'utils/me'
  import { computed } from 'vue'

  const props = defineProps<{
    isIcon?: boolean,
    color?: string
  }>()

  const max = computed(() => {
    function premiumOrDefault() {
      if ($me.data.is_premium) {
        if ($me.data.premium_invites === 0) return '∞'
        if ($me.data.premium_invites !== 0)
          return `${$me.data?.level?.invite_limit} +${$me.data.premium_invites}`
      }

      return $me.data?.level?.invite_limit || '0'
    }

    return $isAUnlimitedInvite.value ? '∞' : premiumOrDefault()
  })
</script>

<template>
  <div v-if="props.isIcon" class="text-center">
    <svg
      class="inline-block w-10"
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      fill="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        fill-rule="evenodd"
        d="M12 6a3.5 3.5 0 1 0 0 7 3.5 3.5 0 0 0 0-7Zm-1.5 8a4 4 0 0 0-4 4 2 2 0 0 0 2 2h7a2 2 0 0 0 2-2 4 4 0 0 0-4-4h-3Zm6.82-3.096a5.51 5.51 0 0 0-2.797-6.293 3.5 3.5 0 1 1 2.796 6.292ZM19.5 18h.5a2 2 0 0 0 2-2 4 4 0 0 0-4-4h-1.1a5.503 5.503 0 0 1-.471.762A5.998 5.998 0 0 1 19.5 18ZM4 7.5a3.5 3.5 0 0 1 5.477-2.889 5.5 5.5 0 0 0-2.796 6.293A3.501 3.501 0 0 1 4 7.5ZM7.1 12H6a4 4 0 0 0-4 4 2 2 0 0 0 2 2h.5a5.998 5.998 0 0 1 3.071-5.238A5.505 5.505 0 0 1 7.1 12Z"
        clip-rule="evenodd"
      />
    </svg>
    <div>
      {{ $me.data.team_count || '0' }}
    </div>
  </div>
  <div v-else>
    {{ $t('InvitedTotal') }}: {{ $me.data.team_count || '0' }}
  </div>
</template>
