<script lang="ts">
  export default {
    name: 'PageTariffsCustomer',
  }
</script>

<script setup lang="ts">
  import Info from 'components/info/index.vue'
  import Tariff from 'components/tariff/index.vue'
  import { onMounted, ref, Ref } from 'vue'
  import { $requests } from 'queries/index'
  import { ERoles, ESizes } from 'types/enums'
  import { $modals } from 'utils/modals'
  import { $t } from 'i18n/index'
  import { $config } from 'utils/configuration'

  const loading = ref(true)
  const list: Ref<Plan[]> = ref([])

  async function getData() {
    loading.value = true

    try {
      const response = await $requests.plans.get()
      if (response) list.value = response
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  function openModalPromocode() {
    if (loading.value) return
    $modals.promocode.show()
  }

  onMounted(() => {
    getData()
  })
</script>

<template>
  <div class="page">
    <div v-if="$config.role === ERoles.client">
      <button class="btn btn-primary" :disabled="loading" @click="openModalPromocode">
        <span>{{ $t('Promocode') }}</span>
      </button>
    </div>

    <div class="text-center text-3xl my-4">{{ $t('TariffTitle') }}</div>

    <div class="grid mobile:grid-cols-1 pc:grid-cols-4 mobile:gap-2 pc:gap-6">
      <Tariff v-for="tariff in list" :key="tariff.code" :data="tariff" />
    </div>

    <Info :size="ESizes.SM" :message="$t('info.tariffs')" />
  </div>
</template>
