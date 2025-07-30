package domain

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	UID          string
	BuyUID       string
	CreatedAt    time.Time
	CurrencyCode string
	Amount       int64
}

type CurrencyCode string

const (
	CurrencyUSDT CurrencyCode = "USDT"
)

var AllCurrencyCode = []CurrencyCode{
	CurrencyUSDT,
}

func (e CurrencyCode) IsValid() bool {
	return slices.Contains(AllCurrencyCode, e)
}

type CustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Customer struct {
	Account string `json:"account"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

type OrderCreateReq struct {
	Amount             float64       `json:"amount"`
	Callbacks          Callbacks     `json:"callbacks"`
	Comment            string        `json:"comment"`
	CurrencyCode       CurrencyCode  `json:"currency"`
	CustomFields       []CustomField `json:"customFields"`
	Customer           Customer      `json:"customer"`
	ExpirationDateTime time.Time     `json:"expirationDateTime"`
	PaymentSystemID    uuid.UUID     `json:"paymentSystemID"`
}

type Callbacks struct {
	FailureURL string `json:"failureURL"`
	ReturnURL  string `json:"returnURL"`
	SuccessURL string `json:"successURL"`
}

type CreateOrderResp struct {
	BillId uuid.UUID `json:"billId"`
	URL    string    `json:"url"`
}

type OrderInfo struct {
	Amount             float64       `json:"amount"`
	AmountPaid         float64       `json:"amountPaid"`
	BillID             uuid.UUID     `json:"billId"`
	ChangedDateTime    time.Time     `json:"changedDateTime"`
	Comment            string        `json:"comment"`
	CreationDateTime   time.Time     `json:"creationDateTime"`
	CurrencyCode       CurrencyCode  `json:"currency"`
	CustomFields       []CustomField `json:"customFields"`
	Customer           Customer      `json:"customer"`
	ExpirationDateTime time.Time     `json:"expirationDateTime"`
	MerchantID         uuid.UUID     `json:"merchantId"`
	PayURL             string        `json:"payUrl"`
	PaymentSystemID    uuid.UUID     `json:"paymentSystemID"`
	Status             OrderStatus   `json:"status"`
}

type OrderStatus string

const (
	OrderStatusCanceled OrderStatus = "CANCELED"
	OrderStatusExpired  OrderStatus = "EXPIRED"
	OrderStatusPaid     OrderStatus = "PAID"
	OrderStatusRejected OrderStatus = "REJECTED"
	OrderStatusWaiting  OrderStatus = "WAITING"
)

var AllOrderStatus = []OrderStatus{
	OrderStatusCanceled,
	OrderStatusExpired,
	OrderStatusPaid,
	OrderStatusRejected,
	OrderStatusWaiting,
}

func (e OrderStatus) IsValid() bool {
	return slices.Contains(AllOrderStatus, e)
}

type OrderIn struct {
	MerchantID         string         `json:"merchantId"`
	BillID             string         `json:"billId"`
	PaymentSystemID    string         `json:"paymentSystemID"`
	Amount             float64        `json:"amount"`
	CurrencyCode       CurrencyCode   `json:"currency"`
	Status             OrderStatus    `json:"status"`
	Comment            string         `json:"comment"`
	ChangedDateTime    string         `json:"changedDateTime"`
	CreationDateTime   string         `json:"creationDateTime"`
	ExpirationDateTime string         `json:"expirationDateTime"`
	PayURL             string         `json:"payUrl"`
	Customer           *Customer      `json:"customer"`
	CustomFields       []*CustomField `json:"customFields"`
}
