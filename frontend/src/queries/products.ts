import { $request } from 'utils/axios'
import { $me } from 'utils/me'
import { $store } from 'utils/store'

export async function get({ category = '' }): ResposePartial<Product[]> {
  try {
    const response = await $request({
      query: `query ($category: String){
        products (category:$category) {
          code,
          category,
          period,
          price,
          retail_price,
          cv,
          currency_code,
          precision,
          priority,
          multiplier,
          limit,
          count
        }
      }`,
      cache: true,
      variables: {
        category,
      },
    })

    const products = response.data.data?.products || []

    $store.set('products', products)

    return response.data.data?.products || []
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function buyProduct(input: BuyProductReq): Promise<Partial<BuyRes>> {
  try {
    const response = await $request({
      query: `mutation ($input: BuyProductReq!) {
        buy_product (input: $input) {
          uid
        }
      }`,
      variables: {
        input: {
          code: input.code,
        },
      },
    })

    if (response.data.data?.buy_product?.uid) {
      $me.update()
    }

    return response.data.data?.buy_product || []
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function buyTariff(input: BuyReq): Promise<BuyRes> {
  try {
    const response = await $request({
      query: `mutation ($input: BuyReq!) {
          buy (input: $input) {
            uid
          }
        }`,
      variables: {
        input,
      },
    })

    if (response.data.data?.buy?.uid) {
      $me.update()
    }

    return response.data.data?.buy || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function buyDistributor(): Promise<DistRes> {
  try {
    const response = await $request({
      query: `mutation {
        dist {
          uid
        }
      }`,
    })

    if (response.data.data?.dist?.uid) {
      $me.update()
    }

    return response.data.data?.dist || {}
  } catch (error) {
    return Promise.reject(error)
  }
}
