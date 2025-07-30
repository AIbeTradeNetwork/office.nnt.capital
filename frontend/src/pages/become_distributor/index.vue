<script lang="ts">
  export default {
    name: 'PageAuth',
  }
</script>

<script setup lang="ts">
  import { $modals } from 'utils/modals'
  import { $t } from 'i18n/index'
  import { $formatPriceByLocale } from 'utils/formats'
  import { $config } from 'utils/configuration'
  import { $store } from 'utils/store'

  function buy() {
    $modals.paySystem.show({
      type: 'distributor',
      what: $t('BuyDistributor'),
      price: $store.get('cfg').distributor_price,
      currency_code: $store.get('cfg').default_currency_code,
      precision: 2,
      code: '',
    })
  }
</script>

<template>
  <div class="hero h-full -padding-x">
    <div class="hero-overlay bg-opacity-50"></div>
    <div class="hero-content">
      <div class="prose prose-sm text-white">
        <h1 class="mb-5 text-4xl drop-shadow font-bold text-white">
          {{ $t('becomeDistributor.title') }}
        </h1>
        <p class="drop-shadow">{{ $t('becomeDistributor.text_1', { name: $config.appName }) }}</p>
        <ul>
          <li>
            <div class="mb-2 text-lg font-semibold drop-shadow">
              {{ $t('becomeDistributor.list_1.title') }}
            </div>
          </li>
          <li>
            <div class="mb-2 text-lg font-semibold drop-shadow">
              {{ $t('becomeDistributor.list_2.title') }}
            </div>
            <div>
              <blockquote class="text-sm italic text-gray">
                <p class="italic drop-shadow">
                  {{
                    $t('becomeDistributor.list_2.text_1', {
                      amount: $formatPriceByLocale({ count: 50 }),
                    })
                  }}
                </p>
              </blockquote>
              <blockquote class="text-sm italic text-gray">
                <p class="drop-shadow">{{ $t('becomeDistributor.list_2.text_2') }}</p>
              </blockquote>
              <blockquote class="text-sm italic text-gray">
                <p class="drop-shadow">{{ $t('becomeDistributor.list_2.text_3') }}</p>
              </blockquote>
            </div>
          </li>
          <li>
            <p class="mb-2 text-lg font-semibold drop-shadow">
              {{ $t('becomeDistributor.list_3.title') }}
            </p>

            <blockquote class="text-sm italic text-gray">
              <p class="italic drop-shadow">
                {{
                  $t('becomeDistributor.list_3.text_1', {
                    amount: $formatPriceByLocale({ count: 50 }),
                  })
                }}
              </p>
            </blockquote>
            <blockquote class="text-sm italic text-gray">
              <p class="drop-shadow">
                {{ $t('becomeDistributor.list_3.text_2') }}
              </p>
            </blockquote>
          </li>
        </ul>

        <button class="btn btn-primary my-6" @click="buy">
          {{ $t('becomeDistributor.name') }}
        </button>

        <p class="drop-shadow">
          {{ $t('becomeDistributor.text_2', { name: $config.appName }) }}
        </p>
      </div>
    </div>
  </div>
</template>
