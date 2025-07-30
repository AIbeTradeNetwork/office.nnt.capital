<script lang="ts">
  export default {
    name: 'ModalRankInfo',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $t } from 'i18n/index'
  import { CABTCounts, EABTRanks } from 'types/enums'
  import { $formatK } from 'utils/formats'
  import { $modals } from 'utils/modals'
  import { computed, onMounted } from 'vue'

  const rewardFriendsPurchases = {
    [EABTRanks.Novice]: [5],
    [EABTRanks.Farmer]: [7],
    [EABTRanks.Supervisor]: [7, 2],
    [EABTRanks.Expert]: [8, 5],
    [EABTRanks.Master]: [8, 6],
    [EABTRanks.Champion]: [9, 6],
    [EABTRanks.Veteran]: [9, 8],
    [EABTRanks.Legionary]: [10, 10],
    [EABTRanks.Grandmaster]: [12, 10],
    [EABTRanks.Archon]: [15, 10],
    [EABTRanks['archon-I']]: [15, 10],
    [EABTRanks['archon-II']]: [15, 10],
    [EABTRanks['archon-III']]: [15, 10],
    [EABTRanks['archon-IV']]: [15, 10],
    [EABTRanks['archon-V']]: [15, 10],
    [EABTRanks['archon-VI']]: [15, 10],
    [EABTRanks['archon-VII']]: [15, 10],
    [EABTRanks['archon-VIII']]: [15, 10],
    [EABTRanks['archon-IX']]: [15, 10],
    [EABTRanks['archon-X']]: [15, 10],
  }

  const toggleRangId = 9

  const getRankId = computed(() => {
    if (!('rankId' in $modals.ranks.data)) return -1
    const r = $modals.ranks.data['rankId']
    if (typeof r !== 'number' || r <= -1) return -1
    return $modals.ranks.data['rankId']
  })

  const isTier1 = computed(() => {
    return getRankId.value > -1 && getRankId.value < toggleRangId
  })

  const isTier2 = computed(() => {
    return getRankId.value > -1 && getRankId.value >= toggleRangId
  })

  const getRankList = computed(() => {
    const out = []
    const tier2Id = CABTCounts[EABTRanks['archon-I']].id
    Object.keys(CABTCounts).forEach((key) => {
      if (
        (isTier1.value && CABTCounts[key].id < tier2Id) ||
        (isTier2.value && CABTCounts[key].id >= tier2Id)
      ) {
        out.push(key)
        return
      }
    })
    return out
  })
</script>

<template>
  <Modal :modal="$modals.ranks" :z="50" :title="$t('Ranks')">
    <div class="space-y-4">
      <div
        v-for="rank in getRankList"
        :key="rank"
        class="flex items-center"
        :class="`${isTier2 ? 'gap-6 pl-2' : 'gap-4'}`"
      >
        <div class="relative w-[50px] self-start" :class="`${isTier2 ? 'mt-6' : 'mt-2'}`">
          <img
            v-if="CABTCounts[rank].img_top"
            class="absolute max-w-none w-[80px] -top-[25px] left-[50%] -translate-x-2/4 pointer-events-none opacity-80 z-0"
            :src="CABTCounts[rank].img_top"
            alt="coinTop"
          />
          <img
            v-if="CABTCounts[rank].img_bottom"
            class="absolute w-full h-full -bottom-[15px] left-0 pointer-events-none opacity-80 z-0"
            :src="CABTCounts[rank].img_bottom"
            alt="coinBottom"
          />
          <img v-if="isTier1" class="relative z-0" :src="CABTCounts[rank].img" alt="coin" />
          <img v-if="isTier2" class="relative z-0" :src="CABTCounts[rank].img_premium" alt="coin" />
        </div>
        <div>
          <div>
            <b>{{ rank }}</b>
            : {{ $formatK(CABTCounts[rank].min) }} - {{ $formatK(CABTCounts[rank].max) || '∞' }} ABT
          </div>
          <div class="opacity-80 text-sm">
            {{ $t('InvitedLimit') }}:
            {{ CABTCounts[rank].invitingMax > -1 ? CABTCounts[rank].invitingMax : '∞' }}
          </div>
          <div class="opacity-80 text-sm">
            <template v-if="isTier1">
              {{ $t('coin.Farming') }}: {{ CABTCounts[rank].farming }} ABT/{{ $t('Days', 1) }}
            </template>
            <template v-if="isTier2">
              {{ $t('coin.Farming') }}: {{ CABTCounts[rank].farming * 3 }} ABT/{{ $t('Days', 1) }}
            </template>
          </div>
          <div class="opacity-80 text-sm mt-2">
            <span class="font-bold">{{ $t('coin.rewardFriendsPurchases') }}:</span>
            <div v-for="(item, index) in rewardFriendsPurchases[rank]" :key="index">
              {{ $t('Level') }} {{ index + 1 }}:
              <span class="pl-1">{{ item }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Modal>
</template>
