<script lang="ts">
  export default {
    name: 'LocaleSelector',
  }
</script>

<script setup lang="ts">
  import Dropdown from 'components/dropdown/index.vue'
  import { $i18nGlobal, $i18nOptions, $setLocale } from 'i18n/index'
  import { computed, ref } from 'vue'
  const currentLocale = ref($i18nGlobal.locale)
  const currentLocaleOptions = computed(() => {
    return $i18nOptions.localesOptions.find((item) => {
      return item.locale === currentLocale.value
    })
  })

  function setLocale(locale: string) {
    $setLocale(locale)
    ;(document.activeElement.closest('button.btn') as HTMLButtonElement).blur()
  }
</script>

<template>
  <Dropdown>
    <template #trigger>
      <button
        tabindex="0"
        class="btn btn-ghost btn-circle avatar [&_>svg]:w-7"
        v-html="currentLocaleOptions.icon"
      ></button>
    </template>
    <template #drop>
      <div class="p-2 menu-horizontal">
        <template v-for="item in $i18nOptions.localesOptions">
          <div v-if="item.name !== currentLocale" :key="item.name">
            <button
              class="inline-flex justify-center content-center btn btn-circle btn-ghost [&_>svg]:w-7"
              @click="setLocale(item.locale)"
              v-html="item.icon"
            ></button>
          </div>
        </template>
      </div>
    </template>
  </Dropdown>
</template>
