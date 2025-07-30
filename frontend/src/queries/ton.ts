import { $authentication } from 'utils/authentication'
import { $request } from 'utils/axios'

export async function getTonPayload(): Promise<string> {
  try {
    const response = await $request({
      query: `query {
        get_ton_payload
      }`,
    })

    return response.data.data?.get_ton_payload || ''
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function setTonWallet(input: string): Promise<string> {
  try {
    const response = await $request({
      query: `mutation ($input: String!) {
        set_ton_wallet (input: $input)
      }`,
      variables: {
        input,
      },
    })

    return response.data.data?.set_ton_wallet || ''
  } catch (error) {
    return Promise.reject(error)
  }
}
