import { $request } from 'utils/axios'

export async function get(): ResposePartial<Rank[]> {
  try {
    const response = await $request({
      query: `query {
        ranks {
          code,
          min_cv,
          team_condition {
            rank_code,
            team_type,
            is_ref,
            count,
          }
          priority,
          bin_bonus,
          bin_bonus_week_limit,
          ref_bonus {
            plan_code,
            percent,
          }
          match_bonus,
          first_bonus {
            amount,
            months,
          }
          approve_bonus {
            amount,
            months,
          }
        }
      }`,
      cache: true,
    })

    return response.data.data?.ranks || []
  } catch (error) {
    return Promise.reject(error)
  }
}
