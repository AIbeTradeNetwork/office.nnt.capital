<script lang="ts">
  export default {
    name: 'PageRoadmap',
  }
</script>

<script setup lang="ts">
  import Info from 'coin/components/info.vue'
  import Roadmap from 'coin/components/timeline.vue'
</script>

<template>
  <div class="page flex flex-col gap-10 items-center">
    <Info />
    <Roadmap />
  </div>
</template>
