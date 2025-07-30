<script lang="ts">
  export default {
    name: 'Dropdown',
  }
</script>

<script setup lang="ts">
  import { nextTick, Ref, ref } from 'vue'

  const show = ref(false)
  const dropdownAction: Ref<HTMLElement> = ref()
  const dropdownContent: Ref<HTMLElement> = ref()

  const props = withDefaults(
    defineProps<{
      // position?: 'left' | 'top' | 'right' | 'bottom'
      // align: 'start' | 'end'
      closeOnAction?: boolean
      open?: boolean
      closeBnt?: boolean
      class?: string
    }>(),
    {
      closeOnAction: true,
      open: false,
      closeBnt: false,
      class: '',
    },
  )

  // const position = computed(() => {
  //   if (props.position === 'top') return 'dropdown-top'
  //   if (props.position === 'left') return 'dropdown-left'
  //   if (props.position === 'right') return 'dropdown-right'
  //   return 'dropdown-bottom'
  // })

  // const align = computed(() => {
  //   if (props.align === 'end') return 'dropdown-end'
  //   return ''
  // })

  function setDropDownPosition() {
    dropdownContent.value.style.position = `absolute`

    const indent = 10
    const dropdownActionRect: DOMRect = dropdownAction.value.getBoundingClientRect()
    const dropdownContentRect: DOMRect = dropdownContent.value.getBoundingClientRect()

    const windowWidth = window.innerWidth
    const windowHeight = window.innerHeight

    const outesideLeft = dropdownActionRect.left - dropdownContentRect.width - indent < 0
    const outesideRight = dropdownActionRect.left + dropdownContentRect.width + indent > windowWidth

    const outesideTop = dropdownActionRect.top - dropdownContentRect.height - indent < 0
    const outesideBottom =
      dropdownActionRect.bottom + dropdownContentRect.height + indent > windowHeight

    if (outesideRight && outesideLeft && outesideBottom && outesideTop) {
      dropdownContent.value.style.left = `0`
      dropdownContent.value.style.right = `0`
      dropdownContent.value.style.top = `0`
      dropdownContent.value.style.bottom = '0'
      dropdownContent.value.style.maxWidth = '100%'
      dropdownContent.value.style.maxHeight = '100%'
      dropdownContent.value.style.minWidth = 'initial'
    }

    function checkLeftRight() {
      if (outesideLeft && outesideRight) {
        dropdownContent.value.style.left = `0`
        dropdownContent.value.style.right = `0`
        dropdownContent.value.style.maxWidth = '100%'
        dropdownContent.value.style.minWidth = 'initial'
        return
      }

      if (outesideRight) {
        dropdownContent.value.style.right = `${windowWidth - dropdownActionRect.right}px`
        return
      }

      dropdownContent.value.style.left = `${dropdownActionRect.left}px`
    }

    function checkTopBottom() {
      if (outesideTop && outesideBottom) {
        dropdownContent.value.style.top = `0`
        dropdownContent.value.style.bottom = `0`
        dropdownContent.value.style.maxHeight = '100%'
        dropdownContent.value.style.minHeight = 'initial'
        return
      }

      if (outesideBottom) {
        dropdownContent.value.style.bottom = `${windowHeight - dropdownActionRect.bottom}px`
        return
      }

      dropdownContent.value.style.top = `${dropdownActionRect.top + dropdownActionRect.height}px`
    }

    dropdownContent.value.removeAttribute('style')
    dropdownContent.value.style.position = `fixed`
    checkLeftRight()
    checkTopBottom()
  }

  function closeEvent(event) {
    const onInner = event?.target?.closest('.dropdown-content')
    if (!props.closeOnAction) {
      if (onInner) return
    }
    event.preventDefault()
    event.stopPropagation()
    show.value = false
    document.removeEventListener('click', closeEvent)
  }

  function setCloseEvent() {
    document.addEventListener('click', closeEvent)
  }

  function toggleDropdownContent() {
    show.value = !show.value
    if (show.value) {
      nextTick(() => {
        setDropDownPosition()
        setTimeout(() => setCloseEvent())
      })
    }
  }
</script>

<template>
  <div
    ref="dropdownAction"
    class="cursor-pointer"
    :class="`dropdown ${props.open ? 'dropdown-open' : ''}`"
    @click="toggleDropdownContent"
  >
    <slot name="trigger"></slot>

    <Teleport to="body">
      <div
        v-if="show"
        ref="dropdownContent"
        tabindex="0"
        class="dropdown-content drop-shadow-custom bg-base-100 rounded-box z-10 min-w-max"
        :class="props.class"
      >
        <button
          v-if="props.closeBnt"
          class="btn btn-sm btn-circle absolute right-2 top-2"
          @click="toggleDropdownContent"
        >
          ✕
        </button>
        <div :class="`${props.closeBnt ? 'pt-12' : ''}`">
          <slot name="drop"></slot>
        </div>
      </div>
    </Teleport>
  </div>
</template>
