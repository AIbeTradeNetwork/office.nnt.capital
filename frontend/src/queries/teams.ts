import { $request } from 'utils/axios'

export async function getRef(input: { user_uid: string }): Promise<Partial<TeamUser[]>> {
  try {
    const response = await $request({
      query: `query ($user_uid: String) {
        team_ref (user_uid: $user_uid) {
          uid,
          email,
          nickname,
          role,
          plan {
            code,
            end_at
          },
          rank {
            code,
            end_at
          },
          team_count,
          created_at,
        }
      }`,
      variables: {
        ...input,
      },
      cache: true,
    })

    return response.data.data?.team_ref || []
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function getBin(input: { user_uid: string; rows: Int }): Promise<Partial<TeamUser[]>> {
  try {
    const response = await $request({
      query: `query ($user_uid: String, $rows: Int!) {
        team_bin (user_uid: $user_uid, rows: $rows) {
          uid,
          ref_uid,
          email,
          nickname,
          role,
          plan {
            code,
          },
          rank {
            code,
          },
          place {
            row,
            col,
            match_uid
          }
        }
      }`,
      variables: {
        ...input,
      },
      cache: true,
    })

    return response.data.data?.team_bin || []
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function getMatch(input: { user_uid: string }): Promise<TeamUser[]> {
  try {
    const response = await $request({
      query: `query ($user_uid: String) {
        team_match (user_uid: $user_uid) {
          uid,
          email,
          nickname,
          role,
          plan {
            code,
            end_at
          },
          rank {
            code,
            end_at
          },
          team_count,
        }
      }`,
      variables: {
        ...input,
      },
      cache: true,
    })

    return response.data.data?.team_match || []
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function setSwitchKnee(inputs: { type: TeamType }): Promise<UserConfig> {
  try {
    const response = await $request({
      query: `mutation ($type: TeamType!) {
        set_team_type (type: $type) {
          team_type
        }
      }`,
      variables: {
        ...inputs,
      },
    })

    return response.data.data?.set_team_type || []
  } catch (error) {
    return Promise.reject(error)
  }
}
