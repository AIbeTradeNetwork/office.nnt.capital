<script lang="ts">
  export default {
    name: 'ModalPromocode',
  }
</script>

<script setup lang="ts">
  import Modal from 'components/modals/modal/index.vue'
  import { $t } from 'i18n/index'
  import { $requests } from 'queries/index'
  import { ERoles } from 'types/enums'
  import { $config } from 'utils/configuration'
  import { $modals } from 'utils/modals'
  import { $notify } from 'utils/notify'
  import { ref } from 'vue'

  const loading = ref(false)
  const promocode = ref('')

  async function sendPromocode() {
    if (loading.value) return

    loading.value = true

    try {
      const response = await $requests.me.setRef(promocode.value)

      if (response) {
        $notify.show({
          type: 'success',
          text: $t('info.PromocodeApplied'),
        })
        $modals.promocode.close()
      }
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  function apply() {
    if ($config.role !== ERoles.client) {
      $notify.show({
        type: 'error',
        text: $t('info.OnlyForClients'),
      })
      return
    }

    sendPromocode()
  }

  function deny() {
    $modals.promocode.close()
  }
</script>

<template>
  <Modal :modal="$modals.promocode" :z="50">
    <template #title>
      <h3 class="font-bold text-lg flex items-center text-info">
        {{ $t('Promocode') }}
      </h3>
    </template>

    <div class="min-w-[20rem]">
      <template v-if="loading">
        <div class="flex justify-center items-center">{{ $t('Loading') }}...</div>
      </template>

      <template v-else>
        <div class="form m-0 mb-4">
          <div class="form-control">
            <input
              v-model.trim="promocode"
              type="text"
              class="w-full input input-bordered"
              name="login_login"
            />
          </div>
        </div>
      </template>
    </div>

    <template #actions>
      <div clas="flex gap-2 mt-4">
        <button class="btn btn-sm btn-primary btn-success" @click="apply">
          {{ $t('Apply') }}
        </button>

        <button class="btn btn-sm btn-outline ml-2" @click="deny">
          {{ $t('Close') }}
        </button>
      </div>
    </template>
  </Modal>
</template>
