<script setup lang="ts">
  import { $i18nGlobal, $t } from 'i18n/index'
  import { computed } from 'vue'

  const howFarmList = computed(() => {
    return [
      {
        title: $t('coin.farm_info_1.title'),
        text: $t('coin.farm_info_1.text'),
      },
      {
        title: $t('coin.farm_info_2.title'),
        text: $t('coin.farm_info_2.text'),
      },
      {
        title: $t('coin.farm_info_3.title'),
        text: $t('coin.farm_info_3.text'),
      },
      {
        title: $t('coin.farm_info_4.title'),
        text: $t('coin.farm_info_4.text'),
      },
      //   {
      //     title: $t('coin.farm_info_5.title'),
      //     text: $t('coin.farm_info_5.text'),
      //   },
    ]
  })
</script>
<template>
  <div class="max-w-[600px] space-y-3">
    <h1 class="coin-title">
      <span class="text-gradient">{{ $t('coin.menu.how_farm') }}</span>
    </h1>

    <!-- <div class="flex place-content-center items-stretch mobile:gap-3 pc:gap-12 flex-wrap"> -->
    <!-- <div
      class="pc:flex-0 mobile:flex-[0_0_100%] mobile:max-w-[260px] inline-flex flex-col items-center justify-center gap-4"
      >
      <div class="coin-box flex w-full">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 122.88 92.5" class="w-12">
          <path
            fill-rule="evenodd"
            stroke="currentColor"
            fill="currentColor"
            d="M3.1,86.3H14.57A2.14,2.14,0,0,1,14,84.8v-47a2.15,2.15,0,0,1,2.15-2.15h90.1a2.16,2.16,0,0,1,2.15,2.15V84.8a2.19,2.19,0,0,1-.61,1.5h12a3.1,3.1,0,0,1,0,6.2H3.1a3.1,3.1,0,0,1,0-6.2Zm85-64.3H82v9.56H77.73V19.86a2.13,2.13,0,0,1,2.14-2.14h8.26A6.28,6.28,0,1,1,88.07,22ZM44,23.5v8.06H39.71V25.64H34.8a6.28,6.28,0,1,1,0-4.29h7.1a2.12,2.12,0,0,1,1.51.63A2.15,2.15,0,0,1,44,23.5ZM33.06,0a6.28,6.28,0,1,0,5.75,8.81H51.37V31.56h4.28V6.68a2.13,2.13,0,0,0-2.14-2.15H39.09a6.27,6.27,0,0,0-6-4.53ZM67.85.1a6.28,6.28,0,0,0-2.14,12.19V31.56H70V12.29A6.28,6.28,0,0,0,67.85.1ZM51.12,48a10.17,10.17,0,0,1,1.56-5.44H70.2A10.15,10.15,0,0,1,71.76,48a10.32,10.32,0,0,1-20.64,0Zm40.11.39v2.16a3.69,3.69,0,0,1-3.67,3.68H83.79a22.82,22.82,0,0,1-2.27,6L84,62.69a3.68,3.68,0,0,1,0,5.19l-3,3.06a3.69,3.69,0,0,1-5.2,0l-2.29-2.29a22.65,22.65,0,0,1-6.21,2.69v3.12a3.7,3.7,0,0,1-3.68,3.68H59.28a3.7,3.7,0,0,1-3.68-3.68V71.38a22.67,22.67,0,0,1-6.27-2.67L47.1,70.94a3.68,3.68,0,0,1-5.19,0l-3.06-3.06a3.68,3.68,0,0,1,0-5.19l2.4-2.41a22.63,22.63,0,0,1-2.32-6.09H35.32a3.69,3.69,0,0,1-3.67-3.68V46.18a3.68,3.68,0,0,1,3.5-3.66H47.46a14.87,14.87,0,0,0-.89,5.06c0,.26,0,.51,0,.77s0,.51,0,.76a14.79,14.79,0,1,0,29.58,0c0-.25,0-.51,0-.76s0-.51,0-.77a14.87,14.87,0,0,0-.89-5.06H87.73a3.68,3.68,0,0,1,3.5,3.66v2.17ZM17.64,39H105V84.2H17.64V39ZM57.9,87.83H65A1.21,1.21,0,0,1,66.18,89h0A1.21,1.21,0,0,1,65,90.24H57.9A1.21,1.21,0,0,1,56.7,89h0a1.2,1.2,0,0,1,1.2-1.2Z"
          />
        </svg>
        <div class="mx-3">
          <div class="text-gradient">Balance</div>
          <div>1.00010</div>
        </div>
      </div>

      <div class="coin-box flex w-full">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-12"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99"
          />
        </svg>
        <div class="mx-3">
          <div class="text-gradient">Mining</div>
          <div>1.00010 ABT/час</div>
        </div>
      </div>

      <div class="coin-box flex relative w-full">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 122.88 122.85" class="w-10">
          <path
            stroke="currentColor"
            fill="currentColor"
            d="M28.89,20.82h65a8.42,8.42,0,0,1,8.39,8.4V94A8.4,8.4,0,0,1,99.8,99.9h0a8.36,8.36,0,0,1-5.92,2.47h-65a8.43,8.43,0,0,1-8.4-8.4V29.22A8.37,8.37,0,0,1,23,23.3h0a8.37,8.37,0,0,1,5.92-2.47ZM45.47,37.51H77.29a8,8,0,0,1,8,8V77.71a8,8,0,0,1-8,8H45.47a8,8,0,0,1-8-8V45.49a8,8,0,0,1,8-8ZM91.88,89.16a2.94,2.94,0,1,1-2.94-2.93,2.94,2.94,0,0,1,2.94,2.93ZM36.75,34a2.94,2.94,0,1,1-2.94-2.93A2.94,2.94,0,0,1,36.75,34ZM88.94,31.1A2.94,2.94,0,1,1,86,34a2.94,2.94,0,0,1,2.93-2.93ZM33.81,86.23a2.94,2.94,0,1,1-2.93,2.93,2.94,2.94,0,0,1,2.93-2.93Zm.7,22.76v13.86h39.34V109Zm11.94,0v13.86H41.28V109Zm11.93,0v13.86H53.21V109Zm11.94,0v13.86H65.15V109Zm11.93,0v13.86H77.08V109Zm11.94,0v13.86H89V109ZM34.51,0V13.86h39.34V0ZM46.45,0V13.86H41.28V0ZM58.38,0V13.86H53.21V0ZM70.32,0V13.86H65.15V0ZM82.25,0V13.86H77.08V0ZM94.19,0V13.86H89V0ZM109,88.81h13.86V94H109V88.81Zm0-11.94h13.86V82H109V76.87Zm0-11.93h13.86v5.17H109V64.94ZM109,53h13.86v5.17H109V53Zm0-11.93h13.86v5.17H109V41.07Zm0-11.94h13.86V34.3H109V29.13ZM0,88.81H13.86V94H0V88.81ZM0,76.87H13.86V82H0V76.87ZM0,64.94H13.86v5.17H0V64.94ZM0,53H13.86v5.17H0V53ZM0,41.07H13.86v5.17H0V41.07ZM0,29.13H13.86V34.3H0V29.13Zm93.87-4.28h-65a4.39,4.39,0,0,0-3.09,1.27h0a4.34,4.34,0,0,0-1.27,3.09V94a4.39,4.39,0,0,0,4.37,4.38h65A4.35,4.35,0,0,0,97,97.07h0A4.35,4.35,0,0,0,98.24,94V29.22a4.39,4.39,0,0,0-4.37-4.37Z"
          />
        </svg>
        <div class="mx-3">
          <div class="text-gradient">Storage</div>
          <div>4h to fill</div>
        </div>
        <button class="btn btn-success">Collect</button>
      </div>

      <div class="bg-black-300 w-full rounded relative pc:-mt-2 mobile:-mt-1 mobile:mb-6">
        <div class="rounded w-[10%] h-full bg-success px-1 py-1"></div>
        <div class="absolute top-0 left-0 w-full h-full flex justify-center items-center">
          10% / 100%
        </div>
      </div>
    </div> -->

    <div class="space-y-3">
      <div
        v-for="item in howFarmList"
        :key="item.title"
        tabindex="0"
        class="collapse collapse-open bg-base-100"
      >
        <div class="collapse-title text-xl font-medium text-gradient">
          {{ item.title }}
        </div>
        <div class="collapse-content">
          <p>
            {{ item.text }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
