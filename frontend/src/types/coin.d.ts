declare interface Claim {
  code: string!
  min_period: Int!
  max_period: Int!
  amount: Int!
  currency_code: string!
  precision: Int!
}

declare interface ClaimBalance {
  claim_code: string!
  claimed_at: Int!
  balance: Int!
  precision: Int!
}

declare interface EABTRanksCountItem {
  code: import('./enums').EABTRanks
  min: number
  max: number
  img: string
}
