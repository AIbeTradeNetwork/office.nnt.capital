<script lang="ts">
  export default {
    name: 'PayTelegramWarning',
  }
</script>

<script setup lang="ts">
  import { $t } from 'i18n/index'
  import { $authentication } from 'utils/authentication'
  import { $copyToClipboard } from 'utils/clipboard'
  import { $config } from 'utils/configuration'

  function createLink() {
    const url = new URL(location.href)
    url.searchParams.set('token', $authentication.getToken())
    url.searchParams.set('modal', 'addingFunds')

    return url.toString()
  }
</script>

<template>
  <div class="rounded-lg border border-error p-4 space-y-4">
    <div class="font-bold text-lg text-warning">{{ $t('info.payment_upgrading.title') }}</div>
    <div class="">{{ $t('info.payment_upgrading.text') }}</div>

    <button class="btn btn-warning" @click="$copyToClipboard(createLink())">
      {{ $t('info.payment_upgrading.copy') }}
    </button>
    <div class="text-info text-sm">
      {{ $t('info.payment_upgrading.if_fill_error') }}:
      <a class="link" :href="`https://t.me/${$config.supportLink}`" target="_blank">
        @{{ $config.supportLink }}
      </a>
    </div>
  </div>
</template>
