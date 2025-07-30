package domain

import "time"

type UserClaimType string

const (
	UserClaimTypeInvite     UserClaimType = "invite"
	UserClaimTypeRef        UserClaimType = "ref"
	UserClaimTypeRefBuy     UserClaimType = "refBuy"
	UserClaimTypeRefBonus   UserClaimType = "refBonus"
	UserClaimTypeClaim      UserClaimType = "claim"
	UserClaimTypePartner    UserClaimType = "partner"
	UserClaimTypePlan       UserClaimType = "plan"
	UserClaimTypeBonus      UserClaimType = "bonus"
	UserClaimTypeBuy        UserClaimType = "buy"
	UserClaimTypeTask       UserClaimType = "task"
	UserClaimTypePrise      UserClaimType = "prise"
	UserClaimTypeDepBonus   UserClaimType = "depBonus"
	UserClaimTypeTradeInRef UserClaimType = "tradeInRef"
	UserClaimTypeCombo      UserClaimType = "combo"
	UserClaimTypeSafe       UserClaimType = "safe"
	UserClaimTypeBoost      UserClaimType = "boost"
)

type UserClaim struct {
	UID          string
	ClaimCode    string
	UserUID      string
	RefUID       string
	Level        int
	CreatedAt    time.Time
	ClaimedAt    time.Time
	Amount       int64
	CurrencyCode string
	Precision    uint8
	Type         UserClaimType
	TaskCode     string
	ComboUID     string
	PartnerCode  string
}
