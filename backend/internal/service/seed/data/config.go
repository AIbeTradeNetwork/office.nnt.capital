package data

import (
	"server/internal/domain"
	"time"
)

func Config() domain.Config {
	return domain.Config{
		PayoutAmountMin:  5000,
		PayoutFeeMin:     100,
		PayoutFeePercent: 100,
		FastStartBonuses: []domain.FastStartBonus{
			{
				Amount: 20000,
				MinCv:  195000,
			},
			{
				Amount: 40000,
				MinCv:  260000,
			},
			{
				Amount: 60000,
				MinCv:  390000,
			},
		},
		FastStartDuration:   30 * 24 * time.Hour,
		ActivityCvAmount:    5000,
		ActivityMaxDuration: 60 * 24 * time.Hour,
		ClientRefBonus: map[string]int64{
			"start":        300,
			"advanced":     400,
			"professional": 600,
			"silver":       600,
			"gold":         600,
			"platinum":     600,
			"brilliant":    600,
		},
		DistributorRefBonus: map[string]int64{
			"start":        500,
			"advanced":     600,
			"professional": 800,
			"silver":       800,
			"gold":         800,
			"platinum":     800,
			"brilliant":    800,
		},
		DistributorPrice:          5000,
		DefaultCurrencyCode:       "usd",
		AllowTeamTypeSwitchAmount: 9500,
		CoinCode:                  "UDEX",
		CoinRefBonus:              1000000000,
		CoinRefLinePercent:        []int64{2500, 2000, 1000, 500, 500, 500, 500, 500, 500, 500},
		InitializedAt:             time.Now().UTC(),
	}
}
