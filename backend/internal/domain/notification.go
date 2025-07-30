package domain

import "time"

type Notification struct {
	UID       string
	Texts     map[string]string
	CreatedAt time.Time
	SentAt    time.Time
	ToUserUID []string
	FlowUID   string
}
