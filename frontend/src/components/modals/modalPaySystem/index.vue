<script lang="ts">
  export default {
    name: 'ModalPaySystem',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import PaySystem from 'components/paySystem/index.vue'
  import { $modals } from 'utils/modals'
</script>

<template>
  <Modal :modal="$modals.paySystem" :title="$t('Payment')">
    <PaySystem :data="$modals.paySystem.data" @success="$modals.paySystem.close()" />
  </Modal>
</template>
