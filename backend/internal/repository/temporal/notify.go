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
	notifyErrorSource = "[repository.temporal.notify]"
)

type NotifyWorkflow interface {
	NotifyFlow(context.Context, *domain.Notification, int64) (*domain.Notification, error)
}

func (mr *Repo) NotifyCreate(ctx context.Context, newNotify *domain.Notification) error {
	options := client.StartWorkflowOptions{
		ID:                    newNotify.UID,
		TaskQueue:             workflow.NotifyQueue,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}

	c := *mr.c

	_, err := c.ExecuteWorkflow(ctx, options, "NotifyFlow", newNotify, 1)
	if err != nil {
		fmt.Println(err)
		return domain.NewError(notifyErrorSource).SetCode(domain.ErrNotifyWorkflowExecute).Add(err)
	}

	return nil
}
