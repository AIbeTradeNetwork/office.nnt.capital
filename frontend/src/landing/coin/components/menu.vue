<script setup lang="ts">
  import { $t } from 'i18n/index'
  import { $modals } from 'utils/modals'
  import { ECoinFarmingActions } from 'types/enums'

  function showInfo(type: ECoinFarmingActions) {
    $modals.farmingInfo.show(type)
  }
</script>

<template>
  <div class="sticky w-full max-w-[calc(500px-55px)] bottom-3 z-[2] mobile:p-2 pc:p-4">
    <div
      class="mobile:gap-1 flex justify-between items-center text-base-content px-4 box bg-opacity-90 h-[55px]"
    >
      <div
        class="inline-flex flex-col items-center gap-1 cursor-pointer"
        @click="$router.push({ name: 'Shop' })"
      >
        <!-- <span class="inline-flex items-center text-sm">&nbsp;</span> -->
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 feather feather-shopping-cart"><circle cx="9" cy="21" r="1"></circle><circle cx="20" cy="21" r="1"></circle><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path></svg>
        <!-- <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-5"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M15.75 10.5V6a3.75 3.75 0 1 0-7.5 0v4.5m11.356-1.993 1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 0 1-1.12-1.243l1.264-12A1.125 1.125 0 0 1 5.513 7.5h12.974c.576 0 1.059.435 1.119 1.007ZM8.625 10.5a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm7.5 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
        />
      </svg> -->

        <span class="text-[10px]">{{ $t('Boosts') }}</span>
      </div>

      <!-- <div
        class="inline-flex flex-col items-center gap-1 cursor-pointer"
        @click="showInfo(ECoinFarmingActions.farming)"
      >
        <svg
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          xmlns:xlink="http://www.w3.org/1999/xlink"
          viewBox="0 0 122.49 122.88"
          xml:space="preserve"
          fill="currentColor"
          class="w-5 h-5"
        >
          <path
            fill-rule="evenodd"
            clip-rule="evenodd"
            class="st0"
            d="M101.12,37.47c14.95,18.54,22.23,40.44,21.28,60.48c-6.91-16.93-17.64-34.09-31.87-49.9l-4.77,4.77 c-0.54,0.54-1.42,0.54-1.96,0L68.72,37.75c-0.54-0.54-0.54-1.42,0-1.96l4.63-4.63C57.16,17.12,39.68,6.67,22.54,0.2 c20.2-1.52,42.5,5.45,61.44,20.33l2.09-2.09c0.54-0.54,1.42-0.54,1.96,0l15.08,15.08c0.54,0.54,0.54,1.42,0,1.96L101.12,37.47 L101.12,37.47z M68.16,42.51l12.22,12.22l-65.64,65.64c-3.36,3.36-8.86,3.36-12.22,0l0,0c-3.36-3.36-3.36-8.86,0-12.22L68.16,42.51 L68.16,42.51z"
          />
        </svg>
        <span class="text-[10px]">{{ $t('coin.farming.Farming') }}</span>
      </div> -->

      <div
        class="inline-flex flex-col items-center gap-1 cursor-pointer"
        @click="showInfo(ECoinFarmingActions.invite)"
      >
        <!-- <span class="inline-flex items-center text-sm">
        +1.0
        <img width="12px" :src="symbolSrc" alt="symbol" class="ml-[2px]" />
      </span> -->
        <svg
          class="w-5 h-5"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          fill="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            fill-rule="evenodd"
            d="M9 4a4 4 0 1 0 0 8 4 4 0 0 0 0-8Zm-2 9a4 4 0 0 0-4 4v1a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2v-1a4 4 0 0 0-4-4H7Zm8-1a1 1 0 0 1 1-1h1v-1a1 1 0 1 1 2 0v1h1a1 1 0 1 1 0 2h-1v1a1 1 0 1 1-2 0v-1h-1a1 1 0 0 1-1-1Z"
            clip-rule="evenodd"
          />
        </svg>
        <span class="text-[10px]">{{ $t('coin.farming.Invite') }}</span>
      </div>

      <router-link
        :to="{ name: 'Autotrade' }"
        class="box-shadow inline-flex flex-col justify-center items-center gap-1 cursor-pointer rounded-full bg-base-100 h-[85px] w-[85px] drop-shadow-xl"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-6 -mt-2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M2.25 18 9 11.25l4.306 4.306a11.95 11.95 0 0 1 5.814-5.518l2.74-1.22m0 0-5.94-2.281m5.94 2.28-2.28 5.941"
          />
        </svg>

        <span class="text-[10px]">{{ $t('autotrade.title') }}</span>
      </router-link>

      <div
        class="inline-flex flex-col items-center gap-1 cursor-pointer"
        @click="showInfo(ECoinFarmingActions.team)"
      >
        <!-- <span class="inline-flex items-center text-sm">
        +0.33
        <img width="12px" :src="symbolSrc" alt="symbol" class="ml-[2px]" />
      </span> -->
        <svg
          class="w-5 h-5"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          fill="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            fill-rule="evenodd"
            d="M12 6a3.5 3.5 0 1 0 0 7 3.5 3.5 0 0 0 0-7Zm-1.5 8a4 4 0 0 0-4 4 2 2 0 0 0 2 2h7a2 2 0 0 0 2-2 4 4 0 0 0-4-4h-3Zm6.82-3.096a5.51 5.51 0 0 0-2.797-6.293 3.5 3.5 0 1 1 2.796 6.292ZM19.5 18h.5a2 2 0 0 0 2-2 4 4 0 0 0-4-4h-1.1a5.503 5.503 0 0 1-.471.762A5.998 5.998 0 0 1 19.5 18ZM4 7.5a3.5 3.5 0 0 1 5.477-2.889 5.5 5.5 0 0 0-2.796 6.293A3.501 3.501 0 0 1 4 7.5ZM7.1 12H6a4 4 0 0 0-4 4 2 2 0 0 0 2 2h.5a5.998 5.998 0 0 1 3.071-5.238A5.505 5.505 0 0 1 7.1 12Z"
            clip-rule="evenodd"
          />
        </svg>
        <span class="text-[10px]">{{ $t('coin.farming.Team') }}</span>
      </div>

      <div
        class="inline-flex flex-col items-center gap-1 cursor-pointer"
        @click="showInfo(ECoinFarmingActions.tasks)"
      >
        <!-- <span class="inline-flex items-center text-sm">
        +15.0
        <img width="12px" :src="symbolSrc" alt="symbol" class="ml-[2px]" />
      </span> -->
      <svg
          class="w-5 h-5 text-green-500"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          fill="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            fill-rule="evenodd"
            d="M12 6a3.5 3.5 0 1 0 0 7 3.5 3.5 0 0 0 0-7Zm-1.5 8a4 4 0 0 0-4 4 2 2 0 0 0 2 2h7a2 2 0 0 0 2-2 4 4 0 0 0-4-4h-3Zm6.82-3.096a5.51 5.51 0 0 0-2.797-6.293 3.5 3.5 0 1 1 2.796 6.292ZM19.5 18h.5a2 2 0 0 0 2-2 4 4 0 0 0-4-4h-1.1a5.503 5.503 0 0 1-.471.762A5.998 5.998 0 0 1 19.5 18ZM4 7.5a3.5 3.5 0 0 1 5.477-2.889 5.5 5.5 0 0 0-2.796 6.293A3.501 3.501 0 0 1 4 7.5ZM7.1 12H6a4 4 0 0 0-4 4 2 2 0 0 0 2 2h.5a5.998 5.998 0 0 1 3.071-5.238A5.505 5.505 0 0 1 7.1 12Z"
            clip-rule="evenodd"
          />
        </svg>

        <span class="text-[10px]">{{ $t('coin.farming.Tasks') }}</span>
      </div>
    </div>
  </div>
</template>
