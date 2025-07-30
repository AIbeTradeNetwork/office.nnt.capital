package data

import (
	"server/internal/domain"
	"time"
)

func UserPlans() []domain.UserPlan {
	return []domain.UserPlan{
		{
			UID:      "SYSTEM",
			UserUID:  "SYSTEM",
			BuyUID:   "",
			PlanCode: "brilliant",
			StartAt:  time.Now().UTC().Add(-24 * time.Hour),
			EndAt:    time.Now().UTC().Add(20 * 365 * 24 * time.Hour),
			Priority: 800,
		},
	}
}
