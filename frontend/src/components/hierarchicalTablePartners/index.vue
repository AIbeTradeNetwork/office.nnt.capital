<script lang="ts">
export default {
  name: 'HierarchicalTreeFriends',
}
</script>

<script setup lang="ts">
import MyLine from './myLine.vue'
import CommandLine from './commandLine.vue'
import { onMounted, ref } from 'vue'
import { $requests } from 'queries/index'
import { $me } from 'utils/me'
import { $historBack } from 'utils/history'

const loading = ref<boolean>(false)
const childrens = ref<FriendUser[]>([])
const colspan = 3
const limit = 15
let page = 0
const isEnd = ref(true)

async function getRefs() {
  if (!page || page < 0) page = 0

  loading.value = true

  try {
    childrens.value.push(...(await $requests.friends.list($me.data.uid, limit, limit * page)))

    if (childrens.value.length >= $me.data.team_count) {
      isEnd.value = true
    } else {
      isEnd.value = false
    }

    page++
  } catch (error) {
    console.error(error)
  }

  loading.value = false
}

onMounted(() => {
  getRefs()
})
</script>

<template>
<div class="">
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
          <th>{{ $t('Registration') }}</th>
          <th></th>
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
              <span class="badge badge-secondary ml-2 font-normal">
                {{ $me.data.team_count }}
              </span>
            </div>
          </th>
        </tr>
      </thead>

      <tbody>
        <tr v-if="!childrens.length">
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

          <tr v-if="!isEnd" id="0">
            <td :colspan="colspan">
              <div class="flex mobile:justify-start pc:justify-center">
                <div v-if="loading" class="font-bold p-4">{{ $t('Loading') }}</div>
                <button v-else class="btn btn-primary btn-sm" @click="getRefs">
                  {{ $t('ShowMore') }}
                </button>
              </div>
            </td>
          </tr>
        </template>
      </tbody>
    </table>
  </div>
</div>
</template>
