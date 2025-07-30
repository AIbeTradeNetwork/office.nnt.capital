package data

import (
	"math/big"
	"server/internal/domain"
	"time"
)

func UserPlaces() []domain.UserPlace {
	return []domain.UserPlace{
		{
			UserUID:   "SYSTEM",
			MatchUID:  "",
			Row:       big.NewInt(1),
			Col:       big.NewInt(1),
			CreatedAt: time.Now().UTC(),
		},
	}
}
