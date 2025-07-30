<script lang="ts">
  export default {
    name: 'Exchanges',
  }
</script>

<script setup lang="ts">
  import Exchanges from 'components/exchanges/index.vue'
  import Info from 'components/info/index.vue'
  import { EAlerts, ESizes } from 'types/enums'
</script>

<template>
  <div class="page">
    <Exchanges />

    <div class="grid mobile:gap-2 pc:gap-4 mobile:grid-cols-1 pc:grid-cols-3">
      <Info
        :type="EAlerts.BASE"
        :size="ESizes.SM"
        :title="$t('info.pageExchanges.registration_on_bybit')"
        link="assets/pdfs/registration_on_bybit.pdf"
      />
      <Info
        :type="EAlerts.BASE"
        :size="ESizes.SM"
        :title="$t('info.pageExchanges.create_wallet_ММ_iOS')"
        link="assets/pdfs/create_wallet_mm_ios.pdf"
      />
      <Info
        :type="EAlerts.BASE"
        :size="ESizes.SM"
        :title="$t('info.pageExchanges.connecting_exchange_and_creating_bot')"
        link="assets/pdfs/connecting_exchange_and_creating_bot.pdf"
      />
    </div>
  </div>
</template>
