import { $request } from 'utils/axios'
import { $store } from 'utils/store'

export async function get(paging?: Paging): Promise<Partial<Key[]>> {
  if (!paging) {
    paging = {
      limit: 0,
      skip: 0,
    }
  }

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BOTS,
      },
      query: `query ($limit: Int!, $skip: Int!) {
        keys (limit: $limit, skip: $skip) {
          name,
          uid,
          exchange_code,
          key,
        }
      }`,
      variables: {
        ...paging,
      },
      cache: true,
    })

    const data = response.data.data?.keys || []

    $store.set('keys', data)

    return data
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function add(input: {
  exchange_code: KeyAddReq['exchange_code']
  name: KeyAddReq['name']
  key: KeyAddReq['key']
  secret: KeyAddReq['secret']
}): Promise<Partial<Key>> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BOTS,
      },
      query: `mutation ($input: KeyAddReq!) {
        add_key (input: $input) {
          uid
        }
      }`,
      variables: {
        input,
      },
    })

    $store.updateKeys()

    return response.data.data?.add_key
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function edit(o: {
  uid: string
  name: KeyEditReq['name']
  key: KeyEditReq['key']
  secret: KeyEditReq['secret']
}): Promise<Partial<Key>> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BOTS,
      },
      query: `mutation ($uid: String!, $input: KeyEditReq!) {
        edit_key (uid: $uid, input: $input) {
          uid
        }
      }`,
      variables: {
        uid: o.uid,
        input: {
          name: o.name,
          key: o.key,
          secret: o.secret,
        },
      },
    })

    $store.updateKeys()

    return response.data.data?.edit_key
  } catch (error) {
    return Promise.reject(error)
  }
}
