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
	autofarmErrorSource = "[repository.temporal.autofarm]"
)

type AutofarmWorkflow interface {
	AutofarmFlow(context.Context, *domain.UserProduct) error
}

func (mr *Repo) AutofarmCreate(ctx context.Context, newAutofarm *domain.UserProduct) error {
	options := client.StartWorkflowOptions{
		ID:                    newAutofarm.UID,
		TaskQueue:             workflow.AutofarmQueue,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}

	c := *mr.c

	_, err := c.ExecuteWorkflow(ctx, options, "AutofarmFlow", newAutofarm)
	if err != nil {
		fmt.Println(err)
		return domain.NewError(depositErrorSource).SetCode(domain.ErrAutofarmWorkflowExecute).Add(err)
	}

	return nil
}
