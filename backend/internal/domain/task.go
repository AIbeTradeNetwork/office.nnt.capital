package domain

import "time"

type TaskReq struct {
	Code   string `json:"code"`
	UserId int64  `json:"userId"`
	Action bool   `json:"action"`
}

type ClaimReq struct {
	PartnerCode string `json:"partnerCode"`
	Code        string `json:"code"`
	UserUID     string `json:"userUid"`
	Amount      string `json:"amount"`
}

type Task struct {
	Code         string
	Texts        map[string]string
	IsActive     bool
	StartAt      time.Time
	EndAt        time.Time
	CurrencyCode string
	Precision    uint8
	Amount       int64
	Priority     int64
	Link         string
	Completed    bool
	Locale       string
	Count        int64
	Limit        int64
	IsApprove    bool
	RefUID       string
	RefCount     int64
}
