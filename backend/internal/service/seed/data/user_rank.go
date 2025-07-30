package data

import (
	"math/big"
	"server/internal/domain"
	"time"
)

func UserRanks() []domain.UserRank {
	return []domain.UserRank{
		{
			UID:      "SYSTEM1",
			Type:     domain.UserRankTypeBuy,
			UserUID:  "SYSTEM",
			MatchUID: "",
			Row:      big.NewInt(1),
			Col:      big.NewInt(1),
			BuyUID:   "",
			RankCode: "brilliant",
			StartAt:  time.Now().UTC().Add(-24 * time.Hour),
			EndAt:    time.Now().UTC().Add(20 * 365 * 24 * time.Hour),
			Priority: 800,
		},
		{
			UID:      "SYSTEM2",
			Type:     domain.UserRankTypeAuto,
			UserUID:  "SYSTEM",
			MatchUID: "",
			Row:      big.NewInt(1),
			Col:      big.NewInt(1),
			BuyUID:   "",
			RankCode: "legend",
			StartAt:  time.Now().UTC().Add(-24 * time.Hour),
			EndAt:    time.Now().UTC().Add(20 * 365 * 24 * time.Hour),
			Priority: 1600,
		},
	}
}
