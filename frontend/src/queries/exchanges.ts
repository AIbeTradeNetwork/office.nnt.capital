import { $request } from 'utils/axios'
import { $store } from 'utils/store'

export async function get(paging?: Paging): Promise<Partial<Exchange[]>> {
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
        exchanges (limit: $limit, skip: $skip) {
            name,
            code,
            is_active,
            link
        }
      }`,
      variables: {
        ...paging,
      },
      cache: true,
    })

    const data = response.data.data?.exchanges || []

    $store.set('exchanges', data)

    return data
  } catch (error) {
    return Promise.reject(error)
  }
}
