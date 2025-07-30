package domain

type UserLevel struct {
	Code        string
	Balance     int64
	InviteLimit int64
	ClaimAmount int64
	LevelBonus  []int64
	Tier        int64
}

var UserLevels = []*UserLevel{
	{
		Code:        "novice",
		Balance:     0,
		InviteLimit: 0,
		ClaimAmount: 300000000,
		LevelBonus:  []int64{500},
		Tier:        1,
	},
	{
		Code:        "farmer",
		Balance:     50000000000,
		InviteLimit: 0,
		ClaimAmount: 400000000,
		LevelBonus:  []int64{700},
		Tier:        1,
	},
	{
		Code:        "supervisor",
		Balance:     150000000000,
		InviteLimit: 0,
		ClaimAmount: 500000000,
		LevelBonus:  []int64{700, 200},
		Tier:        1,
	},
	{
		Code:        "expert",
		Balance:     500000000000,
		InviteLimit: 0,
		ClaimAmount: 600000000,
		LevelBonus:  []int64{800, 500},
		Tier:        1,
	},
	{
		Code:        "master",
		Balance:     1000000000000,
		InviteLimit: 0,
		ClaimAmount: 700000000,
		LevelBonus:  []int64{800, 600},
		Tier:        1,
	},
	{
		Code:        "champion",
		Balance:     2500000000000,
		InviteLimit: 0,
		ClaimAmount: 800000000,
		LevelBonus:  []int64{900, 600},
		Tier:        1,
	},
	{
		Code:        "veteran",
		Balance:     5000000000000,
		InviteLimit: 0,
		ClaimAmount: 900000000,
		LevelBonus:  []int64{900, 800},
		Tier:        1,
	},
	{
		Code:        "legionary",
		Balance:     8000000000000,
		InviteLimit: 0,
		ClaimAmount: 1000000000,
		LevelBonus:  []int64{1000, 1000},
		Tier:        1,
	},
	{
		Code:        "grandmaster",
		Balance:     15000000000000,
		InviteLimit: 0,
		ClaimAmount: 1400000000,
		LevelBonus:  []int64{1200, 1000},
		Tier:        1,
	},
	{
		Code:        "archon",
		Balance:     30000000000000,
		InviteLimit: 0,
		ClaimAmount: 1700000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        1,
	},
	{
		Code:        "archon-I",
		Balance:     60000000000000,
		InviteLimit: 0,
		ClaimAmount: 12000000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        2,
	},
	{
		Code:        "archon-II",
		Balance:     150000000000000,
		InviteLimit: 0,
		ClaimAmount: 15000000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        2,
	},
	{
		Code:        "archon-III",
		Balance:     300000000000000,
		InviteLimit: 0,
		ClaimAmount: 18000000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        2,
	},
	{
		Code:        "archon-IV",
		Balance:     600000000000000,
		InviteLimit: 0,
		ClaimAmount: 21000000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        2,
	},
	{
		Code:        "archon-V",
		Balance:     1500000000000000,
		InviteLimit: 0,
		ClaimAmount: 24000000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        2,
	},
	{
		Code:        "archon-VI",
		Balance:     3000000000000000,
		InviteLimit: 0,
		ClaimAmount: 27000000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        2,
	},
	{
		Code:        "archon-VII",
		Balance:     6000000000000000,
		InviteLimit: 0,
		ClaimAmount: 30000000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        2,
	},
	{
		Code:        "archon-VIII",
		Balance:     12000000000000000,
		InviteLimit: 0,
		ClaimAmount: 42000000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        2,
	},
	{
		Code:        "archon-IX",
		Balance:     20000000000000000,
		InviteLimit: 0,
		ClaimAmount: 51000000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        2,
	},
	{
		Code:        "archon-X",
		Balance:     40000000000000000,
		InviteLimit: 0,
		ClaimAmount: 200000000000,
		LevelBonus:  []int64{1500, 1000},
		Tier:        2,
	},
}

func GetLevelByCode(code string) *UserLevel {
	for _, l := range UserLevels {
		if l.Code == code {
			return l
		}
	}
	return nil
}
