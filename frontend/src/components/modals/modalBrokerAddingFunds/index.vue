<script lang="ts">
  export default {
    name: 'ModalAddingFunds',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
  import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
  import { $t } from 'i18n/index'
  import { $me } from 'utils/me'
  import { $formatInt, $formatPriceByLocale } from 'utils/formats'
  import { ECurrencies } from 'types/enums'
  import { Decimal } from 'decimal.js'
  import { $requests } from 'queries/index'
  import { $notify } from 'utils/notify'
  import { $store } from 'utils/store'

  const loading = ref(false)

  const agree = ref(false)

  const cfg = ref<Partial<Broker.Config>>({
    min_deposit: '0',
    max_deposit: '0',
    deposit_fee_percent: '0',
  })

  const balance = computed(() => $formatInt($me.data.balance))

  const refillWallet = computed(() => ({
    min: Number(cfg.value.min_deposit),
    max: Math.min(balance.value, Number(cfg.value?.max_deposit)),
  }))

  const inputValue = ref(refillWallet.value.min)

  const getFee = computed(() => {
    return new Decimal(inputValue.value || 0)
      .mul(new Decimal(cfg.value.deposit_fee_percent))
      .div(new Decimal(100))
      .toNumber()
  })

  const getTotal = computed(() => {
    return new Decimal(inputValue.value || 0).minus(new Decimal(getFee.value)).toNumber()
  })

  const isShowSend = computed(() => {
    return (
      agree.value &&
      inputValue.value >= refillWallet.value.min &&
      inputValue.value <= refillWallet.value.max
    )
  })

  const isNotEnoughMin = computed(() => {
    return inputValue.value < refillWallet.value.min
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

  function close() {
    $modals.brokerAddingFunds.close()
  }

  function addAmout(val: number) {
    inputValue.value = new Decimal(inputValue.value).plus(new Decimal(val)).toNumber()
    nextTick(() => checkRefillWalletValue())
  }

  async function send() {
    if (!isShowSend.value) return
    if (loading.value) return

    loading.value = true

    try {
      const response = await $requests.broker.deposit({
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

  onMounted(async () => {
    loading.value = true

    try {
      cfg.value = await $requests.broker.config()
      checkRefillWalletValue()
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  })

  onBeforeUnmount(() => {})
</script>

<template>
  <Modal :modal="$modals.brokerAddingFunds" :title="$t('AddingFunds')">
    <div class="min-w-[20rem]">
      <div v-if="loading">
        <p class="text-2xl text-center py-2">{{ $t('Loading') }}</p>
        <p class="text-xl text-error text-center py-4">{{ $t('DoNotClose') }}</p>
      </div>

      <div v-else class="form">
        <div>
          <div class="label flex">
            <div class="flex-[0_0_50%]">
              <div class="text-base label-text flex-[0_0_50%]">{{ $t('Balance') }}:</div>
              <div class="text-xs italic opacity-80">
                {{ $t('DepositAmountMin') }}:
                {{ $formatPriceByLocale({ count: refillWallet.min }) }}
              </div>
            </div>
            <div class="font-bold">
              {{ $formatPriceByLocale({ count: balance }) }}
            </div>
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
              @change="checkRefillWalletValue"
              @blur="checkRefillWalletValue"
            />
          </label>
          <div class="mt-2 flex justify-center gap-2">
            <button class="btn btn-square" @click="addAmout(1)">+1</button>
            <button class="btn btn-square" @click="addAmout(5)">+5</button>
            <button class="btn btn-square" @click="addAmout(10)">+10</button>
            <button class="btn btn-square" @click="addAmout(balance)">
              {{ $t('Max') }}
            </button>
          </div>
        </div>

        <div>
          <div class="label flex">
            <div class="flex-[0_0_50%]">
              <div class="text-base label-text">
                {{ $t('Commission') }}: {{ cfg.deposit_fee_percent }}%
              </div>
              <!-- <div class="text-xs italic opacity-80">Но не менее 3,00&nbsp;UDEX</div> -->
            </div>
            <div class="font-bold">
              {{
                $formatPriceByLocale({
                  count: getFee,
                })
              }}
            </div>
          </div>

          <div class="label flex gap-4">
            <div>
              <div class="text-base label-text">{{ $t('BalanceNext') }}:</div>
            </div>
            <div class="font-bold text-lg text-warning">
              {{
                $formatPriceByLocale({
                  count: getTotal,
                })
              }}
            </div>
          </div>
        </div>
      </div>

      <div v-if="isNotEnoughMin" class="border border-warning rounded-md p-2 mb-4">
        {{ $t('autotrade.info.notEnoughMin') }}
        <button
          class="btn btn-link p-0 btn-sm text-success m-0 ml-1"
          @click="$modals.addingFunds.show"
        >
          {{ $t('AddingFunds') }}
        </button>
      </div>

      <div class="form-control mb-4">
        <label class="label justify-start cursor-pointer">
          <input v-model="agree" type="checkbox" class="checkbox mr-2" />
          <span class="label-text">
            {{ $t('ApplyAgreement') }}.
            <a class="link link-info" href="#" @click="$modals.documents.show">
              {{ $t('Show') }} {{ $t('Documents').toLowerCase() }}
            </a>
          </span>
        </label>
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
