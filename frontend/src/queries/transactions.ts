import { $request } from 'utils/axios'

export async function transactions(paging: Paging): Promise<Partial<Transaction[]>> {
  try {
    const response = await $request({
      query: `query ($limit: Int!, $skip: Int!) {
        transactions (limit: $limit, skip: $skip) {
          from_uid,
          type,
          amount,
          buy_uid,
          charged_at,
        }
      }`,
      variables: {
        ...paging,
      },
      cache: true,
    })

    return response.data.data?.transactions || []
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function buys(paging: Paging): Promise<Partial<Buy[]>> {
  try {
    const response = await $request({
      query: `query ($limit: Int!, $skip: Int!) {
        buys (limit: $limit, skip: $skip) {
          uid,
          user_uid,
          type,
          paid_at,
          amount,
          cv,
          plan_code,
          currency_code
        }
      }`,
      variables: {
        ...paging,
      },
      cache: true,
    })

    return response.data.data?.buys || []
  } catch (error) {
    return Promise.reject(error)
  }
}

// export async function buyProduct(input: BuyProductReq): Promise<BuyRes[]> {
//   try {
//     const response = await $request({
//       query: `mutation ($input: BuyProductReq!) {
//           buy_product (input: $input) {
//             url,
//             uid
//           }
//         }`,
//       variables: {
//         input,
//       },
//     })

//     return response.data.data?.buy || {}
//   } catch (error) {
//     return Promise.reject(error)
//   }
// }

export async function setBuyDistributor(): Promise<DistRes> {
  try {
    const response = await $request({
      query: `mutation () {
          dist{
            uid
          }
        }`,
    })

    return response.data.data?.dist || {}
  } catch (error) {
    return Promise.reject(error)
  }
}
