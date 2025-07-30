<script lang="ts">
  export default {
    name: 'BalanceSimple',
  }
</script>

<script setup lang="ts">
  // import Dropdown from 'components/dropdown/index.vue'
  import { $formatInt, $formatPriceByLocale } from 'utils/formats'
  import { $modals } from 'utils/modals'
  import { ref, onMounted } from 'vue'
  import { coin } from 'queries/claim'

  const udexBalance = ref(0)
  const udexPrecision = ref(9)

  onMounted(async () => {
    const balanceObj = await coin.getBalance()
    udexBalance.value = balanceObj.balance || 0
    udexPrecision.value = balanceObj.precision || 9
  })
</script>

<template>
  <div>
    <div
      class="flex mobile:flex-col items-center pc:gap-2 cursor-pointer mt-1"
      @click="$modals.balance.show"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-6 h-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M21 12a2.25 2.25 0 0 0-2.25-2.25H15a3 3 0 1 1-6 0H5.25A2.25 2.25 0 0 0 3 12m18 0v6a2.25 2.25 0 0 1-2.25 2.25H5.25A2.25 2.25 0 0 1 3 18v-6m18 0V9M3 12V9m18 0a2.25 2.25 0 0 0-2.25-2.25H5.25A2.25 2.25 0 0 0 3 9m18 0V6a2.25 2.25 0 0 0-2.25-2.25H5.25A2.25 2.25 0 0 0 3 6v3"
        />
      </svg>

      <span class="mobile:hidden">
        {{ (udexBalance / Math.pow(10, udexPrecision)).toFixed(udexPrecision) }} UDEX
      </span>

      <div class="pc:hidden text-xs">{{ $t('Balance') }}</div>
    </div>

    <div class="mobile:hidden mobile:mt-4">
      <button
        class="btn btn-link p-0 btn-xs text-success m-0 mr-2"
        @click="$modals.addingFunds.show"
      >
        {{ $t('AddingFunds') }}
      </button>
      <button class="btn btn-link p-0 btn-xs text-info m-0" @click="$modals.payout.show">
        {{ $t('Payout') }}
      </button>
    </div>
  </div>
</template>
