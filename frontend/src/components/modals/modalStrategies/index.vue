<script lang="ts">
  export default {
    name: 'ModalStrategiesList',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
  import StrategiesList from 'components/strategies/index.vue'
  import { EStrategyType } from 'types/enums'
  import { ref } from 'vue'

  const activeTab = ref<keyof typeof EStrategyType>(EStrategyType.classic)

  function close() {
    $modals.strategies.close()
  }
</script>

<template>
  <Modal :modal="$modals.strategies" :title="$t('Strategies')" size="screen">
    <div role="tablist" class="inline-flex tabs tabs-boxed mb-4">
      <a
        v-for="item in EStrategyType"
        :key="item"
        role="tab"
        class="tab"
        :class="`${activeTab === item ? 'tab-active' : ''}`"
        @click="activeTab = item"
      >
        <span class="first-letter:uppercase">{{ item }}</span>
      </a>
    </div>

    <StrategiesList :key="activeTab" :type="EStrategyType[activeTab]" @selected="close" />
  </Modal>
</template>
