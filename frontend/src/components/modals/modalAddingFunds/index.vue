<script lang="ts">
  export default {
    name: 'ModalAddingFunds',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import PayTelegramWarning from 'components/payTelegramWarning/index.vue'

  import { $modals } from 'utils/modals'
  import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
  import { $t } from 'i18n/index'
  import { $me } from 'utils/me'
  import { $formatInt, $formatPriceByLocale } from 'utils/formats'
  import { $store } from 'utils/store'
  import TonButton from 'components/ton/button.vue'
  import { useTonConnect } from 'utils/wallets/tonconnect'
  import { ECurrencies } from 'types/enums'
  import { $notify } from 'utils/notify'
  import { Decimal } from 'decimal.js'
  import { $useTelegram } from 'utils/telegram'

  const tonConnect = useTonConnect()

  const inTelegramAppIOS = $useTelegram().inTelegramAppIOS

  const isWalletConnected = computed(() =>
    inTelegramAppIOS ? false : tonConnect.wallet.isConnected().value,
  )

  const loading = ref(false)

  const refillWallet = {
    min: 1,
    max: 99999999999,
  }

  const refillWalletValue = ref(refillWallet.min)

  const Rates = {
    data: ref(null),
    loaded: ref(false),
    timeOutRates: null as ReturnType<typeof setTimeout>,
    getPrice: computed(() => {
      if (!Rates.data.value) return 0
      return Rates.data.value?.rates?.TON?.prices?.USDT
    }),
    async update() {
      try {
        const json = await fetch('https://tonapi.io/v2/rates?tokens=ton&currencies=usdt')
        const parsed = await json.json()
        Rates.data.value = parsed
        Rates.loaded.value = true
      } catch (error) {
        $notify.show({
          error: error,
        })
      }
    },
    async start() {
      await Rates.update()
      Rates.timeOutRates = setTimeout(async () => this.start(), 2000)
    },
    clear() {
      if (Rates.timeOutRates) clearTimeout(Rates.timeOutRates)
    },
  }

  function checkRefillWalletValue(el: HTMLInputElement) {
    const _value = Number(el.value)
    if (_value > refillWallet.max) return (refillWalletValue.value = refillWallet.max)
    if (_value < refillWallet.min) return (refillWalletValue.value = refillWallet.min)
  }

  function close() {
    $modals.addingFunds.close()
  }

  async function send() {
    loading.value = true

    if (!isWalletConnected.value) return

    try {
      const response = await tonConnect.wallet.transferToAPersonalAccount({
        amount: refillWalletValue.value,
      })

      if (response && response.boc) {
        $modals.info.show({
          title: $t('info.paymentsMessages.success'),
          text: $t('info.paymentsMessages.transactionHasBeenSent'),
        })

        await $me.update()

        close()
      }
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  onMounted(async () => {
    Rates.start()
    $store.updateCfg()
  })

  onBeforeUnmount(() => {
    Rates.clear()
  })
</script>

<template>
  <Modal :modal="$modals.addingFunds" :title="$t('AddingFunds')">
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
            </div>
            <div class="font-bold">
              {{ $formatPriceByLocale({ count: $formatInt($me.data.balance) }) }}
            </div>
          </div>
        </div>

        <div v-if="!inTelegramAppIOS" class="form-control">
          <label class="label">
            <span class="text-base label-text">
              {{ $t('DepositAmount') }} {{ ECurrencies.ton.toUpperCase() }}:
            </span>
          </label>
          <label class="input input-bordered flex items-center gap-2">
            <input
              v-model="refillWalletValue"
              type="number"
              :step="1"
              :min="1"
              class="grow bg-transparent"
              name="refill_value"
              @change="checkRefillWalletValue($event.target as HTMLInputElement)"
              @blur="checkRefillWalletValue($event.target as HTMLInputElement)"
            />
          </label>
          <div class="mt-2 flex justify-center gap-2">
            <button class="btn btn-square" @click="refillWalletValue += 1">+1</button>
            <button class="btn btn-square" @click="refillWalletValue += 5">+5</button>
            <button class="btn btn-square" @click="refillWalletValue += 10">+10</button>
          </div>
        </div>

        <div v-if="!inTelegramAppIOS && Rates.loaded">
          <div>{{ $t('BalanceNext') }}:</div>
          <div class="text-lg text-warning">
            ≈
            {{
              $formatPriceByLocale({
                count: new Decimal(Rates.getPrice.value)
                  .mul(new Decimal(refillWalletValue || 0))
                  .toNumber(),
              })
            }}
          </div>
        </div>

        <div v-if="inTelegramAppIOS" class="form-control">
          <PayTelegramWarning />
        </div>

        <div v-if="!inTelegramAppIOS" class="form-control">
          <label class="label">
            <span class="text-base label-text">{{ $t('Wallet') }}:</span>
          </label>

          <TonButton />
        </div>
      </div>

      <hr class="opacity-50" />

      <!-- Editing -->
      <div class="flex gap-2 justify-end mt-6">
        <button :disabled="loading || !isWalletConnected" class="btn btn-success" @click="send">
          {{ $t('Sent') }}
        </button>
        <button class="btn btn-outline" @click="close">{{ $t('Close') }}</button>
      </div>
    </div>
  </Modal>
</template>
