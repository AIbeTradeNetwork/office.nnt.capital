<script lang="ts">
  export default {
    name: 'Purchases',
  }
</script>

<script setup lang="ts">
  import Info from 'components/info/index.vue'
  import { $requests } from 'queries/index'
  import { ESizes } from 'types/enums'
  import { $formatDate } from 'utils/date'
  import { $formatInt, $formatPriceByLocale } from 'utils/formats'
  import { onMounted, ref } from 'vue'

  const loading = ref(false)
  const page = ref(0)
  const limit = 20
  const isLastPage = ref(false)
  const list = ref<Buy[]>([])

  async function getData() {
    if (loading.value) return

    loading.value = true

    try {
      const response = await $requests.transactions.buys({
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
            <th>{{ $t('Plan') }}</th>
            <th>CV</th>
            <th>{{ $t('Amount') }}</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in list" :key="item.uid">
            <td>{{ $formatDate(item.paid_at, 'date|time') }}</td>
            <td>{{ item.uid }}</td>
            <td>{{ item.user_uid }}</td>
            <td>
              <div class="badge badge-outline">{{ item.type }}</div>
            </td>
            <td>{{ item.plan_code }}</td>
            <td>{{ $formatInt(item.cv) }}</td>
            <td>
              <b>
                {{
                  $formatPriceByLocale({
                    count: $formatInt(item.amount as number, {
                      currency_code: item.currency_code,
                    }),
                    currency: item.currency_code,
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

    <div class="grid mobile:gap-2 pc:gap-4 mobile:grid-cols-1 pc:grid-cols-3">
      <Info :size="ESizes.SM" :message="$t('info.finance')" />
    </div>
  </div>
</template>
