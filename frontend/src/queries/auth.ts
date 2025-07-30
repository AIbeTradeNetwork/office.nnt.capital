import { $authentication } from 'utils/authentication'
import { $request } from 'utils/axios'

export async function login(input: AuthLogin): Promise<AuthRes> {
  try {
    const response = await $request({
      query: `mutation ($input: AuthLogin!) {
        login (input: $input) {
          auth_token
          refresh_token
        }
      }`,
      variables: {
        input,
      },
    })

    const data = response.data.data?.login

    if (data) $authentication.logIn(data)

    return data || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function telegram(input: TelegramAuth): Promise<AuthRes> {
  try {
    const response = await $request({
      query: `mutation ($input: TelegramAuth!) {
        telegram (input: $input) {
          auth_token
          refresh_token
        }
      }`,
      variables: {
        input,
      },
    })

    const data = response.data.data?.telegram

    if (data) $authentication.logIn(data)

    return data || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function telegramApp(input: string): Promise<AuthRes> {
  try {
    const response = await $request({
      query: `mutation ($input: String!) {
        telegram_app (input: $input) {
          auth_token
          refresh_token
        }
      }`,
      variables: {
        input,
      },
    })

    const data = response.data.data?.telegram_app

    if (data) $authentication.logIn(data)

    return data || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function refresh(): Promise<AuthRes> {
  try {
    const response = await $request({
      query: `mutation ($input: AuthRefresh!) {
        refresh (input: $input) {
          auth_token,
          refresh_token
        }
      }`,
      variables: {
        input: {
          refresh_token: $authentication.getRefreshToken(),
        },
      },
    })

    const refresh = response.data.data.refresh
    if (refresh?.refresh_token) {
      $authentication.setRefreshToken(refresh.refresh_token)
    }

    if (refresh?.auth_token) {
      $authentication.setToken(refresh.auth_token)
    }

    return refresh || {}
  } catch (error) {
    $authentication.logOut()
    return Promise.reject(error)
  }
}

export async function register(input: AuthRegister): Promise<AuthRes> {
  try {
    const response = await $request({
      query: `mutation ($input: AuthRegister!) {
        register (input: $input) {
          auth_token,
          refresh_token
        }
      }`,
      variables: {
        input,
      },
    })

    const data = response.data.data?.register
    if (data) $authentication.logIn(data)

    return data || {}
  } catch (error) {
    return Promise.reject(error)
  }
}
