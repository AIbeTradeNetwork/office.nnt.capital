package domain

import "time"

type TonTransaction struct {
	UID          string
	TxUID        string
	TxLT         uint64
	Addr         string
	Amount       int64
	Fee          int64
	Sender       string
	CurrencyCode string
	Precision    uint8
	UserUID      string
	Comment      string
	CreatedAt    time.Time
}
