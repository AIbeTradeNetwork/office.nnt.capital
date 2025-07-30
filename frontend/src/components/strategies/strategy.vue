<script lang="ts">
  export default {
    name: 'Strategies',
  }
</script>

<script setup lang="ts">
  import { $t } from 'i18n/index'
  import { $formatInt } from 'utils/formats'
  import { EStrategyType } from 'types/enums'
  import { $modals } from 'utils/modals'
  import { $startTimer } from 'utils/timeout'
  import { onMounted, ref } from 'vue'
  import { $riskLevelBadge } from 'utils/colors'

  const emits = defineEmits(['selected'])

  const { data, soon } = defineProps<{
    data: Strategy
    soon?: boolean
  }>()

  const timeOut = ref('00:00:00')
  const timer = ref<ReturnType<typeof setTimeout>>()

  // function showModalInfo(item: Strategy) {
  //   $modals.info.show({
  //     title: item.name,
  //     text: item.code,
  //   })
  // }

  function showFullDescription(item: Strategy) {
    $modals.info.show({
      title: item.name,
      text: $t(`strategies.${data.code}`),
    })
  }

  function select(strategy: Strategy) {
    if (soon) return

    emits('selected')

    $modals.bot.show(strategy)
  }

  function startTimeOut() {
    if (data.type !== EStrategyType.ico) return

    timer.value = $startTimer(data.start_at, timeOut)
  }

  onMounted(() => {
    startTimeOut()
  })
</script>

<template>
  <div v-if="data" class="box relative p-0 overflow-hidden border border-base-300 bg-base-200">
    <div class="flex h-full items-stretch">
      <div class="p-4 text-start flex flex-col justify-between flex-1">
        <div>
          <div class="flex justify-between items-center space-x-2">
            <div class="text-lg font-bold flex-[0_1_70%]">{{ data.name }}</div>
            <div v-if="data.is_new" class="text-sm badge badge-error font-mono">NEW!</div>
            <div v-if="data.is_free" class="text-sm badge badge-success font-mono">FREE!</div>
          </div>

          <div class="text-sm flex flex-wrap gap-2 mt-2">
            <template v-if="(data.symbol_codes || []).length">
              <p
                v-for="symbol_code in data.symbol_codes"
                :key="symbol_code"
                class="badge badge-warning badge-md text-xs"
              >
                {{ symbol_code }}
              </p>
            </template>

            <template v-else>
              <p class="badge badge-warning badge-md">
                {{ 'All coins' }}
              </p>
            </template>

            <p v-if="data.exchange_code" class="badge badge-outline badge-md">
              <!-- {{ $t('Exchange') }}:&nbsp; -->
              <b>{{ data.exchange_code }}</b>
            </p>

            <p v-if="data.category" class="badge badge-outline badge-md">
              <!-- {{ $t('Category') }}:&nbsp; -->
              <b>{{ $t(data.category) }}</b>
            </p>

            <p
              v-if="data.risk_level"
              class="badge badge-outline badge-md"
              :class="$riskLevelBadge[data.risk_level]"
            >
              {{ $t('RiskLevel') }}:&nbsp;
              <b>{{ data.risk_level }}</b>
            </p>

            <p v-if="data.pos_profit > 0" class="badge badge-outline badge-md">
              <!-- {{ $t('PotentialProfitability') }}:&nbsp; -->
              {{ $t('before') }}&nbsp;+
              <b>{{ $formatInt(data.pos_profit) }}%</b>
            </p>

            <p v-if="data.min_deposit > 0" class="badge badge-outline badge-md">
              <!-- {{ $t('PotentialProfitability') }}:&nbsp; -->
              {{ $t('MinDeposit') }}:&nbsp;
              <b>{{ $formatInt(data.min_deposit) }}</b>
            </p>

            <p v-if="data.max_deposit > 0" class="badge badge-outline badge-md">
              <!-- {{ $t('PotentialProfitability') }}:&nbsp; -->
              {{ $t('MaxDeposit') }}:&nbsp;
              <b>{{ $formatInt(data.max_deposit) }}</b>
            </p>

            <p v-if="data.share_profit > 0" class="badge badge-info badge-outline badge-md">
              {{ $t('ShareProfit') }}:&nbsp;
              <b>{{ $formatInt(data.share_profit) }}%</b>
            </p>
          </div>

          <p class="text-sm mt-2 opacity-60 overflow-hidden text-ellipsis line-clamp-2">
            {{ $t(`strategies.${data.code}`) }}
          </p>
        </div>

        <!-- <div class="flex flex-wrap gap-2 text-sm mt-2">
        <div>
            <span class="opacity-60">Risk:</span>
            {{ item.risk }}
        </div>
        <div>
            <span class="opacity-60">Par1:</span>
            {{ item.param1 }}
        </div>
        <div>
            <span class="opacity-60">Par2:</span>
            {{ item.param2 }}
        </div>
        </div> -->

        <div class="flex justify-between items-center gap-4 mt-4">
          <div>
            <div v-if="soon" class="badge badge-outline badge-accent">{{ $t('Soon') }}</div>
          </div>

          <div class="flex gap-4">
            <button
              v-if="data"
              class="btn btn-info btn-outline btn-sm"
              @click="showFullDescription(data)"
            >
              {{ $t('Description') }}
            </button>
            <button :disabled="soon" class="btn btn-primary btn-sm" @click="select(data)">
              {{ $t('Select') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
