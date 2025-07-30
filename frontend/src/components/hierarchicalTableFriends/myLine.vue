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
          <!-- <p>
            Приглашенных:
            <span class="text-primary">11</span>
          </p> -->
        </div>
      </div>
    </td>
    <td>{{ $formatDate($me.data.created_at) }}</td>
    <td></td>
    <!-- <td>
      <div v-if="$me.data?.plan">
        <b>{{ $me.data.plan.code }}</b>
        {{ $t('before') }}: {{ $formatDate($me.data.plan.end_at) }}
      </div>
      <div v-else>-</div>
      <span class="badge badge-outline text-brilliant">Brilliant</span>
    </td>
    <td>
      <div v-if="$me.data?.rank">
        <b>{{ $me.data.rank.code }}</b>
        {{ $t('before') }}: {{ $formatDate($me.data.rank.end_at) }}
      </div>
      <div v-else>-</div>
    </td>
    <td>
      <div v-if="$me.data?.activity">
        {{ $t('before') }}: {{ $formatDate($me.data.activity.end_at) }}
      </div>
      <div v-else>-</div>
    </td> -->
  </tr>
</template>
