package domain

import "time"

type FastStartBonus struct {
	Amount int64
	MinCv  int64
}

type Config struct {
	PayoutAmountMin           int64
	PayoutFeeMin              int64
	PayoutFeePercent          int64
	FastStartBonuses          []FastStartBonus
	FastStartDuration         time.Duration
	ActivityCvAmount          int64
	ActivityMaxDuration       time.Duration
	ClientRefBonus            map[string]int64
	DistributorRefBonus       map[string]int64
	DistributorPrice          int64
	DefaultCurrencyCode       string
	AllowTeamTypeSwitchAmount int64
	DefaultRefUid             string    // Default referrer for new users without invite link
	InitializedAt             time.Time // Date of initialization
	CoinCode                  string    // Code of claiming coin
	CoinRefBonus              int64     // Amount of coins to charge to inviter for invited referral
	CoinToRefBonus            int64     // Amount of coins to charge to invited referral
	CoinRefLinePercent        []int64   // Percent for charging coins for claiming to referrals by lines
	UnlimInvite               bool      // Unlimited invite for all
	TonLastLT                 uint64    // Last processed TON transaction LT
	TonWallet                 string    // Treasury TON wallet
	RefSafeUID                string    // Default safe for referrals
	Tier1SafeUID              string    // Safe for buyers of 6th digit
	Tier2SafeUID              string    // Safe for buyers of 7th digit
	Tier1CoinSafeUID          string    // Safe for buyers of 6th digit
	Tier2CoinSafeUID          string    // Safe for buyers of 7th digit
}
