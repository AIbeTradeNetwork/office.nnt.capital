<script setup lang="ts">
  import { $t } from 'i18n/index'
  import Info from 'components/info/index.vue'
  import { EAlerts, ESizes } from 'types/enums'
  import { $isNotOurCompany, $isRefInviteEnds } from 'utils/checks'
  import { $modals } from 'utils/modals'
  import { $me } from 'utils/me'

  // const informationsList = computed(() => {
  //   return [
  //     {
  //       title: $t('coin.info_1.title'),
  //       text: $t('coin.info_1.text'),
  //     },
  //     {
  //       title: $t('coin.info_2.title'),
  //       text: $t('coin.info_2.text'),
  //     },
  //     {
  //       title: $t('coin.info_3.title'),
  //       text: $t('coin.info_3.text'),
  //     },
  //     {
  //       title: $t('coin.info_4.title'),
  //       text: $t('coin.info_4.text'),
  //     },
  //     {
  //       title: $t('coin.info_5.title'),
  //       text: $t('coin.info_5.text'),
  //     },
  //   ]
  // })
</script>
<template>
  <div
    v-if="!$isNotOurCompany && !$isRefInviteEnds"
    class="z-[2] absolute top-0 left-0 w-full h-full flex items-center justify-center"
  >
    <div></div>
    <Info
      class=""
      :type="EAlerts.BASE"
      :size="ESizes.BASE"
      :title="$t('info.isOurCompany.title')"
      :message="$t('info.isOurCompany.text')"
    >
      <template #after>
        <div class="">
          <button class="btn btn-primary mt-2 mr-2" @click="$modals.promocode.show()">
            <span>{{ $t('Promocode') }}</span>
          </button>
          <a class="btn btn-secondary mt-2" href="https://t.me/abtminerchat" target="_blank">
            <span>{{ $t('Chat') }}</span>
          </a>
        </div>
      </template>
    </Info>
  </div>

  <!-- $isNotOurCompany $isRefInviteEnds -->
  <div
    v-if="!$isNotOurCompany && $isRefInviteEnds"
    class="z-[2] absolute top-0 left-0 w-full h-full flex items-center justify-center"
  >
    <Info
      :type="EAlerts.BASE"
      :size="ESizes.BASE"
      :title="$t('info.isRefInviteEnds.title')"
      :message="$t('info.isRefInviteEnds.text', { promocode: $me.data.lim_ref_uid || '000000' })"
    >
      <template #after>
        <div class="">
          <button class="btn btn-primary mt-2 mr-2" @click="$modals.promocode.show()">
            <span>{{ $t('Promocode') }}</span>
          </button>
          <a class="btn btn-success mt-2" href="https://t.me/abtminerchat" target="_blank">
            <span>{{ $t('Chat') }}</span>
          </a>
        </div>
      </template>
    </Info>
  </div>
</template>
