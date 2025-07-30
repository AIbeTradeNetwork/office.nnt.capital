package domain

type UserBalance struct {
	UserUID      string
	CurrencyCode string
	Precision    uint8
	Amount       int64
}
