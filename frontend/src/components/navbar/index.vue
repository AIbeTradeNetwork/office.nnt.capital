<script lang="ts">
  export default {
    name: 'NavBar',
  }
</script>

<script setup lang="ts">
  import BalanceNav from 'components/balance/nav.vue'
  import KneeSwitcher from 'components/kneeSwitcher/index.vue'
  import LocaleSelector from 'components/localeSelector/index.vue'
  // import EventBell from 'components/eventBell/index.vue'
  import SettingsSelector from 'components/settingsSelector/index.vue'
  import UserAvatarAndInfo from 'components/userAvatarAndInfo/index.vue'
  import Notifications from 'components/notifications/index.vue'
  import { $screen } from 'utils/screen'
  import { $config } from 'utils/configuration'
  import { ERoles } from 'types/enums'
  // import { Router } from 'routes/index'
  // import { computed } from 'vue'

  // const pageNameTitle = computed(() => {
  //   return Router.currentRoute.value.name
  // })
</script>

<template>
  <div
    class="navbar pc:h-[70px] sticky top-0 flex justify-between items-stretch flex-[0_0_auto] z-10 bg-base-100 backdrop-blur border-b border-white/30 bg-opacity-50 mobile:px-2 pc:px-6"
  >
    <div class="flex items-center safearea-t">
      <span
        class="pc:hidden tooltip tooltip-bottom before:text-xs before:content-[attr(data-tip)]"
        data-tip="Menu"
      >
        <label
          aria-label="Open menu"
          for="drawer"
          class="btn btn-square btn-ghost drawer-button text-lg pc:hidden"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            class="inline-block h-10 w-10 stroke-current mobile:h-6 pc:w-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16"
            ></path>
          </svg>
        </label>
      </span>

      <div class="pc:hidden divider divider-horizontal mx-0 opacity-30"></div>
      <UserAvatarAndInfo />

      <template v-if="$config.role === ERoles.distributor">
        <div
          v-if="$screen.isPc.value"
          class="divider divider-horizontal mobile:hidden opacity-30"
        ></div>
        <KneeSwitcher v-if="$screen.isPc.value" class="mobile:hidden" />
      </template>

      <div class="mobile:hidden divider divider-horizontal opacity-30"></div>
      <BalanceNav class="mobile:hidden" />
    </div>

    <div class="flex items-center gap-2">
      <!-- <EventBell /> -->
      <BalanceNav class="pc:hidden" />
      <Notifications />
      <LocaleSelector v-if="$screen.isPc.value" class="mobile:hidden" />
      <SettingsSelector v-if="$screen.isPc.value" class="mobile:hidden" />
    </div>
  </div>
</template>
