<script setup lang="ts">
  import { $modals } from 'utils/modals'
  import { defineAsyncComponent } from 'vue'
  const Info = defineAsyncComponent(() => import('./modalInfo/index.vue'))
  const Waning = defineAsyncComponent(() => import('./modalWarning/index.vue'))
  const Error = defineAsyncComponent(() => import('./modalError/index.vue'))
  // const PDF = defineAsyncComponent(() => import( './modalPDF/index.vue'))
  const Settings = defineAsyncComponent(() => import('./modalSettings/index.vue'))
  // const ExchangeAdd = defineAsyncComponent(() => import( './modalExchange/add.vue'))
  // const ExchangeList = defineAsyncComponent(() => import( './modalExchange/list.vue'))
  const Bot = defineAsyncComponent(() => import('./modalBot/index.vue'))
  const Strategies = defineAsyncComponent(() => import('./modalStrategies/index.vue'))
  const Documents = defineAsyncComponent(() => import('./modalDocuments/index.vue'))
  const Social = defineAsyncComponent(() => import('./modalSocial/index.vue'))
  const PaySystem = defineAsyncComponent(() => import('./modalPaySystem/index.vue'))
  const Keys = defineAsyncComponent(() => import('./modalKeys/index.vue'))
  const Question = defineAsyncComponent(() => import('./modalQuestion/index.vue'))
  const Payout = defineAsyncComponent(() => import('./modalPayout/index.vue'))
  const BrokerPayout = defineAsyncComponent(() => import('./modalBrokerPayout/index.vue'))
  const Promocode = defineAsyncComponent(() => import('./modalPromocode/index.vue'))
  const AddingFunds = defineAsyncComponent(() => import('./modalAddingFunds/index.vue'))
  const BrokerAddingFunds = defineAsyncComponent(() => import('./modalBrokerAddingFunds/index.vue'))
  const Balance = defineAsyncComponent(() => import('./modalBalance/index.vue'))
  const Prize = defineAsyncComponent(() => import('./modalPrize/index.vue'))

  const GoToPro = defineAsyncComponent(() => import('./modalGotopro/index.vue'))
  // import { config } from './config'
  // const cfg = config
</script>

<template>
  <!-- <PDF v-if="$modals.PDF.active" /> -->
  <Balance v-if="$modals.balance.active" />

  <Settings v-if="$modals.settings.active" />

  <!-- <ExchangeAdd v-if="$modals.exchange.active"  /> -->
  <!-- <ExchangeList v-if="$modals.exchanges.active"  /> -->

  <Bot v-if="$modals.bot.active" />
  <Strategies v-if="$modals.strategies.active" />

  <Info v-if="$modals.info.active" />
  <Waning v-if="$modals.warning.active" />
  <Error v-if="$modals.error.active" />

  <Social v-if="$modals.userSocial.active" />

  <Keys v-if="$modals.keys.active" />

  <Question v-if="$modals.question.active" />

  <BrokerPayout v-if="$modals.brokerPayout.active" />
  <Payout v-if="$modals.payout.active" />

  <Promocode v-if="$modals.promocode.active" />

  <PaySystem v-if="$modals.paySystem.active" />

  <BrokerAddingFunds v-if="$modals.brokerAddingFunds.active" />
  <AddingFunds v-if="$modals.addingFunds.active" />

  <Documents v-if="$modals.documents.active" />

  <Prize v-if="$modals.prize.active" />

  <GoToPro v-if="$modals.goToPro.active" />
</template>

<style scoped></style>
