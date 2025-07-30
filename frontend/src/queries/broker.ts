import { BrokerTransactionType, ECurrencies } from 'types/enums'
import { $request } from 'utils/axios'
import { $store } from 'utils/store'

export async function config(): Promise<Broker.Config> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BROKER,
      },
      query: `query {
        config {
          deposit_fee_percent,
          withdraw_fee_percent,
          max_deposit,
          min_deposit,
          min_withdraw,
          max_withdraw,
        }
      }`,
      cache: true,
    })

    return response.data.data?.config
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function balance(currency = ECurrencies.usdHidden): Promise<Broker.Balance> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BROKER,
      },
      cache: true,
      query: `query ($currency: String!) {
        balance (currency: $currency) {
          currency,
          amount,
          plus,
          minus,
        }
      }`,
      variables: { currency },
    })

    const res = response.data.data?.balance

    $store.set('brokerBalance', res)

    return
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function transactions(o: {
  currency?: string
  type?: BrokerTransactionType
  limit?: Int
  skip?: Int
}): Promise<Broker.Transaction[]> {
  if (!o.limit) o.limit = 10
  if (!o.skip) o.skip = 0

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BROKER,
      },
      query: `query ($currency: String, $type: TransactionType, $limit: Int, $skip: Int) {
        transactions (currency: $currency, type: $type, limit: $limit, skip: $skip) {
          uid,
          currency,
          percent,
          amount,
          type,
          trade_uid,
          created_at,
        }
      }`,
      cache: true,
      variables: { ...o },
    })

    return response.data.data?.transactions || []
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function deposit({ amount, currency }: Broker.DepositIn): Promise<Broker.Balance> {
  if (!amount) return Promise.reject(new Error('Amount is required'))
  if (!currency) return Promise.reject(new Error('Currency is required'))

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BROKER,
      },
      query: `mutation ($input: DepositIn!) {
        deposit (input: $input) {
          currency,
          amount,
          plus,
          minus,
        }
      }`,
      cache: true,
      variables: {
        input: {
          amount,
          currency,
        },
      },
    })

    return response.data.data?.deposit
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function withdraw({ amount, currency }: Broker.WithdrawIn): Promise<Broker.Balance> {
  if (!amount) return Promise.reject(new Error('Amount is required'))
  if (!currency) return Promise.reject(new Error('Currency is required'))

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL_BROKER,
      },
      query: `mutation ($input: WithdrawIn!) {
        withdraw (input: $input) {
          currency,
          amount,
          plus,
          minus,
        }
      }`,
      cache: true,
      variables: {
        input: {
          amount,
          currency,
        },
      },
    })

    return response.data.data?.withdraw
  } catch (error) {
    return Promise.reject(error)
  }
}
