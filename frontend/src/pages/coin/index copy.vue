<script lang="ts">
  export default {
    name: 'Dashboard',
  }
</script>

<script setup lang="ts">
  // import ChartDefault from 'components/charts/default/index.vue'
  import { EPlans, EPlansVIP, ERangs } from 'types/enums'
  import { $config } from 'utils/configuration'

  const currentPlanVip = EPlans.advanced
  const currentRang = EPlans.advanced
  const nextRang = EPlans.advanced + 1

  function enumToArray(enumObj: any) {
    return Object.keys(enumObj)
      .filter((key) => isNaN(Number(key)))
      .map((key) => enumObj[key])
  }
</script>

<template>
  <div class="page">
    <!-- <div class="box">
      <div class="title flex justify-between">
        My qualification
        <router-link class="text-base link" :to="{ name: 'Index' }">
          Qualification history
        </router-link>
      </div>
    </div> -->

    <div class="box">
      <div class="stats w-full">
        <div class="stat content-between">
          <div class="stat-title">Приглашенных за месяц</div>
          <div class="stat-value flex gap-4 items-center text-secondary text-[60px] py-4">
            57
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-14"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M18 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0ZM3 19.235v-.11a6.375 6.375 0 0 1 12.75 0v.109A12.318 12.318 0 0 1 9.374 21c-2.331 0-4.512-.645-6.374-1.766Z"
              />
            </svg>
          </div>
          <div class="stat-desc">21% больше чем в прошлом месяце</div>
        </div>

        <div class="stat content-between">
          <div class="stat-title">Приглашенных за неделю</div>
          <div class="stat-value flex gap-4 items-center text-secondary text-[60px] py-4">
            57
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-14"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M18 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0ZM3 19.235v-.11a6.375 6.375 0 0 1 12.75 0v.109A12.318 12.318 0 0 1 9.374 21c-2.331 0-4.512-.645-6.374-1.766Z"
              />
            </svg>
          </div>
          <div class="stat-desc">21% больше чем на прошлой неделе</div>
        </div>

        <div class="stat content-between">
          <div class="stat-title">Приглашенные командой за месяц</div>
          <div class="stat-value flex gap-4 items-center text-secondary text-[60px] py-4">
            30
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-14"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z"
              />
            </svg>
          </div>
          <div class="stat-desc">21% больше чем в прошлом месяце</div>
        </div>
      </div>
    </div>

    <div class="flex gap-6">
      <div class="box flex flex-col items-center">
        <div class="title mb-4">До следующего ранга:</div>
        <div
          class="mx-auto text-primary radial-progress text-4xl"
          style="--value: 65; --size: 12rem; --thickness: 1.5rem"
          role="progressbar"
        >
          65%
        </div>
        <div class="text-center mt-4">
          Отсалось
          <span class="text-secondary">30%</span>
          до ранга
          <span class="font-bold text-primary">{{ ERangs[nextRang] }}</span>
        </div>
      </div>

      <div class="flex flex-col gap-6 flex-1">
        <div class="box flex items-center justify-between">
          <div>
            <div class="title mb-4">Текущий план VIP:</div>

            <div class="flex justify-between">
              <ul class="steps -ml-4">
                <li
                  v-for="item in enumToArray(EPlansVIP)"
                  :key="item"
                  class="step"
                  :class="currentPlanVip >= item ? 'step-primary' : ''"
                >
                  {{ EPlansVIP[item] }}
                </li>
              </ul>
            </div>
          </div>

          <div class="stat inline-grid p-0 w-auto">
            <div class="stat-title">Осталось дней:</div>
            <div class="stat-value text-secondary text-[60px] py-2">30</div>
          </div>
        </div>

        <div class="box">
          <div class="title mb-4">Текущий ранг:</div>
          <ul class="steps counter-steps -ml-4">
            <template v-for="item in enumToArray(ERangs)">
              <li
                v-if="!(item in EPlansVIP)"
                :key="item"
                class="step"
                :class="currentRang >= item ? 'step-primary' : ''"
              >
                {{ ERangs[item] }}
              </li>
            </template>
          </ul>
        </div>
      </div>
    </div>

    <!-- <ChartDefault /> -->
  </div>
</template>
