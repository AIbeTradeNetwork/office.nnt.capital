package data

import (
	"server/internal/domain"
	"time"
)

func Products() []domain.Product {
	return []domain.Product{
		{
			Code:         "premium_month_usd",
			Category:     "premium",
			Period:       30 * 24 * time.Hour,
			CurrencyCode: "usd",
			Precision:    2,
			RetailPrice:  359,
			Price:        359,
			Cv:           359,
			Priority:     100,
			IsActive:     true,
		},
		{
			Code:         "premium_month_udex",
			Category:     "premium",
			Period:       30 * 24 * time.Hour,
			CurrencyCode: "UDEX",
			Precision:    9,
			RetailPrice:  35000000000,
			Price:        35000000000,
			Cv:           70000000,
			Priority:     100,
			IsActive:     true,
		},
		{
			Code:         "premium_year_usd",
			Category:     "premium",
			Period:       365 * 24 * time.Hour,
			CurrencyCode: "usd",
			Precision:    2,
			RetailPrice:  4308,
			Price:        2999,
			Cv:           2999,
			Priority:     200,
			IsActive:     true,
		},
		{
			Code:         "premium_year_udex",
			Category:     "premium",
			Period:       365 * 24 * time.Hour,
			CurrencyCode: "UDEX",
			Precision:    9,
			RetailPrice:  420000000000,
			Price:        300000000000,
			Cv:           600000000,
			Priority:     200,
			IsActive:     true,
		},
	}
}
