import { ECurrencies } from 'types/enums'
import { $request } from 'utils/axios'
import { $config } from 'utils/configuration'
import { $store } from 'utils/store'

export async function get(): Promise<Partial<Cfg>> {
  try {
    const response = await $request({
      query: `query {
        cfg {
          payout_amount_min,
          payout_fee_min,
          payout_fee_percent,
          distributor_price,
          default_currency_code,
          default_ref_uid,
          unlim_invite,
          ton_wallet,
        }
      }`,
    })

    const data = response.data.data?.cfg

    if (data) {
      $config.tonAddress = data.ton_wallet
      $store.set('cfg', data)
    }

    return data
  } catch (error) {
    return Promise.reject(error)
  }
}
