package workflow

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"server/internal/domain"
	"time"
)

const (
	DepositQueue        = "deposit"
	DepositBuildVersion = "v1"

	depositErrorSource = "[workflow.deposit]"
)

type DepositWorkflow struct {
	user UserService
}

func NewDepositWorkflow(u UserService) *DepositWorkflow {
	return &DepositWorkflow{u}
}

func (dwf *DepositWorkflow) DepositFlow(ctx workflow.Context, newDeposit *domain.Deposit) error {
	// workflow settings
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Hour,
		MaximumAttempts:    10,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failures by default, this is just an example.
		RetryPolicy: retryPolicy,
		// VersioningIntent: temporal.VersioningIntentCompatible,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, dwf.user.DepositCreate, newDeposit).Get(ctx, newDeposit)
	if err != nil {
		return domain.NewError(depositErrorSource).SetCode(domain.ErrDepositWorkflowGetLastDeposits).Add(err)
	}

	err = workflow.ExecuteActivity(ctx, dwf.user.DepositTransactionCreate, newDeposit).Get(ctx, newDeposit)
	if err != nil {
		return domain.NewError(depositErrorSource).SetCode(domain.ErrDepositWorkflowGetLastDeposits).Add(err)
	}

	return nil
}
