package domain

import "time"

type TransactionType string

const (
	TransactionTypeUndefined        TransactionType = ""
	TransactionTypeRefBonus         TransactionType = "ref"
	TransactionTypeRefBuy           TransactionType = "refBuy"
	TransactionTypeBinBonus         TransactionType = "bin"
	TransactionTypeMatchBonus       TransactionType = "match"
	TransactionTypeFirstRankBonus   TransactionType = "firstRank"
	TransactionTypeApproveRankBonus TransactionType = "approveRank"
	TransactionTypeFastStartBonus   TransactionType = "fastStart"
	TransactionTypePayout           TransactionType = "payout"
	TransactionTypeDeposit          TransactionType = "deposit"
	TransactionTypeBuy              TransactionType = "buy"
	TransactionTypeDist             TransactionType = "dist"
	TransactionTypeTask             TransactionType = "task"
	TransactionTypePrise            TransactionType = "prise"
	TransactionTypeTradeIn          TransactionType = "tradeIn"
	TransactionTypeTradeOut         TransactionType = "tradeOut"
	TransactionTypeTradeInFee       TransactionType = "tradeInFee"
	TransactionTypeTradeOutFee      TransactionType = "tradeOutFee"
	TransactionTypeTradeInRef       TransactionType = "tradeInRef"
	TransactionTypeCombo            TransactionType = "combo"
	TransactionTypeSafe             TransactionType = "safe"
)

type TransactionMsgCode string

const (
	TransactionMsgCodeUndefined       TransactionMsgCode = ""
	TransactionMsgCodeWeekLimit       TransactionMsgCode = "weekLimit"
	TransactionMsgCodeMonthLimit      TransactionMsgCode = "monthLimit"
	TransactionMsgCodeNotEnoughRank   TransactionMsgCode = "notEnoughRank"
	TransactionMsgCodeActivityExpired TransactionMsgCode = "activityExpired"
	TransactionMsgCodeAmountNotEnough TransactionMsgCode = "amountNotEnough"
)

type Transaction struct {
	UID         string
	UserUID     string
	FromUID     string
	Percent     int64
	Level       int
	Type        TransactionType
	RankCode    string
	TaskCode    string
	ComboUID    string
	Amount      int64
	PosAmount   int64
	FullAmount  int64
	Coefficient int64
	BuyUID      string
	PayoutUID   string
	DepositUID  string
	CreatedAt   time.Time
	ChargedAt   time.Time
	MsgCodes    []TransactionMsgCode
}
