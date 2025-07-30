<script lang="ts">
  export default {
    name: 'StartegiesPage',
  }
</script>

<script setup lang="ts">
  import StrategiesList from 'components/strategies/index.vue'
  import { EStrategyType } from 'types/enums'
  import { $store } from 'utils/store'
  import { onMounted } from 'vue'

  onMounted(async () => {
    $store.updateStrategies(EStrategyType.classic)
  })
</script>

<template>
  <div class="page">
    <StrategiesList :type="EStrategyType.classic" />
  </div>
</template>
