<script lang="ts">
  export default {
    name: 'Plan',
  }
</script>

<script setup lang="ts">
  import { $formatPriceByLocale, $formatInt } from 'utils/formats'
  import { $formatDay } from 'utils/date'
  import { $modals } from 'utils/modals'
  import { $t } from 'i18n/index'
  import { EPlans } from 'types/enums'
  import { $store } from 'utils/store'
  import { computed } from 'vue'
  import { $me } from 'utils/me'

  const { data } = defineProps<{
    data: Plan
  }>()

  const isRetailPrice = computed(() => {
    return $store.get('cfg').default_ref_uid === $me.data._ref_uid
  })

  const getPrice = computed(() => {
    if (isRetailPrice.value) return data.retail_price
    return data.price
  })

  const isPriorityMin = computed(() => {
    return data.priority <= 300
  })

  const isPriorityMax = computed(() => {
    return data.priority > 300
  })

  // function getColor(code: string) {
  //   if (code === 'start') $planColors['sstart']
  //   return $planColors[code]
  // }

  function buy() {
    $modals.paySystem.show({
      type: 'tariff',
      what: $t('BuyTariff') + ' ' + $t(`tariffs.${data.code}.name`),
      price: getPrice.value,
      currency_code: data.currency_code as unknown as string,
      precision: 2,
      code: data.code as unknown as string,
    })
  }
</script>

<template>
  <div
    class="flex flex-col justify-between relative rounded-lg drop-shadow-custom bg-base-100 text-center"
  >
    <div class="p-4">
      <div v-if="isPriorityMax" class="absolute top-2 right-2 badge badge-secondary">VIP</div>

      <h2 :class="`text-2xl font-semibold text-center`">
        <!-- :style="`color: ${getColor(data.code)}`"" -->
        {{ $t(`tariffs.${data.code}.name`) }}
      </h2>

      <div v-if="isPriorityMax" class="badge badge-error bg-red-400">
        {{ $t(`tariffs.${data.code}.sale`) }}
      </div>
      <!-- <p class="mt-2 text-sm text-base-content">
        {{ $t(`tariffs.${data.code}.description`) }}
      </p> -->
    </div>
    <div>
      <div class="border-y border-y-solid border-gray-600 bg-base-200">
        <div class="p-4">
          <h4 class="text-xs font-semibold tracking-wide uppercase text-base-content text-left">
            {{ $t('WHATS_INCLUDED') }}:
          </h4>
          <ul class="mt-4 space-y-4">
            <template v-if="isPriorityMin">
              <!-- <li class="flex space-x-3">
                <svg
                  class="flex-shrink-0 h-5 w-5 text-green-500"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path
                    fill-rule="evenodd"
                    d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                    clip-rule="evenodd"
                  ></path>
                </svg>
                <span class="text-xs leading-normal text-left">
                  {{ $t(`tariffs.${data.code}.features.bot`) }}
                </span>
              </li> -->
              <li v-if="data.bot_count > 0" class="flex space-x-3">
                <svg
                  class="flex-shrink-0 h-5 w-5 text-green-500"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path
                    fill-rule="evenodd"
                    d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                    clip-rule="evenodd"
                  ></path>
                </svg>
                <span class="text-xs leading-normal text-left">
                  {{
                    $t(`tariffs.features.bot`, {
                      count: data.bot_count,
                    })
                  }}
                </span>
              </li>
            </template>
            <template v-else>
              <li class="flex space-x-3">
                <svg
                  class="flex-shrink-0 h-5 w-5 text-green-500"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path
                    fill-rule="evenodd"
                    d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                    clip-rule="evenodd"
                  ></path>
                </svg>
                <span class="text-xs leading-normal text-left">
                  {{ $t(`tariffs.${data.code}.features.prev`) }}
                </span>
              </li>
            </template>

            <li class="flex space-x-3">
              <svg
                class="flex-shrink-0 h-5 w-5 text-green-500"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                  clip-rule="evenodd"
                ></path>
              </svg>
              <span class="text-xs leading-normal text-left">
                {{
                  $t(`tariffs.${data.code}.features.cv`, {
                    count: $formatInt(data.cv),
                  })
                }}
              </span>
            </li>
            <li class="flex space-x-3">
              <svg
                class="flex-shrink-0 h-5 w-5 text-green-500"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                  clip-rule="evenodd"
                ></path>
              </svg>
              <span class="text-xs leading-normal text-left">
                {{
                  $t(`tariffs.${data.code}.features.period`, {
                    period: $formatDay(data.period),
                  })
                }}
                {{ $t('Days', $formatDay(data.period)) }}
              </span>
            </li>

            <li v-if="data.max_deposit > 0 && isPriorityMin" class="flex space-x-3">
              <svg
                class="flex-shrink-0 h-5 w-5 text-green-500"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                  clip-rule="evenodd"
                ></path>
              </svg>
              <span class="text-xs leading-normal text-left">
                {{
                  $t(`tariffs.features.deposit`, {
                    amount: $formatPriceByLocale({
                      count: $formatInt(data.max_deposit),
                      currency: data.currency_code,
                    }),
                  })
                }}
              </span>
            </li>
          </ul>
        </div>
      </div>

      <div class="p-4">
        <p class="text-2xl font-extrabold relative py-2">
          <span
            v-if="!isRetailPrice"
            class="absolute bottom-full -mb-2 left-1/2 -translate-x-1/2 line-through text-sm font-medium text-error"
          >
            {{
              $formatPriceByLocale({
                count: $formatInt(data.retail_price),
                currency: data.currency_code,
              })
            }}
          </span>
          {{
            $formatPriceByLocale({
              count: $formatInt(getPrice),
              currency: data.currency_code,
            })
          }}
          <!-- <span class="text-base font-medium text-base-content">/mo</span> -->
        </p>

        <button class="btn mt-4 w-full btn-primary" @click="buy">
          {{ $t('Buy') }}
        </button>
      </div>
    </div>
  </div>
</template>
