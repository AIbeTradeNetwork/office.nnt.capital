<script setup lang="ts">
  import { useTonConnect } from 'utils/wallets/tonconnect'
  import { computed } from 'vue'
  const tonConnect = useTonConnect()

  const addr = computed(() => {
    const substSize = 5
    const addrString = tonConnect.wallet.getNonBounceableAddress()
    return addrString.replace(addrString.substring(substSize, addrString.length - substSize), '...')
  })

  const isConnected = computed(() => {
    return tonConnect.wallet.isConnected().value
  })
</script>
<template>
  <div v-if="isConnected">
    <div>
      <div class="flex items-start gap-2 max-w-[300px]">
        <div class="flex-[0_0_50px] h-[50px] overflow-hidden rounded-full">
          <img
            :src="tonConnect.wallet.info().imageUrl"
            :alt="tonConnect.wallet.info().openMethod"
          />
        </div>
        <div class="overflow-hidden">
          <div class="font-bold">{{ tonConnect.wallet.info().appName }}</div>
          <div class="truncate">{{ addr }}</div>
          <div class="flex gap-2">
            <a
              href="javascript:void(0)"
              class="text-info hover:underline whitespace-nowrap"
              @click.prevent="tonConnect.wallet.open"
            >
              {{ $t('OpenWallet') }}
            </a>
            <a
              href="javascript:void(0)"
              class="hover:text-error whitespace-nowrap"
              @click.prevent="tonConnect.wallet.disconnect"
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
                  d="M5.636 5.636a9 9 0 1 0 12.728 0M12 3v9"
                />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div v-if="!isConnected" class="border border-error rounded-md p-2 mb-4">
    {{ $t('info.wallet.wallets_count') }}
  </div>

  <button v-if="!isConnected" class="btn w-full btn-success" @click="tonConnect.modal.open">
    <svg class="w-8" viewBox="0 0 56 56" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M28 56C43.464 56 56 43.464 56 28C56 12.536 43.464 0 28 0C12.536 0 0 12.536 0 28C0 43.464 12.536 56 28 56Z"
        fill="#0098EA"
      />
      <path
        d="M37.5603 15.6277H18.4386C14.9228 15.6277 12.6944 19.4202 14.4632 22.4861L26.2644 42.9409C27.0345 44.2765 28.9644 44.2765 29.7345 42.9409L41.5381 22.4861C43.3045 19.4251 41.0761 15.6277 37.5627 15.6277H37.5603ZM26.2548 36.8068L23.6847 31.8327L17.4833 20.7414C17.0742 20.0315 17.5795 19.1218 18.4362 19.1218H26.2524V36.8092L26.2548 36.8068ZM38.5108 20.739L32.3118 31.8351L29.7417 36.8068V19.1194H37.5579C38.4146 19.1194 38.9199 20.0291 38.5108 20.739Z"
        fill="white"
      />
    </svg>

    {{ $t('ConnectWallet') }}
  </button>
</template>
