<script lang="ts">
  export default {
    name: 'FormAuth',
  }
</script>

<script setup lang="ts">
  import { $requests } from 'queries/index'
  import { $initTelegramAuthBtnScript } from 'utils/auth/telegram'
  import { $config } from 'utils/configuration'
  import { ref, reactive, onMounted } from 'vue'

  const loading = ref(false)

  const form: AuthLogin = reactive({
    login: '',
    password: '',
  })

  async function send() {
    if (loading.value) return
    loading.value = true
    try {
      await $requests.auth.login(form)
    } catch (error) {
      console.error(error)
    }

    loading.value = false
  }

  onMounted(() => {
    $initTelegramAuthBtnScript(loading)
  })
</script>

<template>
  <div
    class="flex w-full justify-center transition-shadow duration-100 [transform:translate3d(0,0,0)]"
    :class="loading ? 'pointer-events-none opacity-80' : ''"
  >
    <div class="flex flex-col gap-4 bg-base-200 w-full p-6 rounded-md drop-shadow-custom max-w-lg">
      <h1 class="text-3xl font-semibold text-center">
        {{ $t('Authorization') }}
      </h1>

      <div class="form m-0 mb-4">
        <div class="form-control">
          <label class="label">
            <span class="text-base label-text">{{ $t('LoginOrEmail') }}</span>
          </label>
          <input
            v-model.trim="form.login"
            type="text"
            class="w-full input input-bordered"
            name="login_login"
            enterkeyhint="send"
            @keyup.enter="send"
          />
        </div>

        <div class="form-control">
          <label class="label">
            <span class="text-base label-text">{{ $t('Password') }}</span>
          </label>
          <input
            v-model.trim="form.password"
            enterkeyhint="send"
            type="password"
            class="w-full input input-bordered"
            name="login_password"
            @keyup.enter="send"
          />
        </div>
      </div>

      <button class="btn btn-primary" @click="send">
        <span v-show="loading" class="loading loading-spinner"></span>
        {{ $t('SignIn') }}
      </button>

      <div class="flex items-center w-full text-sm mt-2">
        <hr class="w-full" />
        <p class="px-3">{{ $t('OR') }}</p>
        <hr class="w-full" />
      </div>

      <div id="telegram-login" class="flex justify-center"></div>

      <div class="flex items-center w-full mt-2">
        <hr class="w-full" />
        <p class="px-3 text-sm">{{ $t('OR') }}</p>
        <hr class="w-full" />
      </div>

      <div class="flex justify-center text-sm">
        {{ $t('NotAlready') }}
        <router-link :to="{ name: 'Registration' }" class="link link-primary ml-2">
          {{ $t('Registration') }}
        </router-link>
      </div>
    </div>
  </div>
</template>
