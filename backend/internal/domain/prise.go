package domain

type PriseType string

const (
	PriseTypeUDEX             PriseType = "UDEX"
	PriseTypePremiumMonth     PriseType = "premium_month"
	PriseTypePremiumYear      PriseType = "premium_year"
	PriseTypePremiumMonthUdex PriseType = "premium_month_udex"
	PriseTypePremiumYearUdex  PriseType = "premium_year_udex"
	PriseTypeAutofarmMonth    PriseType = "autofarm_month"
	PriseTypeAutofarmYear     PriseType = "autofarm_year"
)

type PriseReq struct {
	Type      PriseType
	Amount    int64
	ToUserUID []string
}
