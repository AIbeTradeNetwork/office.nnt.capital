package domain

import "time"

type Claim struct {
	Code         string
	MaxPeriod    time.Duration
	MinPeriod    time.Duration
	Amount       int64
	CurrencyCode string
	Precision    uint8
}
