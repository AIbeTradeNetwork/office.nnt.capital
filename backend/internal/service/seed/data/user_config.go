package data

import (
	"server/internal/domain"
)

func UserConfigs() []domain.UserConfig {
	return []domain.UserConfig{
		{
			UserUID:      "SYSTEM",
			TeamType:     "",
			LastTeamType: domain.UserTeamTypeRight,
			AllowSwitch:  false,
		},
	}
}
