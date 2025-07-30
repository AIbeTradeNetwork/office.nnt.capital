<script lang="ts">
  export default {
    name: 'ModalInfo',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $t } from 'i18n/index'
  import { $modals } from 'utils/modals'
  import { $textSplit } from 'utils/text'
</script>

<template>
  <Modal :modal="$modals.info" :z="50">
    <template #title>
      <h3 class="font-bold text-lg flex items-center text-info">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-6 h-6 mr-1"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z"
          />
        </svg>

        {{ $modals.info.title || $t('Information') }}
      </h3>
    </template>

    <div class="space-y-4">
      <p v-for="par in $textSplit($modals.info.text)" :key="par">
        <template v-if="par.match(/&b&/gm)">
          <b class="text-lg text-info">{{ par.replace(/&b&/gm, '') }}</b>
        </template>
        <template v-else>
          {{ par }}
        </template>
      </p>
    </div>
  </Modal>
</template>
