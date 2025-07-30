<script lang="ts">
  export default {
    name: 'ModalGotopro',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $t } from 'i18n/index'
  import { $modals } from 'utils/modals'
</script>

<template>
  <Modal :modal="$modals.goToPro" :z="50">
    <h3 class="font-bold text-lg flex items-center text-info">
      {{ $t('autotrade_pro.gotopro.title') }}
    </h3>

    <div class="mt-4">
      <p>{{ $t('autotrade_pro.gotopro.text') }}</p>
    </div>

    <div class="mt-4 text-center">
      <router-link
        :to="{ name: 'AutotradePRO' }"
        class="btn text-white bg-gold/80 pc:hover:bg-gold/40 mobile:active:bg-gold/40"
      >
        {{ $t('autotrade_pro.go') }}
      </router-link>
    </div>
  </Modal>
</template>
