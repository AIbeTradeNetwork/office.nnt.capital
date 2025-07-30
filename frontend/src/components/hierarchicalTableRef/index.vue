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
  import { $historBack } from 'utils/history'

  const loading = ref<boolean>(false)
  const childrens = ref<TeamUser[]>([])
  const colspan = 5

  const isClient = computed(() => $config.role === ERoles.client)
  const topSpace = computed(() =>
    isClient.value ? 'top-0' : 'top-[calc(theme(spacing.16)+theme(spacing.2))]',
  )

  async function loadData() {
    loading.value = true

    try {
      childrens.value = await $requests.teams.getRef({ user_uid: $me.data.uid })
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
      <!-- min-h-[calc(100%-theme(spacing.4))] -->
      <table class="table table-pin-rows whitespace-nowrap [&_th]:bg-base-300">
        <!-- <thead>
          <tr class="uppercase">
            <th>Мой ментор</th>
            
          </tr>
        </thead>

        <tbody>
          <MentorLine />
        </tbody> -->

        <thead>
          <tr class="bg-base-300 uppercase">
            <th>{{ $t('Me') }}</th>
            <th></th>
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
            <th :colspan="colspan">
              <div class="flex items-center">
                <!-- {{ $t('MyTeam') }}: -->
                <button class="btn btn-success btn-xs mr-2 pc:hidden" @click="$historBack">
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
                {{ $t('PersonallyInvited') }}:
                <span class="badge badge-secondary ml-2 font-normal">{{ childrens.length }}</span>
              </div>
            </th>
          </tr>
        </thead>

        <tbody>
          <tr v-if="loading">
            <td :colspan="colspan">
              <div class="text-center font-bold p-4">{{ $t('Loading') }}</div>
            </td>
          </tr>

          <tr v-else-if="!childrens.length">
            <td :colspan="colspan">
              <div class="pc:text-center font-bold p-6">{{ $t('ListEmpty') }}</div>
            </td>
          </tr>

          <template v-else>
            <CommandLine
              v-for="user in childrens"
              :key="user.uid"
              :lineid="0"
              :user="user"
              :show-ref-count="true"
            />
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>
