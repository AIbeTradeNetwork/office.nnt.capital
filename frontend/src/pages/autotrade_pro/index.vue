<script lang="ts">
  export default {
    name: 'PageAuth',
  }
</script>

<script setup lang="ts">
  import { $textSplit } from 'utils/text'
  import trade_img from 'assets/autotradePRO/img-white.png'
  import binance from 'assets/autotradePRO/binance.svg'
  import bybit from 'assets/autotradePRO/bybit.svg'
  import okx from 'assets/autotradePRO/okx.svg'
  import bing_x from 'assets/autotradePRO/bing-x.svg'
  import bitget from 'assets/autotradePRO/biget.svg'
</script>

<template>
  <div class="py-8 px-4 max-w-[768px] mx-auto">
    <h2 class="text-center mb-8 text-2xl font-bold">{{ $t('autotrade_pro.name') }}</h2>

    <h3 class="text-center my-8 text-xl font-bold">
      {{ $t('autotrade_pro.part_1.description.title') }}
    </h3>
    <p v-for="par in $textSplit($t('autotrade_pro.part_1.description.text'))" :key="par">
      {{ par }}
    </p>

    <div class="overflow-hidden rounded-lg my-8 drop-shadow-xl">
      <img class="-mb-1" :src="trade_img" alt="" />
    </div>

    <h3 class="text-center my-8 text-xl font-bold">
      {{ $t('autotrade_pro.part_1.procedure.title') }}
    </h3>
    <p v-for="par in $textSplit($t('autotrade_pro.part_1.procedure.text'))" :key="par">{{ par }}</p>

    <h3 class="text-center my-8 text-xl font-bold">
      {{ $t('autotrade_pro.part_1.exchanges.title') }}
    </h3>
    <div class="flex flex-wrap gap-4 items-stretch my-8">
      <a
        href="https://www.binance.com"
        class="inline-flex items-center py-2 px-4 rounded-md bg-white/10 h-[50px]"
        target="_blank"
      >
        <img class="h-[30px]" :src="binance" alt="" />
      </a>
      <a
        href="https://bybit.com"
        class="inline-flex items-center py-2 px-4 rounded-md bg-white/10 h-[50px]"
        target="_blank"
      >
        <img class="h-[30px]" :src="bybit" alt="" />
      </a>
      <a
        href="https://www.okx.com"
        class="inline-flex items-center py-2 px-4 rounded-md bg-white/10 h-[50px]"
        target="_blank"
      >
        <img class="h-[30px]" :src="okx" alt="" />
      </a>

      <a
        href="https://bingx.com"
        class="inline-flex items-center py-2 px-4 rounded-md bg-white/10 h-[50px]"
        target="_blank"
      >
        <img class="h-[100px]" :src="bing_x" alt="" />
      </a>
      <a
        href="https://www.bitget.com"
        class="inline-flex items-center py-2 px-4 rounded-md bg-white/10 h-[50px]"
        target="_blank"
      >
        <img class="h-[30px]" :src="bitget" alt="" />
      </a>
    </div>
    <p v-for="par in $textSplit($t('autotrade_pro.part_1.exchanges.text'))" :key="par">{{ par }}</p>

    <h3 class="text-center my-8 text-xl font-bold">
      {{ $t('autotrade_pro.part_1.mutual.title') }}
    </h3>
    <p v-for="par in $textSplit($t('autotrade_pro.part_1.mutual.text'))" :key="par">{{ par }}</p>

    <h3 class="text-center my-8 text-xl font-bold">
      {{ $t('autotrade_pro.part_1.refusal.title') }}
    </h3>
    <p v-for="par in $textSplit($t('autotrade_pro.part_1.refusal.text'))" :key="par">{{ par }}</p>

    <p
      class="my-8 text-xl font-bold [&_>_a]:underline [&_>_a]:text-info"
      v-html="$t('autotrade_pro.part_1.link.end')"
    ></p>
  </div>
</template>
