package data

import (
	"server/internal/domain"
	"time"
)

func Users() []domain.User {
	return []domain.User{
		{
			UID:       "SYSTEM",
			RefUID:    "SYSTEM",
			Nickname:  "root-system",
			Email:     "root@system.com",
			CreatedAt: time.Now().UTC(),
		},
	}
}
