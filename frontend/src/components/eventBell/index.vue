<script lang="ts">
  export default {
    name: 'EventBell',
  }
</script>

<script setup lang="ts">
  import Dropdown from 'components/dropdown/index.vue'
</script>

<template>
  <Dropdown>
    <template #trigger>
      <div tabindex="0" role="button" class="btn btn-ghost btn-circle">
        <div class="indicator">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20">
            <path
              d="M16 7a5.38 5.38 0 0 0-4.46-4.85C11.6 1.46 11.53 0 10 0S8.4 1.46 8.46 2.15A5.38 5.38 0 0 0 4 7v6l-2 2v1h16v-1l-2-2zm-6 13a3 3 0 0 0 3-3H7a3 3 0 0 0 3 3z"
              fill="currentColor"
            ></path>
          </svg>
          <span class="badge badge-sm indicator-item text-1">+ 99</span>
        </div>
      </div>
    </template>
    <template #drop>
      <div class="card-body">
        <span class="font-bold text-lg">Что-то</span>
        <span class="text-info">Что-то еще</span>
        <div class="card-actions">
          <button class="btn btn-primary btn-block">btn</button>
        </div>
      </div>
    </template>
  </Dropdown>
</template>
