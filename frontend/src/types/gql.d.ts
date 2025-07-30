declare interface GQLRespose<T> {
  data: { [key: string]: T }
  error: [
    {
      message: string
      locations: [{ line: number; column: number }]
      path: [string, number]
    },
  ]
}

declare type ResposePartial<T> = Promise<Partial<T>>

declare interface Locale {
  ru: string
  en: string
}

declare interface AuthLogin {
  login: string
  password: string
}

declare interface TelegramAuth {
  id: string
  first_name: string
  last_name: string
  username: string
  photo_url: string
  auth_date: string
  hash: string
  ref_uid: string
}

declare interface AuthRefresh {
  refresh_token: string
}

declare type AuthRegister = {
  email: string
  nickname: string
  password: string
  repassword: string
  ref_uid: string
  locale?: string
}

declare interface AuthRes {
  auth_token: string
  refresh_token: string
}

declare interface UserPlan {
  code: string
  start_at: Int
  end_at: Int
  priority: Int
}

declare interface UserRank {
  code: string
  start_at: Int
  end_at: Int
  priority: Int
}

declare interface UserPlace {
  row: string
  col: string
  created_at: Int
  match_uid: string
}

declare type UserRole = 'client' | 'distributor'

declare interface UserLevel {
  code: string
  balance: Int
  invite_limit: Int
}

declare interface User {
  uid: string
  ref_uid: string
  _ref_uid: string
  email: string
  nickname: string
  created_at: Int
  plans: UserPlan[]
  plan: UserPlan
  role: UserRole
  left: Int
  right: Int
  ranks: UserRank[]
  rank: UserRank
  place: UserPlace
  activity: UserActivity
  config: UserConfig
  balance: Int
  lim_ref_uid: string
  level: UserLevel
  team_count: number
  unlim_invite?: boolean
  ton_wallet: string
  products: UserProduct[]
  is_premium: boolean
  premium_until: Int
  premium_invites: Int
  premium_multiplier: Int
  is_autofarm: boolean
  autofarm_until: Int
  active_boost: UserProduct
}

declare interface UserProduct {
  code: string
  category: string
  start_at: Int
  end_at: Int
  priority: Int
  multiplier: Int
}

declare interface UserConfig {
  team_type: TeamType
}

declare type TeamType = import('./common').ETeamType

declare interface Plan {
  code: import('./enums').EPlans
  period: Int
  price: Int
  cv: Int
  currency_code: import('./enums').ECurrencies
  rank_code: string
  rank_period: Int
  priority: Int
  retail_price: Int
  bot_count: Int
  max_deposit: Int
}

declare interface RankCondition {
  rank_code: string
  team_type: TeamType
  is_ref: boolean
  count: Int
}

declare interface RankGetBonus {
  amount: Int
  months: Int
}

declare interface RankRefBonus {
  plan_code: string
  percent: Int
}

declare interface Rank {
  code: keyof typeof import('./enums').ERangs
  min_cv: Int
  team_condition: RankCondition[]
  priority: Int
  bin_bonus: Int
  bin_bonus_week_limit: Int
  ref_bonus: RankRefBonus[]
  match_bonus: Int[]
  first_bonus: RankGetBonus
  approve_bonus: RankGetBonus
}

declare interface BuyReq {
  code: string
}

declare interface BuyProductReq {
  code: string
}

declare interface BuyRes {
  url: string
  uid: string
}

declare interface Currency {
  code: string
  precision: Int
}

declare interface UserActivity {
  start_at: Int
  end_at: Int
  cv: Int
}

declare interface TeamUser extends User {
  team_count: Int
  created_at: Int
}

declare type TransactionType =
  | undefined
  | 'ref'
  | 'bin'
  | 'match'
  | 'firstRank'
  | 'approveRank'
  | 'fastStart'

declare interface Transaction {
  user_uid: string
  from_uid: string
  percent: Int
  level: Int
  type: TransactionType
  rank_code: string
  amount: Int
  pos_amount: Int
  full_amount: Int
  buy_uid: string
  payout_uid: string
  created_at: Int
  charged_at: Int
  msg_codes: string[]
}

declare type BuyType = undefined | 'client' | 'distributor'

declare interface Buy {
  uid: string
  user_uid: string
  ref_uid: string
  match_uid: string
  row: string
  col: string
  type: BuyType
  created_at: Int
  paid_at: Int
  approved_at: Int
  charged_at: Int
  refunded_at: Int
  plan_code: string
  currency_code: import('./enums').ECurrencies
  amount: Int
  cv: Int
}

declare interface DistRes {
  url: string
  uid: string
}

declare interface Payout {
  uid: string
  amount: Int
  commission: Int
  currency_code: string
  method_code: string
  account_number: string
  account_name: string
  created_at: Int
  approved_at: Int
  charged_at: Int
  cancelled_at: Int
  reason: string
}

declare interface PayoutReq {
  currency_code: string!
  method_code: string!
  account_number: string!
  account_name: string!
}

declare interface INotification {
  uid: string
  to_user_uid: string[]
  texts: {
    ru: string
    en: string
  }
  created_at: Int
}

declare interface NotifyReq {
  to_user_uid: string[]
  texts: {
    ru: string
    en: string
  }
}

declare interface UserProduct {
  code: string
  category: string
  start_at: Int
  end_at: Int
  priority: Int
}

declare interface Product {
  code: string
  category: string
  period: Int
  price: Int
  retail_price: Int
  cv: Int
  currency_code: string
  precision: Int
  priority: Int
  multiplier: Int
  limit: Int
  count: Int
}

declare interface Task {
  code: string
  texts: Locale
  currency_code: string
  precision: Int
  amount: Int
  link: string
  completed: boolean
  is_approve: boolean
}

declare interface Combo {
  uid: string
  name: string
  currency_code: string
  precision: Int
  amount: Int
  prise_code: string
  limit: Int
  count: Int
  start_at: Int
  end_at: Int
}

declare interface FriendUser {
  uid: string
  ref_uid: string
  email: string
  nickname: string
  created_at: Int
  team_count: Int
}

declare type UserSafeVariantType = import('./enums').ECurrencies

declare interface UserSafeVariant {
  type: UserSafeVariantType
  amount: Int
}

declare interface UserSafe {
  uid: string
  safe_uid: string
  secret: string
  created_at: Int
  claimed_at: Int
  variant: UserSafeVariant
}
