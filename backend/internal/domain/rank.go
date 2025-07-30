package domain

type RankTeamCondition struct {
	RankCode string
	TeamType UserTeamType
	IsRef    bool
	Count    uint64
}

type RankGetBonus struct {
	Amount int64
	Months int
}

type Rank struct {
	Code          string
	MinCv         int64
	TeamCondition []RankTeamCondition
	Priority      int64
	BinBonus      int64
	BinBonusLimit int64
	RefBonus      map[string]int64
	MatchBonus    []int64
	FirstBonus    RankGetBonus
	ApproveBonus  RankGetBonus
}
