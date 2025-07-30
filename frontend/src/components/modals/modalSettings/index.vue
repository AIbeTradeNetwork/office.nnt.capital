<script lang="ts">
  export default {
    name: 'ModalSettings',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
</script>

<template>
  <Modal :modal="$modals.settings" :title="$t('modalSettings.title')">
    Settings
  </Modal>
</template>
