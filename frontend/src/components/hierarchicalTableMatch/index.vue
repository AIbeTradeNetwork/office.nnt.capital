<script lang="ts">
  export default {
    name: 'HierarchicalTree',
  }
</script>

<script setup lang="ts">
  import MyLine from './myLine.vue'
  import CommandLine from './commandLine.vue'
  import { computed, onMounted, ref } from 'vue'
  import { $requests } from 'queries/index'
  import { $me } from 'utils/me'
  import { $config } from 'utils/configuration'
  import { ERoles } from 'types/enums'

  const loading = ref<boolean>(false)
  const childrens = ref<TeamUser[]>([])

  const isClient = computed(() => $config.role === ERoles.client)
  const topSpace = computed(() =>
    isClient.value ? 'top-0' : 'top-[calc(theme(spacing.16)+theme(spacing.2))]',
  )

  async function loadData() {
    loading.value = true

    try {
      childrens.value = await $requests.teams.getMatch({ user_uid: $me.data.uid })
      loading.value = false
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  onMounted(async () => {
    await loadData()
  })
</script>

<template>
  <div class="" :class="`${topSpace}`">
    <div class="box w-full h-full p-0 overflow-auto rounded-none">
      <table class="table table-pin-rows whitespace-nowrap [&_th]:bg-base-300">
        <thead>
          <tr class="bg-base-300 uppercase">
            <th>{{ $t('Me') }}</th>
            <th>{{ $t('Plan') }}</th>
            <th>{{ $t('Rank') }}</th>
            <th>{{ $t('Activity') }}</th>
          </tr>
        </thead>

        <tbody>
          <MyLine />
        </tbody>

        <thead>
          <tr class="bg-base-300 uppercase">
            <th colspan="4">{{ $t('MyTeam') }}</th>
          </tr>
        </thead>

        <tbody>
          <tr v-if="loading" colspan="">
            <td colspan="4">
              <div class="text-center font-bold p-4">{{ $t('Loading') }}</div>
            </td>
          </tr>

          <tr v-else-if="!childrens.length">
            <td colspan="4">
              <div class="pc:text-center font-bold p-6">{{ $t('ListEmpty') }}</div>
            </td>
          </tr>

          <template v-else>
            <CommandLine v-for="user in childrens" :key="user.uid" :lineid="0" :user="user" />
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>
