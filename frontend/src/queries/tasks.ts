import { $request } from 'utils/axios'

export async function get(): ResposePartial<Task[]> {
  try {
    const response = await $request({
      query: `query {
        tasks {
          code,
          texts {
            ru,
            en,
          },
          currency_code,
          precision,
          amount,
          link,
          completed,
          is_approve,
        }
      }`,
      cache: true,
    })

    return response.data.data?.tasks || []
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function approve(code: string): ResposePartial<Pick<Task, 'code'>> {
  try {
    const response = await $request({
      query: `mutation ($code: String!) {
        approve_task (code: $code) {
          code,
        }
      }`,
      variables: {
        code,
      },
    })

    return response.data.data?.approve_task
  } catch (error) {
    return Promise.reject(error)
  }
}

// export async function set(input: BuyProductReq): Promise<Partial<BuyRes>> {
//   try {
//     const response = await $request({
//       query: `mutation ($input: BuyProductReq!) {
//         buy_product (input: $input) {
//           uid
//         }
//       }`,
//       variables: {
//         input: {
//           code: input.code,
//         },
//       },
//     })

//     if (response.data.data?.buy_product?.uid) {
//       $me.update()
//     }

//     return response.data.data?.buy_product || []
//   } catch (error) {
//     return Promise.reject(error)
//   }
// }
