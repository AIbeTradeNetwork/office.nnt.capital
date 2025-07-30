<script lang="ts">
  export default {
    name: 'LayoutAuth',
  }
</script>

<script setup lang="ts">
  import Background from 'components/background/index.vue'
  import Logo from 'assets/logos/logo.vue'
  import LocaleSelector from 'components/localeSelector/index.vue'
</script>

<template>
  <Background />

  <div
    class="layout-auth w-full min-h-lvh relative flex flex-col items-center overflow-hidden overflow-y-auto scroll-smooth z-10"
  >
    <div class="m-auto pc:min-w-96">
      <router-view v-slot="{ Component }">
        <div class="flex justify-center">
          <Logo class="h-[100px]" />
        </div>
        <div class="relative -mb-[3.6rem] mr-[.6rem] z-10 flex justify-end">
          <LocaleSelector />
        </div>
        <component :is="Component" />
      </router-view>
    </div>
  </div>
</template>
