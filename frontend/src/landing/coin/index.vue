<script lang="ts">
  export default {
    name: 'FarmUDEX',
  }
</script>

<script setup lang="ts">
  // import Menu from './components/menu.vue'
  // import Timeline from './components/timeline.vue'
  // import Info from './components/info.vue'
  // import Autotrade from './components/autotrade.vue'
  // import bg from './assets/bg.png'
  // import HowFarm from './components/howFarm.vue'
  // import ModalShop from './components/modalShop/index.vue'

  import { $modals } from 'utils/modals'
  import { computed, defineAsyncComponent } from 'vue'
  import { Router } from 'routes/index'

  import('./coin.css')

  const Home = defineAsyncComponent(() => import('./components/home.vue'))
  const Autotrade = defineAsyncComponent(() => import('./components/autotrade.vue'))

  const ModalRanks = defineAsyncComponent(() => import('components/modals/modalRankInfo/index.vue'))
  const ModalCoinInfo = defineAsyncComponent(() => import('./components/modalCoinInfo/index.vue'))
  const ModalSubscription = defineAsyncComponent(
    () => import('./components/modalSubscription/index.vue'),
  )
  const activeComponent = computed(() => {
    if (Router.currentRoute.value.name === 'FarmUDEX') return Home
    if (Router.currentRoute.value.name === 'Autotrade') return Autotrade
    return undefined
  })

  // import { EMenu } from './enums'
  // import { computed } from 'vue'
  // import { menu } from './utils'

  // const activeComponent = computed(() => {
  //   if (menu.active.value === EMenu.timeline) return Timeline
  //   if (menu.active.value === EMenu.info) return Info
  //   return Home
  // })
</script>

<template>
  <div
    class="relative flex justify-center items-stretch coin z-[3] mobile:min-h-[calc(100dvh-75px)] pc:min-h-[calc(100dvh-70px)]"
  >
    <div class="flex justify-center items-stretch max-w-[620px] w-full">
      <component :is="activeComponent"></component>
    </div>
  </div>

  <!-- <ModalShop v-if="$modals.coinShop.active" /> -->
  <ModalCoinInfo v-if="$modals.farmingInfo.active" />
  <ModalSubscription v-if="$modals.farmingSubscription.active" />
  <ModalRanks v-if="$modals.ranks.active" />
</template>
