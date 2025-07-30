<script lang="ts">
  export default {
    name: 'ModalQuestion',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'

  function apply() {
    $modals.question.onConfirm()
    $modals.question.close()
  }

  function deny() {
    $modals.question.onDeny()
    $modals.question.close()
  }
</script>

<template>
  <Modal :modal="$modals.question" :z="50">
    <template #title>
      <h3 class="font-bold text-lg flex items-center text-warning">
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
            d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"
          />
        </svg>

        {{ $modals.question.title || $t('Warning') }}
      </h3>
    </template>

    {{ $modals.question.text || '' }}

    <template #actions>
      <div clas="flex gap-2">
        <button class="btn btn-sm btn-outline btn-warning" @click="apply">
          {{ $t('Apply') }}
        </button>

        <button class="btn btn-sm btn-primary ml-2" @click="deny">
          {{ $t('Close') }}
        </button>
      </div>
    </template>
  </Modal>
</template>
