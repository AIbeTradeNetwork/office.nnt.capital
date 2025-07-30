<script lang="ts">
  export default {
    name: 'ModalError',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
  import { useSlots } from 'vue'
  const slots = useSlots()
</script>

<template>
  <Modal :modal="$modals.error" :z="50">
    <template #title>
      <h3 class="font-bold text-lg flex items-center text-error">
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
            d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
          />
        </svg>

        {{ $modals.error.title || $t('Error') }}
      </h3>
    </template>

    <div class="mt-4 break-words">{{ $modals.error.text }}</div>
  </Modal>
</template>
