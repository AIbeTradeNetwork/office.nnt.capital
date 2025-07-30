package data

import (
	"server/internal/domain"
	"time"
)

func UserActivities() []domain.UserActivity {
	return []domain.UserActivity{
		{
			UserUID:  "SYSTEM",
			StartAt:  time.Now().UTC().Add(-24 * time.Hour),
			EndAt:    time.Now().UTC().Add(20 * 365 * 24 * time.Hour),
			CvAmount: 0,
		},
	}
}
