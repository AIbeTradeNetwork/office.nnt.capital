<script setup lang="ts">
  import { $config } from 'utils/configuration'
  import Modals from 'components/modals/index.vue'
  import Notify from 'components/notify/index.vue'
  // import { config } from './config'
  // const cfg = config
</script>

<template>
  <component :is="$route.meta.layout">
    <template v-if="$route.meta.keepAlive">
      <router-view v-slot="{ Component }">
        <KeepAlive v-if="$route.meta.keepAlive">
          <component :is="Component" />
        </KeepAlive>
      </router-view>
    </template>

    <template v-else>
      <router-view v-slot="{ Component }">
        <component :is="Component" />
      </router-view>
    </template>
  </component>

  <Modals />

  <Notify />

  <div v-if="$config.isDevMode" class="fixed top-0 right-0 text-xs badge badge-error z-[100]">
    DEV
  </div>
  <!-- <div class="glass">{{ cfg.version }}</div> -->
</template>
