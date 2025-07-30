<script lang="ts">
  export default {
    name: 'Clock',
  }
</script>

<script setup lang="ts">
  import { $i18nGlobal } from 'i18n/index'
  import { ref } from 'vue'

  const dateTime = ref('')

  function updateTime() {
    dateTime.value = new Date()
      .toLocaleString($i18nGlobal.locale.value, {
        month: 'short',
        weekday: 'short',
        day: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
      })
      .replace(/,|\./g, '')
    setTimeout(updateTime, 5000)
  }

  updateTime()
</script>

<template>
  <div class="flex text-2xl p-2">{{ dateTime }}</div>
</template>
