import { $request } from 'utils/axios'
import { $formatToInt } from 'utils/formats'
import { $store } from 'utils/store'

export async function get(paging?: Paging): Promise<Partial<Bot[]>> {
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
        bots (limit: $limit, skip: $skip) {
          uid,
          key_uid,
          exchange_code,
          strategy_code,
          trade_percent,
          is_active,
          symbol_code,
          trade_type,
          trade_limit,
          trade_reinvest,
          error,
          error_code,
          leverage,
          profit_all,
          profit_month,
          profit_month_prev,
        }
      }`,
      variables: {
        ...paging,
      },
    })

    const data = response.data.data?.bots || []

    $store.set('bots', data)

    return data
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function traders(bot_uid: string, paging?: Paging): Promise<Partial<Trade[]>> {
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
      query: `query ($bot_uid: String!, $limit: Int!, $skip: Int!) {
        bot_trades (bot_uid: $bot_uid, limit: $limit, skip: $skip) {
          uid,
          signal_uid,
          symbol_code,
          price_open,
          price_close,
          commission,
          take_profit,
          stop_loss,
          created_at,
          closed_at,
          error,
          error_code,
          profit,
          side,
          qty,
        }
      }`,
      variables: {
        bot_uid,
        ...paging,
      },
    })

    const data = response.data.data?.bot_trades || []

    return data
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function add(input: BotAddReq): Promise<Partial<Bot>> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BOTS,
      },
      query: `mutation ($input: BotAddReq!) {
        add_bot (input: $input) {
          uid
        }
      }`,
      variables: {
        input: <BotAddReq>{
          is_active: input.is_active,
          key_uid: input.key_uid,
          strategy_code: input.strategy_code,
          trade_percent: input.trade_percent,
          symbol_code: input.symbol_code,
          trade_type: input.trade_type,
          trade_limit: $formatToInt(input.trade_limit),
          trade_reinvest: input.trade_reinvest,
          leverage: input.leverage > -1 ? input.leverage : undefined,
        },
      },
    })

    return response.data.data?.add_bot
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function edit(uid: string, input: BotEditReq): Promise<Partial<Bot>> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BOTS,
      },
      query: `mutation ($uid: String!, $input: BotEditReq!) {
        edit_bot (uid: $uid, input: $input) {
          uid
        }
      }`,
      variables: {
        uid,
        input: <BotEditReq>{
          key_uid: input.key_uid,
          is_active: input.is_active,
          trade_percent: input.trade_percent,
          trade_type: input.trade_type,
          trade_limit: $formatToInt(input.trade_limit),
          trade_reinvest: input.trade_reinvest,
          leverage: input.leverage > -1 ? input.leverage : undefined,
        },
      },
    })

    return response.data.data?.edit_bot
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function del(uid: string): Promise<Partial<Bot>> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BOTS,
      },
      query: `mutation ($uid: String!) {
        del_bot (uid: $uid) {
          uid
        }
      }`,
      variables: {
        uid,
      },
    })

    return response.data.data?.del_bot
  } catch (error) {
    return Promise.reject(error)
  }
}
