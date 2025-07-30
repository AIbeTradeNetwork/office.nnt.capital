<script lang="ts">
export default {
  name: 'ModalSubscriptions',
}
</script>

<script setup lang="ts">
import Modal from 'components/modals/modal/index.vue'
import { $t } from 'i18n/index'
import { $modals } from 'utils/modals'
import { computed, onMounted, ref } from 'vue'
import { $store } from 'utils/store'
import { $formatInt, $formatPriceByLocale } from 'utils/formats'

const loading = ref(true)
const products = computed(() => $store.get('products'))
const filteredProducts = computed(() => {
  return products.value.filter(product => product.category === 'subscription')
})

const selectedProduct = ref(null)

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

function getSubscriptionDescription(code: string): string {
  const descriptions: { [key: string]: string } = {
    'subscription_researcher': 'Базовый абонемент для всех пользователей',
    'subscription_start': 'Стартовый абонемент для активных пользователей',
    'subscription_advanced': 'Продвинутый абонемент для опытных пользователей',
    'subscription_professional': 'Профессиональный абонемент для экспертов',
    'subscription_ambassador': 'Амбассадорский абонемент для лидеров',
    'subscription_leader': 'Лидерский абонемент для топ-менеджеров',
    'subscription_vip': 'VIP абонемент для элитных пользователей',
  }
  return descriptions[code] || 'Описание недоступно'
}

function buy(product: Product) {
  $modals.paySystem.show({
    type: 'product',
    what: getSubscriptionName(product.code),
    code: product.code,
    currency_code: product.currency_code,
    onSuccess: $modals.farmingSubscription.onSuccess,
  })
}

onMounted(async () => {
  loading.value = true
  await $store.updateProducts()
  loading.value = false
})
</script>

<template>
  <Modal :modal="$modals.farmingSubscription" :z="50" :title="$t('Premium')">
    <div v-if="loading" class="text-center font-bold my-10 text-info">{{ $t('Loading') }}</div>
    <div v-else>
      <div class="grid gap-4 grid-cols-1 sm:grid-cols-2">
        <div
          v-for="product in filteredProducts"
          :key="product.code"
          class="box flex flex-col items-center bg-base-200 relative gap-2 cursor-pointer border-2"
          :class="selectedProduct && selectedProduct.code === product.code ? 'border-primary' : 'border-transparent'"
          @click="selectedProduct = product"
        >
          <div class="text-xl font-bold text-gradient">{{ getSubscriptionName(product.code) }}</div>
          <div class="text-2xl font-bold text-gold">
            {{ $formatPriceByLocale({ count: $formatInt(product.price, { precision: product.precision }), currency: product.currency_code }) }}
          </div>
          <div class="text-sm">
            <span class="opacity-50 mr-2">Партнёров:</span>
            <b class="text-gradient text-lg">{{ product.multiplier }}</b>
          </div>
                     <div class="text-xs text-center opacity-80 mt-2">{{ getSubscriptionDescription(product.code) }}</div>
          <button class="btn btn-sm w-full btn-primary mt-2" @click.stop="buy(product)">
            {{ $t('Buy') }}
          </button>
        </div>
      </div>
      <div v-if="selectedProduct" class="mt-6 p-4 bg-base-100 rounded-lg shadow">
                 <div class="text-lg font-bold mb-2 text-gradient">{{ getSubscriptionName(selectedProduct.code) }}</div>
         <div class="mb-2">{{ getSubscriptionDescription(selectedProduct.code) }}</div>
        <div class="mb-2">
          <span class="opacity-50 mr-2">Партнёров:</span>
          <b class="text-gradient text-lg">{{ selectedProduct.multiplier }}</b>
        </div>
        <div class="mb-2">
          <span class="opacity-50 mr-2">Цена:</span>
          <b class="text-gold">{{ $formatPriceByLocale({ count: $formatInt(selectedProduct.price, { precision: selectedProduct.precision }), currency: selectedProduct.currency_code }) }}</b>
        </div>
        <button class="btn btn-primary w-full mt-2" @click="buy(selectedProduct)">{{ $t('Buy') }}</button>
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.box {
  transition: border 0.2s;
}
.box:hover {
  border-color: var(--color-primary);
}
</style>
