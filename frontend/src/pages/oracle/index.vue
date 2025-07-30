<script lang="ts">
  export default {
    name: 'PageOracle',
  }
</script>

<script setup lang="ts">
  import Info from 'components/info/index.vue'
  import { ESizes } from 'types/enums'
</script>

<template>
  <div class="page">
    <Info :size="ESizes.SM" :message="$t('info.oracle')" />
  </div>
</template>
