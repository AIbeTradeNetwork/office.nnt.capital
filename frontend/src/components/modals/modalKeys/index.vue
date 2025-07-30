<script lang="ts">
  export default {
    name: 'ModalKeys',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $modals } from 'utils/modals'
  import { computed, ref } from 'vue'
  import { $store } from 'utils/store'
  import { $requests } from 'queries/index'

  const modalData = $modals.keys.data
  const loading = ref(false)

  const keys = computed(() => {
    return $store.get('keys').filter((item) => item.exchange_code === modalData.code)
  })

  const editingData = ref<(KeyAddReq | ({ uid: string } & KeyEditReq)) | null>(null)

  function adding() {
    editingData.value = {
      name: '',
      exchange_code: modalData.code,
      key: '',
      secret: '',
    }
  }

  function edit(keyData: Key) {
    editingData.value = {
      name: keyData.name,
      uid: keyData.uid,
      key: keyData.key,
      secret: '',
    }
  }

  async function editApply() {
    if (!editingData.value) return

    loading.value = true

    try {
      if ('uid' in editingData.value) {
        await $requests.keys.edit(editingData.value)
      }
      if ('exchange_code' in editingData.value) {
        await $requests.keys.add(editingData.value)
      }
    } catch (error) {
      console.error(error)
    }

    editingData.value = null

    loading.value = false
  }

  function editCancel() {
    editingData.value = null
  }
</script>

<template>
  <Modal :modal="$modals.keys" :z="50">
    <template #title>
      <h3 class="font-bold text-lg flex items-center gap-2 mb-4">
        {{ $t('Keys') }}
        <button class="btn btn-xs btn-primary" @click="adding">
          <span>{{ $t('Add') }} {{ $t('Key').toLowerCase() }}</span>
        </button>
      </h3>
    </template>

    <!-- List -->
    <div class="min-w-[20rem]">
      <div v-if="loading">
        <div class="text-center font-bold my-10">{{ $t('Loading') }}</div>
      </div>

      <div v-else-if="!editingData && !keys.length">
        <div class="text-center font-bold my-10">{{ $t('NoKeys') }}</div>
      </div>

      <div v-else-if="!editingData && keys.length">
        <div class="flex flex-col gap-3">
          <div v-for="item in keys" :key="item.key" class="flex items-center gap-2">
            <span class="cursor-pointer pc:hover:text-success" @click="edit(item)">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-6 h-6 relative top-[-2px]"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10"
                />
              </svg>
            </span>
            <div>
              <div class="whitespace-nowrap font-bold">{{ $t('Key') }}: {{ item.name }}</div>
              <div class="italic text-xs opacity-80">{{ item.uid }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Editing -->
      <div v-else-if="editingData">
        <div class="form m-0 mb-6">
          <div class="flex items-center w-full text-sm mt-4">
            <hr class="w-full" />
            <p class="px-3 font-bold">{{ $t('Adding') }}</p>
            <hr class="w-full" />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="text-base label-text">
                {{ $t('Key') }} {{ $t('name').toLowerCase() }}
              </span>
            </label>
            <input
              v-model.trim="editingData.name"
              type="text"
              class="w-full input input-bordered"
              name="keyedit_key"
            />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="text-base label-text">{{ $t('Key') }}</span>
            </label>
            <input
              v-model.trim="editingData.key"
              type="text"
              class="w-full input input-bordered"
              name="keyedit_key"
            />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="text-base label-text">{{ $t('SecretKey') }}</span>
            </label>
            <input
              v-model.trim="editingData.secret"
              type="text"
              class="w-full input input-bordered"
              name="keyedit_s_key"
            />
          </div>
        </div>

        <div class="flex gap-2 justify-end">
          <button class="btn btn-sm btn-outline btn-error" @click="editCancel">
            {{ $t('Cancel') }}
          </button>
          <button class="btn btn-sm btn-success" @click="editApply">{{ $t('Apply') }}</button>
        </div>
      </div>
    </div>
  </Modal>
</template>
