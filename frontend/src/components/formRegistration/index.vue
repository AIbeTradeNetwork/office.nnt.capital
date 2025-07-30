<script lang="ts">
  export default {
    name: 'FormReg',
  }
</script>

<script setup lang="ts">
  // import { $getLocaleOptions, $i18nGlobal } from 'i18n/index'
  // import { $requests } from 'queries/index'
  import { $initTelegramAuthBtnScript } from 'utils/auth/telegram'
  import { $ref } from 'utils/ref'
  import { createTelegrammRefLink } from 'utils/shares'
  import { computed, onMounted, ref } from 'vue'

  const loading = ref(false)
  const promocode = ref($ref.getFromCookies())
  const tgAppLinkFull = computed(() => createTelegrammRefLink(promocode.value))

  // const refId = $ref.getFromCookies()
  // const locale = computed(() => {
  //   const mappedOption = $getLocaleOptions($i18nGlobal.locale.value)
  //   return mappedOption.map
  // })

  // const form: AuthRegister = reactive({
  //   email: '',
  //   nickname: '',
  //   ref_uid: refId || '',
  //   password: '',
  //   repassword: '',
  //   locale: locale.value,
  // })

  // async function send() {
  //   if (loading.value) return
  //   loading.value = true

  //   try {
  //     await $requests.auth.register(form)
  //     loading.value = false
  //   } catch (error) {
  //     console.error(error)
  //   }

  //   loading.value = false
  // }

  function savePromocode() {
    $ref.saveToCookies(promocode.value)
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
        {{ $t('Registration') }}
      </h1>
      <div class="form m-0 mb-4">
        <!-- <div class="form-control">
          <label class="label">
            <span class="text-base label-text">{{ $t('Email') }}</span>
          </label>
          <input
            v-model.trim="form.email"
            type="text"
            class="w-full input input-bordered"
            name="reg_email"
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
            type="password"
            class="w-full input input-bordered"
            name="reg_password"
            enterkeyhint="send"
            @keyup.enter="send"
          />
        </div>

        <div class="form-control">
          <label class="label">
            <span class="text-base label-text">
              {{ $t('ConfirmPassword') }}
            </span>
          </label>
          <input
            v-model.trim="form.repassword"
            type="password"
            class="w-full input input-bordered"
            name="reg_repassword"
            enterkeyhint="send"
            @keyup.enter="send"
          />
        </div>

        <div class="form-control">
          <label class="label">
            <span class="text-base label-text">{{ $t('Nickname') }}</span>
          </label>
          <input
            v-model.trim="form.nickname"
            type="text"
            class="w-full input input-bordered"
            name="reg_name"
            enterkeyhint="send"
            @keyup.enter="send"
          />
        </div> -->

        <div class="form-control">
          <label class="label">
            <span class="text-base label-text">
              {{ $t('Promocode') }}
            </span>
          </label>
          <input
            v-model.trim="promocode"
            type="text"
            class="w-full input input-bordered"
            name="reg_ref"
            enterkeyhint="save"
            @keyup="savePromocode"
            @change="savePromocode"
          />
        </div>
      </div>

      <!-- <button class="btn btn-block btn-primary" @click="savePromocode">
        <span v-show="loading" class="loading loading-spinner"></span>
        {{ $t('Save') }} {{ $t('Promocode').toLowerCase() }}
      </button> -->

      <!-- <button class="btn btn-block btn-primary" @click="send">
        <span v-show="loading" class="loading loading-spinner"></span>
        {{ $t('SignUp') }}
      </button> -->

      <!-- <div class="flex items-center w-full text-sm mt-2">
        <hr class="w-full" />
        <p class="px-3">{{ $t('OR') }}</p>
        <hr class="w-full" />
      </div> -->

      <div id="telegram-login" class="flex justify-center"></div>

      <div class="flex items-center w-full text-sm mt-2">
        <hr class="w-full" />
        <p class="px-3">{{ $t('OR') }}</p>
        <hr class="w-full" />
      </div>

      <div class="text-center">
        <a
          :href="tgAppLinkFull"
          target="_blank"
          class="inline-flex items-center gap-1 link link-primary"
        >
          AiBeTrade Bot
          <svg
            viewBox="0 0 48 48"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            class="relative top-[1px] w-4"
          >
            <path
              d="M19 11H37V29"
              stroke="currentColor"
              stroke-width="4"
              stroke-linecap="butt"
              stroke-linejoin="bevel"
            ></path>
            <path
              d="M11.5439 36.4559L36.9997 11"
              stroke="currentColor"
              stroke-width="4"
              stroke-linecap="butt"
              stroke-linejoin="bevel"
            ></path>
          </svg>
        </a>
      </div>

      <div class="flex items-center w-full text-sm mt-2">
        <hr class="w-full" />
        <p class="px-3">{{ $t('OR') }}</p>
        <hr class="w-full" />
      </div>

      <div class="flex justify-center text-sm space-y-2">
        {{ $t('AlreadyHave') }}

        <router-link :to="{ name: 'Authorization' }" class="link link-primary ml-2">
          {{ $t('Authorization') }}
        </router-link>
      </div>
    </div>
  </div>
</template>
