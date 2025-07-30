<script lang="ts">
  export default {
    name: 'ModalExchangeAdd',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import Exchages from 'components/exchanges/index.vue'
  import { $modals } from 'utils/modals'
</script>

<template>
  <Modal :modal="$modals.exchanges" :title="$t('modalExchange.title')">
    <Exchages />

    <template #actions>
      <button class="btn btn-primary flex-1" @click="$modals.exchanges.close">
        {{ $t('Apply') }}
      </button>

      <button class="btn flex-1" @click="$modals.exchanges.close">
        {{ $t('Close') }}
      </button>
    </template>
  </Modal>
</template>
