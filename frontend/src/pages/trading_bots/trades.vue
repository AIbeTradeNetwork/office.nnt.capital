<script lang="ts">
  export default {
    name: 'BotTrades',
  }
</script>

<script setup lang="ts">
  import { $requests } from 'queries/index'
  import { $formatDate } from 'utils/date'
  import { onMounted, ref } from 'vue'
  // import { profitability } from './profitability'
  // import { $notify } from 'utils/notify'
  import { $t } from 'i18n/index'
  import { $modals } from 'utils/modals'
  import { ECurrencies, ETradeSide } from 'types/enums'
  import { getColorSign, getSign } from 'utils/sign'

  const { bot } = defineProps<{
    bot: Bot
  }>()

  const emit = defineEmits<{
    (e: 'loaded'): void
  }>()

  const colspan = 9
  const trades = ref<Trade[]>([])
  const loading = ref(true)
  const limit = 10
  const page = ref(0)
  const isEndList = ref(true)

  async function loadTraders() {
    if (trades.value.length && loading.value) return

    loading.value = true

    try {
      const response = await $requests.bots.traders(bot.uid, { limit, skip: page.value * limit })

      if (response.length < limit) {
        isEndList.value = true
      } else {
        isEndList.value = false
        page.value++
      }

      response.forEach((item) => trades.value.push(item))

      emit('loaded')
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  function showErrorInfo(trade: Trade) {
    $modals.error.show({
      title: `${$t('Error')} / ${trade.error_code}`,
      text: `${$t(`error.${trade.error_code}`)} ${trade.error}`,
    })
  }

  onMounted(() => {
    loadTraders()
  })
</script>

<template>
  <table class="bg-base-300 table table-zebra table-pin-cols rounded-t-none">
    <thead>
      <tr>
        <!-- <td>ID</td> -->
        <!-- <td>{{ $t('SignalUid') }}</td> -->
        <td>{{ $t('SymbolCode') }}</td>
        <td>{{ $t('Quantity') }}</td>
        <td>{{ $t('PriceOpen') }}</td>
        <td>{{ $t('PriceClose') }}</td>
        <td>{{ $t('Profitability') }}: {{ ECurrencies.usdt }}</td>
        <!-- <td>{{ $t('TakeProfit') }}</td> -->
        <!-- <td>{{ $t('StopLoss') }}</td> -->
        <td>{{ $t('CreatedAt') }}</td>
        <td>{{ $t('ClosedAt') }}</td>
        <td class="p-0"></td>
      </tr>
    </thead>

    <tbody v-if="trades.length">
      <tr v-for="trade in trades" :key="trade.uid">
        <!-- <td>{{ trade.uid }}</td> -->
        <!-- <td>{{ trade.signal_uid }}</td> -->
        <td>
          <div class="flex items-center">
            <span
              v-if="trade.side === ETradeSide.buy"
              class="badge badge-outline badge-sm font-mono badge-success"
            >
              {{ trade.side }}
            </span>
            <span
              v-if="trade.side === ETradeSide.sell"
              class="badge badge-outline badge-sm font-mono badge-error"
            >
              {{ trade.side }}
            </span>
            <span class="ml-1">{{ trade.symbol_code || 'All coins' }}</span>
          </div>
        </td>
        <td>{{ trade.qty }}</td>
        <td>{{ parseFloat(trade.price_open + '').toFixed(3) }}</td>
        <td>
          {{ trade.price_close > 0 ? parseFloat(trade.price_close + '').toFixed(3) : '...' }}
        </td>
        <!-- <td :class="profitability.styles(trade)">{{ profitability.value(trade) }}%</td> -->
        <td>
          <span v-if="trade.profit > 0 || trade.profit < 0" :class="getColorSign(trade.profit)">
            {{ getSign(trade.profit) }}{{ parseFloat(trade.profit + '').toFixed(3) }}
          </span>
          <span v-else>{{ '...' }}</span>
        </td>
        <!-- <td>{{ trade.take_profit }}</td> -->
        <!-- <td>{{ trade.stop_loss }}</td> -->
        <td>{{ $formatDate(trade.created_at, 'date|time') }}</td>
        <td>
          {{ trade.closed_at > 0 ? $formatDate(trade.closed_at, 'date|time') : '...' }}
        </td>
        <td class="p-0">
          <button
            v-if="trade.error"
            class="btn btn-xs btn-error btn-circle btn-outline"
            @click="showErrorInfo(trade)"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-4"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"
              />
            </svg>
          </button>
        </td>
      </tr>
      <tr v-if="!isEndList">
        <td :colspan="colspan">
          <div class="flex justify-center">
            <button class="btn btn-primary" :disabled="loading" @click="loadTraders">
              {{ $t('ShowMore') }}
            </button>
          </div>
        </td>
      </tr>
    </tbody>

    <tbody v-else-if="!loading">
      <tr>
        <td :colspan="colspan" class="pc:text-center font-bold py-6">{{ $t('ListEmpty') }}</td>
      </tr>
    </tbody>

    <tbody v-if="loading">
      <tr>
        <td :colspan="colspan" class="text-center font-bold py-6">{{ $t('Loading') }}</td>
      </tr>
    </tbody>
  </table>
</template>
