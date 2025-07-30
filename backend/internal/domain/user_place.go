package domain

import (
	"math/big"
	"time"
)

type UserPlace struct {
	UserUID   string
	MatchUID  string
	Row       *big.Int
	Col       *big.Int
	CreatedAt time.Time
}

type UserPlaceWeb struct {
	UserUID   string
	MatchUID  string
	Row       *big.Int
	Col       *big.Int
	CreatedAt time.Time
	User      User
	Plan      UserPlan
	Rank      UserRank
	Activity  UserActivity
}
