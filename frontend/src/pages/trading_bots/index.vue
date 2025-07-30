<script lang="ts">
  export default {
    name: 'PageTradingBots',
  }
</script>

<script setup lang="ts">
  import Trades from './trades.vue'
  import { EBotTradeType, EStrategyCategory } from 'types/enums'
  import { $modals } from 'utils/modals'
  import { $store } from 'utils/store'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { $t } from 'i18n/index'
  import { $formatInt } from 'utils/formats'
  import { getColorSign, getSign } from 'utils/sign'

  const colspan = 10
  const loading = ref(true)
  const strategies = computed(() => $store.get('strategies'))
  const bots = computed(() => $store.get('bots'))
  const trades: { [key: Bot['uid']]: 'loaded' | 'loading' } = reactive({})

  function showStrategiesList() {
    $modals.strategies.show()
  }

  function edit(bot: Bot) {
    $modals.bot.show(bot)
  }

  function closeBotTrades(uid: Bot['uid']) {
    delete trades[uid]
  }

  function toggle(bot: Bot) {
    if (trades[bot.uid] === 'loading') return

    if (trades[bot.uid] === 'loaded') {
      closeBotTrades(bot.uid)
      return
    }

    trades[bot.uid] = 'loading'
  }

  function getStrategy(bot: Bot) {
    return strategies.value.find((item) => item.code === bot.strategy_code)
  }

  function showErrorInfo(bot: Bot) {
    $modals.error.show({
      title: `${$t('Error')} / ${bot.error_code}`,
      text: `${$t(`error.${bot.error_code}`)} ${bot.error}`,
    })
  }

  onMounted(async () => {
    loading.value = true
    await $store.updateStrategies()
    await $store.updateBots()
    loading.value = false
  })
</script>

<template>
  <div class="page">
    <div>
      <button class="btn btn-primary" @click="showStrategiesList">
        <span>{{ $t('AddBot') }}</span>

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
            d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
          />
        </svg>
      </button>
    </div>

    <div class="box p-0 overflow-x-auto">
      <table class="table table-zebra table-pin-cols">
        <!-- head -->
        <thead>
          <tr>
            <!-- <td>{{ $t('Bot') }}</td> -->
            <!-- <th>{{ $t('Cryptocurrency') }}</th> -->
            <td>{{ $t('Exchange') }}</td>
            <!-- <td>{{ $t('Key') }}</td> -->
            <td>{{ $t('Name') }}</td>
            <td>{{ $t('Type') }}</td>
            <td>{{ $t('Category') }}</td>
            <!-- <td>{{ $t('Currency') }}</td> -->
            <!-- <th>{{ $t('Bid', { currency: $config.currency }) }}</th> -->
            <!-- <th>{{ $t('Limit') }}</th> -->
            <!-- <th>{{ $t('Profit') }}</th> -->
            <td>{{ $t('Settings') }}</td>

            <td>
              {{ $t('Profit') }}:
              <div class="italic">({{ $t('General') }})</div>
            </td>
            <td>
              {{ $t('Profit') }}:
              <div class="italic">({{ $t('Month', 1) }})</div>
            </td>
            <td>
              {{ $t('Profit') }}:
              <div class="italic">({{ $t('Prev') }})</div>
            </td>

            <td>{{ $t('Activity') }}</td>
            <th></th>
          </tr>
        </thead>

        <tbody v-if="!bots.length">
          <tr>
            <td :colspan="colspan" class="pc:text-center font-bold py-6">{{ $t('ListEmpty') }}</td>
          </tr>
        </tbody>

        <template v-for="bot in bots" v-else :key="bot.uid">
          <tbody>
            <tr>
              <!-- <td>{{ bot.uid }}</td> -->

              <td>{{ bot.exchange_code }}</td>
              <!-- <td>{{ bot.key_uid }}</td> -->
              <td>
                {{ getStrategy(bot)?.name }}
              </td>
              <td>
                {{ getStrategy(bot)?.type }}
              </td>
              <td>
                {{ getStrategy(bot)?.category }}
                <template
                  v-if="
                    getStrategy(bot)?.category === EStrategyCategory.futures && bot.leverage > -1
                  "
                >
                  x{{ bot.leverage > 0 ? bot.leverage : 1 }}
                </template>
              </td>
              <!-- <td>{{ bot.symbol_code || 'All coins' }}</td> -->
              <td>
                <template v-if="bot.trade_type === EBotTradeType.percent">
                  {{ $t('Percent') }}: {{ bot.trade_percent }}%
                </template>
                <template v-if="bot.trade_type === EBotTradeType.limit">
                  {{ $t('Limit') }}: {{ $formatInt(bot.trade_limit) }}
                  <template v-if="bot.trade_reinvest">&nbsp;/&nbsp;{{ $t('Reinvest') }}</template>
                </template>
              </td>

              <td :class="getColorSign(bot.profit_all)">
                {{ getSign(bot.profit_all) }}{{ parseFloat(bot.profit_all + '').toFixed(3) }}
              </td>
              <td :class="getColorSign(bot.profit_month)">
                {{ getSign(bot.profit_month) }}{{ parseFloat(bot.profit_month + '').toFixed(3) }}
              </td>
              <td :class="getColorSign(bot.profit_month_prev)">
                {{ getSign(bot.profit_month_prev)
                }}{{ parseFloat(bot.profit_month_prev + '').toFixed(3) }}
              </td>

              <td>
                <div class="flex items-center gap-2">
                  <span
                    class="label-text"
                    :class="`${bot.is_active ? 'text-success' : 'text-error'}`"
                  >
                    {{ bot.is_active ? $t('On') : $t('Off') }}
                  </span>
                  <button
                    v-if="bot.error"
                    class="btn btn-xs btn-error btn-circle btn-outline"
                    @click="showErrorInfo(bot)"
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
                </div>
              </td>

              <th class="w-[42px] p-0 pr-2">
                <div class="flex gap-2">
                  <button
                    class="btn btn-sm btn-primary btn-circle btn-outline"
                    :disabled="trades[bot.uid] === 'loading'"
                    @click="toggle(bot)"
                  >
                    <svg
                      v-if="!trades[bot.uid]"
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
                        d="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 0 1 0 3.75H5.625a1.875 1.875 0 0 1 0-3.75Z"
                      />
                    </svg>
                    <svg
                      v-if="trades[bot.uid] === 'loading'"
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
                        d="M2.25 15a4.5 4.5 0 0 0 4.5 4.5H18a3.75 3.75 0 0 0 1.332-7.257 3 3 0 0 0-3.758-3.848 5.25 5.25 0 0 0-10.233 2.33A4.502 4.502 0 0 0 2.25 15Z"
                      />
                    </svg>

                    <svg
                      v-if="trades[bot.uid] === 'loaded'"
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
                        d="m4.5 18.75 7.5-7.5 7.5 7.5"
                      />
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="m4.5 12.75 7.5-7.5 7.5 7.5"
                      />
                    </svg>
                  </button>

                  <button class="btn btn-sm btn-accent btn-outline btn-circle" @click="edit(bot)">
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
                        d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10"
                      />
                    </svg>
                  </button>
                </div>
              </th>
            </tr>
          </tbody>

          <tbody v-if="trades[bot.uid]">
            <tr>
              <td :colspan="colspan" class="p-0">
                <Trades :bot="bot" @loaded="trades[bot.uid] = 'loaded'" />
              </td>
            </tr>
          </tbody>
        </template>
      </table>
    </div>

    <!-- <div class="flex justify-start">
      <div class="join">
        <button class="join-item btn">1</button>
        <button class="join-item btn btn-active">2</button>
        <button class="join-item btn">3</button>
        <button class="join-item btn">4</button>
      </div>
    </div> -->

    <div class="grid mobile:gap-2 pc:gap-4 mobile:grid-cols-1 pc:grid-cols-3">
      <!-- <Info
        :type="EAlerts.BASE"
        :size="ESizes.SM"
        :title="$t('info.tradingBots.title')"
        :message="$t('info.tradingBots.text')"
      /> -->

      <!-- <Info
        :type="EAlerts.WARNING"
        :size="ESizes.SM"
        :title="$t('info.warning.tradingBots.title')"
        :message="$t('info.warning.tradingBots.text')"
      /> -->
    </div>
  </div>
</template>
