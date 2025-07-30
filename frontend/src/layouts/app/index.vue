<script lang="ts">
  export default {
    name: 'LayoutApp',
  }
</script>

<script setup lang="ts">
  import Background from 'components/background/index.vue'
  import Menu from 'components/menu/index.vue'
  import Navbar from 'components/navbar/index.vue'
  import { $requests } from 'queries/index'
  import { $me } from 'utils/me'
  import { $store } from 'utils/store'
  import { ref, onMounted } from 'vue'

  const loading = ref(true)

  async function loadMainData() {
    loading.value = true
    try {
      $store.updateCfg()
      await $requests.me.getRole()
      await $me.update()
      $store.updateCurrencies()
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  onMounted(() => {
    loadMainData()
  })
</script>

<template>
  <Background />

  <div
    v-if="loading"
    class="fixed top-0 left-0 w-full h-full z-[100] flex items-center justify-center"
  >
    <span class="loading loading-spinner w-36"></span>
  </div>

  <div v-else class="drawer pc:drawer-open w-full mx-auto max-w-[3840px]">
    <input id="drawer" type="checkbox" class="drawer-toggle" />

    <div class="drawer-content flex flex-col z-[1]">
      <Navbar />

      <div class="flex-1">
        <router-view v-slot="{ Component }">
          <component :is="Component" />
        </router-view>
      </div>
    </div>

    <Menu />
  </div>
</template>
