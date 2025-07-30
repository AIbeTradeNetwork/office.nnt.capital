<script lang="ts">
  export default {
    name: 'ModalPayouts',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import PayTelegramWarning from 'components/payTelegramWarning/index.vue'

  import { $modals } from 'utils/modals'
  import { computed, onMounted, ref } from 'vue'
  import { $t } from 'i18n/index'
  import { ECurrencies } from 'types/enums'
  import { $me } from 'utils/me'
  import { $formatInt } from 'utils/formats'
  import { $store } from 'utils/store'
  import { coin } from 'queries/claim'
  import { Router } from 'routes/index'
  import { $useTelegram } from 'utils/telegram'

  const inTelegramAppIOS = $useTelegram().inTelegramAppIOS

  const loading = ref(true)
  const coinBalance = ref<ClaimBalance>()

  function close() {
    $modals.balance.close()
  }

  async function getABTBalance() {
    try {
      const response = await coin.getBalance()
      coinBalance.value = response
    } catch (error) {
      console.error(error)
    }
    return 1
  }

  function openModal(type: 'payout' | 'addingFunds') {
    close()
    if (type === 'payout') $modals.payout.show()
    if (type === 'addingFunds') $modals.addingFunds.show()
  }

  function toPayments() {
    close()
    Router.push({ name: 'Payments' })
  }

  onMounted(async () => {
    loading.value = true
    Promise.allSettled([$store.updateCfg(), $me.update(), getABTBalance()])
    loading.value = false
  })
</script>

<template>
  <Modal :modal="$modals.balance" :title="$t('Balance')">
    <div class="min-w-[20rem]">
      <div v-if="loading">
        <div class="text-center font-bold my-10">{{ $t('Loading') }}</div>
      </div>

      <div v-else>
        <div class="mb-6 space-y-4">


          <div class="flex justify-between font-bold text-lg">
            <div class="flex-[0_0_50%]">{{ ECurrencies.udex.toUpperCase() }}:</div>
            <div class="font-bold text-lg">
              {{
                coinBalance ?
                  $formatInt(coinBalance.balance, { precision: coinBalance?.precision })
                : '-'
              }}
            </div>
          </div>
        </div>

        <div class="space-y-6 mb-6">
          <button class="btn btn-block btn-outline btn-neutral" @click="toPayments">
            {{ $t('PaymentsLog') }}
          </button>
          <button class="btn btn-block btn-outline btn-info" @click="openModal('payout')">
            {{ $t('WithdrawFromBalance') }}
          </button>

          <PayTelegramWarning v-if="inTelegramAppIOS" />

          <button v-else class="btn btn-block btn-success" @click="openModal('addingFunds')">
            {{ $t('TopUpBalance') }}
          </button>
        </div>
      </div>

      <hr class="opacity-50" />

      <!-- Editing -->
      <div class="flex gap-2 justify-end mt-6">
        <button class="btn btn-outline" @click="close">{{ $t('Close') }}</button>
      </div>
    </div>
  </Modal>
</template>
