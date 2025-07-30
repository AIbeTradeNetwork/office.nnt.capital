package domain

import "time"

type Payout struct {
	UID           string
	UserUID       string
	Amount        int64
	Fee           int64
	CurrencyCode  string
	MethodCode    string
	AccountNumber string
	AccountName   string
	CreatedAt     time.Time
	ApprovedAt    time.Time
	ChargedAt     time.Time
	CancelledAt   time.Time
	Reason        string
}
