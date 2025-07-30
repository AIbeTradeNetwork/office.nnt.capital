package domain

import "time"

type UserActivity struct {
	UserUID  string
	StartAt  time.Time
	EndAt    time.Time
	CvAmount int64
}
