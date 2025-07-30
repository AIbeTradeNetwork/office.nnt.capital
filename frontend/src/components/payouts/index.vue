<script lang="ts">
  export default {
    name: 'PagePayments',
  }
</script>

<script setup lang="ts">
  import { $t } from 'i18n/index'
  import { $requests } from 'queries/index'
  import { $formatDate } from 'utils/date'
  import { $formatInt, $formatPriceByLocale } from 'utils/formats'
  import { $modals } from 'utils/modals'
  import { onMounted, ref } from 'vue'

  const loading = ref(false)
  const page = ref(0)
  const limit = 15
  const isLastPage = ref(false)
  const list = ref<Payout[]>([])

  async function getData() {
    if (loading.value) return

    loading.value = true

    try {
      const response = await $requests.payouts.get({
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

  function getStatus(item: Payout) {
    if (item.cancelled_at > 0)
      return {
        class: 'badge badge-outline badge-error',
        text: $t('CancelledAt'),
      }
    if (item.charged_at > 0)
      return {
        class: 'badge badge-outline badge-success',
        text: $t('ChargedAt'),
      }
    if (item.approved_at > 0)
      return {
        class: 'badge badge-outline badge-warning',
        text: $t('ApprovedAt'),
      }
    if (item.created_at > 0)
      return {
        class: 'badge badge-outline',
        text: $t('CreatedAt'),
      }
  }

  function showReason(item: Payout) {
    if (!item.reason) return
    $modals.info.show({
      title: $t('Reason'),
      text: item.reason,
    })
  }

  function addPayout() {
    $modals.payout.onSucess = () => {
      getData()
    }

    $modals.payout.show()
  }

  onMounted(() => {
    getData()
  })
</script>

<template>
  <div>
    <button class="btn btn-primary inline-flex items-center" @click="addPayout">
      <span>{{ $t('CreatePayout') }}</span>

      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-6 h-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
        />
      </svg>
    </button>
  </div>

  <div class="box p-0 overflow-x-auto">
    <table class="table table-zebra">
      <thead>
        <tr>
          <td>UID</td>
          <td>{{ $t('CreatedAt') }}</td>
          <td>{{ $t('MethodCode') }}</td>
          <td>{{ $t('AccountName') }}</td>
          <td>{{ $t('Amount') }}</td>
          <td>{{ $t('Commission') }}</td>
          <td>{{ $t('Total') }}</td>
          <td>{{ $t('Currency') }}</td>
          <td>{{ $t('Status') }}</td>
          <td>{{ $t('Reason') }}</td>
        </tr>
      </thead>

      <tbody v-if="!list.length">
        <tr>
          <td colspan="9">
            <div class="pc:text-center font-bold p-6">{{ $t('ListEmpty') }}</div>
          </td>
        </tr>
      </tbody>

      <tbody>
        <tr v-for="item in list" :key="item.uid">
          <td>{{ item.uid }}</td>
          <td>{{ $formatDate(item.created_at, 'date|time') }}</td>
          <td>{{ item.method_code }}</td>
          <td>{{ item.account_name }}</td>
          <td>{{ $formatInt(item.amount) }}</td>
          <td>{{ $formatInt(item.commission) }}</td>
          <td>{{ $formatInt(item.amount - item.commission).toFixed(2) }}</td>
          <td>{{ item.currency_code }}</td>
          <td>
            <div :class="getStatus(item).class">{{ getStatus(item).text }}</div>
          </td>
          <td>
            <div
              v-if="item.reason"
              class="flex gap-2 items-center cursor-pointer hover:text-primary"
              @click="showReason(item)"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="flex-[0_0_1rem]"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
                />
              </svg>

              <div class="text-ellipsis line-clamp-1 max-w-[100px]">
                {{ item.reason }}
              </div>
            </div>
            <div v-else>-</div>
          </td>
        </tr>
      </tbody>

      <!-- <tbody v-if="!isLastPage">
        <tr>
          <td colspan="9">
            <div class="flex justify-center">
              <button class="btn btn-sm" @click="loadMore">
                <span v-if="loading" class="loading loading-spinner"></span>
                <span v-else>{{ $t('ShowMore') }}</span>
              </button>
            </div>
          </td>
        </tr>
      </tbody> -->
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
</template>
