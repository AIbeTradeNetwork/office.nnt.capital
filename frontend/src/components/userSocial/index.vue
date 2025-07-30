<script lang="ts">
  export default {
    name: 'AvatarNav',
  }
</script>

<script setup lang="ts">
  import { ref } from 'vue'
  import avatar from 'src/assets/demo/social.jpg'
  import { saveImage } from './svgToImg'
  import { $formatPriceByLocale } from 'utils/formats'
  import { $me } from 'utils/me'

  const loading = ref(false)

  const props = withDefaults(
    defineProps<{
      download: boolean
    }>(),
    {
      download: true,
    },
  )

  async function saveImg() {
    if (loading.value) return

    loading.value = true

    try {
      await saveImage('socialBlock')
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }
</script>

<template>
  <div class="bg-base-300 rounded-lg w-[300px]">
    <div
      id="socialBlock"
      class="insta-bg relative p-4 pb-10 overflow-hidden rounded-lg pointer-events-none"
    >
      <div class="relative aspect-[9/16] shadow-2xl">
        <img :src="avatar" class="block aspect-[9/16] object-cover rounded-lg" />
      </div>
      <div
        class="absolute bottom-4 left-5 right-5 uppercase font-bold text-center tracking-wide text-white p-3 px-6 bg-black rounded-md rounded-tl-3xl rounded-br-3xl inline-block shadow-2xl backdrop-blur-md bg-opacity-80"
      >
        {{ $me.data.nickname }}
      </div>

      <div
        class="absolute bottom-[80px] right-0 flex flex-col tracking-wide rounded-xl rounded-r-none overflow-hidden shadow-2xl"
      >
        <div class="bg-black text-white px-6 pr-4 py-1 text-center backdrop-blur-md bg-opacity-80">
          {{ $t('Profitability') }}
        </div>
        <div
          class="bg-white px-6 pr-4 py-4 text-black font-bold text-lg text-center backdrop-blur-md bg-opacity-80"
        >
          {{ $formatPriceByLocale({ count: 100500 }) }}
        </div>
      </div>

      <div
        class="absolute bottom-[188px] right-0 flex flex-col tracking-wide rounded-xl rounded-r-none overflow-hidden shadow-lg"
      >
        <div class="bg-black text-white px-6 pr-4 py-1 text-center backdrop-blur-md bg-opacity-80">
          {{ $t('Rank') }}
        </div>
        <div
          class="bg-white px-6 pr-4 py-4 text-black font-bold text-lg text-center backdrop-blur-md bg-opacity-80"
        >
          Savvy
        </div>
      </div>
    </div>
    <div v-if="props.download" class="p-4 flex justify-center">
      <button class="btn btn-sm btn-primary" @click="saveImg">
        <span v-if="loading" class="loading loading-spinner"></span>
        <svg
          v-else
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
            d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m.75 12 3 3m0 0 3-3m-3 3v-6m-1.5-9H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
          />
        </svg>
        {{ $t('Download') }}
      </button>
    </div>
  </div>
</template>
