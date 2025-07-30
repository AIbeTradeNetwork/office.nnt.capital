<script lang="ts">
  export default {
    name: 'AvatarNav',
  }
</script>

<script setup lang="ts">
  import { $modals } from 'utils/modals'
  import { $me } from 'utils/me'
  import { $config } from 'utils/configuration'

  import Dropdown from 'components/dropdown/index.vue'
  import QualificationsUser from 'components/userQualificationInfo/client.vue'
  import QualificationsDistributor from 'components/userQualificationInfo/distributor.vue'
  import { ERoles } from 'types/enums'
  import Avatar from 'vue-boring-avatars'
  import { $copyToClipboard } from 'utils/clipboard'

  function getQualification() {
    if ($config.role === ERoles.distributor) {
      return QualificationsDistributor
    }
    if ($config.role === ERoles.client) {
      return QualificationsUser
    }
  }

  function showSocial() {
    $modals.userSocial.show()
  }
</script>

<template>
  <div class="h-[48px]">
    <div class="flex">
      <Dropdown class="flex">
        <template #trigger>
          <div class="flex items-center cursor-pointer">
            <button tabindex="0" role="button">
              <div class="flex">
                <div class="rounded-full overflow-hidden avatar w-12">
                  <Avatar class="w-full h-full" variant="marble" :name="$me.data.uid" />
                </div>

                <div
                  class="ml-2 whitespace-nowrap text-left flex flex-col items-start justify-center"
                >
                  <div class="text-left truncate max-w-[150px]">
                    <b>
                      {{ $me.data.nickname || 'Username' }}
                    </b>
                  </div>
                  <div v-if="$me.data.uid" class="text-base-content text-center">
                    <!-- {{ $t('MyPromoCode') }}: -->
                    <span
                      class="whitespace-nowrap underline cursor-pointer text-gradient text-sm"
                      @click.self="$copyToClipboard($me.data.uid)"
                    >
                      {{ $me.data.uid }}
                      <!-- <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke-width="1.5"
                        stroke="currentColor"
                        class="w-4 inline relative -top-[3px]"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          d="M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 0 1-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 0 1 1.5.124m7.5 10.376h3.375c.621 0 1.125-.504 1.125-1.125V11.25c0-4.46-3.243-8.161-7.5-8.876a9.06 9.06 0 0 0-1.5-.124H9.375c-.621 0-1.125.504-1.125 1.125v3.5m7.5 10.375H9.375a1.125 1.125 0 0 1-1.125-1.125v-9.25m12 6.625v-1.875a3.375 3.375 0 0 0-3.375-3.375h-1.5a1.125 1.125 0 0 1-1.125-1.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H9.75"
                        />
                      </svg> -->
                    </span>
                  </div>
                  <!-- <div v-if="$me.data.is_premium">
                    <div class="text-gradient">{{ $t('Premium') }}</div>
                  </div> -->
                </div>
              </div>
            </button>
          </div>
        </template>
        <template #drop><component :is="getQualification()" /></template>
      </Dropdown>

      <!-- <button class="flex self-center ml-4" @click="showSocial">
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
            d="M10.5 19.5h3m-6.75 2.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-15a2.25 2.25 0 0 0-2.25-2.25H6.75A2.25 2.25 0 0 0 4.5 4.5v15a2.25 2.25 0 0 0 2.25 2.25Z"
          />
        </svg>
      </button> -->
    </div>
  </div>
</template>
