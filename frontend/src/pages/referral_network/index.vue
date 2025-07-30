<script lang="ts">
  export default {
    name: 'PageFriends',
  }
</script>

<script setup lang="ts">
  import { computed, ref } from 'vue'
  import HierarchicalTree from 'components/hierarchicalTree/index.vue'
  import HierarchicalTableRef from 'components/hierarchicalTableRef/index.vue'
  import HierarchicalTableMatch from 'components/hierarchicalTableMatch/index.vue'
  import { $config } from 'utils/configuration'
  import { ERoles } from 'types/enums'

  const activeTab = ref(0)
  const isClient = computed(() => $config.role === ERoles.client)
  const activeComponent = computed(() => {
    if (isClient.value) return HierarchicalTableRef
    if (activeTab.value === 0) return HierarchicalTableRef
    if (activeTab.value === 1) return HierarchicalTree
    if (activeTab.value === 2) return HierarchicalTableMatch
    return undefined
  })
</script>

<template>
  <div :class="activeTab === 1 ? 'relative w-full h-full overflow-hidden' : ''">
    <div class="flex">
      <div v-if="!isClient" role="tablist" class="m-4 mobile:m-2 tabs tabs-boxed z-[1]">
        <a
          role="tab"
          class="tab"
          :class="`${activeTab === 0 ? 'tab-active' : ''}`"
          @click="activeTab = 0"
        >
          Referral
        </a>

        <a
          role="tab"
          class="tab"
          :class="`${activeTab === 1 ? 'tab-active' : ''}`"
          @click="activeTab = 1"
        >
          Binary
        </a>

        <a
          role="tab"
          class="tab"
          :class="`${activeTab === 2 ? 'tab-active' : ''}`"
          @click="activeTab = 2"
        >
          Matching
        </a>
      </div>
    </div>

    <KeepAlive>
      <component :is="activeComponent"></component>
    </KeepAlive>
  </div>
</template>
