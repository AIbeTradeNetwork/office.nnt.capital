<script lang="ts">
  export default {
    name: 'Exchanges',
  }
</script>

<script setup lang="ts">
  import { $modals } from 'utils/modals'
  import { $store } from 'utils/store'
  import { computed } from 'vue'

  const { data } = defineProps<{
    data: Exchange
  }>()

  const isActive = computed(() => {
    return !data.is_active
  })

  const keys = computed(() => {
    return $store.get('keys').filter((item) => item.exchange_code === data.code).length
  })

  function getImageUrl() {
    return new URL(`../../assets/exchanges/${data.code}.svg`, import.meta.url).href
  }

  function showModalKeys() {
    if (!isActive.value) return
    $modals.keys.show(data)
  }

  function getColorByCode() {
    if (data.code === 'bybit') return 'text-orange-300'
    return 'btn-info'
  }
</script>

<template>
  <div
    class="relative box p-0 overflow-hidden"
    :class="`${!isActive ? 'pointer-events-none' : ''}`"
  >
    <div
      v-if="!isActive"
      class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-error font-bold text-lg tracking-wide"
    >
      Offline
    </div>

    <div class="flex h-full items-center" :class="`${!isActive ? 'opacity-10' : ''}`">
      <div class="inline-flex items-stretch bg-black-300 h-full py-4 px-3 w-24 [&_img]:object-fill">
        <img :src="getImageUrl()" alt="" />
      </div>

      <div class="flex-1 p-4">
        <div class="text-lg font-bold">{{ data.name }}</div>
        <div>
          {{ $t('KeysActive') }}:
          <span class="text-primary">{{ keys }}</span>
        </div>

        <div class="mt-2 flex flex-col gap-4">
          <a
            v-if="data.link"
            :href="data.link"
            target="_blank"
            class="btn btn-sm btn-outline"
            :class="getColorByCode()"
          >
            {{ data.name }} {{ $t('Registration') }}
          </a>
          <button class="btn btn-primary btn-sm" @click.prevent="showModalKeys()">
            {{ $t('KeysEdit') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
