<script lang="ts">
  export default {
    name: 'MenuComp',
  }
</script>

<script setup lang="ts">
  import { watch } from 'vue'
  // import Logo from 'assets/logos/logo.vue'
  import DocumentsList from 'components/documentsList/index.vue'
  import { Router } from 'routes/index'
  import { $screen } from 'utils/screen'
  import { $config } from 'utils/configuration'
  import { $copyToClipboard } from 'utils/clipboard'
  // import downloadDebuggingFile from './downloadDebuggingFile'
  import KneeSwitcher from 'components/kneeSwitcher/index.vue'
  import LocaleSelector from 'components/localeSelector/index.vue'
  import SettingsSelector from 'components/settingsSelector/index.vue'
  // import BalanceNav from 'components/balance/nav.vue'
  import { ERoles } from 'types/enums'
  import { createRefLink, openShareLink } from 'utils/shares'
  import { $useTelegram } from 'utils/telegram'
  import { checkDevMode } from 'utils/queriesGuards'

  watch(Router.currentRoute, (routeNew, routeOld) => {
    if (routeNew.path !== routeOld.path) closeMenu()
  })

  function clearCacheAndReload() {
    if (window.caches) {
      window.caches.keys().then(async (keyList) => {
        await Promise.all(
          keyList.map((key) => {
            caches.delete(key)
          }),
        )

        location.reload()
      })
    } else {
      location.reload()
    }
  }

  // const isActiveRoute = computed(() => {
  //   return (arr: string[] = []) => {
  //     return arr.includes(Router.currentRoute.value.name as string)
  //   }
  // })

  // const isOpenByRoute = computed(() => {
  //   return (reg: string): boolean => {
  //     return !!(Router.currentRoute.value.path as string)
  //       .toLowerCase()
  //       .match(reg)
  //   }
  // })

  let devClicks = 0
  let devTimer: ReturnType<typeof setTimeout> = null
  function devMode() {
    clearTimeout(devTimer)
    devTimer = setTimeout(() => (devClicks = 0), 300)

    devClicks += 1

    if (devClicks === 20) {
      if ($config.isDevMode) {
        checkDevMode('0')
      } else {
        checkDevMode('1')
      }
      devClicks = 0
    }
  }

  function closeMenu() {
    ;(document.getElementById('drawer') as HTMLInputElement).checked = false
  }
</script>

<template>
  <div
    class="drawer-side z-[2] overflow-y-auto scroll-smooth overscroll-none mobile:no-scrollbar bottom-0"
  >
    <label for="drawer" class="drawer-overlay" aria-label="Close menu"></label>

    <aside
      class="safearea-l safearea-b bg-base-100 pc:bg-opacity-50 w-full max-w-[300px] pc:min-h-full"
    >
      <!-- scrollbar-stable -->
      <div
        class="safearea-t sticky top-0 z-[1] border-b border-white/30 backdrop-blur hover:bg-base-100 hover:bg-opacity-10"
      >
        <a class="flex justify-center h-[69px]" href="/" aria-current="page" aria-label="Logo">
          <img src="/src/assets/logos/logo.png" alt="Logo" class="h-full object-contain" />
        </a>
      </div>

      <div class="h-[calc(100%-66px)] flex flex-col justify-between py-4">
        <!-- <div v-if="$screen.isMobile" class="pc:hidden">
          <div class="mobile:px-6">
            <UserStatusSimply />
          </div>
          <div class="divider my-2"></div>
        </div> -->

        <!-- <div v-if="$screen.isMobile.value" class="px-4 mb-6">
          <BalanceNav />
        </div> -->

        <div v-if="$screen.isMobile.value && $config.role === ERoles.distributor" class="px-4 mb-6">
          <KneeSwitcher />
        </div>

        <div>
          <div class="px-4 mb-4">
            <button
              class="btn btn-outline w-full drop-shadow-custom border border-info"
              @click="openShareLink()"
            >
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
          </div>

          <!-- <div class="px-4 mb-4">
            <button
              class="btn btn-outline w-full drop-shadow-custom"
              @click="$copyToClipboard(createRefLink())"
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
            </button>
          </div> -->
        </div>

        <ul class="menu pt-0">
          <!-- START TEST -->
          <li v-if="$config.isDevMode"></li>
          <li v-if="$config.isDevMode" class="border border-red-600 rounded-lg">
            <details>
              <summary class="group">
                <span>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="size-6"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M4.26 10.147a60.438 60.438 0 0 0-.491 6.347A48.62 48.62 0 0 1 12 20.904a48.62 48.62 0 0 1 8.232-4.41 60.46 60.46 0 0 0-.491-6.347m-15.482 0a50.636 50.636 0 0 0-2.658-.813A59.906 59.906 0 0 1 12 3.493a59.903 59.903 0 0 1 10.399 5.84c-.896.248-1.783.52-2.658.814m-15.482 0A50.717 50.717 0 0 1 12 13.489a50.702 50.702 0 0 1 7.74-3.342M6.75 15a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5Zm0 0v-3.675A55.378 55.378 0 0 1 12 8.443m-7.007 11.55A5.981 5.981 0 0 0 6.75 15.75v-1.5"
                    />
                  </svg>
                </span>
                {{ $t('TradingBots') }}
              </summary>
              <ul>
                <li>
                  <router-link :to="{ name: 'Exchanges' }">
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
                          d="M12 21v-8.25M15.75 21v-8.25M8.25 21v-8.25M3 9l9-6 9 6m-1.5 12V10.332A48.36 48.36 0 0 0 12 9.75c-2.551 0-5.056.2-7.5.582V21M3 21h18M12 6.75h.008v.008H12V6.75Z"
                        />
                      </svg>
                    </span>
                    <span>{{ $t('Exchanges') }}</span>
                  </router-link>
                </li>
                <li>
                  <router-link class="border border-warning" :to="{ name: 'StrategiesIco' }">
                    <span>
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke-width="1.5"
                        stroke="currentColor"
                        class="w-6 text-orange-400"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          d="M15.362 5.214A8.252 8.252 0 0 1 12 21 8.25 8.25 0 0 1 6.038 7.047 8.287 8.287 0 0 0 9 9.601a8.983 8.983 0 0 1 3.361-6.867 8.21 8.21 0 0 0 3 2.48Z"
                        />
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          d="M12 18a3.75 3.75 0 0 0 .495-7.468 5.99 5.99 0 0 0-1.925 3.547 5.975 5.975 0 0 1-2.133-1.001A3.75 3.75 0 0 0 12 18Z"
                        />
                      </svg>
                    </span>
                    {{ $t('Strategies') }} ICO
                  </router-link>
                </li>
                <li>
                  <router-link :to="{ name: 'Strategies' }">
                    <span>
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
                          d="M3.75 3v11.25A2.25 2.25 0 0 0 6 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0 1 18 16.5h-2.25m-7.5 0h7.5m-7.5 0-1 3m8.5-3 1 3m0 0 .5 1.5m-.5-1.5h-9.5m0 0-.5 1.5m.75-9 3-3 2.148 2.148A12.061 12.061 0 0 1 16.5 7.605"
                        />
                      </svg>
                    </span>
                    <span>{{ $t('Strategies') }}</span>
                  </router-link>
                </li>
                <li>
                  <router-link :to="{ name: 'TradingBots' }">
                    <span>
                      <svg
                        stroke="currentColor"
                        class="w-6 h-6"
                        viewBox="0 0 24 24"
                        fill="none"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          fill-rule="evenodd"
                          clip-rule="evenodd"
                          d="M7.5 11C6.94772 11 6.5 11.4477 6.5 12C6.5 12.5523 6.94772 13 7.5 13C8.05228 13 8.5 12.5523 8.5 12C8.5 11.4477 8.05228 11 7.5 11ZM5.5 12C5.5 10.8954 6.39543 10 7.5 10C8.60457 10 9.5 10.8954 9.5 12C9.5 13.1046 8.60457 14 7.5 14C6.39543 14 5.5 13.1046 5.5 12Z"
                          fill="#47495F"
                        />
                        <path
                          fill-rule="evenodd"
                          clip-rule="evenodd"
                          d="M16.5 11C15.9477 11 15.5 11.4477 15.5 12C15.5 12.5523 15.9477 13 16.5 13C17.0523 13 17.5 12.5523 17.5 12C17.5 11.4477 17.0523 11 16.5 11ZM14.5 12C14.5 10.8954 15.3954 10 16.5 10C17.6046 10 18.5 10.8954 18.5 12C18.5 13.1046 17.6046 14 16.5 14C15.3954 14 14.5 13.1046 14.5 12Z"
                          fill="#47495F"
                        />
                        <path
                          fill-rule="evenodd"
                          clip-rule="evenodd"
                          d="M10 15.5C10.2761 15.5 10.5 15.7239 10.5 16L10.5003 16.0027C10.5003 16.0027 10.5014 16.0073 10.5034 16.0122C10.5074 16.022 10.5171 16.0405 10.5389 16.0663C10.5845 16.1202 10.6701 16.1902 10.8094 16.2599C11.0883 16.3993 11.5085 16.5 12 16.5C12.4915 16.5 12.9117 16.3993 13.1906 16.2599C13.3299 16.1902 13.4155 16.1202 13.4611 16.0663C13.4829 16.0405 13.4926 16.022 13.4966 16.0122C13.4986 16.0073 13.4997 16.0027 13.4997 16.0027L13.5 16C13.5 15.7239 13.7239 15.5 14 15.5C14.2761 15.5 14.5 15.7239 14.5 16C14.5 16.5676 14.0529 16.9468 13.6378 17.1543C13.1928 17.3768 12.6131 17.5 12 17.5C11.3869 17.5 10.8072 17.3768 10.3622 17.1543C9.9471 16.9468 9.5 16.5676 9.5 16C9.5 15.7239 9.72386 15.5 10 15.5Z"
                          fill="#47495F"
                        />
                        <path
                          fill-rule="evenodd"
                          clip-rule="evenodd"
                          d="M16 5.5V7.5H8V5.5C8 5.22386 7.77614 5 7.5 5C7.22386 5 7 5.22386 7 5.5V7.5H6C4.61929 7.5 3.5 8.61929 3.5 10V17C3.5 18.3807 4.61929 19.5 6 19.5H18C19.3807 19.5 20.5 18.3807 20.5 17V10C20.5 8.61929 19.3807 7.5 18 7.5H17V5.5C17 5.22386 16.7761 5 16.5 5C16.2239 5 16 5.22386 16 5.5ZM6 8.5C5.17157 8.5 4.5 9.17157 4.5 10V17C4.5 17.8284 5.17157 18.5 6 18.5H18C18.8284 18.5 19.5 17.8284 19.5 17V10C19.5 9.17157 18.8284 8.5 18 8.5H6Z"
                          fill="#47495F"
                        />
                      </svg>
                    </span>
                    <span>{{ $t('TradingBots') }}</span>
                  </router-link>
                </li>

                <li v-if="$config.role === ERoles.client">
                  <router-link class="border border-success" :to="{ name: 'BecomeDistributor' }">
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
                        d="M4.26 10.147a60.438 60.438 0 0 0-.491 6.347A48.62 48.62 0 0 1 12 20.904a48.62 48.62 0 0 1 8.232-4.41 60.46 60.46 0 0 0-.491-6.347m-15.482 0a50.636 50.636 0 0 0-2.658-.813A59.906 59.906 0 0 1 12 3.493a59.903 59.903 0 0 1 10.399 5.84c-.896.248-1.783.52-2.658.814m-15.482 0A50.717 50.717 0 0 1 12 13.489a50.702 50.702 0 0 1 7.74-3.342M6.75 15a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5Zm0 0v-3.675A55.378 55.378 0 0 1 12 8.443m-7.007 11.55A5.981 5.981 0 0 0 6.75 15.75v-1.5"
                      />
                    </svg>
                    {{ $t('becomeDistributor.name') }}
                  </router-link>
                </li>
              </ul>
            </details>
          </li>
          <!-- END TEST -->

          <li>
            <router-link :to="{ name: 'FarmUDEX' }">
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
                    d="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375m16.5 0v3.75m-16.5-3.75v3.75m16.5 0v3.75C20.25 16.153 16.556 18 12 18s-8.25-1.847-8.25-4.125v-3.75m16.5 0c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125"
                  />
                </svg>
              </span>
              <span>{{ $t('Index') }}</span>
            </router-link>
          </li>
          <div class="my-2 border-t opacity-30"></div>

          <li>
            <router-link :to="{ name: 'Shop' }">
              <span>
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-shopping-cart"><circle cx="9" cy="21" r="1"></circle><circle cx="20" cy="21" r="1"></circle><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path></svg>
              </span>
              <span>{{ $t('Boosts') }}</span>
            </router-link>
          </li>
          <div class="my-2 border-t opacity-30"></div>

          <li>
            <router-link :to="{ name: 'Friends' }">
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
          </li>
          <div class="my-2 border-t opacity-30"></div>
          <li>
  <router-link :to="{ name: 'Partners' }">
    <span>
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 0 0-5.356-1.857M17 20H7m10 0v-2a3 3 0 0 0-5.356-1.857M7 20H2v-2a3 3 0 0 1 5.356-1.857M7 20v-2a3 3 0 0 1 5.356-1.857M15 7a3 3 0 1 1-6 0 3 3 0 0 1 6 0Zm6 3a2 2 0 1 1-4 0 2 2 0 0 1 4 0ZM7 10a2 2 0 1 1-4 0 2 2 0 0 1 4 0Z" />
      </svg>
    </span>
    <span>{{ $t('Partners') }}</span>
  </router-link>
</li>
          <div class="my-2 border-t opacity-30"></div>
          <li>
            <router-link :to="{ name: 'Payments' }">
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
                    d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 0 0 2.25-2.25V6.75A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25v10.5A2.25 2.25 0 0 0 4.5 19.5Z"
                  />
                </svg>
              </span>
              <span>{{ $t('Payments') }}</span>
            </router-link>
          </li>
          <div class="my-2 border-t opacity-30"></div>
          <li>
            <router-link :to="{ name: 'Purchases' }">
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
                    d="M12 6v12m-3-2.818.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
                  />
                </svg>
              </span>
              <span>{{ $t('Purchases') }}</span>
            </router-link>
          </li>
          <div class="my-2 border-t opacity-30"></div>

          <li>
            <router-link class="border border-info" :to="{ name: 'AutotradePRO' }">
              <span>
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
                    d="M9.813 15.904 9 18.75l-.813-2.846a4.5 4.5 0 0 0-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 0 0 3.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 0 0 3.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 0 0-3.09 3.09ZM18.259 8.715 18 9.75l-.259-1.035a3.375 3.375 0 0 0-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 0 0 2.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 0 0 2.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 0 0-2.456 2.456ZM16.894 20.567 16.5 21.75l-.394-1.183a2.25 2.25 0 0 0-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 0 0 1.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 0 0 1.423 1.423l1.183.394-1.183.394a2.25 2.25 0 0 0-1.423 1.423Z"
                  />
                </svg>
              </span>
              <span>{{ $t('autotrade_pro.name') }}</span>
            </router-link>
          </li>

          <!-- <li>
            <router-link :to="{ name: 'Analytics' }">
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
                    d="M3.75 3v11.25A2.25 2.25 0 0 0 6 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0 1 18 16.5h-2.25m-7.5 0h7.5m-7.5 0-1 3m8.5-3 1 3m0 0 .5 1.5m-.5-1.5h-9.5m0 0-.5 1.5m.75-9 3-3 2.148 2.148A12.061 12.061 0 0 1 16.5 7.605"
                  />
                </svg>
              </span>
              <span>{{ $t('Analytics') }}</span>
            </router-link>
          </li>

          <li></li> -->

          <!-- <li>
            <router-link :to="{ name: 'Oracle' }">
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
                    d="M8.25 3v1.5M4.5 8.25H3m18 0h-1.5M4.5 12H3m18 0h-1.5m-15 3.75H3m18 0h-1.5M8.25 19.5V21M12 3v1.5m0 15V21m3.75-18v1.5m0 15V21m-9-1.5h10.5a2.25 2.25 0 0 0 2.25-2.25V6.75a2.25 2.25 0 0 0-2.25-2.25H6.75A2.25 2.25 0 0 0 4.5 6.75v10.5a2.25 2.25 0 0 0 2.25 2.25Zm.75-12h9v9h-9v-9Z"
                  />
                </svg>
              </span>
              <span>{{ $t('Oracle') }}</span>
            </router-link>
          </li>

          <li></li> -->

          <!-- <li v-if="$config.role === ERoles.client"></li> -->

          <!-- <li>
            <details id="disclosure-docs" :open="isOpenByRoute('tariffs_')">
              <summary class="group">
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
                {{ $t('Tariffs') }}
              </summary>
              <ul>
                <li>
                  <router-link class="group" :to="{ name: 'TariffsCustomer' }">
                    <span>{{ $t('Customer') }}</span>
                  </router-link>
                </li>
                <li>
                  <router-link class="group" :to="{ name: 'TariffsParntners' }">
                    <span>{{ $t('Parntner') }}</span>
                  </router-link>
                </li>
              </ul>
            </details>
          </li> -->
        </ul>

        <div>
          <div class="h-8"></div>

          <div class="px-2"><DocumentsList /></div>

          <div class="h-4"></div>

          <div class="flex gap-4 px-3">
            <a
              class="block text-base transition-colors p-1"
              target="_blank"
              href="mailto:support@aibetrade.com"
            >
              <svg
                class="w-8 h-8"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M16 12C16 14.2091 14.2091 16 12 16C9.79086 16 8 14.2091 8 12C8 9.79086 9.79086 8 12 8C14.2091 8 16 9.79086 16 12ZM16 12V13.5C16 14.8807 17.1193 16 18.5 16V16C19.8807 16 21 14.8807 21 13.5V12C21 7.02944 16.9706 3 12 3C7.02944 3 3 7.02944 3 12C3 16.9706 7.02944 21 12 21H16"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </a>

            <a
              class="block text-base transition-colors p-1"
              target="_blank"
              href="https://t.me/aibetradesupport"
            >
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 496 512" class="w-8 h-8">
                <path
                  fill="currentColor"
                  d="M248,8C111.033,8,0,119.033,0,256S111.033,504,248,504,496,392.967,496,256,384.967,8,248,8ZM362.952,176.66c-3.732,39.215-19.881,134.378-28.1,178.3-3.476,18.584-10.322,24.816-16.948,25.425-14.4,1.326-25.338-9.517-39.287-18.661-21.827-14.308-34.158-23.215-55.346-37.177-24.485-16.135-8.612-25,5.342-39.5,3.652-3.793,67.107-61.51,68.335-66.746.153-.655.3-3.1-1.154-4.384s-3.59-.849-5.135-.5q-3.283.746-104.608,69.142-14.845,10.194-26.894,9.934c-8.855-.191-25.888-5.006-38.551-9.123-15.531-5.048-27.875-7.717-26.8-16.291q.84-6.7,18.45-13.7,108.446-47.248,144.628-62.3c68.872-28.647,83.183-33.623,92.511-33.789,2.052-.034,6.639.474,9.61,2.885a10.452,10.452,0,0,1,3.53,6.716A43.765,43.765,0,0,1,362.952,176.66Z"
                />
              </svg>
            </a>
          </div>

          <div class="h-5"></div>

          <div class="px-4 pb-6">
            <div>{{ $t('project.addr') }}</div>
            <div>{{ $t('project.copyright', { year: new Date().getFullYear() }) }}</div>
          </div>

          <div class="px-4 opacity-30">
            <div class="cursor-pointer" @click="devMode">
              {{ $t('Build') }}: {{ $config.version }}
            </div>
            <!-- <div>
              <a href="#" class="underline" @click.prevent="downloadDebuggingFile">
                {{ $t('Debugging_file') }}
              </a>
            </div> -->
            <div class="cursor-pointer underline" @click="clearCacheAndReload">
              {{ $t('ClearCache') }}
            </div>
          </div>
        </div>
      </div>
    </aside>

    <div
      v-if="$screen.isMobile.value"
      class="menu-mobile-right safearea-r safearea-t safearea-b fixed top-2 mobile:right-4 bottom-2 flex flex-col items-end justify-between pc:hidden pointer-events-none text-gray-300"
    >
      <button
        for="drawer"
        class="btn-close btn btn-circle btn-ghost btn-sm pointer-events-auto"
        @click="closeMenu"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
      <div class="flex flex-col pointer-events-auto">
        <LocaleSelector />
        <SettingsSelector v-if="!$useTelegram().inTelegramApp" />
      </div>
    </div>
  </div>
</template>

<style scoped>
  @import 'index.css';
</style>
