<script lang="ts">
  export default {
    name: 'KneeSwitcher',
  }
</script>

<script setup lang="ts">
  // import Dropdown from 'components/dropdown/index.vue'
  import { $t } from 'i18n/index'
  import { $requests } from 'queries/index'
  import { ECurrencies, ERoles, ETeamType } from 'types/enums'
  import { $config } from 'utils/configuration'
  import { $formatInt } from 'utils/formats'
  import { $me } from 'utils/me'
  import { $notify } from 'utils/notify'
  import { computed, ref } from 'vue'

  const loading = ref(false)
  const activeType = computed(() => $me.data?.config?.team_type)
  const activeClass = 'border border-primary hover:border-primary'

  async function setKnee(type: TeamType) {
    if (!type || loading.value) return null

    loading.value = true

    try {
      const response = await $requests.teams.setSwitchKnee({
        type,
      })

      await $requests.me.updateConfig()

      if (response.team_type === type) {
        $notify.show({
          type: 'success',
          text: $t('switchKnee.success'),
        })
      } else {
        $notify.show({
          type: 'error',
          text: $t('switchKnee.error'),
        })
      }
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }
</script>

<template>
  <div v-if="$config.role === ERoles.distributor" class="flex items-center gap-4">
    <div class="flex items-center gap-2">
      <button
        class="btn btn-sm btn-outline inline-flex content-center btn-square tooltip tooltip-bottom"
        :class="`${activeType === 'left' ? activeClass : ''}`"
        :data-tip="$t('switchKnee.left.tooltip')"
        @click="setKnee(ETeamType.left)"
      >
        <span v-if="loading" class="loading loading-spinner"></span>
        <template v-else>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="w-4 h-4"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M10.5 19.5 3 12m0 0 7.5-7.5M3 12h18"
            />
          </svg>
        </template>
      </button>

      <span class="text-sm">{{ ECurrencies.cv }}: {{ $formatInt($me.data.left) }}</span>
    </div>

    <button
      class="btn btn-sm btn-square btn-outline tooltip tooltip-bottom"
      :class="`${activeType === 'auto' ? activeClass : ''}`"
      :data-tip="$t('switchKnee.auto.tooltip')"
      @click="setKnee(ETeamType.auto)"
    >
      <span v-if="loading" class="loading loading-spinner"></span>
      <template v-else>A</template>
    </button>

    <div class="flex items-center gap-2">
      <span class="text-sm">{{ ECurrencies.cv }}: {{ $formatInt($me.data.right) }}</span>
      <button
        class="btn btn-sm inline-flex btn-square btn-outline tooltip tooltip-bottom"
        :data-tip="$t('switchKnee.right.tooltip')"
        :class="`${activeType === 'right' ? activeClass : ''}`"
        @click="setKnee(ETeamType.right)"
      >
        <span v-if="loading" class="loading loading-spinner"></span>
        <template v-else>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="w-4 h-4"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3"
            />
          </svg>
        </template>
      </button>
    </div>

    <!-- <Dropdown class="flex">
      <template #action>
        <button class="btn btn-xs inline-flex btn-square btn-outline">
          <svg
            class="w-6 h-6"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              stroke="currentColor"
              stroke-linecap="round"
              stroke-width="2"
              d="M6 12h0m6 0h0m6 0h0"
            />
          </svg>
        </button>
      </template>

      <template #drop>
        <table class="table">
          <thead>
            <tr>
              <th><b>Бинарный оборот</b></th>
              <th><b>Left</b></th>
              <th><b>Right</b></th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Весь период</td>
              <td><b>123</b></td>
              <td><b>123</b></td>
            </tr>
            <tr>
              <td>Квалификационный период</td>
              <td><b>123</b></td>
              <td><b>123</b></td>
            </tr>
            <tr>
              <td>Текущий перелив</td>
              <td><b>123</b></td>
              <td><b>123</b></td>
            </tr>
          </tbody>
        </table>
      </template>
    </Dropdown> -->
  </div>
</template>
