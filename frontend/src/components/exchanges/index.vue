<script lang="ts">
  export default {
    name: 'Exchanges',
  }
</script>

<script setup lang="ts">
  import Info from 'components/info/index.vue'
  import { EAlerts } from 'types/enums'
  import Exchange from './exchange.vue'
  import { $store } from 'utils/store'
  import { onMounted } from 'vue'

  onMounted(async () => {
    try {
      $store.updateExchanges()
      $store.updateKeys()
    } catch (error) {
      console.error(error)
    }
  })
</script>

<template>
  <div class="flex gap-4">
    <Info
      :type="EAlerts.BASE"
      :title="$t('info.pageExchangesIP.title')"
      :message="$t('info.pageExchangesIP.text')"
      :copytext="$t('info.pageExchangesIP.IPS')"
    />
  </div>

  <div class="grid mobile:grid-cols-1 pc:grid-cols-3 mobile:gap-2 pc:gap-4">
    <!-- <button
      class="btn h-full box flex flex-col justify-center items-center"
      @click="showModalExchange"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-12 h-12"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
        />
      </svg>
      <div>{{ $t('AddExchange') }}</div>
    </button> -->
    <template v-for="item in $store.get('exchanges')" :key="item.code">
      <Exchange :data="item" />
    </template>
  </div>
</template>
