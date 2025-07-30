<script lang="ts">
  export default {
    name: 'PageAnalitics',
  }
</script>

<script setup lang="ts">
  import { onMounted, ref } from 'vue'
  import Info from 'components/info/index.vue'
  import { ESizes } from 'types/enums'
  const mounted = ref(false)
  onMounted(() => {
    mounted.value = true
  })
</script>

<template>
  <div class="page">
    <!-- <teleport v-if="mounted" to="#header_actions"> -->
    <button class="btn btn-primary">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-6 h-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
        />
      </svg>
      Add new
    </button>
    <!-- </teleport> -->

    <Info :size="ESizes.SM" :message="$t('info.analitics')" />
  </div>
</template>
