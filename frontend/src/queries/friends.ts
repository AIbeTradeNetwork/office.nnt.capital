import { $request } from 'utils/axios'

export async function count(user_uid: string): Promise<Int> {
  if (!user_uid) return 0

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `query ($user_uid: String!) {
        friends_count (user_uid: $user_uid)
      }`,
      variables: {
        user_uid,
      },
      cache: false,
    })

    return response.data.data?.friends_count || 0
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function list(user_uid: string, limit: Int, skip: Int): Promise<FriendUser[]> {
  if (!limit) limit = 5
  if (!skip || skip < 0) skip = 0

  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `query ($user_uid: String!, $limit: Int!, $skip: Int!) {
        friends (user_uid: $user_uid, limit: $limit, skip: $skip) {
          uid,
          ref_uid,
          email,
          nickname,
          created_at,
          team_count
        }
      }`,
      variables: {
        user_uid,
        limit,
        skip,
      },
      cache: false,
    })

    return response.data.data?.friends || []
  } catch (error) {
    return Promise.reject(error)
  }
}
