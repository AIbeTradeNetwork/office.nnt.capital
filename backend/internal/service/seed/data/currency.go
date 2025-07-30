package data

import "server/internal/domain"

func Currencies() []domain.Currency {
	return []domain.Currency{
		{
			Code:      "cv",
			Precision: 2,
		},
		{
			Code:      "usd",
			Precision: 2,
		},
		{
			Code:      "usdt",
			Precision: 6,
		},
		{
			Code:      "UDEX",
			Precision: 9,
		},
	}
}
