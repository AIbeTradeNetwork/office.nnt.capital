package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type Deposit struct {
	UID           string
	UserUID       string
	Amount        int64
	Fee           int64
	CurrencyCode  string
	Precision     uint8
	Rate          decimal.Decimal
	MethodCode    string
	AccountNumber string
	AccountName   string
	CreatedAt     time.Time
	ApprovedAt    time.Time
	ChargedAt     time.Time
	CancelledAt   time.Time
	Comment       string
	TxUID         string
	TxLT          uint64
}
