<script lang="ts">
  export default {
    name: 'ModalPayouts',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
  import { $requests } from 'queries/index'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { $notify } from 'utils/notify'
  import { $t } from 'i18n/index'
  import { ECurrencies, EPayoutCodes } from 'types/enums'
  import { $me } from 'utils/me'
  import { $formatInt, $formatPriceByLocale } from 'utils/formats'
  import { $store } from 'utils/store'
  import TonButton from 'components/ton/button.vue'
  import { useTonConnect } from 'utils/wallets/tonconnect'

  const tonConnect = useTonConnect()

  const isWalletConnected = computed(() => tonConnect.wallet.isConnected().value)
  const isShowSend = computed(() => isWalletConnected.value && getTotalAmount() > 0)

  const loading = ref(false)
  const currencies = [ECurrencies.usdHidden]
  const cfg = computed(() => $store.get('cfg'))

  const form: PayoutReq = reactive({
    account_name: '',
    account_number: '',
    currency_code: currencies[0],
    method_code: EPayoutCodes.usdt_ton,
  })

  // function checkFormFill() {
  //   // if (!form.account_name || form.account_name.length <= 3) return
  //   if (!form.account_number || form.account_number.length <= 3) return
  //   if (!form.currency_code) return
  //   if (!form.method_code) return
  //   return true
  // }

  async function send() {
    if (loading.value) return

    if (!isWalletConnected.value) return
    if (getTotalAmount() <= 0) return
    if (!isShowSend.value) return

    form.account_number = tonConnect.wallet.getBounceableAddress() || ''

    // if (!checkFormFill()) {
    //   $notify.show({
    //     type: 'error',
    //     text: $t('error.fillAllFields'),
    //   })

    //   return
    // }

    loading.value = true

    try {
      const response = await $requests.payouts.create(form)

      if (response?.uid) {
        $me.update()
        $modals.payout.onSucess()

        $notify.show({
          type: 'success',
          text: $t('info.payoutSent'),
        })

        close()
      }
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  function close() {
    $modals.payout.close()
  }

  function getCommission() {
    const balance = $me.data.balance
    const balanceMin = cfg.value.payout_amount_min
    const feeMin = cfg.value.payout_fee_min
    const feePercent = cfg.value.payout_fee_percent
    if (balance < balanceMin) return 0
    return (balance / 10000) * feePercent < feeMin ? feeMin : (balance / 10000) * feePercent
  }

  function getTotalAmount() {
    const balance = $me.data.balance
    const balanceMin = cfg.value.payout_amount_min
    if (balance < balanceMin) return 0
    return balance - getCommission()
  }

  onMounted(async () => {
    loading.value = true
    await $store.updateCfg()
    loading.value = false
  })
</script>

<template>
  <Modal :modal="$modals.payout" :title="$t('Payout')">
    <div class="min-w-[20rem]">
      <div v-if="loading">
        <div class="text-center font-bold my-10">{{ $t('Loading') }}</div>
      </div>

      <div class="form">
        <div>
          <div class="label flex">
            <div class="flex-[0_0_50%]">
              <div class="text-base label-text flex-[0_0_50%]">{{ $t('Balance') }}:</div>
              <div class="text-xs italic opacity-80">
                {{ $t('PayoutAmountMin') }}:
                {{ $formatPriceByLocale({ count: $formatInt(cfg.payout_amount_min) }) }}
              </div>
            </div>
            <div class="font-bold">
              {{ $formatPriceByLocale({ count: $formatInt($me.data.balance) }) }}
            </div>
          </div>

          <div class="label flex">
            <div class="flex-[0_0_50%]">
              <div class="text-base label-text">
                {{ $t('PayoutFeePercent') }}: {{ $formatInt(cfg.payout_fee_percent) }}%
              </div>
              <div class="text-xs italic opacity-80">
                {{
                  $t('PayoutFeePercentMinMax', {
                    amount: $formatPriceByLocale({ count: $formatInt(cfg.payout_fee_min) }),
                  })
                }}
              </div>
            </div>
            <div class="font-bold">
              {{ $formatPriceByLocale({ count: $formatInt(getCommission()) }) }}
            </div>
          </div>

          <div class="label flex">
            <div class="flex-[0_0_50%]">
              <div class="text-base label-text flex-[0_0_50%]">{{ $t('Total') }}:</div>
              <div class="text-xs italic opacity-80">
                {{ $t('PayoutTotalInfo') }}
              </div>
            </div>
            <div class="font-bold">
              {{
                $formatPriceByLocale({
                  count: $formatInt(getTotalAmount()),
                  currency: ECurrencies.usdt,
                })
              }}
            </div>
          </div>
        </div>

        <div class="border border-error p-2 rounded-md">
          {{ $t('info.paymentsMessages.payoutTimeRestriction') }}
        </div>

        <!-- <div class="form-control">
          <label class="label">
            <span class="text-base label-text">{{ $t('AccountName') }}</span>
          </label>
          <input
            v-model="form.account_name"
            type="text"
            class="w-full input input-bordered"
            name="reg_name"
          />
        </div> -->

        <div class="form-control">
          <!-- <label class="label">
            <span class="text-base label-text">{{ $t('AccountNumber') }}</span>
          </label> -->

          <TonButton />

          <!-- <input
            v-model.trim="form.account_number"
            type="text"
            class="w-full input input-bordered"
            name="reg_email"
          /> -->
        </div>
        <!-- <div class="form-control">
          <label class="label">
            <span class="text-base label-text">{{ $t('Currency') }}</span>
          </label>
          <select v-model="form.currency_code" class="select select-bordered">
            <option v-for="currency in currencies" :key="currency" :value="currency">
              {{ currency }}
            </option>
          </select>
        </div>
        <div class="form-control">
          <label class="label">
            <span class="text-base label-text">{{ $t('MethodCode') }}</span>
          </label>
          <select v-model="form.method_code" class="select select-bordered">
            <option v-for="code in EPayoutCodes" :key="code" :value="code">
              {{ code }}
            </option>
          </select>
        </div> -->
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
