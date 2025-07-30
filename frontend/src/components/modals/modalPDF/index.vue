<script lang="ts">
  export default {
    name: 'ModalPdf',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
</script>

<template>
  <Modal :modal="$modals.PDF" :title="$modals.PDF.title" :z="50">
    <object
      width="1000"
      height="1000"
      type="application/pdf"
      :src="$modals.PDF.src"
    ></object>
  </Modal>
</template>
