<script lang="ts">
  export default {
    name: 'PageDiscovery',
  }
</script>

<script setup lang="ts">
  import { $t } from 'i18n/index'
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
  import { ref } from 'vue'
  import { $textSplit } from 'utils/text'

  const modalData = ref<(typeof apps)[0] | null>()

  const apps = [
    {
      name: $t('discovery.app.nfterrium.name'),
      description: $t('discovery.app.nfterrium.description'),
      img: new URL('../../assets/apps/nfterrium/cover.jpeg', import.meta.url).href,
      url: $t('discovery.app.nfterrium.url'),
      soon: false,
    },
    // {
    //   name: $t('discovery.app.cardgame.name'),
    //   description: $t('discovery.app.cardgame.description'),
    //   img: getUrl('/src/assets/apps/cardgame/Cover.png'),
    //   url: 'https://aibetrade.com',
    //   soon: true,
    // },
    {
      name: '',
      description: '',
      img: undefined,
      url: undefined,
      soon: true,
    },
  ]

  function getUrl(url: string): string {
    return new URL(url, import.meta.url).href
  }

  function openInfo(app: (typeof apps)[0]) {
    if (app.soon) return
    modalData.value = app
    $modals.any.onClose = () => (modalData.value = null)
    $modals.any.show()
  }
</script>

<template>
  <div class="page p-4">
    <div class="flex flex-wrap gap-4">
      <template v-for="app in apps" :key="app.name">
        <div class="w-[128px] cursor-pointer" @click="openInfo(app)">
          <div
            class="w-[128px] h-[128px] bg-black rounded-md overflow-hidden bg-center bg-cover relative flex justify-center items-center"
            :style="`background-image: url(${app.img})`"
          >
            <div v-if="app.soon" class="badge badge-sm badge-info absolute top-1 right-1">
              {{ $t('Soon') }}
            </div>
            <span v-if="app.soon">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-[64px]"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 5.25h.008v.008H12v-.008Z"
                />
              </svg>
            </span>
          </div>
          <div class="text-center mt-2 text-sm">{{ app.name }}</div>
        </div>
      </template>
    </div>

    <Modal
      v-if="$modals.any.active"
      :modal="$modals.any"
      :class="'[&_.modal-box]:max-w-[700px]'"
      :z="50"
    >
      <div class="flex gap-4">
        <div
          class="w-[128px] h-[128px] flex-[0_0_128px] bg-black rounded-md overflow-hidden bg-center bg-cover"
          :style="`background-image: url(${modalData?.img})`"
        ></div>
        <div>
          <div>
            <div class="title m-0 mb-2">{{ modalData?.name || '' }}</div>
            <div class="text">
              <p v-for="par in $textSplit(modalData?.description || '')" :key="par">{{ par }}</p>
            </div>
          </div>
          <template v-if="modalData.soon">
            <div class="btn btn-sm mt-3 btn-info flex-[0_0_100%]">{{ $t('Soon') }}</div>
          </template>
          <template v-else>
            <div v-if="modalData?.url" class="link mt-4">
              <!-- {{ $t('discovery.app.siteUrl') }}: -->
              <a class="btn btn-success" :href="'https://' + modalData?.url" target="_blank">
                {{ modalData?.url }}
                <svg
                  width="12"
                  height="12"
                  viewBox="0 0 48 48"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M19 11H37V29"
                    stroke="currentColor"
                    stroke-width="4"
                    stroke-linecap="butt"
                    stroke-linejoin="bevel"
                  ></path>
                  <path
                    d="M11.5439 36.4559L36.9997 11"
                    stroke="currentColor"
                    stroke-width="4"
                    stroke-linecap="butt"
                    stroke-linejoin="bevel"
                  ></path>
                </svg>
              </a>
            </div>
          </template>
        </div>
      </div>
    </Modal>
  </div>
</template>
