<script lang="ts">
  export default {
    name: 'ModalPayouts',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
  import { $requests } from 'queries/index'
  import { computed, nextTick, onMounted, ref } from 'vue'
  import { $notify } from 'utils/notify'
  import { $t } from 'i18n/index'
  import { ECurrencies } from 'types/enums'
  import { $formatPriceByLocale } from 'utils/formats'
  import { $store } from 'utils/store'
  import { Decimal } from 'decimal.js'
  import { $me } from 'utils/me'

  const loading = ref(true)

  const balance = ref(Number($store.get('brokerBalance').amount))

  const cfg = ref<Partial<Broker.Config>>({
    min_withdraw: '0',
    max_withdraw: '0',
    withdraw_fee_percent: '0',
  })

  const refillWallet = computed(() => ({
    min: Number(cfg.value.min_withdraw),
    max: Math.min(balance.value, Number(cfg.value?.max_withdraw)),
  }))

  const inputValue = ref(refillWallet.value.min)

  const isShowSend = computed(() => {
    return inputValue.value >= refillWallet.value.min && inputValue.value <= refillWallet.value.max
  })

  const getCommission = computed(() => {
    return new Decimal(inputValue.value || 0)
      .mul(new Decimal(cfg.value.withdraw_fee_percent))
      .div(new Decimal(100))
      .toNumber()
  })

  const getTotalAmount = computed(() => {
    return new Decimal(inputValue.value).minus(new Decimal(getCommission.value)).toNumber()
  })

  function checkRefillWalletValue() {
    if (Number(refillWallet.value.min) > Number(refillWallet.value.max)) {
      inputValue.value = Number(refillWallet.value.max)
      return
    }

    if (Number(inputValue.value) > Number(refillWallet.value.max)) {
      inputValue.value = Number(refillWallet.value.max)
      return
    }

    if (Number(inputValue.value) < Number(refillWallet.value.min)) {
      inputValue.value = Number(refillWallet.value.min)
      return
    }
  }

  function addAmout(val: number) {
    inputValue.value = new Decimal(inputValue.value).plus(new Decimal(val)).toNumber()
    nextTick(() => checkRefillWalletValue())
  }

  async function send() {
    if (loading.value) return

    if (getTotalAmount.value <= 0) return
    if (!isShowSend.value) return

    loading.value = true

    try {
      const response = await $requests.broker.withdraw({
        amount: inputValue.value + '',
        currency: ECurrencies.usdHidden,
      })

      if (response.amount) {
        $notify.show({
          type: 'success',
          text: $t('Success'),
        })
        $me.update()
        $store.updateBrokerBalance()
        close()
      } else {
        $notify.show({
          type: 'error',
          text: $t('Error'),
        })
      }
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  function close() {
    $modals.brokerPayout.close()
  }

  onMounted(async () => {
    loading.value = true

    try {
      await $store.updateBrokerBalance()
      cfg.value = await $requests.broker.config()
      checkRefillWalletValue()
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  })
</script>

<template>
  <Modal :modal="$modals.brokerPayout" :title="$t('Payout')">
    <div class="min-w-[20rem]">
      <div v-if="loading">
        <div class="text-center font-bold my-10">{{ $t('Loading') }}</div>
      </div>

      <div v-else class="form">
        <div>
          <div class="label flex">
            <div class="flex-[0_0_50%]">
              <div class="text-base label-text flex-[0_0_50%]">{{ $t('Balance') }}:</div>
              <div class="text-xs italic opacity-80">
                {{ $t('PayoutAmountMin') }}:
                {{ $formatPriceByLocale({ count: refillWallet.min }) }}
              </div>
            </div>
            <div class="font-bold">
              {{ $formatPriceByLocale({ count: balance }) }}
            </div>
          </div>

          <div class="form-control">
            <label class="label">
              <span class="text-base label-text">
                {{ $t('DepositAmount') }} {{ ECurrencies.usd.toUpperCase() }}:
              </span>
            </label>
            <label class="input input-bordered flex items-center gap-2">
              <input
                v-model="inputValue"
                type="number"
                :step="1"
                :min="refillWallet.min"
                :max="refillWallet.max"
                class="grow bg-transparent"
                name="refill_value"
                @change="checkRefillWalletValue()"
                @blur="checkRefillWalletValue()"
              />
            </label>
            <div class="mt-2 flex justify-center gap-2">
              <div class="mt-2 flex justify-center gap-2">
                <button class="btn btn-square" @click="addAmout(1)">+1</button>
                <button class="btn btn-square" @click="addAmout(5)">+5</button>
                <button class="btn btn-square" @click="addAmout(10)">+10</button>
                <button class="btn btn-square" @click="addAmout(balance)">
                  {{ $t('Max') }}
                </button>
              </div>
            </div>
          </div>

          <div class="label flex mt-6">
            <div class="flex-[0_0_50%]">
              <div class="text-base label-text">
                {{ $t('PayoutFeePercent') }}: {{ cfg.withdraw_fee_percent }}%
              </div>
            </div>
            <div class="font-bold">
              {{ $formatPriceByLocale({ count: getCommission }) }}
            </div>
          </div>

          <div class="label flex">
            <div class="flex-[0_0_50%]">
              <div class="text-base label-text flex-[0_0_50%]">{{ $t('Total') }}:</div>
            </div>
            <div class="font-bold text-lg text-warning">
              {{
                $formatPriceByLocale({
                  count: getTotalAmount,
                  currency: ECurrencies.usd,
                })
              }}
            </div>
          </div>
        </div>
      </div>

      <hr class="opacity-50" />

      <!-- Editing -->
      <div class="flex gap-2 justify-end mt-6">
        <button :disabled="loading || !isShowSend" class="btn btn-success" @click="send">
          {{ $t('Sent') }}
        </button>
        <button class="btn btn-outline" @click="close">{{ $t('Close') }}</button>
      </div>
    </div>
  </Modal>
</template>
