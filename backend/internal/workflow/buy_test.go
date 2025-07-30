package workflow_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"math/big"
	rand2 "math/rand"
	"server/internal/service/auth"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"

	"server/internal/domain"
	"server/internal/repository"
	"server/internal/repository/mongodb"
	"server/internal/service/buy"
	"server/internal/service/seed"
	"server/internal/service/team"
	"server/internal/workflow"
)

var (
	one    = big.NewInt(1)
	two    = big.NewInt(2)
	maxRow = big.NewInt(15)
)

type BuySuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	ctx context.Context
	db  *mongodb.Repo

	seed *seed.Service
	buy  *buy.Service
	team *team.Service
	auth *auth.Service

	flow *workflow.BuyWorkflow

	env *testsuite.TestWorkflowEnvironment
}

func (s *BuySuite) SetupTest() {
	s.ctx = context.Background()
	s.env = s.NewTestWorkflowEnvironment()

	dbRepo, _ := repository.NewDbRepo(s.ctx)
	s.db = dbRepo

	s.seed = seed.NewSeedService(s.db)

	s.buy = buy.NewBuyService(s.db, nil, nil, s.auth)
	s.auth = auth.NewAuthService(s.db, nil)
	s.team = team.NewTeamService(s.db, nil, s.auth)

	s.flow = workflow.NewBuyWorkflow(s.buy, s.team)
	s.registerTestWorkflow()
}

//func TestSuite(t *testing.T) {
//	suite.Run(t, new(BuySuite))
//}

func (s *BuySuite) Test_BuyFlow() {
	err := s.db.DropTest(s.ctx)
	s.NoError(err)

	// Ensure indexes
	err = s.db.EnsureIndexes(s.ctx)
	s.NoError(err)

	err = s.seed.Seed(s.ctx)
	s.NoError(err)
	//err = s.seedTestData()
	//s.NoError(err)

	refUid := "SYSTEM"

	for i := 1; i <= 1000; i++ {

		dur := time.Duration(i) * 5 * time.Minute

		userUid := domain.GenUID(8)
		userNickname := domain.GenNickname()
		err = s.db.UserCreate(s.ctx, &domain.User{
			UID:       userUid,
			RefUID:    refUid,
			Nickname:  userNickname,
			Email:     fmt.Sprintf("%s@system.com", userNickname),
			CreatedAt: time.Now().Add(dur).UTC(),
		})
		s.NoError(err)

		if rand2.Intn(10) >= 5 {
			refUid = userUid
		}

		planCodes := map[int]string{
			0: "start",
			1: "advanced",
			2: "professional",
			3: "silver",
			4: "gold",
			5: "platinum",
			6: "brilliant",
		}
		randPlanInd := rand2.Intn(6)
		randPlan := planCodes[randPlanInd]

		plan, err := s.db.PlanGetByCode(s.ctx, randPlan)
		s.NoError(err)

		maxRowSub := new(big.Int)
		maxRowSub.Sub(maxRow, one)
		randRow, _ := rand.Int(rand.Reader, maxRowSub)
		randRow.Add(randRow, two)
		maxCol := new(big.Int).Exp(two, new(big.Int).Set(randRow).Sub(randRow, one), nil)
		randCol, _ := rand.Int(rand.Reader, maxCol)
		randCol.Add(randCol, one)

		if rand2.Intn(10) >= 5 {
			placeEx, _ := s.db.UserPlaceGetByUserUID(s.ctx, userUid)
			if placeEx == nil {
				distBuy, err := s.team.DistBuy(s.ctx, userUid)
				s.NoError(err)

				orderIn := &domain.OrderIn{
					BillID: distBuy.PayUID,
				}

				distBuyPaid, err := s.team.DistPaid(s.ctx, orderIn)
				s.NoError(err)
				s.NotEmpty(distBuyPaid.PaidAt)
			}
		}

		buyUid := domain.GenUID(12)
		payUuid, _ := uuid.NewUUID()
		newBuy := &domain.Buy{
			UID:          buyUid, //fmt.Sprintf("buy-000%d", i),
			UserUID:      userUid,
			PlanCode:     plan.Code,
			CurrencyCode: plan.CurrencyCode,
			Amount:       plan.Price,
			Cv:           plan.Cv,
			CreatedAt:    time.Now().Add(dur).UTC(),
			PayUID:       payUuid.String(),
		}
		paidInfo := &domain.BuySignalPaid{
			UID:    buyUid, //fmt.Sprintf("buy-000%d", i),
			Amount: plan.Price,
		}

		s.env = s.NewTestWorkflowEnvironment()
		s.registerTestWorkflow()

		s.env.RegisterDelayedCallback(func() {
			s.env.SignalWorkflow("paid", paidInfo)
		}, time.Millisecond)

		s.env.ExecuteWorkflow(s.flow.BuyFlow, newBuy)

		s.True(s.env.IsWorkflowCompleted())
		s.NoError(s.env.GetWorkflowError())
	}

}

func (s *BuySuite) seedTestData() error {
	var users []domain.User
	var userPlaces []domain.UserPlace
	var userRanks []domain.UserRank
	var userActivities []domain.UserActivity
	var userConfigs []domain.UserConfig
	i := big.NewInt(2)
	j := big.NewInt(1)
	for i.Set(two); i.Cmp(maxRow) <= 0; i.Add(i, one) {
		ji := new(big.Int).Set(i)
		ji.Sub(ji, one)
		jEnd := new(big.Int).Exp(two, ji, nil)
		for j.Set(one); j.Cmp(jEnd) <= 0; j.Add(j, one) {
			refUid := "SYSTEM"
			if i.Cmp(two) > 0 {
				row := new(big.Int).Set(i)
				col := new(big.Int).Set(j)
				row.Sub(row, one)
				var r *big.Int
				col, r = new(big.Int).DivMod(col, two, new(big.Int))
				if len(r.Bits()) > 0 {
					col.Add(col, one)
				}
				refUid = fmt.Sprintf("SUBSYSTEM-%s-%s", row.String(), col.String())
			}
			users = append(users, domain.User{
				UID:       fmt.Sprintf("SUBSYSTEM-%s-%s", new(big.Int).Set(i).String(), new(big.Int).Set(j).String()),
				RefUID:    refUid,
				Nickname:  fmt.Sprintf("sub-system-%s-%s", new(big.Int).Set(i).String(), new(big.Int).Set(j).String()),
				Email:     fmt.Sprintf("sub-system-%s-%s@system.com", new(big.Int).Set(i).String(), new(big.Int).Set(j).String()),
				CreatedAt: time.Now().UTC(),
			})
			if i.Cmp(big.NewInt(11)) < 0 {
				userPlaces = append(userPlaces, domain.UserPlace{
					UserUID:   fmt.Sprintf("SUBSYSTEM-%s-%s", new(big.Int).Set(i).String(), new(big.Int).Set(j).String()),
					MatchUID:  refUid,
					Row:       new(big.Int).Set(i),
					Col:       new(big.Int).Set(j),
					CreatedAt: time.Now().UTC(),
				})
				userRanks = append(userRanks, domain.UserRank{
					UID:      domain.GenUID(12),
					Type:     domain.UserRankTypeBuy,
					UserUID:  fmt.Sprintf("SUBSYSTEM-%s-%s", new(big.Int).Set(i).String(), new(big.Int).Set(j).String()),
					MatchUID: refUid,
					Row:      new(big.Int).Set(i),
					Col:      new(big.Int).Set(j),
					BuyUID:   "",
					RankCode: "brilliant",
					StartAt:  time.Now().UTC().Add(-24 * time.Hour),
					EndAt:    time.Now().UTC().Add(30 * 24 * time.Hour),
					Priority: 400,
				})
				userActivities = append(userActivities, domain.UserActivity{
					UserUID:  fmt.Sprintf("SUBSYSTEM-%s-%s", new(big.Int).Set(i).String(), new(big.Int).Set(j).String()),
					StartAt:  time.Now().UTC().Add(-24 * time.Hour),
					EndAt:    time.Now().UTC().Add(730 * 30 * 24 * time.Hour),
					CvAmount: 0,
				})
				teamType := domain.UserTeamTypeLeft
				q := big.NewInt(0)
				r := big.NewInt(0)
				q.DivMod(j, big.NewInt(2), r)
				if r.Int64() == 1 {
					teamType = domain.UserTeamTypeRight
				}
				userConfigs = append(userConfigs, domain.UserConfig{
					UserUID:      fmt.Sprintf("SUBSYSTEM-%s-%s", new(big.Int).Set(i).String(), new(big.Int).Set(j).String()),
					TeamType:     domain.UserTeamTypeUndefined,
					LastTeamType: teamType,
					AllowSwitch:  true,
				})
			}
		}
	}

	var err error
	for _, u := range users {
		err = s.db.UserCreate(s.ctx, &u)
		s.NoError(err)
	}

	for _, up := range userPlaces {
		err = s.db.UserPlaceCreate(s.ctx, &up)
		s.NoError(err)
	}

	for _, ur := range userRanks {
		err = s.db.UserRankCreate(s.ctx, &ur)
		s.NoError(err)
	}

	for _, ua := range userActivities {
		err = s.db.UserActivityCreate(s.ctx, &ua)
		s.NoError(err)
	}

	for _, uc := range userConfigs {
		err = s.db.UserConfigCreate(s.ctx, &uc)
		s.NoError(err)
	}

	return nil
}

func (s *BuySuite) registerTestWorkflow() {
	s.env.RegisterWorkflow(s.flow.BuyFlow)
	s.env.RegisterActivity(s.buy.Init)
	s.env.RegisterActivity(s.buy.Plan)
	s.env.RegisterActivity(s.buy.Approved)
	s.env.RegisterActivity(s.buy.Paid)
	s.env.RegisterActivity(s.buy.PlanAdd)
	s.env.RegisterActivity(s.buy.RankAddFromPlan)
	s.env.RegisterActivity(s.buy.PlanEnd)
	s.env.RegisterActivity(s.buy.RankEnd)
	s.env.RegisterActivity(s.team.ChargeRefBonus)
	s.env.RegisterActivity(s.team.ChargeBinBonus)
	s.env.RegisterActivity(s.team.ChargeMatchBonus)
	s.env.RegisterActivity(s.team.ChargeFirstRankBonus)
	s.env.RegisterActivity(s.team.ChargeApproveRankBonus)
	s.env.RegisterActivity(s.team.ChargeFastStartBonus)
	s.env.RegisterActivity(s.buy.Charged)
	s.env.RegisterActivity(s.buy.Refund)
	s.env.RegisterActivity(s.buy.Cancelled)
	s.env.RegisterActivity(s.team.UserPlaceGet)
	s.env.RegisterActivity(s.team.PlaceGetAllUp)
	s.env.RegisterActivity(s.team.PlaceGetRefUpByBuy)
	s.env.RegisterActivity(s.team.PlaceRefGetAllUp)
	s.env.RegisterActivity(s.buy.BuySetRowCol)
	s.env.RegisterActivity(s.team.BuyClientSetRowCol)
	s.env.RegisterActivity(s.team.UserRankSetRowCol)
	s.env.RegisterActivity(s.team.CalculateNextRank)
	s.env.RegisterActivity(s.team.CalculateActivity)
	s.env.RegisterActivity(s.team.ChargeCoinRefBonus)
}
