package workflow_test

import (
	"context"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
	"log"
	"server/internal/repository"
	"server/internal/repository/mongodb"
	"server/internal/repository/ton"
	"server/internal/service/seed"
	"server/internal/service/user"
	"server/internal/workflow"
	"testing"
	"time"
)

type DepSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	ctx context.Context
	ton *ton.Repo
	db  *mongodb.Repo

	user *user.Service
	seed *seed.Service

	flow *workflow.DepositWorkflow

	env *testsuite.TestWorkflowEnvironment
}

func (s *DepSuite) SetupTest() {
	s.ctx = context.Background()
	s.env = s.NewTestWorkflowEnvironment().
		SetTestTimeout(60 * time.Second).
		SetWorkflowRunTimeout(60 * time.Second)

	tonRepo, _ := repository.NewTonRepo(s.ctx)
	s.ton = tonRepo

	dbRepo, _ := repository.NewDbRepo(s.ctx)
	s.db = dbRepo

	s.seed = seed.NewSeedService(s.db)

	s.user = user.NewUserService(nil, nil, tonRepo, nil, nil, nil)

	s.flow = workflow.NewDepositWorkflow(s.user)
	s.registerTestWorkflow()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DepSuite))
}

func (s *DepSuite) Test_DepFlow() {
	err := s.db.DropTest(s.ctx)
	s.NoError(err)

	// Ensure indexes
	err = s.db.EnsureIndexes(s.ctx)
	s.NoError(err)
	err = s.seed.Seed(s.ctx)
	log.Println(err)
	s.env.ExecuteWorkflow(s.flow.DepositFlow, nil)
}

func (s *DepSuite) registerTestWorkflow() {
	s.env.RegisterWorkflow(s.flow.DepositFlow)
	s.env.RegisterActivity(s.user.GetLastTonDeposits)
}
