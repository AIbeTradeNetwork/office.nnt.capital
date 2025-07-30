package domain

import "time"

type Dist struct {
	UID          string
	UserUID      string
	RefUID       string
	Amount       int64
	CurrencyCode string
	CreatedAt    time.Time
	PaidAt       time.Time
	PayUID       string
}
