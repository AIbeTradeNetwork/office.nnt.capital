<script lang="ts">
  export default {
    name: 'PagePayments',
  }
</script>

<script setup lang="ts">
  import Info from 'components/info/index.vue'
  import { $requests } from 'queries/index'
  import { EAlerts, ESizes } from 'types/enums'
  import { $formatDate } from 'utils/date'
  import { $formatInt, $formatPriceByLocale } from 'utils/formats'
  import { onMounted, ref } from 'vue'
  import Payouts from 'components/payouts/index.vue'
  import { $store } from 'utils/store'

  const loading = ref(false)
  const page = ref(0)
  const limit = 20
  const isLastPage = ref(false)
  const list = ref<Transaction[]>([])

  async function getData() {
    if (loading.value) return

    loading.value = true

    try {
      const response = await $requests.transactions.transactions({
        limit,
        skip: page.value > 1 ? page.value * limit : 0,
      })

      if (response.length) {
        list.value.push(...response)
        page.value++
      }

      if (!response.length || response.length < limit) isLastPage.value = true
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  function loadMore() {
    if (loading.value || isLastPage.value) return
    getData()
  }

  onMounted(() => {
    getData()
  })
</script>

<template>
  <div class="page">
    <div class="box p-0 overflow-x-auto">
      <table class="table table-zebra">
        <thead>
          <tr>
            <th>{{ $t('Date') }}</th>
            <th>{{ 'UID' }}</th>
            <th>{{ $t('FromUID') }}</th>
            <th>{{ $t('Type') }}</th>
            <th>{{ $t('Amount') }}</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in list" :key="item.buy_uid">
            <td>{{ $formatDate(item.charged_at, 'date|time') }}</td>
            <td>{{ item.buy_uid }}</td>
            <td>{{ item.from_uid }}</td>
            <td>
              <div class="badge badge-outline">{{ item.type }}</div>
            </td>
            <td>
              <b>
                {{
                  $formatPriceByLocale({
                    count: $formatInt(item.amount),
                  })
                }}
              </b>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="flex justify-center">
      <div class="join">
        <button v-if="!isLastPage" class="join-item btn" @click="loadMore">
          <span v-if="loading" class="loading loading-spinner"></span>
          <span v-else>{{ $t('ShowMore') }}</span>
        </button>
      </div>
    </div>

    <Payouts />
    <div class="grid mobile:gap-2 pc:gap-4 mobile:grid-cols-1 pc:grid-cols-3">
      <!-- <Info
        :type="EAlerts.WARNING"
        :size="ESizes.SM"
        :title="$t('info.pagePayments.rules.title')"
        :message="
          $t('info.pagePayments.rules.message', {
            feepercent: `${$formatInt($store.get('cfg').payout_fee_percent)}%`,
            feemin: `${$formatPriceByLocale({ count: $formatInt($store.get('cfg').payout_fee_min) })}`,
            amountmin: `${$formatPriceByLocale({ count: $formatInt($store.get('cfg').payout_amount_min) })}`,
          })
        "
      /> -->
      <Info :size="ESizes.SM" :message="$t('info.pagePayments.info')" />
    </div>
  </div>
</template>
