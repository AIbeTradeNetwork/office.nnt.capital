<script lang="ts">
  export default {
    name: 'UsersList',
  }
</script>

<script setup lang="ts">
  import { ref } from 'vue'
  import Avatar from 'vue-boring-avatars'

  const emit = defineEmits<{
    (e: 'selected', user: User): void
  }>()

  const userLists = ref<User[]>([])

  function addUser(user: User) {
    if (!user) return
    const lastUserUid = userLists.value[userLists.value.length - 1]?.uid
    if (lastUserUid === user.uid) return
    userLists.value.push(user)
  }

  function removeUsers(user: User) {
    if (!user) return
    const index = userLists.value.findIndex((item) => item.uid === user.uid)

    if (index > -1) {
      userLists.value.splice(index + 1, userLists.value.length)
    }
  }

  function select(item) {
    removeUsers(item)
    emit('selected', item)
  }

  function stepBack() {
    select(userLists.value[userLists.value.length - 2])
  }

  defineExpose({ addUser })
</script>

<template>
  <div
    class="select-none pc:absolute mobile:fixed safearea-b bottom-2 left-0 right-0 avatar-group overflow-hidden overflow-x-auto scroll-smooth scrollbar-stable mobile:px-2 pc:px-4"
  >
    <div
      v-if="userLists.length > 1"
      class="box rounded-full flex-[0_0_48px] h-[48px] px-0 py-0 cursor-pointer hover:bg-base-300"
      @click="stepBack"
    >
      <div class="flex w-full h-full items-center justify-center select-none">
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
            d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3"
          />
        </svg>
      </div>
    </div>

    <div
      v-for="item in [...userLists].reverse()"
      :key="item.uid"
      class="flex items-center [&:first-child_>_.arrow]:hidden"
    >
      <transition-group appear name="list">
        <div :key="item.uid" class="arrow">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="w-6 h-6 opacity-50 -rotate-90"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M9 8.25H7.5a2.25 2.25 0 0 0-2.25 2.25v9a2.25 2.25 0 0 0 2.25 2.25h9a2.25 2.25 0 0 0 2.25-2.25v-9a2.25 2.25 0 0 0-2.25-2.25H15M9 12l3 3m0 0 3-3m-3 3V2.25"
            />
          </svg>
        </div>

        <div
          :key="item.uid"
          class="box rounded-full px-0 py-1 pr-4 cursor-pointer hover:bg-base-300"
          @click="select(item)"
        >
          <div class="flex items-center select-none">
            <Avatar class="w-12" variant="marble" :name="item.uid" />
            <div class="text-xs whitespace-nowrap">
              <p>{{ item.nickname }}</p>
              <p>{{ item.email }}</p>
            </div>
          </div>
        </div>
      </transition-group>
    </div>
  </div>
</template>
