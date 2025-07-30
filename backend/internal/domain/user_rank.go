package domain

import (
	"math/big"
	"time"
)

type UserRankType string

const (
	UserRankTypeUndefined UserRankType = ""
	UserRankTypeBuy       UserRankType = "buy"
	UserRankTypeAuto      UserRankType = "auto"
)

type UserRank struct {
	UID      string
	Type     UserRankType
	UserUID  string
	MatchUID string
	Row      *big.Int
	Col      *big.Int
	BuyUID   string
	RankCode string
	StartAt  time.Time
	EndAt    time.Time
	Priority int64
}
