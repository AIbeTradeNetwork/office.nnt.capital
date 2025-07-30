<script lang="ts">
  export default {
    name: 'QualificationsClient',
  }
</script>

<script setup lang="ts">
  import { EHideClasses } from 'types/enums'
  import { $isNotOurCompany, $isNotTMeMail } from 'utils/checks'
  import { $copyToClipboard } from 'utils/clipboard'
  import { $formatDate } from 'utils/date'
  import { hideByRoute } from 'utils/hideByRoute'
  import { $me } from 'utils/me'
  import { $modals } from 'utils/modals'
</script>

<template>
  <div class="p-6 flex flex-col gap-4">
    <div v-if="hideByRoute(EHideClasses.autotrade)" class="flex items-center text-sm">
      <div class="mobile:max-w-[55px] truncate">
        {{ $t(`roles.${$me.data?.role.toLowerCase()}`) }}
      </div>
      <!-- <span class="mx-1">/</span>
      <div v-if="$me.data.plan">{{ $me.data.plan.code }}</div>
      <div v-else>
        <router-link :to="{ name: 'Tariffs' }" class="btn btn-outline btn-xs" @click.stop>
          {{ $t('Plan') }}
        </router-link>
      </div> -->
    </div>

    <template v-if="$isNotOurCompany">
      <div v-if="$me.data.ref_uid" class="text-base-content">
        {{ $t('ReferralPromoCode') }}:
        <span class="text-secondary">
          {{ $me.data.ref_uid }}
        </span>
      </div>
    </template>
    <template v-else>
      <button class="btn btn-primary" @click="$modals.promocode.show()">
        <span>{{ $t('Promocode') }}</span>
      </button>
    </template>

    <div class="text-base-content">
      {{ $t('MyPromoCode') }}:
      <span
        class="whitespace-nowrap underline cursor-pointer text-gradient"
        @click="$copyToClipboard($me.data.uid)"
      >
        {{ $me.data.uid }}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-4 h-4 inline relative -top-[3px]"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 0 1-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 0 1 1.5.124m7.5 10.376h3.375c.621 0 1.125-.504 1.125-1.125V11.25c0-4.46-3.243-8.161-7.5-8.876a9.06 9.06 0 0 0-1.5-.124H9.375c-.621 0-1.125.504-1.125 1.125v3.5m7.5 10.375H9.375a1.125 1.125 0 0 1-1.125-1.125v-9.25m12 6.625v-1.875a3.375 3.375 0 0 0-3.375-3.375h-1.5a1.125 1.125 0 0 1-1.125-1.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H9.75"
          />
        </svg>
      </span>
    </div>

    <div v-if="$isNotTMeMail" class="text-base-content">
      {{ $t('Email') }}:
      <span class="">
        {{ $me.data.email }}
      </span>
    </div>

    <div v-if="hideByRoute(EHideClasses.autotrade)" class="-mt-4">
      <div class="text-base-content flex items-center gap-2">
        {{ $t('Plan') }}:
        <span v-if="$me.data.plan" class="font-bold">{{ $me.data.plan.code }}</span>
        <!-- <router-link v-else :to="{ name: 'Tariffs' }" class="btn btn-xs btn-outline">
          {{ $t('Select') }}
          <svg
            class="w-3"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="m16.2 19 4.8-7-4.8-7H3l4.8 7L3 19h13.2Z"
            />
          </svg>
        </router-link> -->
      </div>

      <div v-if="$me.data.plan" class="text-base-content">
        {{ $t('StartDate') }}:
        <span class="text-success">
          {{ $formatDate($me.data.plan.start_at, 'Pp') }}
        </span>
      </div>

      <div v-if="$me.data.plan" class="text-base-content">
        {{ $t('EndDate') }}:
        <span class="text-success">
          {{ $formatDate($me.data.plan.end_at, 'Pp') }}
        </span>
      </div>

      <!-- <div>
      <button class="btn btn-xs btn-outline btn-primary">
        {{ $t('ToExtend') }}
        <svg
          class="w-3"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <path
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="m16 10 3-3m0 0-3-3m3 3H5v3m3 4-3 3m0 0 3 3m-3-3h14v-3"
          />
        </svg>
      </button>
    </div> -->
    </div>

    <div v-if="$me.data.is_premium" class="text-base-content">
      {{ $t('Premium') }}
      <span>{{ $t('before') }} {{ $formatDate($me.data.premium_until, 'Pp') }}</span>
    </div>

    <div v-if="$me.data.is_autofarm" class="text-base-content">
      {{ $t('Autofarming') }}
      <span>{{ $t('before') }} {{ $formatDate($me.data.autofarm_until, 'Pp') }}</span>
    </div>

    <div v-if="$me.data.active_boost && $me.data.active_boost.end_at" class="text-base-content">
      {{ $t('Boost') }}
      <span>{{ $t('before') }} {{ $formatDate($me.data.active_boost.end_at, 'Pp') }}</span>
    </div>
  </div>
</template>
