package data

import (
	"server/internal/domain"
	"time"
)

func UserAuths() []domain.UserAuth {
	return []domain.UserAuth{
		{
			UserUID:   "SYSTEM",
			Type:      domain.AuthTypePassword,
			Token:     "$2a$10$9R4ka5ZRlgCxpu8hzAGNTOf/kWGPw/OE20O5flWqtbnPX..9hbPhq",
			CreatedAt: time.Time{},
		},
	}
}
