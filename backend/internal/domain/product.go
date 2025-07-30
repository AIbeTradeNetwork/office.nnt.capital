package domain

import "time"

// Subscription types
const (
	SubscriptionResearcher   = "subscription_researcher"
	SubscriptionStart        = "subscription_start"
	SubscriptionAdvanced     = "subscription_advanced"
	SubscriptionProfessional = "subscription_professional"
	SubscriptionAmbassador   = "subscription_ambassador"
	SubscriptionLeader       = "subscription_leader"
	SubscriptionVIP          = "subscription_vip"
)

// SubscriptionInfo содержит информацию об абонементе
type SubscriptionInfo struct {
	Code           string
	Name           string
	PartnersLimit  int64
	CommissionRate float64 // Комиссионные выплаты в процентах
	RoyaltyRate    float64 // Роялти от доходов клиентов в процентах
	Price          int64   // Цена в USD
	RenewalPrice   int64   // Цена продления в USD
	DurationMonths int     // Длительность в месяцах
	IsBasic        bool    // Базовый абонемент (Researcher)
}

// GetSubscriptionInfo возвращает информацию об абонементе по коду
func GetSubscriptionInfo(code string) *SubscriptionInfo {
	subscriptions := map[string]*SubscriptionInfo{
		SubscriptionResearcher: {
			Code:           SubscriptionResearcher,
			Name:           "Researcher",
			PartnersLimit:  1,
			CommissionRate: 15.0,
			RoyaltyRate:    15.0,
			Price:          0,
			RenewalPrice:   0,
			DurationMonths: 0,
			IsBasic:        true,
		},
		SubscriptionStart: {
			Code:           SubscriptionStart,
			Name:           "Start",
			PartnersLimit:  3,
			CommissionRate: 16.0,
			RoyaltyRate:    14.0,
			Price:          25000, // $250.00 в центах
			RenewalPrice:   12500, // $125.00 в центах
			DurationMonths: 12,
			IsBasic:        false,
		},
		SubscriptionAdvanced: {
			Code:           SubscriptionAdvanced,
			Name:           "Advanced",
			PartnersLimit:  8,
			CommissionRate: 18.0,
			RoyaltyRate:    12.0,
			Price:          50000, // $500.00 в центах
			RenewalPrice:   25000, // $250.00 в центах
			DurationMonths: 18,
			IsBasic:        false,
		},
		SubscriptionProfessional: {
			Code:           SubscriptionProfessional,
			Name:           "Professional",
			PartnersLimit:  20,
			CommissionRate: 20.0,
			RoyaltyRate:    10.0,
			Price:          100000, // $1,000.00 в центах
			RenewalPrice:   50000,  // $500.00 в центах
			DurationMonths: 24,
			IsBasic:        false,
		},
		SubscriptionAmbassador: {
			Code:           SubscriptionAmbassador,
			Name:           "Ambassador",
			PartnersLimit:  50,
			CommissionRate: 22.0,
			RoyaltyRate:    8.0,
			Price:          150000, // $1,500.00 в центах
			RenewalPrice:   75000,  // $750.00 в центах
			DurationMonths: 30,
			IsBasic:        false,
		},
		SubscriptionLeader: {
			Code:           SubscriptionLeader,
			Name:           "Leader",
			PartnersLimit:  100,
			CommissionRate: 24.0,
			RoyaltyRate:    6.0,
			Price:          250000, // $2,500.00 в центах
			RenewalPrice:   125000, // $1,250.00 в центах
			DurationMonths: 36,
			IsBasic:        false,
		},
		SubscriptionVIP: {
			Code:           SubscriptionVIP,
			Name:           "VIP",
			PartnersLimit:  300,
			CommissionRate: 26.0,
			RoyaltyRate:    4.0,
			Price:          750000, // $7,500.00 в центах
			RenewalPrice:   375000, // $3,750.00 в центах
			DurationMonths: 48,
			IsBasic:        false,
		},
	}

	return subscriptions[code]
}

type Product struct {
	Code         string
	Name         string
	Description  string
	Category     string
	Period       time.Duration
	CurrencyCode string
	Precision    uint8
	RetailPrice  int64
	Price        int64
	Cv           int64
	Priority     int64
	Multiplier   int64
	IsActive     bool
	Limit        int64
	Count        int64
	Sort         int64
}
