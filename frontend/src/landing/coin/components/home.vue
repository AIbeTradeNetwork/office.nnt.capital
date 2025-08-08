<script setup lang="ts">
  // import coinSrc from 'coin/assets/coin.png?url'
  import symbolSrc from 'coin/assets/symbol.png?url'
  import InCompanyWarning from './inCompanyWarning.vue'
  // import InTasksWarning from './inTasksWarning.vue'
  import IvitedTotal from 'components/invitedTotal/index.vue'
  import PartnersTotal from 'components/partnersTotal/index.vue'
  import Menu from './menu.vue'

  import { intervalToDuration } from 'date-fns'
  import { coin } from 'queries/claim'
  import {
    computed,
    getCurrentInstance,
    nextTick,
    onBeforeUnmount,
    onMounted,
    ref,
    watch,
  } from 'vue'
  import { drawCoinCanvas, useConffeti, useCounter, useTimeout } from '../utils'
  import { Decimal } from 'decimal.js'
  import { $t } from 'i18n/index'
  import { $modals } from 'utils/modals'
  import { ECoinFarmingActions } from 'types/enums'
  import ABTRank from '../ranks'
  import { Router } from 'routes/index'
  import { $me } from 'utils/me'
  import { $request } from 'utils/axios'

  import '../arrowGame.css'
  import { $config } from 'utils/configuration'

  import { useSafeGame } from '../safeGame'
  import { useArrowGame } from '../arrowGame'
  import { $textSplit } from 'utils/text'
  // import { $copyToClipboard } from 'utils/clipboard'
  import coinClatterSrc from 'coin/assets/coin-clatter-svgrepo-com.svg?url'

  const instance = getCurrentInstance()
  const counter = useCounter(instance)
  const conffeti = useConffeti(instance)
  const timer = useTimeout(instance)

  // const infoList = computed(() => {
  //   return [
  //     // {
  //     //   title: $t('coin.info.start.title'),
  //     //   data: $t('coin.info.start.text'),
  //     // },
  //     {
  //       title: $t('coin.info.max.title'),
  //       data: $t('coin.info.max.text'),
  //     },
  //     // {
  //     //   title: $t('coin.info.presale.title'),
  //     //   data: $t('coin.info.presale.text'),
  //     // },
  //     {
  //       title: $t('coin.info.volume.title'),
  //       data: $t('coin.info.volume.text'),
  //     },
  //     // {
  //     //   title: $t('coin.info.listing.title'),
  //     //   data: $t('coin.info.presale.text'),
  //     // },
  //     // {
  //     //   title: $t('coin.info.comments.title'),
  //     //   data: $t('coin.info.presale.text'),
  //     // },
  //   ]
  // })

  const loading = ref(true)
  const loadingInner = ref(false)
  const now = ref(Date.now())
  const farmedAmount = ref(0)
  const farmedAmountTime = 1000
  const checkBoxPremium = ref($me.data.is_premium)

  const autoFarming = {
    loading: ref(false),
    active: ref($me.data.is_autofarm || false),
  }

  const claimData = {
    claim: ref<Claim>(),
    balance: ref<ClaimBalance>(),
  }

  const isNew = computed(() => (claimData.balance.value.claimed_at > 0 ? false : true))

  const isFarmedFull = computed(() =>
    loading.value ? false : farmedAmount.value >= claimData.claim.value.amount,
  )

  const getMinFarmTime = computed(
    () => (loading.value ? 0 : claimData.claim.value.min_period / 10 ** 9) * 1000,
  )
  const getMaxFarmTime = computed(
    () => (loading.value ? 0 : claimData.claim.value.max_period / 10 ** 9) * 1000,
  )

  const getCoinImg = new URL('../assets/coin.png', import.meta.url).href

  const getCoinTier2TopImg = computed(() => {
    // return new URL('coin/assets/abt_ranks/tier-2/t10.png', import.meta.url).href
    return ABTRank.currentRank.value.img_top
  })

  const getCoinTier2BottomImg = computed(() => {
    // return new URL('coin/assets/abt_ranks/tier-2/b4.png', import.meta.url).href
    return ABTRank.currentRank.value.img_bottom
  })

  const getShortTimeTextValue = (val: number, type: 'h' | 'm' | 's') => {
    const t = {
      h: $t('shorts.hours'),
      m: $t('shorts.minutes'),
      s: $t('shorts.seconds'),
    }

    return val ? `${Math.floor(val)}${t[type]} ` : `00${t[type]} `
  }

  // const getIntervalToDurationMin = computed(() => {
  //   const q = intervalToDuration({
  //     start: now.value,
  //     end: claimData.balance.value.claimed_at + getMinFarmTime.value,
  //   })

  //   // const year = q.years ? `${Math.floor(q.years)}y ` : '00y '
  //   // const month = q.months ? `${Math.floor(q.months)}mo ` : '00mo '
  //   // const day = q.days ? `${Math.floor(q.days)}d ` : '00d '
  //   const hour = getShortTimeTextValue(q.hours, 'h')
  //   const minute = getShortTimeTextValue(q.minutes, 'm')
  //   const second = getShortTimeTextValue(q.seconds, 's')

  //   return `${hour}${minute}${second}`
  // })

  const getIntervalToDurationMax = computed(() => {
    const q = intervalToDuration({
      start: now.value,
      end: claimData.balance.value.claimed_at + getMaxFarmTime.value,
    })

    // const year = q.years ? `${Math.floor(q.years)}y ` : '00y '
    // const month = q.months ? `${Math.floor(q.months)}mo ` : '00mo '
    // const day = q.days ? `${Math.floor(q.days)}d ` : '00d '
    const hour = getShortTimeTextValue(q.hours, 'h')
    const minute = getShortTimeTextValue(q.minutes, 'm')
    const second = getShortTimeTextValue(q.seconds, 's')

    return `${hour}${minute}${second}`
  })

  const isMinFarmTimeDone = computed(() => {
    return now.value >= claimData.balance.value.claimed_at + getMinFarmTime.value
  })

  const isMaxFarmTimeDone = computed(
    () => now.value >= claimData.balance.value.claimed_at + getMaxFarmTime.value,
  )

  const getFarmAmountPeriod = computed(() => {
    return (
      (claimData.claim.value.amount * claimData.claim.value.min_period) /
      claimData.claim.value.max_period
    )
  })

  const percentDone = computed(() => {
    const val = ((farmedAmount.value / claimData.claim.value.amount) * 100).toFixed(0)
    if (Number(val) <= 0) return 0
    if (Number(val) >= 100) return 100
    return Number(val)
  })

  async function getBalance() {
    try {
      const response = await coin.getBalance()
      claimData.balance.value = response
    } catch (error) {
      console.error(error)
    }
    return 1
  }

  async function getClaim() {
    try {
      const response = await coin.getClaim()
      claimData.claim.value = response
    } catch (error) {
      console.error(error)
    }
    return 1
  }

  async function setClaimData() {
    if (loadingInner.value) return
    loadingInner.value = true

    try {
      const response = await coin.setClaim()
      if ('errors' in response) {
        loadingInner.value = false
        return
      }
    } catch (error) {
      loadingInner.value = true
      console.error(error)
      return
    }

    farmedAmount.value = 0
    claimData.balance.value.claimed_at = Date.now()

    try {
      await getBalance()
      await getClaim()
    } catch (error) {
      console.error(error)
    }

    updateABTRank()

    loadingInner.value = false
    return 1
  }

  function formatAmount(amount: number) {
    return (amount / 10 ** claimData.claim.value.precision).toFixed(claimData.claim.value.precision)
  }

  function updateNow() {
    now.value = Date.now()
  }

  function updateFarmedAmount() {
    const shift = now.value - claimData.balance.value.claimed_at

    const amount = new Decimal(claimData.claim.value.amount)
    const maxPeriod = new Decimal(getMaxFarmTime.value)

    const farmed = amount.div(maxPeriod).mul(shift)

    farmedAmount.value = farmed.toNumber()
  }

  function checkConffeti() {
    if (!isNew.value && (isMinFarmTimeDone.value || isMaxFarmTimeDone.value)) {
      nextTick(() => {
        conffeti.start()
      })
      return
    }

    conffeti.stop()
  }

  function drapProgress() {
    drawCoinCanvas(instance, percentDone.value)
    //   drawDonutCanvas(instance, percentDone.value)
  }

  function updateABTRank() {
    nextTick(() => {
      ABTRank.setClaimed(
        formatAmount(claimData.balance.value.balance),
        claimData.claim.value.precision,
      )
    })
  }

  function initTimer() {
    timer.functions.push(updateNow, drapProgress, updateFarmedAmount, checkConffeti)
    timer.start(farmedAmountTime)
  }

  // watch(() => ), () => {
  //   if (!isNew.value && !isFarmedFull.value) {
  //     nextTick(() => {
  //       conffeti.start()
  //     })
  //     retuisMinFarmTimeDone || isMaxFarmTimeDonern
  //   }

  //   conffeti.stop()
  // })

  function showInfo(type: ECoinFarmingActions) {
    $modals.farmingInfo.show(type)
  }

  function checkTogglePremium() {
    $modals.farmingSubscription.show('premium')

    if ($me.data.is_premium) {
      checkBoxPremium.value = true
    } else {
      setTimeout(() => {
        if ($me.data.is_premium) {
          checkBoxPremium.value = true
        } else {
          checkBoxPremium.value = false
        }
      }, 300)
    }
  }

  function checkToggleAutoFarm() {
    $modals.farmingSubscription.show('autofarm')

    if ($me.data.is_autofarm) {
      autoFarming.active.value = true
    } else {
      setTimeout(() => {
        if ($me.data.is_autofarm) {
          autoFarming.active.value = true
        } else {
          autoFarming.active.value = false
        }
      }, 300)
    }
  }

  const AutoUpdate = {
    timer: null as ReturnType<typeof setTimeout>,
    async update() {
      await getBalance()
    },
    async start() {
      await AutoUpdate.update()
      AutoUpdate.timer = setTimeout(async () => this.start(), 5000)
      return 1
    },
    stop() {
      if (AutoUpdate.timer) clearTimeout(AutoUpdate.timer)
    },
  }

  const games = {
    active: undefined,
    start(game) {
      this.stop()
      if (!game) return
      this.active = game()
      setTimeout(() => this.active.startGame(instance), 300)
    },
    async stop() {
      if (!this.active) return
      this.active.stopGame()
      this.active = undefined
    },
  }

  function showRanks() {
    $modals.ranks.data['rankId'] = ABTRank.currentRank.value.id || -1
    $modals.ranks.show()
  }

  
  function formatDepositAmount(amount: any) {
  return Number(amount || 0).toLocaleString(undefined, { maximumFractionDigits: 2 })
}

  function getCurrentSubscriptionName(): string {
    // Получаем текущую подписку пользователя
    const subscription = $me.data.products?.find(p => p.category === 'subscription')
    if (!subscription) return 'Researcher'
    
    const names: { [key: string]: string } = {
      'subscription_researcher': 'Researcher',
      'subscription_start': 'Start',
      'subscription_advanced': 'Advanced',
      'subscription_professional': 'Professional',
      'subscription_ambassador': 'Ambassador',
      'subscription_leader': 'Leader',
      'subscription_vip': 'VIP'
    }
    return names[subscription.code] || 'Researcher'
  }

  const totalDeposit = ref(0)

  onMounted(async () => {
    loading.value = true

    await AutoUpdate.start()

    try {
      await getClaim()
    } catch (error) {
      console.error(error)
    }

    updateABTRank()

    loading.value = false

    nextTick(() => {
      initTimer()
    })

    setTimeout(() => {
      // games.start()
      // startArrowGame(instance)
      // startSafeGame(instance)
    }, 300)

    try {
      const response = await $request({
        query: `query($userUid: String!) { depositTotal(userUid: $userUid) }`,
        variables: { userUid: $me.data.uid }
      })
      totalDeposit.value = response.data.data.depositTotal
    } catch (e) {
      totalDeposit.value = 0
    }
  })

  onBeforeUnmount(() => {
    AutoUpdate.stop()
    counter.stop()
    timer.stop()
    games.stop()
    // stopArrowGame()
    // conffeti.stop()
  })
</script>

<template>
  <div class="flex flex-col justify-center items-center w-full">
    <div v-if="loading" class="w-full flex justify-center items-center">
      <span class="text-gradient text-2xl uppercase font-bold">{{ $t('Loading') }}</span>
    </div>

    <template v-else>
      <div
        class="relative w-full h-full flex flex-col justify-around mobile:p-2 pc:p-4 space-y-2 overflow-hidden"
      >
        <!-- <div class="text-3xl text-center">
          <span class="text-gradient">ABT Coin</span>
        </div> -->

        <!-- <h1 class="coin-title">
          <span class="text-gradient font-bold">ABT Coin</span>
        </h1> -->

        <InCompanyWarning />
        <!-- <InTasksWarning /> -->

        <!-- <label class="flex cursor-pointer gap-2 items-center justify-center">
          <span class="label-text text-lg font-bold text-gradient">{{ $t('Premium') }}</span>
          <input
            v-model="checkBoxPremium"
            type="checkbox"
            class="toggle theme-controller checked:toggle-success toggle-lg"
            :true-value="true"
            :false-value="false"
            @change="checkTogglePremium"
          />
          <span class="label-text text-xs">
            {{ $me.data.is_premium ? $t('On') : $t('Off') }}
          </span>
        </label> -->

        <!-- <div class="text-base-content text-center font-bold mb-6">
          {{ $t('MyPromoCode') }}:
          <span
            class="whitespace-nowrap underline cursor-pointer text-gradient"
            @click="$copyToClipboard($me.data.uid)"
          >
            {{ $me.data.uid }}
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-4 h-4 inline relative -top-[3px]"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 0 1-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 0 1 1.5.124m7.5 10.376h3.375c.621 0 1.125-.504 1.125-1.125V11.25c0-4.46-3.243-8.161-7.5-8.876a9.06 9.06 0 0 0-1.5-.124H9.375c-.621 0-1.125.504-1.125 1.125v3.5m7.5 10.375H9.375a1.125 1.125 0 0 1-1.125-1.125v-9.25m12 6.625v-1.875a3.375 3.375 0 0 0-3.375-3.375h-1.5a1.125 1.125 0 0 1-1.125-1.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H9.75"
              />
            </svg>
          </span>
        </div> -->

        <div>
          <div class="flex justify-center items-center mt-4">
            <!-- {{ $t('CurrentRank') }}: -->
            <b class="text-gradient uppercase text-4xl">
              {{ ABTRank.currentRank.value.code || '-' }}
            </b>
          </div>

          <div
            v-if="
              ($me.data.products || []).filter((p) => p.code.match(/boost_x\d{1,}_lifetime_/))
                .length
            "
            class="flex justify-center gap-1 items-center my-2"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="#aa00d9"
              class="w-6 relative top-[1px]"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M16.5 6v.75m0 3v.75m0 3v.75m0 3V18m-9-5.25h5.25M7.5 15h3M3.375 5.25c-.621 0-1.125.504-1.125 1.125v3.026a2.999 2.999 0 0 1 0 5.198v3.026c0 .621.504 1.125 1.125 1.125h17.25c.621 0 1.125-.504 1.125-1.125v-3.026a2.999 2.999 0 0 1 0-5.198V6.375c0-.621-.504-1.125-1.125-1.125H3.375Z"
              />
            </svg>
            <span class="text-gradient">{{ $t('products.category.boost_x50_lifetime') }}</span>
          </div>

          <div class="flex justify-center items-center mt-6">
            <img width="35px" :src="symbolSrc" alt="symbol" class="relative -top-[1px]" />
            <div class="pl-2">
              <div class="text-2xl font-bold text-left font-mono">
                {{ formatAmount(claimData.balance.value.balance) }}
              </div>
            </div>
          </div>

          <div class="mt-1">
            <div class="flex justify-between items-center -mb-1">
              <span>
                {{ ABTRank.currentRank.value.code || '-' }} /
                {{ ABTRank.currentRank.value.min > -1 ? ABTRank.currentRank.value.min : '∞' }}
              </span>

              <button class="btn btn-ghost btn-sm btn-circle" @click="showRanks">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </button>
              <span>
                {{ ABTRank.nextRank.value.code || '-' }} /
                {{ ABTRank.currentRank.value.max || '∞' }}
              </span>
              <!-- {{ $t('CurrentRank') }}:  -->
            </div>
            <progress
              :key="ABTRank.getPercent.value"
              class="progress progress-success"
              :value="ABTRank.getPercent.value"
              max="100"
            ></progress>
          </div>
        </div>

        <div>
          <div class="relative">
            <div>
              <div class="flex justify-center">
                <!-- <template v-if="!isMinFarmTimeDone">
              {{ getIntervalToDurationMin }}
            </template> 
                <div class="opacity-80 font-mono text-sm">
                  <template v-if="!isMaxFarmTimeDone">
                    {{ getIntervalToDurationMax }}
                  </template>
                  <template v-else></template>
                </div>-->

                <div></div>
              </div>
            </div>

            <div class="absolute top-0 left-0 flex flex-col items-center cursor-pointer h-20 w-20" @click="Router.push({ name: 'Friends' })">
              <span class="text-xs font-bold mb-1 text-center">{{$t('Clients')}}</span>
              <IvitedTotal
                :is-icon="true"
                color="#38bdf8"
              />
            </div>

            <div class="absolute cursor-pointer top-0 right-0 text-center flex flex-col items-center h-20 w-20" @click="Router.push({ name: 'ReferralNetworkPartners' })">
              <span class="text-xs font-bold mb-1 text-center">{{$t('Partners')}}</span>
              <PartnersTotal
                :is-icon="true"
                color="#22c55e"
              />
            </div>


            <div
              class="absolute bottom-[50px] right-0 text-center flex flex-col items-center h-20 w-20"
            >
              <span class="text-xs font-bold mb-1 text-center">{{ $t('MyDeposit') }}</span>
              <img :src="coinClatterSrc" alt="deposit" class="w-10 h-10 text-white" />
              <div class="text-center">{{ formatDepositAmount(totalDeposit) }}</div>
            </div>

            <div class="absolute bottom-[50px] left-0 flex flex-col items-center text-center space-y-1 cursor-pointer h-20 w-20" @click="checkTogglePremium">
              <span class="text-xs font-bold mb-1 text-center">{{ $t('PremiumSlots') }}</span>
              <svg
                class="w-10 opacity-90 inline"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 122.88 107.7"
              >
                <path
                  fill="currentColor"
                  d="M21.08,83.82h24a15.89,15.89,0,0,1,4.79-11.49A16.33,16.33,0,0,1,77.72,83.82H101.4l12.49-34.67a6.65,6.65,0,0,0,2.1.3,6.9,6.9,0,1,0-6.9-6.89,6.78,6.78,0,0,0,1.3,4l-7.09,5.9c-10,8.19-16.19,4.39-14.29-8.29l1.1-7.2a4.87,4.87,0,0,0,1.2.1,6.88,6.88,0,1,0-4.3-1.5l-1.69,2.7c-8.4,12.59-14.59,7.79-17-4.69L63.94,13.29A6.79,6.79,0,0,0,68.13,7,7.16,7.16,0,0,0,61,0a6.9,6.9,0,0,0-6.89,6.89,6.75,6.75,0,0,0,5.09,6.6l-2.9,11.59C54,35.67,51,56.84,37.56,38.76l-2.19-3a7,7,0,0,0,2.89-5.6,6.89,6.89,0,1,0-6.89,6.89,7.72,7.72,0,0,0,1.5-.2l.5,4.7c.9,6.39,2,15-5.3,14.59-3.59-.2-5-1.4-7.79-3.4l-7.89-5.6A6.8,6.8,0,0,0,13.79,43a6.9,6.9,0,1,0-6.9,6.89,7,7,0,0,0,2.5-.5L21.08,83.82ZM61,3.1a3.8,3.8,0,1,1-3.8,3.79A3.8,3.8,0,0,1,61,3.1Zm54.75,35.46a4.15,4.15,0,1,1-4.2,4.1,4.14,4.14,0,0,1,4.2-4.1Zm-109,0a4.15,4.15,0,1,1,0,8.29A4.13,4.13,0,0,1,2.7,42.66a4.05,4.05,0,0,1,4.09-4.1ZM31.27,26a4.12,4.12,0,0,1,4.1,4.1,4.15,4.15,0,1,1-8.3,0,4.27,4.27,0,0,1,4.2-4.1Zm59.84,0a4.15,4.15,0,1,1,0,8.3,4.15,4.15,0,0,1,0-8.3Zm-70,67.44H48.25a14.62,14.62,0,0,0,1.6,1.9,16.34,16.34,0,0,0,11.49,4.79,16.77,16.77,0,0,0,11.59-4.79,14.62,14.62,0,0,0,1.6-1.9h27V107.7H21.08V93.41Z"
                />
                <path
                  fill="currentColor"
                  d="M61.34,72.53A11.29,11.29,0,1,1,50.05,83.82,11.28,11.28,0,0,1,61.34,72.53Z"
                />
              </svg>
              <div class="text-center">{{ getCurrentSubscriptionName() }}</div>
            </div>

            <div
              ref="coinBlock"
              class="coinBlock relative my-0 w-[250px] h-[250px] mx-auto [&>_>*]:select-none z-0"
            >
              <img
                ref="coinImg"
                :src="getCoinImg"
                alt="coin"
                class="absolute w-[190px] top-[30px] left-[30px] pointer-events-auto cursor-pointer"
              />
            </div>

          </div>
        </div>


        <!-- <div class="space-y-4">
          <div class="box">
            <div class="text-2xl text-gradient font-bold">{{ $t('coin.Store') }}</div>

            <div class="space-y-1 mt-2">
              <div v-if="!isMinFarmTimeDone" class="flex-[0_0_100%] opacity-80 font-mono text-sm">
                {{ $t('coin.store.output.min') }} {{ getIntervalToDurationMin }}
              </div>
              <div v-else-if="!isMaxFarmTimeDone" class="flex-[0_0_100%] opacity-80 font-mono text-sm">
                {{ $t('coin.store.output.max') }} {{ getIntervalToDurationMax }}
              </div>
              <div v-else-if="isNew">{{ $t('coin.store.output.new_full') }}</div>
              <div v-else-if="!isNew">{{ $t('coin.store.output.full') }}</div>
            </div>

            <div class="mt-4 mb-2">
              <progress
                class="progress h-4 w-full progress-coin-color-2"
                :value="percentDone"
                max="100"
              ></progress>
            </div>

            <div class="flex flex-wrap">
              <div class="flex flex-[0_0_100%] items-center justify-items-start gap-2">
                <span class="text-lg">{{ $t('coin.Farming') }}</span>
                <img class="relative top-[1px]" width="20px" :src="coinSrc" alrt="coin" />
                <span v-if="isNew" ref="counterElement" class="relative top-[1px] text-xl font-mono">
                  {{ formatAmount(getFarmAmountPeriod) }}
                </span>
                <span v-else-if="!isNew && !isFarmedFull" class="relative top-[1px] text-xl font-mono">
                  {{ formatAmount(farmedAmount) }}
                </span>
                <span v-else-if="!isNew && isFarmedFull" class="relative top-[1px] text-xl font-mono">
                  {{ formatAmount(claimData.claim.value.amount) }}
                </span>
              </div>
            </div>

            <p class="mt-4 text-sm text-mono">
              {{ $t('coin.store.output.farming_info') }}
            </p>
          </div>

          <div>
            <div ref="conffetiElement" class="relative w-full"></div>
            <button
              class="btn relative flex-wrap w-full p-4 h-auto overflow-hidden"
              :class="`${isMinFarmTimeDone || isMaxFarmTimeDone ? 'bg-[var(--coin-color-2)] border-[var(--coin-color-2)] text-white' : 'opacity-70'}`"
              @click="setClaimData"
            >
              {{ $t('coin.btn.toBalance') }}
            </button>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div v-for="item in infoList" :key="item.title" class="coin-box w-full">
              <h3 class="text-gradient font-bold">{{ item.title }}</h3>
              <p>{{ item.data }}</p>
            </div>
          </div>
        </div> -->
      </div>

      <Menu />
    </template>
  </div>
</template>
