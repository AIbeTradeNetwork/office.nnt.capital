<script lang="ts">
  export default {
    name: 'MyLine',
  }
</script>

<script setup lang="ts">
  import { $formatDate } from 'utils/date'
  import { $me } from 'utils/me'
  import Avatar from 'vue-boring-avatars'
</script>

<template>
  <tr>
    <td>
      <div class="flex items-center gap-2">
        <div class="avatar placeholder">
          <div class="bg-neutral text-neutral-content rounded-full w-12">
            <Avatar class="w-full h-full" variant="marble" :name="$me.data.uid" />
          </div>
        </div>
        <div>
          <p>{{ $me.data.nickname }}</p>
          <p>{{ $me.data.email }}</p>
        </div>
      </div>
    </td>
    <td>{{ $formatDate($me.data.created_at) }}</td>
    <td></td>
  </tr>
</template>