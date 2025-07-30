<script lang="ts">
  export default {
    name: 'ModalExchangesList',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
  import Exchanges from 'components/exchanges/index.vue'
  function close() {
    $modals.exchanges.close()
  }
</script>

<template>
  <Modal :title="$t('modalExchanges.title')" :close="close">
    <Exchanges />
  </Modal>
</template>
