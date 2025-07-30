package domain

import "time"

type UserPlan struct {
	UID      string
	UserUID  string
	BuyUID   string
	PlanCode string
	StartAt  time.Time
	EndAt    time.Time
	Priority int64
}
