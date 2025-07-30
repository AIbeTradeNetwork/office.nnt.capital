<script lang="ts">
  export default {
    name: 'Notify',
  }
</script>

<script setup lang="ts">
  import { $notify } from 'utils/notify'

  const list = $notify.list
  const getType = (type: INotify['type']) => {
    if (type === 'error') return 'alert-error'
    if (type === 'info') return 'alert-info'
    if (type === 'warning') return 'alert-warning'
    if (type === 'success') return 'alert-success'

    return ''
  }

  const getTypeStroke = (type: INotify['type']) => {
    if (type) return 'stroke-current'

    return 'stroke-info'
  }
</script>

<template>
  <div class="fixed flex flex-col gap-2 bottom-2 left-1/2 -translate-x-1/2 z-[1000]">
    <transition-group appear name="list">
      <div
        v-for="item in list"
        :key="item.id"
        role="alert"
        class="alert flex gap-2 rounded-md p-2 drop-shadow-custom"
        :class="getType(item.type)"
      >
        <svg
          class="shrink-0 w-6 h-6"
          :class="getTypeStroke(item.type)"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <path
            v-if="item.type === 'info'"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
          <path
            v-else-if="item.type === 'success'"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
          />
          <path
            v-else-if="item.type === 'warning'"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
          <path
            v-else-if="item.type === 'error'"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
          />
          <path
            v-else
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>

        <div>
          <h3 v-if="item.title" class="font-bold text-left">{{ item.title }}</h3>
          <div v-if="item.title && item.text" class="text-xs text-left">{{ item.text }}</div>
          <span v-else-if="!item.title && item.text">{{ item.text }}</span>
        </div>

        <div class="flex gap-2">
          <button
            v-if="item.decline"
            class="btn btn-outline btn-sm btn-error"
            @click="item.decline"
          >
            {{ $t('Deny') }}
          </button>
          <button v-if="item.accept" class="btn btn-sm" @click="item.accept">
            {{ $t('Accept') }}
          </button>
          <button v-if="item.see" class="btn btn-sm" @click="item.see">{{ $t('See') }}</button>
        </div>
      </div>
    </transition-group>
  </div>
</template>
