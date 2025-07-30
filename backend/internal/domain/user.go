package domain

import "time"

type User struct {
	UID         string
	LimRefUID   string
	RefUID      string
	Nickname    string
	Email       string    `bson:"email"`
	FirstName   string
	LastName    string
	PhotoUrl    string
	Locale      string
	TgID        int64
	TgUsername  string
	CreatedAt   time.Time
	UnlimInvite bool
	TonWallet   string
	TeamCount   int64
}

type UserTg struct {
	UID    string
	TgUID  string
	Locale string
}
