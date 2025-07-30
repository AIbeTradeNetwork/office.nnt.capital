package data

import (
	"server/internal/domain"
	"time"
)

func Claims() []domain.Claim {
	return []domain.Claim{
		{
			Code:         "UDEX",
			MaxPeriod:    8 * time.Hour,
			MinPeriod:    2 * time.Hour,
			Amount:       300000000,
			CurrencyCode: "UDEX",
			Precision:    9,
		},
	}
}
