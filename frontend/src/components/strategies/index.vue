<script lang="ts">
  export default {
    name: 'Strategies',
  }
</script>

<script setup lang="ts">
  import { $store } from 'utils/store'
  import Strategy from './strategy.vue'
  import { EStrategyType } from 'types/enums'
  import { computed } from 'vue'

  const { type } = defineProps<{
    type: EStrategyType
  }>()

  // const soonStrategies: Strategy[] = []

  const soonStrategiesIco: Strategy[] = [
    {
      code: 'ICO All',
      name: 'ICO all coins',
      start_at: new Date('2024/03/16 15:40').getTime(),
      symbol_code: 'All coins',
      type: 'ico',
      category: 'spot',
      exchange_code: 'bybit',
      exchange_codes: [],
      pos_profit: 50000,
      risk_level: 'high',
      min_deposit: 0,
      max_deposit: 0,
      symbol_codes: [],
      is_new: false,
      share_profit: 0,
      is_free: false,
      min_leverage: 0,
      max_leverage: 0,
      fix_leverage: 0,
      is_reinvest: false,
      is_trade_limit: false,
      is_trade_percent: false,
    },
  ]

  const list = computed(() => {
    return $store.get('strategies').filter((item) => item.type === type)
  })

  defineEmits(['selected'])
</script>

<template>
  <div class="grid mobile:grid-cols-1 pc:grid-cols-3 mobile:gap-2 pc:gap-4">
    <template v-if="type === EStrategyType.ico">
      <Strategy :data="soonStrategiesIco[0]" :soon="true" @selected="$emit('selected')" />
    </template>

    <template v-for="item in list" :key="item.code">
      <Strategy :data="item" @selected="$emit('selected')" />
    </template>

    <!-- <template v-if="type === EStrategyType.classic">
      <Strategy :data="soonStrategies[0]" :soon="true" @selected="$emit('selected')" />
      <Strategy :data="soonStrategies[1]" :soon="true" @selected="$emit('selected')" />
    </template> -->
  </div>
</template>
