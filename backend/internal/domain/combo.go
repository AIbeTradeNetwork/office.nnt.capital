package domain

import "time"

type Combo struct {
	UID          string
	Code         string
	Name         string
	CurrencyCode string
	Precision    uint8
	Amount       int64
	PriseCode    string
	Limit        int
	Count        int
	StartAt      time.Time
	EndAt        time.Time
	IsActive     bool
}
