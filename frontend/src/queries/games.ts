import { $request } from 'utils/axios'

export async function set_combo(code: string): Promise<Combo | undefined> {
  if (!code) code = ''

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `mutation ($code: String!) {
        set_combo (code: $code) {
          uid,
          name,
          currency_code,
          precision,
          amount,
          prise_code,
          limit,
          count,
          start_at,
          end_at,
        }
      }`,
      variables: {
        code,
      },
    })

    return response.data.data?.set_combo || ''
  } catch (error: any) {
    return Promise.reject(error)
  }
}

export async function user_safe(): Promise<UserSafe> {
  try {
    const response = await $request({
      query: `query {
        user_safe {
          uid,
          safe_uid,
          secret,
          created_at,
          claimed_at,
          variant {
            type,
            amount,
          }
        }
      }`,
    })

    return response.data.data?.user_safe
  } catch (error: any) {
    return Promise.reject(error)
  }
}

export async function hack_safe(code: string): Promise<UserSafe> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `mutation ($code: String!) {
        hack_safe(code:$code) {
          uid,
          safe_uid,
          secret,
          created_at,
          claimed_at,
          variant {
            type,
            amount,
          }
        }
      }`,
      variables: {
        code,
      },
    })

    return response.data.data?.hack_safe
  } catch (error: any) {
    return Promise.reject(error)
  }
}
