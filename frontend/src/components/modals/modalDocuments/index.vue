<script lang="ts">
  export default {
    name: 'ModalDocuments',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
  import DocumentsList from 'components/documentsList/index.vue'
</script>

<template>
  <Modal :modal="$modals.documents" :title="$t('Documents')" :z="50">
    <DocumentsList />
  </Modal>
</template>
