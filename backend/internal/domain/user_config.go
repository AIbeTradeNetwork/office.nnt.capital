package domain

type UserConfig struct {
	UserUID      string
	TeamType     UserTeamType
	LastTeamType UserTeamType
	AllowSwitch  bool
}
