import { $request } from 'utils/axios'
import { $store } from 'utils/store'

export async function get(): Promise<Partial<INotification[]>> {
  try {
    const response = await $request({
      query: `query {
        notifications {
          uid,
          texts {
            ru,
            en,
          },
          created_at,
        }
      }`,
    })

    const data = (response.data.data?.notifications || []) as INotification[]

    $store.set('notifications', data)

    return data
  } catch (error) {
    return Promise.reject(error)
  }
}
