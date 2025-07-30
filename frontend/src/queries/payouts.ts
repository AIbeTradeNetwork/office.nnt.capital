import { $request } from 'utils/axios'

export async function get(paging: Paging): Promise<Partial<Payout[]>> {
  try {
    const response = await $request({
      query: `query ($limit: Int!, $skip: Int!) {
        payouts (limit: $limit, skip: $skip) {
          uid,
          amount,
          commission,
          currency_code,
          method_code,
          account_number,
          account_name,
          created_at,
          approved_at,
          charged_at,
          cancelled_at,
          reason,
        }
      }`,
      variables: {
        ...paging,
      },
      cache: true,
    })

    return response.data.data?.payouts || []
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function create(input: PayoutReq): Promise<Partial<Payout>> {
  try {
    const response = await $request({
      query: `mutation ($input: PayoutReq!) {
        payout (input: $input) {
          uid,
        }
      }`,
      variables: {
        input,
      },
    })

    return response.data.data?.payout || []
  } catch (error) {
    return Promise.reject(error)
  }
}
