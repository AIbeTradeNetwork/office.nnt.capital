<script lang="ts">
  export default {
    name: 'ModalInfo',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import IvitedTotal from 'components/invitedTotal/index.vue'
  // import TaskMarafon from './taskMarafon.vue'

  import { $getLocaleOptions, $i18nGlobal, $t } from 'i18n/index'
  import { ECoinFarmingActions } from 'types/enums'
  import { $modals } from 'utils/modals'
  import { openShareLink } from 'utils/shares'
  import { $textSplit } from 'utils/text'
  import { Ref, computed, onMounted, ref } from 'vue'
  import { $me } from 'utils/me'
  import { $requests } from 'queries/index'
  import { $formatInt, $formatPriceByLocale } from 'utils/formats'

  const showCompletedTasks = ref(false)
  const taskLoading = ref(false)

  const data = computed(() => {
    if ($modals.farmingInfo.type === ECoinFarmingActions.farming) {
      return {
        title: $t('coin.farming.info.farming.title'),
        text: $t('coin.farming.info.farming.text'),
      }
    }

    if ($modals.farmingInfo.type === ECoinFarmingActions.invite) {
      return {
        title: $t('coin.farming.info.invite.title'),
        text: $t('coin.farming.info.invite.text'),
      }
    }

    if ($modals.farmingInfo.type === ECoinFarmingActions.team) {
      return {
        title: $t('coin.farming.info.team.title'),
        text: $t('coin.farming.info.team.text'),
      }
    }

    if ($modals.farmingInfo.type === ECoinFarmingActions.subscription) {
      return {
        title: $t('coin.farming.info.subscription.title'),
        text: $t('coin.farming.info.subscription.text'),
      }
    }

    if ($modals.farmingInfo.type === ECoinFarmingActions.tasks) {
      return {
        title: $t('tasks.title'),
        text: $t('tasks.text'),
      }
    }

    return {
      title: '',
      text: '',
    }
  })

  const tasks: {
    loading: Ref<boolean>
    list: Task[]
    translates(task: Task): string
    getTasks(): Promise<void>
    reload(): void
  } = {
    loading: ref(true),
    translates(task: Task) {
      const options = {
        amount: $formatPriceByLocale({
          count: $formatInt(task.amount, {
            precision: task.precision,
          }),
          currency: task.currency_code,
        }),
      }
      return $t(`tasks.task.${task.code.replace('-', '_')}`, options)
    },
    list: [],
    async getTasks() {
      tasks.loading.value = true
      tasks.list = ((await $requests.tasks.get()) || []).sort((a, b) => (a.completed ? -1 : 1))
      tasks.loading.value = false
    },
    reload() {
      tasks.getTasks()
    },
  }

  function close() {
    $modals.farmingInfo.close()
  }

  function loadTasks() {
    if ($modals.farmingInfo.type === ECoinFarmingActions.tasks) {
      tasks.getTasks()
    }
  }

  async function setApproveTask(task: Task) {
    if (!task) return
    if (taskLoading.value) return
    if (!task.is_approve) {
      taskLoading.value = true

      try {
        const response = await $requests.tasks.approve(task.code)
        if (response?.code && task.link) {
          loadTasks()
          setTimeout(() => {
            window.open(task.link, '_blank')
          }, 300)
        }
      } catch (error) {
        console.error(error)
      }

      taskLoading.value = false
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

  onMounted(() => {
    loadTasks()
  })
</script>

<template>
  <Modal :modal="$modals.farmingInfo" :z="50" :title="data.title">
    <div v-if="$modals.farmingInfo.type === ECoinFarmingActions.team" class="my-6 space-y-4">
      <b>{{ $t('PersonallyInvited') }}</b>
      : {{ $me.data.team_count }}
    </div>

    <div class="space-y-4">
      <p v-for="par in $textSplit(data.text)" :key="par">{{ par }}</p>
    </div>

    <div v-if="$modals.farmingInfo.type === ECoinFarmingActions.invite" class="mt-6 space-y-4">
      <!-- <template v-else>
        <button
          class="btn btn-outline w-full drop-shadow-custom"
          @click="$copyToClipboard(createRefLink())"
        >
          {{ $t('ReferralLink') }}
        </button>
      </template> -->

      <IvitedTotal class="my-3" />

      <div class="mb-2">{{ $t('ReferralLink') }}:</div>

      <button class="btn btn-outline w-full drop-shadow-custom" @click="openShareLink()">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 496 512" class="w-6">
          <path
            fill="currentColor"
            d="M248,8C111.033,8,0,119.033,0,256S111.033,504,248,504,496,392.967,496,256,384.967,8,248,8ZM362.952,176.66c-3.732,39.215-19.881,134.378-28.1,178.3-3.476,18.584-10.322,24.816-16.948,25.425-14.4,1.326-25.338-9.517-39.287-18.661-21.827-14.308-34.158-23.215-55.346-37.177-24.485-16.135-8.612-25,5.342-39.5,3.652-3.793,67.107-61.51,68.335-66.746.153-.655.3-3.1-1.154-4.384s-3.59-.849-5.135-.5q-3.283.746-104.608,69.142-14.845,10.194-26.894,9.934c-8.855-.191-25.888-5.006-38.551-9.123-15.531-5.048-27.875-7.717-26.8-16.291q.84-6.7,18.45-13.7,108.446-47.248,144.628-62.3c68.872-28.647,83.183-33.623,92.511-33.789,2.052-.034,6.639.474,9.61,2.885a10.452,10.452,0,0,1,3.53,6.716A43.765,43.765,0,0,1,362.952,176.66Z"
          />
        </svg>
        {{ $t('ReferralLink') }}
        <svg
          width="12"
          height="12"
          viewBox="0 0 48 48"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M19 11H37V29"
            stroke="currentColor"
            stroke-width="4"
            stroke-linecap="butt"
            stroke-linejoin="bevel"
          ></path>
          <path
            d="M11.5439 36.4559L36.9997 11"
            stroke="currentColor"
            stroke-width="4"
            stroke-linecap="butt"
            stroke-linejoin="bevel"
          ></path>
        </svg>
      </button>
      <!-- <button
        class="btn btn-outline w-full drop-shadow-custom"
        @click="$copyToClipboard(createRefLink('web'))"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-6"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582m15.686 0A11.953 11.953 0 0 1 12 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0 1 21 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0 1 12 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 0 1 3 12c0-1.605.42-3.113 1.157-4.418"
          />
        </svg>

        {{ $t('WebVersion') }}
      </button> -->
    </div>

    <div v-if="$modals.farmingInfo.type === ECoinFarmingActions.team" class="mt-6 space-y-4">
      <router-link :to="{ name: 'Friends' }" class="btn btn-primary w-full">
        <span>
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
              d="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582m15.686 0A11.953 11.953 0 0 1 12 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0 1 21 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0 1 12 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 0 1 3 12c0-1.605.42-3.113 1.157-4.418"
            />
          </svg>
        </span>
        <span>{{ $t('MyFriends') }}</span>
      </router-link>
    </div>

    <div v-if="$modals.farmingInfo.type === ECoinFarmingActions.subscription" class="mt-6">
      <!-- <router-link
        class="btn btn-primary w-full"
        :to="{ name: 'Tariffs' }"
        @click="$modals.farmingInfo.close()"
      >
        <span>
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
              d="M6.429 9.75 2.25 12l4.179 2.25m0-4.5 5.571 3 5.571-3m-11.142 0L2.25 7.5 12 2.25l9.75 5.25-4.179 2.25m0 0L21.75 12l-4.179 2.25m0 0 4.179 2.25L12 21.75 2.25 16.5l4.179-2.25m11.142 0-5.571 3-5.571-3"
            />
          </svg>
        </span>
        <span>{{ $t('Plans') }}</span>
      </router-link> -->
    </div>

    <div v-if="$modals.farmingInfo.type === ECoinFarmingActions.tasks" class="">
      <p>{{ $t('info.tasks.warning_1') }}</p>

      <div v-if="tasks.loading.value" class="text-center font-bold my-10 pb-6 text-info">
        {{ $t('Loading') }}
      </div>

      <div v-else>
        <!-- <TaskMarafon class="border-t pt-5 mt-6" /> -->
        <div class="form-control pt-5">
          <label class="label cursor-pointer flex justify-end gap-2">
            <span class="label-text">{{ $t('ShowCompletedTasks') }}</span>
            <input
              type="checkbox"
              :checked="showCompletedTasks"
              class="checkbox"
              @change="showCompletedTasks = !showCompletedTasks"
            />
          </label>
        </div>
        <template v-for="task in tasks.list">
          <div
            v-if="showCompletedTasks ? true : !task.completed"
            :key="task.code"
            class="flex items-center mobile:flex-col gap-4 border-t pt-5 mt-6"
          >
            <div class="self-start">
              {{ getAwTask(task) }} {{ $t(task.texts[$getLocaleOptions().map]) }}
            </div>

            <a
              href="javascript:void(0)"
              class="btn mobile:self-end"
              :class="`${task.completed ? 'btn-outline' : 'btn-success'}`"
              @click.prevent="setApproveTask(task)"
            >
              <template v-if="taskLoading">{{ $t('Loading') }}</template>
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
        </template>

        <div class="flex gap-2 justify-end mt-10">
          <!-- <button :disabled="tasks.loading.value" class="btn btn-info" @click="tasks.reload">
            {{ $t('UpdateIt') }}
          </button> -->
          <button class="btn btn-outline" @click="close">{{ $t('Close') }}</button>
        </div>
      </div>
    </div>
  </Modal>
</template>
