declare type StrategyType = keyof typeof import('./enums').EStrategyType
declare type StrategyCategory = keyof typeof import('./enums').EStrategyCategory | ''

declare interface Bot {
  uid: string
  key_uid: string
  exchange_code: string
  strategy_code: string
  symbol_code: string
  is_active: boolean
  trade_percent: Int
  trade_type: BotTradeType
  trade_limit: Int
  trade_reinvest: boolean
  started_at: Int
  created_at: Int
  error: string
  error_code: string
  leverage: Int
  profit_all: number
  profit_month: number
  profit_month_prev: number
}

declare interface BotAddReq {
  key_uid: string
  strategy_code: string
  symbol_code: string
  trade_percent: Int
  is_active: boolean
  trade_type: BotTradeType
  trade_limit: Int
  trade_reinvest: boolean
  leverage?: Int
}

declare interface BotEditReq {
  key_uid: string
  trade_percent: Int
  is_active: boolean
  trade_type: BotTradeType
  trade_limit: Int
  trade_reinvest: boolean
  leverage?: Int
}

declare interface Exchange {
  name: string
  code: string
  is_active: boolean
  link: string
}

declare interface Key {
  name: string
  uid: string
  exchange_code: string
  key: string
}

declare interface KeyAddReq {
  exchange_code: string
  name: string
  key: string
  secret: string
}

declare interface KeyEditReq {
  name: string
  key: string
  secret: string
}

declare type RiskLevel = 'low' | 'medium' | 'high'

declare interface Strategy {
  name: string
  code: string
  symbol_code: string
  type: StrategyType
  start_at: Int
  exchange_code: string
  exchange_codes: string[]
  risk_level: RiskLevel
  pos_profit: Int
  category: StrategyCategory
  min_deposit: Int
  max_deposit: Int
  min_leverage: Int
  max_leverage: Int
  fix_leverage: Int
  symbol_codes: string[]
  share_profit: Int
  is_new: boolean
  is_free: boolean
  is_trade_percent: boolean
  is_trade_limit: boolean
  is_reinvest: boolean
}

declare interface Trade {
  uid: string
  bot_uid: string
  signal_uid: string
  exchange_code: string
  strategy_code: string
  symbol_code: string
  price_open: Int
  price_close: Int
  commission: Int
  take_profit: Int
  stop_loss: Int
  created_at: Int
  closed_at: Int
  error: string
  error_code: string
  profit: number
  side: string
  qty: number
}

declare interface Cfg {
  payout_amount_min: Int
  payout_fee_min: Int
  payout_fee_percent: Int
  distributor_price: Int
  default_currency_code: string
  default_ref_uid: string
  unlim_invite: boolean
  ton_wallet: string
}

declare type BotTradeType =
  | import('./enums').EBotTradeType.limit
  | import('./enums').EBotTradeType.percent
