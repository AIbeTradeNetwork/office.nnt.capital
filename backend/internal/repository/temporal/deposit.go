package temporal

import (
	"context"
	"fmt"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"server/internal/domain"
	"server/internal/workflow"
)

const (
	depositErrorSource = "[repository.temporal.deposit]"
)

type DepositWorkflow interface {
	DepositFlow(context.Context, *domain.Deposit) (*domain.Deposit, error)
}

func (mr *Repo) DepositCreate(ctx context.Context, newDeposit *domain.Deposit) error {
	options := client.StartWorkflowOptions{
		ID:                    newDeposit.UID,
		TaskQueue:             workflow.DepositQueue,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}

	c := *mr.c

	_, err := c.ExecuteWorkflow(ctx, options, "DepositFlow", newDeposit)
	if err != nil {
		fmt.Println(err)
		return domain.NewError(depositErrorSource).SetCode(domain.ErrDepositWorkflowExecute).Add(err)
	}

	return nil
}
