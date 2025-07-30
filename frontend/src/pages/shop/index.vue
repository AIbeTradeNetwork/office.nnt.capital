<script lang="ts">
  export default {
    name: 'PageTariffsCustomer',
  }
</script>

<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { $modals } from 'utils/modals'
  import { $t } from 'i18n/index'
  import { $store } from 'utils/store'
  import { ECurrencies } from 'types/enums'
  import { $formatInt, $formatPriceByLocale } from 'utils/formats'
  import { $formatDay } from 'utils/date'
  import { $historBack } from 'utils/history'
  import { $textSplit } from 'utils/text'

  import coinNftSrcMini from 'assets/nft/coin-mini.gif'
  import { $me } from 'utils/me'
  import { $localStorage } from 'utils/localStorage'

  const loading = ref(true)
  const selectedABT = ref(false)
  const products = computed(() => $store.get('products'))

  function getSubscriptionName(code: string): string {
    const names: { [key: string]: string } = {
      'subscription_researcher': 'Researcher',
      'subscription_start': 'Start',
      'subscription_advanced': 'Advanced',
      'subscription_professional': 'Professional',
      'subscription_ambassador': 'Ambassador',
      'subscription_leader': 'Leader',
      'subscription_vip': 'VIP',
    }
    return names[code] || code
  }

  const voucherNFTAmbassador = {
    regexp: new RegExp(/boost_x\d{1,2}_lifetime_.+$/g),
    isProduct: computed(() => (product) => product.code.match(voucherNFTAmbassador.regexp)),
    isBuyed: computed(
      () => (product) =>
        ($me.data.products.filter((p) => p.code.match(voucherNFTAmbassador.regexp)) || []).length,
    ),
    getName: computed(() => (product) => {
      if (voucherNFTAmbassador.isProduct.value(product)) {
        return $t('products.category.boost_x50_lifetime')
      }
      return null
    }),
  }

  const bootAnchor = {
    regexp: new RegExp(/boost_archon$/g),
    isProduct: computed(() => (product) => product.code.match(voucherNFTAmbassador.regexp)),
    isBuyed: computed(
      () => (product) =>
        ($me.data.products.filter((p) => p.code.match(bootAnchor.regexp)) || []).length,
    ),
  }

  const isInf = computed(
    () => (product) =>
      product.period > 0 ?
        Math.trunc((Number($formatDay(product.period)) || 0) / 365) >= 100
      : null,
  )

  const text = computed(() => {
    return function (product: Product) {
      const text = []

      if (product.category.match(/boost_archon/)) {
        text.push(...$textSplit($t('products.category.boost_archon_text')))
        return text
      }

      if (product.category.match(/boost_x/) && voucherNFTAmbassador.isProduct.value(product)) {
        text.push(
          ...$textSplit(
            $t('products.category.boost_x_lifetime.text', {
              count: product.multiplier,
              day: '∞',
            }),
          ),
        )
        return text
      }

      if (product.category.match(/boost_x/)) {
        text.push(
          ...$textSplit(
            $t('products.category.boost_x.text', {
              count: product.multiplier,
              day: $formatDay(product.period),
            }),
          ),
        )
        return text
      }

      if (product.category === 'premium') {
        text.push(
          ...$textSplit(
            $t('products.text.premium.text', {
              multiplier: product.multiplier,
            }),
          ),
        )
        return text
      }

      if (product.category === 'autofarm') {
        text.push(...$textSplit($t('products.text.autofarm.text')))
        return text
      }

      if (product.category.match(/safe_\dth_digit/g)) {
        const digit = Number(product.category.replace('safe_', '').replace('th_digit', '')) || 0
        text.push(
          ...$textSplit($t('products.text.safe_key.text', { digit: digit, digitprev: digit - 1 })),
        )
        return text
      }

      return text
    }
  })

  const filteredProducts = computed(() => {
    return products.value.filter((product) => {
      if (selectedABT.value && product.currency_code === ECurrencies.udex) return true
      if (!selectedABT.value && product.currency_code === ECurrencies.usdHidden) return true
      return null
    })
  })

  function buy(product: Product) {
    // if (voucherNFTAmbassador.isBuyed.value(product)) return

    const name = () => {
      if (product.code.match(/_week_/)) return $t('Week', 1) + ' '
      if (product.code.match(/_month_/)) return $t('Month', 1) + ' '
      if (product.code.match(/_year_/)) return $t('Year', 1) + ' '
      return ''
    }

    $modals.paySystem.show({
      type: 'product',
      what: name() + '' + $t(`products.category.${product.category}`),
      code: product.code,
      currency_code: selectedABT.value ? ECurrencies.udex : ECurrencies.usdHidden,
      onSuccess: () => {
        $store.updateProducts()
      },
    })
  }

  async function init() {
    loading.value = true
    await $store.updateProducts()
    loading.value = false
  }

  onMounted(() => {
    init()
  })
</script>

<template>
  <div class="page relative">
    <div v-if="loading" class="text-center font-bold my-10 text-info">
      {{ $t('Loading') }}
    </div>

    <template v-else>
      <div class="absolute top-4 right-3">
        <button class="btn btn-outline btn-sm" @click="$historBack">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="relative -top-[1px] size-4 -ml-1 -mr-1"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M9 15 3 9m0 0 6-6M3 9h12a6 6 0 0 1 0 12h-3"
            />
          </svg>
          {{ $t('Back') }}
        </button>
      </div>

      <div class="flex items-center justify-center mt-1">
        <div class="text-center text-3xl">{{ $t('Boosts') }}</div>
      </div>

      <!-- Удалён переключатель валют -->
      <!--
      <div class="form-control mb-8">
        <label class="cursor-pointer label inline-flex self-center">
          <span
            class="label-text font-bold text-lg mr-2"
            :class="`${!selectedABT ? 'text-gold' : ''}`"
          >
            {{ ECurrencies.usd.toUpperCase() }}
          </span>
          <input v-model="selectedABT" type="checkbox" class="toggle toggle-lg toggle-primary" />
          <span
            class="label-text font-bold text-lg ml-2"
            :class="`${selectedABT ? 'text-gold' : ''}`"
          >
            {{ ECurrencies.abt.toUpperCase() }}
          </span>
        </label>
      </div>
      -->

      <div class="grid mobile:grid-cols-2 pc:grid-cols-4 mobile:gap-2 pc:gap-4">
        <div
          v-for="product in filteredProducts"
          :key="product.code"
          class="box flex flex-col justify-between w-full bg-base-200 relative gap-4"
        >
          <div>
            <div class="flex items-center flex-wrap gap-2 mb-2">
              <div class="font-bold uppercase flex-0 flex gap-2 items-center">
                <div v-if="voucherNFTAmbassador.isProduct.value(product)">
                  <div class="overflow-hidden rounded-full max-w-10">
                    <img :src="coinNftSrcMini" alt="" class="w-full" />
                  </div>
                </div>
                {{
                  voucherNFTAmbassador.getName.value(product) !== null ?
                    voucherNFTAmbassador.getName.value(product)
                  : product.category === 'subscription' ?
                    getSubscriptionName(product.code)
                  : $t(`products.category.${product.category}`)
                }}
              </div>
            </div>

            <div v-if="isInf(product) !== null" class="text-sm">
              <span class="opacity-50 mr-1">{{ $t('TimeToUse') }}:</span>
              <b class="text-gradient text-lg">
                <template v-if="isInf(product)">
                  {{ isInf(product) ? $t('Infinite') : $t('Forever') }}
                </template>
                <template v-else-if="product.period && product.period > 0">
                  {{ $formatDay(product.period) }} {{ $t('Days', $formatDay(product.period)) }}
                </template>
                <template v-else-if="product.category === 'subscription'">
                  {{ $t('Forever') }}
                </template>
                <template v-else>-</template>
              </b>
            </div>

            <div v-if="product.category === 'subscription'" class="text-sm">
              <span class="opacity-50 mr-2">Партнёров:</span>
              <b class="text-gradient text-lg">{{ product.multiplier }}</b>
            </div>
            <!-- <div v-if="product.limit > 0" class="text-sm">
              <span class="opacity-50 mr-2">{{ $t('Quantity') }}:</span>
              <b class="text-gradient text-lg">{{ product.count }} / {{ product.limit }}</b>
            </div> -->

            <ul class="mx-auto space-y-2 mt-4" role="">
              <li v-for="t in text(product)" :key="t" class="flex text-xs opacity-80">
                {{ t }}
              </li>
            </ul>
          </div>

          <div>
            <div
              v-if="product.price < product.retail_price"
              class="text-error whitespace-nowrap line-through -mb-1"
            >
              {{
                $formatPriceByLocale({
                  count: $formatInt(product.retail_price, {
                    precision: product.precision,
                  }),
                  currency: product.currency_code,
                })
              }}
            </div>

            <div class="text-xl font-bold text-gold">
              {{
                $formatPriceByLocale({
                  count: $formatInt(product.price, {
                    precision: product.precision,
                  }),
                  currency: product.currency_code,
                })
              }}
            </div>

            <template
              v-if="
                (voucherNFTAmbassador.isProduct.value(product) &&
                  voucherNFTAmbassador.isBuyed.value(product)) ||
                (bootAnchor.isBuyed.value(product) && bootAnchor.isProduct.value(product))
              "
            >
              <button class="btn btn-sm btn-disabled btn-block mt-4">
                {{ $t('Buyed') }}
              </button>
            </template>

            <template v-else>
              <button class="btn btn-sm btn-primary btn-block mt-4" @click="buy(product)">
                {{ $t('Buy') }}
              </button>
            </template>
          </div>

          <div
            v-if="product.price < product.retail_price"
            class="absolute top-0 right-2 rotate-6 badge badge-outline badge-info whitespace-nowrap -translate-y-1/4 font-bold"
          >
            {{ $t('Sale') }}
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
