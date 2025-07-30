<script setup lang="ts">
  import { $getLocaleOptions, $i18nGlobal, $t } from 'i18n/index'
  import Info from 'components/info/index.vue'
  import { EAlerts, ELocales, ESizes } from 'types/enums'
  import { computed, onMounted, Ref, ref } from 'vue'
  import { $formatInt, $formatPriceByLocale } from 'utils/formats'
  import { $requests } from 'queries/index'

  const loading = ref(true)
  const list: Ref<Task[]> = ref([])

  const codesForActivate = (() => {
    if ($i18nGlobal.locale.value === ELocales.ru_RU) {
      return ['sub-1001609461642']
    }

    if ($i18nGlobal.locale.value === ELocales.en_US) {
      return ['sub-1001967803227']
    }

    return []
  })()
  const isShowModal = computed(() => list.value.filter((task) => task.is_approve).length)

  async function setApproveTask(task: Task) {
    if (!task) return
    if (loading.value) return
    if (!task.is_approve) {
      loading.value = true

      try {
        const response = await $requests.tasks.approve(task.code)
        if (response?.code && task.link) {
          getTasks()
          setTimeout(() => {
            window.open(task.link, '_blank')
          }, 300)
        }
      } catch (error) {
        console.error(error)
      }

      loading.value = false
    } else {
      setTimeout(() => {
        window.open(task.link, '_blank')
      }, 300)
    }
  }

  function getAwTask(task: Task) {
    if (!task) return 0
    return $formatPriceByLocale({
      count: $formatInt(task.amount, {
        precision: task.precision,
      }),
      currency: task.currency_code,
    })
  }

  async function getTasks() {
    loading.value = true
    const responce = (await $requests.tasks.get()) || []
    list.value = responce.filter((task) => codesForActivate.includes(task.code))
    loading.value = false
  }

  function reloadTasks() {
    getTasks()
  }

  onMounted(() => {
    getTasks()
  })
</script>
<template>
  <div
    v-if="isShowModal"
    class="z-[2] absolute top-0 left-0 w-full h-full flex items-center justify-center"
  >
    <Info
      :type="EAlerts.BASE"
      :size="ESizes.BASE"
      :show-icon="false"
      :title="$t('info.farming_activation_task.title')"
      :message="$t('info.farming_activation_task.text')"
      class="text-center mt-[15%] py-10"
    >
      <template #after>
        <div class="pb-3">
          <div
            v-for="task in list"
            :key="task.code"
            class="flex items-center mobile:flex-col gap-4 pt-3 mt-4"
          >
            <div class="">
              <div>{{ getAwTask(task) }}</div>
              <div class="font-bold">{{ $t(task.texts[$getLocaleOptions().map]) }}</div>
            </div>

            <a
              href="javascript:void(0)"
              class="btn"
              :class="`${task.completed ? 'btn-outline' : 'btn-success'}`"
              @click.prevent="setApproveTask(task)"
            >
              <template v-if="loading">{{ $t('Loading') }}</template>

              <template v-else>
                <svg
                  v-if="!task.completed"
                  data-v-eb552580=""
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 496 512"
                  class="w-6"
                >
                  <path
                    data-v-eb552580=""
                    fill="currentColor"
                    d="M248,8C111.033,8,0,119.033,0,256S111.033,504,248,504,496,392.967,496,256,384.967,8,248,8ZM362.952,176.66c-3.732,39.215-19.881,134.378-28.1,178.3-3.476,18.584-10.322,24.816-16.948,25.425-14.4,1.326-25.338-9.517-39.287-18.661-21.827-14.308-34.158-23.215-55.346-37.177-24.485-16.135-8.612-25,5.342-39.5,3.652-3.793,67.107-61.51,68.335-66.746.153-.655.3-3.1-1.154-4.384s-3.59-.849-5.135-.5q-3.283.746-104.608,69.142-14.845,10.194-26.894,9.934c-8.855-.191-25.888-5.006-38.551-9.123-15.531-5.048-27.875-7.717-26.8-16.291q.84-6.7,18.45-13.7,108.446-47.248,144.628-62.3c68.872-28.647,83.183-33.623,92.511-33.789,2.052-.034,6.639.474,9.61,2.885a10.452,10.452,0,0,1,3.53,6.716A43.765,43.765,0,0,1,362.952,176.66Z"
                  ></path>
                </svg>

                {{ task.completed ? $t('Done') : $t('Make') }}

                <svg
                  data-v-eb552580=""
                  width="12"
                  height="12"
                  viewBox="0 0 48 48"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    data-v-eb552580=""
                    d="M19 11H37V29"
                    stroke="currentColor"
                    stroke-width="4"
                    stroke-linecap="butt"
                    stroke-linejoin="bevel"
                  ></path>
                  <path
                    data-v-eb552580=""
                    d="M11.5439 36.4559L36.9997 11"
                    stroke="currentColor"
                    stroke-width="4"
                    stroke-linecap="butt"
                    stroke-linejoin="bevel"
                  ></path>
                </svg>
              </template>
            </a>
          </div>

          <div class="">
            <div class="border-t my-8 opacity-50"></div>
            <div class="btn btn-primary cursor-pointer" @click="reloadTasks">
              {{ $t('UpdateIt') }} {{ $t('tasks.title').toLowerCase() }}
            </div>
          </div>
        </div>
      </template>
    </Info>
  </div>
</template>
