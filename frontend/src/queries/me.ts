import { ERoles } from 'types/enums'
import { $request } from 'utils/axios'
import { $config } from 'utils/configuration'
import { $me } from 'utils/me'

const fragmentPlans = `
  plans {
    code,
    start_at,
    end_at,
    priority,
  },
`

const fragmentRanks = `
  ranks {
    code,
    start_at,
    end_at,
    priority,
  },
`

const fragmentMain = `
  role,
  uid,
  ref_uid,
  email,
  nickname,
  created_at,
  balance,
  team_count,
  plan {
    code,
    start_at,
    end_at,
    priority,
  },
  place {
    row,
    col,
    created_at,
    match_uid,
  },
  lim_ref_uid,
  level {
    code,
    balance,
    invite_limit,
  },
  unlim_invite,
  ton_wallet,
  is_premium,
  premium_until,
  premium_invites,
  premium_multiplier,
  is_autofarm,
  autofarm_until,
  active_boost {
    code,
    category,
    start_at,
    end_at,
    multiplier,
  },
  products {
    code,
    category,
    start_at,
    end_at,
    multiplier
  }
`

const fragmentDistributor = `
  left,
  right,
  rank {
    code,
    start_at,
    end_at,
    priority,
  },
  place {
    created_at,
  },
  activity {
    start_at,
    end_at,
    cv,
  },
  config {
    team_type,
  },
`

export async function getRole(): Promise<Pick<User, 'role'>> {
  try {
    const response = await $request({
      query: `query {
        me {
          role
        }
      }`,
    })

    if (response.data.data.me) {
      $me.set(response.data.data.me)
      $config.role = response.data.data.me.role
    }

    return response.data.data?.me || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function get(): ResposePartial<User> {
  try {
    const response = await $request({
      query: `query {
        me {
          ${fragmentMain}
          ${$config.role === ERoles.distributor ? fragmentDistributor : ''}
        }
      }`,
      cache: true,
    })

    const data = response.data.data?.me
    if (data) $me.set(data)
    if (response.data.data.me?.role) $config.role = response.data.data.me.role

    return data || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function updateConfig(): Promise<Pick<User, 'config'>> {
  try {
    const response = await $request({
      query: `query {
        me {
          config {
            team_type
          },
        }
      }`,
    })

    const data = response.data.data?.me
    if (data) $me.set(data)

    return data || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function updateBalance(): Promise<Pick<User, 'balance'>> {
  try {
    const response = await $request({
      query: `query {
        me {
          balance
        }
      }`,
    })

    const data = response.data.data?.me
    if (data) $me.set(data)

    return data || {}
  } catch (error) {
    return Promise.reject(error)
  }
}

export async function setRef(uid: string): Promise<string> {
  if (!uid) return

  try {
    const response = await $request({
      query: `mutation ($uid: String!) {
        set_ref (uid: $uid)
      }`,
      variables: {
        uid,
      },
    })

    const data = response.data.data?.set_ref || ''

    if (data) $me.set({ ref_uid: data })

    return data
  } catch (error) {
    return Promise.reject(error)
  }
}
