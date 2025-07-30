<script lang="ts">
  export default {
    name: 'StartegiesIcoPage',
  }
</script>

<script setup lang="ts">
  import { EAlerts, ESizes, EStrategyType } from 'types/enums'
  import Info from 'components/info/index.vue'
  import StrategiesList from 'components/strategies/index.vue'
  import { onMounted, ref } from 'vue'
  import { $store } from 'utils/store'

  onMounted(async () => {
    $store.updateStrategies(EStrategyType.ico)
  })
</script>

<template>
  <div class="page">
    <Info
      :type="EAlerts.BASE"
      :size="ESizes.BASE"
      :title="$t('info.strategies_ico.title')"
      :message="$t('info.strategies_ico.text')"
    />

    <StrategiesList :type="EStrategyType.ico" />
  </div>
</template>
