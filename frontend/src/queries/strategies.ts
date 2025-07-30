import { $request } from 'utils/axios'
import { $store } from 'utils/store'

export async function get(o: {
  type: StrategyType
  paging?: Paging
  cache?: boolean
}): Promise<Partial<Strategy[]>> {
  if (!o.paging) {
    o.paging = {
      limit: 0,
      skip: 0,
    }
  }

  if (typeof o.cache === 'undefined') {
    o.cache = true
  }

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BOTS,
      },
      query: `query ($type: StrategyType!, $limit: Int!, $skip: Int!) {
        strategies (type: $type, limit: $limit, skip: $skip) {
          name,
          code,
          symbol_code,
          type,
          start_at,
          exchange_code,
          exchange_codes,
          risk_level,
          pos_profit,
          category,
          min_deposit,
          max_deposit,
          min_leverage,
          max_leverage,
          fix_leverage,
          symbol_codes,
          share_profit,
          is_new,
          is_free,
          is_trade_percent,
          is_trade_limit,
          is_reinvest,
        }
      }`,
      variables: {
        type: o.type,
        ...o.paging,
      },
      cache: o.cache,
    })

    const data = (response.data.data?.strategies || []) as Strategy[]

    const strategies = $store.get('strategies')

    data.forEach((item) => {
      const findedIndex = strategies.findIndex((strategy) => strategy.code === item.code)

      if (findedIndex >= 0) {
        strategies[findedIndex] = item
        return
      }

      strategies.push(item)
    })

    $store.set('strategies', strategies)

    return data
  } catch (error) {
    return Promise.reject(error)
  }
}
