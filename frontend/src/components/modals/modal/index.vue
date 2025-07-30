<script lang="ts">
  export default {
    name: 'Modal',
  }
</script>

<script setup lang="ts">
  import { Router } from 'routes/index'
  import { useSlots } from 'vue'

  const zOptions = {
    10: 'z-10',
    20: 'z-20',
    30: 'z-30',
    40: 'z-40',
    50: 'z-50',
  } as const

  const slots = useSlots()

  const props = withDefaults(
    defineProps<{
      clearPadding?: boolean
      modal?: any
      title?: string
      z?: keyof typeof zOptions
      backdrop?: boolean
      backdropClose?: boolean
      closeBtn?: boolean
      class?: string
      size?: 'default' | 'screen'
      width?: 'full' | 'auto'
    }>(),
    {
      clearPadding: false,
      modal: null,
      title: '',
      z: 10,
      backdrop: true,
      backdropClose: true,
      closeBtn: true,
      class: '',
      size: 'default',
      width: 'auto',
    },
  )

  const fullSizeClass =
    props.size === 'screen' ? 'h-vhd w-vhd max-w-full w-full h-full scale-100 rounded-none' : ''

  function hideAction() {
    props.modal.close()
  }

  function getZIndex() {
    return document.querySelectorAll('.modal').length + 100
  }

  Router.afterEach((to, from) => {
    if (to.name !== from.name) hideAction()
  })
</script>

<template>
  <div
    class="modal safearea-full opacity-100 pointer-events-auto"
    :class="`${zOptions[props.z]} ${props.class}`"
    :style="'z-index:' + getZIndex()"
    role="dialog"
  >
    <div
      class="modal-box mobile:w-full max-h-full px-4 py-2 w-auto"
      :class="`
        ${fullSizeClass}
        ${props.clearPadding ? 'p-0' : 'py-4'}
        ${props.width === 'full' ? 'w-full' : ''}
      `"
    >
      <div class="mb-4 pr-12 font-bold text-lg">
        <h3 v-if="props.title" class="">
          {{ props.title }}
        </h3>
        <slot v-else name="title"></slot>
      </div>

      <div><slot></slot></div>

      <div v-if="slots.actions" class="modal-action mt-5">
        <slot name="actions"></slot>
      </div>

      <button
        v-if="closeBtn"
        class="btn btn-sm btn-circle absolute right-2 top-2"
        @click="hideAction"
      >
        ✕
      </button>
    </div>
    <div
      v-if="backdrop"
      class="modal-backdrop cursor-pointer"
      @click="props.backdropClose ? hideAction() : null"
    ></div>
  </div>
</template>
