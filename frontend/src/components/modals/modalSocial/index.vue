<script lang="ts">
  export default {
    name: 'ModalSocial',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import UserSocial from 'components/userSocial/index.vue'
  import { $modals } from 'utils/modals'
</script>

<template>
  <Modal :modal="$modals.userSocial" class="scale-100" :clear-padding="true">
    <UserSocial />
  </Modal>
</template>
