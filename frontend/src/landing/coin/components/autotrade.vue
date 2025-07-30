<script setup lang="ts">
  import Menu from './menu.vue'

  import { BrokerTransactionType, ECurrencies } from 'types/enums'
  import { $historBack } from 'utils/history'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { $notify } from 'utils/notify'
  import { $requests } from 'queries/index'
  import { $t } from 'i18n/index'
  import { $formatDate } from 'utils/date'
  import { $modals } from 'utils/modals'
  import { $store } from 'utils/store'
  import { $formatPriceByLocale } from 'utils/formats'
  import { $config } from 'utils/configuration'

  const dataBalances = {
    loading: ref(true),
    balance: computed(() => $store.get('brokerBalance')),
  }

  const dataTransactions = reactive({
    loading: true,
    list: [] as Partial<Broker.Transaction>[],
    cols: 3,
    isEndList: false,
    limit: 8,
    page: 0,
    type: BrokerTransactionType.all,
    types: BrokerTransactionType,
  })

  const informations = computed(() => {
    return {
      balance: {
        title: $t('autotrade.AccountBalance'),
        amount: Number(dataBalances.balance.value.amount || 0),
        currency: ECurrencies.usd,
        plus: dataBalances.balance.value.plus,
        minus: dataBalances.balance.value.minus,
      },
      monthly: {
        title: $t('autotrade.MonthlyProfit'),
        amount: 345.92,
        currency: 'USD',
        inclined: 'error',
        percent: 1.3,
      },
      profit: {
        title: $t('autotrade.AllTimeProfit'),
        amount: 345.92,
        currency: 'USD',
        inclined: 'success',
        percent: 1.3,
      },
      loss: {
        title: $t('autotrade.AllTimeLoss'),
        amount: 345.92,
        currency: 'USD',
        inclined: 'error',
        percent: 1.3,
      },
    }
  })

  function setFilter(type: BrokerTransactionType) {
    dataTransactions.type = type
    getTransactions(true)
  }

  async function getTransactions(clear?: boolean) {
    dataTransactions.loading = true

    if (clear) {
      dataTransactions.page = 0
      dataTransactions.list = []
    }

    try {
      await $requests.broker.config()

      const response = await $requests.broker.transactions({
        type: dataTransactions.type || undefined,
        limit: dataTransactions.limit,
        skip: dataTransactions.page * dataTransactions.limit,
        currency: ECurrencies.usdHidden,
      })

      dataTransactions.list = [...dataTransactions.list, ...response]

      if (response && (!response.length || response.length < dataTransactions.limit)) {
        dataTransactions.isEndList = true
      } else {
        dataTransactions.isEndList = false
        dataTransactions.page++
      }
    } catch (error) {
      $notify.show({
        error: error,
      })
    }

    dataTransactions.loading = false

    return 1
  }

  function addingFunds() {
    $modals.goToPro.show()
    // $modals.brokerAddingFunds.show()
  }

  function payout() {
    $modals.brokerPayout.show()
  }

  function showInfo() {
    $modals.info.show({
      title: 'ABTBits',
      text: $t('autotrade.info.text'),
    })
  }

  onMounted(async () => {
    await $store.updateBrokerBalance()
    getTransactions()
  })
</script>

<template>
  <div class="flex flex-col justify-center items-center w-full">
    <div
      class="relative w-full h-full flex flex-col justify-between mobile:p-2 mobile:pt-6 pc:p-4 mobile:pb-[20px] gap-10 overflow-hidden"
    >
      <div class="flex justify-between">
        <div>
          <div class="text-xl">{{ $t('autotrade.title') }}</div>
        </div>

        <button class="btn btn-outline btn-sm" @click="$historBack">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="relative -top-[1px] size-4 -ml-1 -mr-1"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M9 15 3 9m0 0 6-6M3 9h12a6 6 0 0 1 0 12h-3"
            />
          </svg>
          {{ $t('Back') }}
        </button>
      </div>

      <div class="flex flex-col items-center gap-2">
        <div>{{ $t('autotrade.info.title') }}</div>
        <button class="btn btn-outline btn-info inline btn-sm" @click="showInfo">
          {{ $t('autotrade.info.more') }}
        </button>
      </div>

      <div class="text-center">
        <router-link
          :to="{ name: 'AutotradePRO' }"
          class="btn text-white btn-sm bg-gold/80 pc:hover:bg-gold/40 mobile:active:bg-gold/40"
        >
          {{ $t('autotrade_pro.go') }}
        </router-link>
      </div>

      <div>
        <div class="flex justify-center flex-wrap p-0 overflow-hidden">
          <div class="flex-[0_0_50%] relative text-center py-6 pt-0">
            <div class="coin-text-gradient text-2xl">{{ informations.balance.title }}</div>
            <div class="text-lg space-x-1">
              <span class="font-bold text-4xl">{{ informations.balance.amount }}</span>
              <span class="text-sm">{{ informations.balance.currency }}</span>
            </div>
            <div class="flex justify-center items-center gap-2">
              <div class="inline-flex items-center text-sm text-success">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="size-4 mr-1"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M2.25 18 9 11.25l4.306 4.306a11.95 11.95 0 0 1 5.814-5.518l2.74-1.22m0 0-5.94-2.281m5.94 2.28-2.28 5.941"
                  />
                </svg>
                +{{ informations.balance.plus }}
              </div>
              <div class="inline-flex items-center text-sm text-error">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="size-4 mr-1"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M2.25 6 9 12.75l4.286-4.286a11.948 11.948 0 0 1 4.306 6.43l.776 2.898m0 0 3.182-5.511m-3.182 5.51-5.511-3.181"
                  />
                </svg>
                {{ informations.balance.minus }}
              </div>
            </div>
          </div>
        </div>
        <!-- <div class="box flex justify-center flex-wrap p-0 overflow-hidden">
          <div
            v-for="(k, index) in Object.keys(informations)"
            :key="informations[k].title"
            class="flex-[0_0_50%] relative text-center py-6"
            :class="index > 1 ? 'dark:bg-black/50' : ''"
          >
            <div class="coin-text-gradient leading-5">{{ informations[k].title }}</div>
            <div class="text-lg space-x-1">
              <span class="text-xl font-bold">{{ informations[k].count }}</span>
              <span class="text-sm">{{ informations[k].currency }}</span>
            </div>
            <div
              class="inline-flex items-center text-sm"
              :class="`${
                informations[k].inclined === 'success' ? 'text-success'
                : informations[k].inclined === 'error' ? 'text-error'
                : ''
              }`"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-4 mr-1"
              >
                <path
                  v-if="informations[k].inclined === 'success'"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M2.25 18 9 11.25l4.306 4.306a11.95 11.95 0 0 1 5.814-5.518l2.74-1.22m0 0-5.94-2.281m5.94 2.28-2.28 5.941"
                />
                <path
                  v-if="informations[k].inclined === 'error'"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M2.25 6 9 12.75l4.286-4.286a11.948 11.948 0 0 1 4.306 6.43l.776 2.898m0 0 3.182-5.511m-3.182 5.51-5.511-3.181"
                />
              </svg>
              -0.5%
            </div>
          </div>
        </div> -->

        <div class="space-y-2 mt-4">
          <div class="flex gap-2 justify-center">
            <button class="flex-1 btn btn-outline btn-success btn-md" @click="addingFunds">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-6"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15M12 9l-3 3m0 0 3 3m-3-3h12.75"
                ></path>
              </svg>
              {{ $t('AddingFunds') }}
            </button>
            <button class="flex-1 btn btn-outline btn-md" @click="payout">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-6"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M8.25 9V5.25A2.25 2.25 0 0 1 10.5 3h6a2.25 2.25 0 0 1 2.25 2.25v13.5A2.25 2.25 0 0 1 16.5 21h-6a2.25 2.25 0 0 1-2.25-2.25V15m-3 0-3-3m0 0 3-3m-3 3H15"
                ></path>
              </svg>
              {{ $t('Withdraw') }}
            </button>
          </div>
        </div>
      </div>

      <div class="flex-1 flex flex-col justify-end">
        <div class="flex justify-end items-center gap-2 mb-2">
          <button
            v-for="type in dataTransactions.types"
            :key="type"
            class="btn btn-xs"
            :class="`${dataTransactions.type === type ? 'btn-primary' : 'btn-outline'}`"
            @click="setFilter(type)"
          >
            {{ $t(`autotrade.transactiosType.filter.${type || 'all'}`) }}
          </button>
          <button class="btn btn-xs btn-square" @click="getTransactions(true)">
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
                d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99"
              />
            </svg>
          </button>
        </div>

        <div class="box p-0 overflow-x-auto flex-1">
          <table class="table table-sm table-zebra">
            <thead>
              <tr>
                <th>{{ $t('Date') }}</th>
                <th>{{ $t('Type') }}</th>
                <th>{{ $t('Amount') }}</th>
              </tr>
            </thead>

            <tbody>
              <tr v-for="item in dataTransactions.list" :key="item.uid">
                <td>{{ $formatDate(item.created_at, 'Pp') }}</td>
                <td>
                  <div class="badge badge-outline">
                    {{ $t(`autotrade.transactiosType.filter.${item.type}`) }}
                  </div>
                </td>
                <td>
                  {{
                    $formatPriceByLocale({ count: Number(item.amount), currency: item.currency })
                  }}
                </td>
              </tr>

              <tr v-if="dataTransactions.loading" colspan="">
                <td :colspan="dataTransactions.cols">
                  <div class="text-center font-bold p-4">{{ $t('Loading') }}</div>
                </td>
              </tr>

              <tr v-else-if="!dataTransactions.list.length">
                <td :colspan="dataTransactions.cols">
                  <div class="text-center font-bold p-4">{{ $t('ListEmpty') }}</div>
                </td>
              </tr>

              <tr v-if="!dataTransactions.loading">
                <td :colspan="dataTransactions.cols">
                  <div class="flex justify-center">
                    <button
                      class="btn btn-sm btn-outline btn-primary"
                      :disabled="dataTransactions.isEndList"
                      @click="getTransactions(false)"
                    >
                      {{ $t('ShowMore') }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <Menu />
  </div>
</template>
