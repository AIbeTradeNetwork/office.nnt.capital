package data

import (
	"server/internal/domain"
	"time"
)

func Tasks() []domain.Task {
	return []domain.Task{
		{
			Code:         "sub-1002156085803",
			IsActive:     true,
			StartAt:      time.Time{},
			EndAt:        time.Time{},
			CurrencyCode: "UDEX",
			Precision:    9,
			Amount:       1000000000,
			Priority:     100,
		},
	}
}
