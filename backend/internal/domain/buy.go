package domain

import (
	"math/big"
	"time"
)

type BuyType string

const (
	BuyTypeUndefined   BuyType = ""
	BuyTypeClient      BuyType = "client"
	BuyTypeDistributor BuyType = "distributor"
)

type Buy struct {
	UID          string
	UserUID      string
	RefUID       string
	MatchUID     string
	Row          *big.Int
	Col          *big.Int
	Type         BuyType
	CreatedAt    time.Time
	PaidAt       time.Time
	ApprovedAt   time.Time
	ChargedAt    time.Time
	CancelledAt  time.Time
	RefundedAt   time.Time
	PlanCode     string
	ProductCode  string
	CurrencyCode string
	Amount       int64
	Cv           int64
	FlowUID      string
	PayUID       string
}

type BuyReq struct {
	PlanCode string
	UserUID  string
}

type BuyProductReq struct {
	ProductCode string
	UserUID     string
}

type BuySignalPaid struct {
	UID    string
	Amount int64
}

type BuySignalRefund struct {
	UID    string
	Amount int64
}
