declare namespace Broker {
  export interface Balance {
    currency: string
    amount: string
    plus: string
    minus: string
  }

  export interface Config {
    deposit_fee_percent: string
    withdraw_fee_percent: string
    max_deposit: string
    min_deposit: string
    min_withdraw: string
    max_withdraw: string
  }

  export interface DepositIn {
    currency: string
    amount: string
  }

  export interface Trade {
    uid: string
    side: string
    category: string
    symbol: string
    qty: string
    price_open: string
    price_close: string
    profit: string
    exchange_code: string
    exchange_uid: string
    created_at: Int
    closed_at: Int
  }

  export interface TradeIn {
    side: string
    category: string
    symbol: string
    currency: string
    qty: string
    profit: string
    price_open: string
    price_close: string
    strategy_code: string
    exchange_code: string
    exchange_uid: string
    created_at: string
    closed_at: string
  }

  export interface Transaction {
    uid: string
    currency: string
    percent: string
    amount: string
    type: string
    trade_uid: string
    created_at: Int
  }

  export interface WithdrawIn {
    currency: string
    amount: string
  }
}
