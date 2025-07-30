<script lang="ts">
  export default {
    name: 'LayoutEmpty',
  }
</script>

<script setup lang="ts"></script>

<template>
  <router-view v-slot="{ Component }">
    <component :is="Component" />
  </router-view>
</template>
