import { $request } from 'utils/axios'
import { $store } from 'utils/store'

export async function get(): ResposePartial<Currency[]> {
  try {
    const response = await $request({
      query: `query {
        currencies {
          code
          precision
        }
      }`,
      cache: true,
    })

    const data = response.data.data?.currencies || []

    $store.set('currencies', data)

    return data
  } catch (error) {
    return Promise.reject(error)
  }
}
