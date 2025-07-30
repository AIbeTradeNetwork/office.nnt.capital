<script lang="ts">
  export default {
    name: 'UserCard',
  }
</script>

<script setup lang="ts">
  import { ESizes } from 'types/enums'
  import { defineProps, ref, toRefs } from 'vue'

  const props = defineProps<{
    size: ESizes
  }>()

  const { size } = toRefs(props)

  const isSM = ref(size.value === ESizes.SM)
</script>

<template>
  <div class="card w-96 bg-base-100" :class="{ 'card-compact': isSM }">
    <figure class="p-4 pt-3 pb-0">
      <div class="avatar">
        <div class="w-24 h-24 rounded-full">
          <img src="https://daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg" />
        </div>
      </div>
    </figure>
    <div v-if="!isSM" class="card-body">
      <h2 class="card-title justify-center">Qwerty Qwertovna</h2>
      <p class="text-center">SOme text. SOme text. SOme text. SOme text.SOme text.</p>
      <div class="card-actions justify-end"></div>
    </div>
  </div>
</template>
