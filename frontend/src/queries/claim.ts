import { $request } from 'utils/axios'

const coinName = 'UDEX'

export async function getClaim(code?: string): Promise<Claim> {
  try {
    const response = await $request({
      query: `query ($code: String!) {
        get_claim(code: $code) {
          code,
          min_period,
          max_period,
          amount,
          currency_code,
          precision,
        }
      }`,
      variables: {
        code: code || coinName,
      },
    })

    return response.data.data?.get_claim || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function getBalance(code?: string): Promise<ClaimBalance> {
  try {
    const response = await $request({
      query: `query ($code: String!) {
        get_claim_balance(code: $code) {
          claim_code,
          claimed_at,
          balance,
          precision,
        }
      }`,
      variables: {
        code: code || coinName,
      },
    })

    return response.data.data?.get_claim_balance || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

async function setClaim(code?: string): Promise<ClaimBalance> {
  try {
    const response = await $request({
      query: `mutation ($code: String!) {
        set_claim (code: $code) {
          claim_code,
          claimed_at,
          balance,
          precision,
        }
      }`,
      variables: {
        code: code || coinName,
      },
    })

    return response.data?.data?.set_claim || response.data || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

const coin = {
  getBalance,
  getClaim,
  setClaim,
}

export { coin }
