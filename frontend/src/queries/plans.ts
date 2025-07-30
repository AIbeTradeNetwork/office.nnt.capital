import { $request } from 'utils/axios'

export async function get(): ResposePartial<Plan[]> {
  try {
    const response = await $request({
      query: `query {
        plans {
          code
          period
          price
          cv
          currency_code
          rank_code
          rank_period
          priority,
          retail_price,
          bot_count,
          max_deposit,
        }
      }`,
      cache: true,
    })

    return response.data.data?.plans || []
  } catch (error) {
    return Promise.reject(error)
  }
}
