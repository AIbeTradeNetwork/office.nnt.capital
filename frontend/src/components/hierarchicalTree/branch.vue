<script lang="ts">
  export default {
    name: 'Branch',
  }
</script>

<script setup lang="ts">
  import Branch from './branch.vue'
  import BranchContent from './branchContent.vue'

  defineProps<{
    branchData: {
      data: TeamUser
      childrens: []
    }
  }>()
</script>

<template>
  <BranchContent :user="branchData.data" />

  <template v-if="branchData.childrens.length">
    <ul>
      <li v-for="(branchInner, index) in branchData.childrens" :key="index">
        <Branch :branch-data="branchInner" />
      </li>
    </ul>
  </template>
</template>
