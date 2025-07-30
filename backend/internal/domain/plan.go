package domain

import "time"

type Plan struct {
	Code         string
	Period       time.Duration
	RetailPrice  int64
	Price        int64
	Cv           int64
	CurrencyCode string
	PaymentWait  time.Duration
	ChargeWait   time.Duration
	RankCode     string
	RankPeriod   time.Duration
	Priority     int64
	BotCount     int64
	MaxDeposit   int64
	CoinBonus    int64
}
