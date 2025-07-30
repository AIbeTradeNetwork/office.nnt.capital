<script lang="ts">
  export default {
    name: 'Info',
  }
</script>

<script setup lang="ts">
  import { EAlerts, ESizes } from 'types/enums'
  import { $copyToClipboard } from 'utils/clipboard'
  import { $modals } from 'utils/modals'
  import { $textSplit } from 'utils/text'
  import { ref } from 'vue'

  const props = withDefaults(
    defineProps<{
      type?: EAlerts
      size?: ESizes
      title?: string
      message: string
      copytext?: string
      link?: string
      showIcon?: boolean
    }>(),
    {
      type: EAlerts.BASE,
      title: '',
      size: ESizes.MD,
      message: 'Message',
      copytext: '',
      link: '',
      showIcon: true,
    },
  )

  const isSizeSM = ref(props.size === ESizes.SM)

  function classOfCopyTxt() {
    return 'cursor-pointer'
  }

  function classOfType() {
    if (props.type === EAlerts.INFO) return 'alert-info'
    if (props.type === EAlerts.WARNING) return 'alert-warning'
    if (props.type === EAlerts.SUCCESS) return 'alert-success'
    if (props.type === EAlerts.ERORR) return 'alert-error'
    return ''
  }

  function classOfSize() {
    if (isSizeSM.value) return 'text-sm content-center'
    return ''
  }

  function actions() {
    if (props.link) {
      const link = document.createElement('a')
      link.href = new URL(`../../${props.link}`, import.meta.url).href
      link.target = '_blank'
      link.click()
      return
    }

    if (props.copytext) {
      copy()
      return
    }

    if (isSizeSM.value) {
      $modals.info.show({
        title: props.title,
        text: props.message,
      })
      return
    }
  }

  function copy() {
    if (!props.copytext) return
    $copyToClipboard(props.copytext)
  }
</script>

<template>
  <div
    role="alert"
    class="alert"
    :class="`${classOfCopyTxt()} ${props.showIcon ? 'gap-2' : 'flex'} ${classOfType()} ${classOfSize()} ${isSizeSM ? 'cursor-pointer' : ''}`"
    @click="actions"
  >
    <template v-if="props.showIcon && props.copytext">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="shrink-0 w-6 h-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 0 1-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 0 1 1.5.124m7.5 10.376h3.375c.621 0 1.125-.504 1.125-1.125V11.25c0-4.46-3.243-8.161-7.5-8.876a9.06 9.06 0 0 0-1.5-.124H9.375c-.621 0-1.125.504-1.125 1.125v3.5m7.5 10.375H9.375a1.125 1.125 0 0 1-1.125-1.125v-9.25m12 6.625v-1.875a3.375 3.375 0 0 0-3.375-3.375h-1.5a1.125 1.125 0 0 1-1.125-1.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H9.75"
        />
      </svg>
    </template>

    <template
      v-else-if="props.showIcon && (props.type === EAlerts.INFO || props.type === EAlerts.BASE)"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        class="stroke-current shrink-0 w-6 h-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        ></path>
      </svg>
    </template>

    <template v-else-if="props.showIcon && props.type === EAlerts.ERORR">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        class="stroke-current shrink-0 w-6 h-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        ></path>
      </svg>
    </template>

    <template v-else-if="props.showIcon && props.type === EAlerts.SUCCESS">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="stroke-current shrink-0 h-6 w-6"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
    </template>

    <template v-else-if="props.showIcon && props.type === EAlerts.WARNING">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="stroke-current shrink-0 h-6 w-6"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
        />
      </svg>
    </template>

    <div class="space-y-1 w-full min-w-[1px] break-words">
      <slot name="before"></slot>

      <div
        v-if="props.title"
        class="font-bold text-md"
        :class="`${isSizeSM ? 'text-ellipsis line-clamp-2' : ''}`"
      >
        {{ props.title }}
      </div>

      <div
        :class="`${isSizeSM && props.title ? 'hidden' : ''} ${isSizeSM ? 'text-ellipsis line-clamp-2' : ''}`"
      >
        <div class="space-y-4">
          <p v-for="par in $textSplit(props.message)" :key="par">{{ par }}</p>
        </div>
      </div>

      <slot name="in"></slot>

      <template v-if="props.copytext">
        <div
          :class="`${isSizeSM && props.title ? 'hidden' : ''} ${isSizeSM ? 'text-ellipsis line-clamp-1' : ''}`"
        >
          {{ props.copytext }}
        </div>
        <div class="mt-2 font-bold underline">{{ $t('PressToCopy') }}</div>
      </template>

      <slot name="after"></slot>
    </div>

    <button v-if="isSizeSM" class="btn btn-circle btn-sm btn-ghost">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-6 h-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25"
        />
      </svg>
    </button>
  </div>
</template>
