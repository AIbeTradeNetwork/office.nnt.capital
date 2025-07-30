<script lang="ts">
  export default {
    name: 'Background',
  }
</script>

<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { $theme } from 'utils/theme'
  import lightBg from 'assets/backgrounds/light.jpg'
  import darkBg from 'assets/backgrounds/dark.jpg'

  const themeImgBg = ref(lightBg)
  const bgColorStyle = ref('')
  const currentTheme = computed(() => $theme.current.value)

  function setBgColorStyle() {
    if (currentTheme.value === $theme.themes.light) {
      return (bgColorStyle.value = 'background-color: rgba(255, 255, 255, 0.8)')
    }

    return (bgColorStyle.value = '')
  }

  function setImg() {
    if (currentTheme.value === $theme.themes.light) return (themeImgBg.value = lightBg)
    return (themeImgBg.value = darkBg)
  }

  watch(
    currentTheme,
    () => {
      setBgColorStyle()
      setImg()
    },
    {
      immediate: true,
    },
  )
</script>

<template>
  <div class="pointer-events-none">
    <img class="fixed top-0 left-0 botom-0 h-full w-full" :src="themeImgBg" alt="" />
    <div
      class="fixed top-0 left-0 botom-0 h-full w-full backdrop-blur-3xl"
      :style="bgColorStyle"
    ></div>
  </div>
</template>
