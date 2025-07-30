<script lang="ts">
  export default {
    name: 'PaySystem',
  }
</script>

<script setup lang="ts">
  import TonButton from 'components/ton/button.vue'

  import { $formatInt, $formatPriceByLocale } from 'utils/formats'
  import { $modals } from 'utils/modals'
  import { computed, onMounted, ref } from 'vue'
  import { $t } from 'i18n/index'
  import { $store } from 'utils/store'
  import { useTonConnect } from 'utils/wallets/tonconnect'
  import { $notify } from 'utils/notify'
  import { $requests } from 'queries/index'
  import { $me } from 'utils/me'
  import { Decimal } from 'decimal.js'
  import { ECurrencies } from 'types/enums'
  import { coin } from 'queries/claim'

  const { data } = defineProps<{
    data: ITransaction & { onSuccess?: (...args: any[]) => void }
  }>()

  const emits = defineEmits(['success', 'error'])

  const loading = ref(false)
  const status = ref<string>()
  const balanceABT = ref<ClaimBalance | null>(null)
  const isUDEX = computed(() => data.currency_code === ECurrencies.udex)

  let tonConnect: ReturnType<typeof useTonConnect>

  if (!isUDEX.value) {
    tonConnect = useTonConnect()
  }

  const isWalletConnected = computed(() => {
    if (isUDEX.value) return false
    if (!tonConnect) return false
    return tonConnect.wallet.isConnected().value
  })

  const products = computed(() => $store.get('products'))

  const selectedProduct = computed(() => {
    if (data.type !== 'product') return data
    return products.value.find((item) => item.code === (data?.code || ''))
  })

  const balanceCurrent = computed(() => {
    if (isUDEX.value) {
      return balanceABT.value?.balance || 0
    }

    return $me.data.balance
  })

  const balanceAfter = computed(() => {
    if (!selectedProduct.value) return 0
    let bal: Decimal
    if (isUDEX.value) {
      bal = new Decimal(balanceABT.value?.balance || 0)
    } else {
      bal = new Decimal($me.data.balance)
    }

    const min = new Decimal(selectedProduct.value.price)

    return bal.minus(min).toNumber()
  })

  function showSuccess(text) {
    status.value = text

    $modals.info.show({
      text,
    })
  }

  function showError(text) {
    status.value = text

    $modals.error.show({
      text,
    })
  }

  async function pay() {
    loading.value = true

    try {
      let response

      if (data.type === 'product') {
        response = await $requests.products.buyProduct({ code: selectedProduct.value.code })
      }

      if (data.type === 'tariff') {
        response = await $requests.products.buyTariff({ code: data.code })
      }

      if (data.type === 'distributor') {
        response = await $requests.products.buyDistributor()
      }

      if (response?.uid) {
        emits('success')
        showSuccess($t('info.paymentsMessages.success'))
        setTimeout(() => {
          $me.update()
          if (data.onSuccess) data.onSuccess()
        }, 2000)
      }
    } catch (error) {
      emits('error')
      if (error instanceof Error) return showError(error.message)
      showError($t('info.paymentsMessages.failed'))
    } finally {
      loading.value = false
    }
  }

  async function init() {
    loading.value = true

    try {
      await $store.updateCfg()
      if (isUDEX.value) balanceABT.value = await coin.getBalance()
      if (data.type === 'product') await $store.updateProducts()
    } catch (error) {
      $notify.show({ error: error })
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    init()
  })
</script>

<template>
  <div class="min-w-[300px]">
    <div v-if="selectedProduct" class="my-4 space-y-2">
      <!-- <div class="text-info mb-6">
      {{ $t('info.paySystem.payRequaredUSDT') }}
    </div> -->

      <div>
        <div>{{ $t('Transaction') }}:</div>
        <div class="text-lg text-warning">
          {{ data.what || selectedProduct.category }}
        </div>
      </div>

      <div>
        <div>{{ $t('BalanceCurrent') }}:</div>
        <div class="text-lg text-warning">
          <template v-if="isUDEX">
            {{
              $formatPriceByLocale({
                count: $formatInt(balanceCurrent, { precision: balanceABT?.precision }),
                currency: ECurrencies.udex,
              })
            }}
          </template>
          <template v-else>
            {{ $formatPriceByLocale({ count: $formatInt(balanceCurrent) }) }}
          </template>
        </div>
      </div>

      <div>
        <div>{{ $t('ChargeAmount') }}:</div>
        <div class="text-lg text-warning">
          {{
            $formatPriceByLocale({
              count: $formatInt(selectedProduct.price, {
                precision: selectedProduct.precision,
              }),
              currency: selectedProduct.currency_code,
            })
          }}
        </div>
      </div>

      <div>
        <div>{{ $t('BalanceAfter') }}:</div>
        <div class="text-lg text-warning">
          {{
            $formatPriceByLocale({
              count: $formatInt(balanceAfter <= 0 ? 0 : balanceAfter, {
                precision: selectedProduct.precision,
              }),
              currency: selectedProduct.currency_code,
            })
          }}
        </div>
      </div>
    </div>

    <div v-if="balanceAfter <= 0" class="border border-error my-4 p-4 rounded-md">
      {{ $t('info.BalanceNotEnoughFunds.part_1') }}
      <span v-if="!isUDEX">{{ $t('info.BalanceNotEnoughFunds.part_2') }}</span>
      <a
        v-if="!isUDEX"
        href="javascript:void(0)"
        class="text-success ml-1"
        @click.prevent="$modals.addingFunds.show"
      >
        {{ $t('AddingFunds') }}
      </a>
    </div>

    <template v-if="!loading">
      <TonButton v-if="!isUDEX" />

      <button
        v-if="isUDEX || isWalletConnected"
        :disabled="loading || balanceAfter <= 0"
        class="btn w-full btn-primary mt-6"
        @click="pay"
      >
        {{ $t('Pay') }}
        <!-- <span class="ml-1">
          {{
            $formatPriceByLocale({
              count: $formatInt(selectedProduct.price, {
                precision: selectedProduct.precision,
              }),
              currency: selectedProduct.currency_code,
            })
          }}
        </span> -->
      </button>
    </template>

    <div v-else>
      <p class="text-2xl text-center py-2">{{ $t('Loading') }}</p>
      <p class="text-xl text-error text-center py-4">{{ $t('DoNotClose') }}</p>
    </div>

    <!-- <template v-if="!loading">
      <button class="btn btn-primary btn-block mt-4" @click="pay">
        {{ $t('Pay') }}
      </button>

      <button class="btn btn-outline btn-block mt-4" @click="openWallet">
        {{ $t('Wallet') }}
      </button>
    </template>

    <template v-else>
      <p class="text-2xl text-center py-2">{{ $t('Loading') }}</p>
      <p class="text-xl text-error text-center py-4">{{ $t('DoNotClose') }}</p>
    </template> -->

    <!-- <button
      v-if="$device.isMobile"
      class="btn btn-primary btn-block mt-4"
      @click="generateMobileUrl"
    >
      {{ $t('OpenWallet') }}
    </button>
    <canvas v-if="$device.isPc" ref="qrzone"></canvas> -->
  </div>
</template>
