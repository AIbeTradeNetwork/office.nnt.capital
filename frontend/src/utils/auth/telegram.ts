import { $requests } from 'queries/index'
import { $authentication } from 'utils/authentication'
import { $cookies } from 'utils/cookies'
import { $modals } from 'utils/modals'
import { $notify } from 'utils/notify'
import { $ref } from 'utils/ref'
import { Ref } from 'vue'

/**
 * Sends a Telegram authentication request.
 *
 * @param {TelegramAuth} form - The Telegram authentication form containing the user's information.
 * @return {Promise<boolean>} A promise that resolves to true if the authentication is successful, or rejects with 0 if an error occurs.
 */
async function sendTelegramAuth(form: TelegramAuth) {
  try {
    await $requests.auth.telegram(form)
    return true
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

async function sendTelegramAppAuth(input: string) {
  try {
    await $requests.auth.telegramApp(input)
    return true
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

/**
 * Initializes the Telegram btn authentication script.
 *
 * @param {Ref<boolean>} loading - The loading state reference.
 */
export function $initTelegramAuthBtnScript(loading: Ref<boolean>) {
  const container = document.getElementById('telegram-login')

  if (!container) return
  ;(window as any).onTelegramAuth = async (user) => {
    if (loading.value) return

    loading.value = true

    const form: TelegramAuth = {
      id: user.id ? user.id + '' : '',
      first_name: user.first_name || '',
      last_name: user.last_name || '',
      username: user.username || '',
      photo_url: user.photo_url || '',
      auth_date: user.auth_date ? user.auth_date + '' : '',
      hash: user.hash || '',
      ref_uid: $ref.getFromCookies() || '',
    }

    await sendTelegramAuth(form)

    loading.value = false
  }

  const script = document.createElement('script')
  script.setAttribute('src', 'https://telegram.org/js/telegram-widget.js?22')
  script.async = true
  script.setAttribute('data-telegram-login', 'nntcapitalbot')
  script.setAttribute('data-size', 'large')
  script.setAttribute('data-radius', '8')
  script.setAttribute('data-onauth', 'onTelegramAuth(user)')
  script.setAttribute('data-request-access', 'write')
  container.appendChild(script)
}

export async function $initTelegramAuthFrameAppGuard() {
  const WebApp = (window as any)?.Telegram?.WebApp
  if (!WebApp || !WebApp.initData) return true

  const initDataUnsafe = WebApp.initDataUnsafe
  if (initDataUnsafe.start_param) $ref.saveToCookies(initDataUnsafe.start_param)

  if (!WebApp.isExpanded) WebApp.expand()

  // if ($authentication.getToken()) return true

  if (initDataUnsafe.hash) {
    // const form: TelegramAuth = {
    //   id: initDataUnsafe.user.id ? initDataUnsafe.user.id + '' : '',
    //   first_name: initDataUnsafe.user.first_name || '',
    //   last_name: initDataUnsafe.user.last_name || '',
    //   username: initDataUnsafe.user.username || '',
    //   photo_url: initDataUnsafe.user.photo_url || '',
    //   auth_date: initDataUnsafe.auth_date ? initDataUnsafe.auth_date + '' : '',
    //   hash: initDataUnsafe.hash || '',
    //   ref_uid: initDataUnsafe.start_param || '',
    // }

    // $modals.info.show({
    //   text: WebApp.initData,
    // })

    await sendTelegramAppAuth(WebApp.initData)
  }

  return true
}
