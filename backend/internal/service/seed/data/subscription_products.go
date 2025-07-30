package data

import (
	"server/internal/domain"
	"time"
)

func SubscriptionProducts() []domain.Product {
	return []domain.Product{
		{
			Code:         domain.SubscriptionResearcher,
			Name:         "Researcher",
			Description:  "Базовый абонемент для всех пользователей",
			Category:     "subscription",
			Period:       0, // Бессрочный
			CurrencyCode: "usd",
			Precision:    2,
			RetailPrice:  0,
			Price:        0,
			Cv:           0,
			Priority:     1,
			Multiplier:   1,
			IsActive:     true,
			Limit:        0,
			Count:        0,
			Sort:         1,
		},
		{
			Code:         domain.SubscriptionStart,
			Name:         "Start",
			Description:  "Стартовый абонемент для активных пользователей",
			Category:     "subscription",
			Period:       12 * 30 * 24 * time.Hour, // 12 месяцев
			CurrencyCode: "usd",
			Precision:    2,
			RetailPrice:  25000, // $250.00 в центах
			Price:        25000,
			Cv:           0,
			Priority:     2,
			Multiplier:   3,
			IsActive:     true,
			Limit:        0,
			Count:        0,
			Sort:         2,
		},
		{
			Code:         domain.SubscriptionAdvanced,
			Name:         "Advanced",
			Description:  "Продвинутый абонемент для опытных пользователей",
			Category:     "subscription",
			Period:       18 * 30 * 24 * time.Hour, // 18 месяцев
			CurrencyCode: "usd",
			Precision:    2,
			RetailPrice:  50000, // $500.00 в центах
			Price:        50000,
			Cv:           0,
			Priority:     3,
			Multiplier:   8,
			IsActive:     true,
			Limit:        0,
			Count:        0,
			Sort:         3,
		},
		{
			Code:         domain.SubscriptionProfessional,
			Name:         "Professional",
			Description:  "Профессиональный абонемент для экспертов",
			Category:     "subscription",
			Period:       24 * 30 * 24 * time.Hour, // 24 месяца
			CurrencyCode: "usd",
			Precision:    2,
			RetailPrice:  100000, // $1,000.00 в центах
			Price:        100000,
			Cv:           0,
			Priority:     4,
			Multiplier:   20,
			IsActive:     true,
			Limit:        0,
			Count:        0,
			Sort:         4,
		},
		{
			Code:         domain.SubscriptionAmbassador,
			Name:         "Ambassador",
			Description:  "Амбассадорский абонемент для лидеров",
			Category:     "subscription",
			Period:       30 * 30 * 24 * time.Hour, // 30 месяцев
			CurrencyCode: "usd",
			Precision:    2,
			RetailPrice:  150000, // $1,500.00 в центах
			Price:        150000,
			Cv:           0,
			Priority:     5,
			Multiplier:   50,
			IsActive:     true,
			Limit:        0,
			Count:        0,
			Sort:         5,
		},
		{
			Code:         domain.SubscriptionLeader,
			Name:         "Leader",
			Description:  "Лидерский абонемент для топ-менеджеров",
			Category:     "subscription",
			Period:       36 * 30 * 24 * time.Hour, // 36 месяцев
			CurrencyCode: "usd",
			Precision:    2,
			RetailPrice:  250000, // $2,500.00 в центах
			Price:        250000,
			Cv:           0,
			Priority:     6,
			Multiplier:   100,
			IsActive:     true,
			Limit:        0,
			Count:        0,
			Sort:         6,
		},
		{
			Code:         domain.SubscriptionVIP,
			Name:         "VIP",
			Description:  "VIP абонемент для элитных пользователей",
			Category:     "subscription",
			Period:       48 * 30 * 24 * time.Hour, // 48 месяцев
			CurrencyCode: "usd",
			Precision:    2,
			RetailPrice:  750000, // $7,500.00 в центах
			Price:        750000,
			Cv:           0,
			Priority:     7,
			Multiplier:   300,
			IsActive:     true,
			Limit:        0,
			Count:        0,
			Sort:         7,
		},
	}
}
