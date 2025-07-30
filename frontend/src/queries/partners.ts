import { $request } from 'utils/axios'

export async function count(user_uid: string): Promise<Int> {
  if (!user_uid) return 0

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `query ($user_uid: String!) {
        partners_count (user_uid: $user_uid)
      }`,
      variables: {
        user_uid,
      },
      cache: false,
    })

    return response.data.data?.partners_count || 0
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function addPartner(user_uid: string, partner_uid: string): Promise<boolean> {
  if (!user_uid || !partner_uid) return false

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `mutation ($user_uid: String!, $partner_uid: String!) {
        add_partner (user_uid: $user_uid, partner_uid: $partner_uid)
      }`,
      variables: {
        user_uid,
        partner_uid,
      },
      cache: false,
    })

    return response.data.data?.add_partner || false
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function removePartner(user_uid: string, partner_uid: string): Promise<boolean> {
  if (!user_uid || !partner_uid) return false

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `mutation ($user_uid: String!, $partner_uid: String!) {
        remove_partner (user_uid: $user_uid, partner_uid: $partner_uid)
      }`,
      variables: {
        user_uid,
        partner_uid,
      },
      cache: false,
    })

    return response.data.data?.remove_partner || false
  } catch (error) {
    return Promise.reject(error)
  }
} 