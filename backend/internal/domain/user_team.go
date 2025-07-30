package domain

type UserTeamType string

const (
	UserTeamTypeUndefined UserTeamType = ""
	UserTeamTypeLeft      UserTeamType = "left"
	UserTeamTypeRight     UserTeamType = "right"
)

type TeamUser struct {
	User      *User
	Place     *UserPlace
	Ranks     []*UserRank
	Plans     []*UserPlan
	Activity  *UserActivity
	TeamCount int64
}
