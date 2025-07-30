package temporal

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"server/internal/domain"
	"server/internal/workflow"
)

const (
	buyErrorSource = "[repository.temporal.buy]"
)

type BuyWorkflow interface {
	BuyFlow(context.Context, *domain.Buy) (*domain.Buy, error)
}

func (mr *Repo) BuyCreate(ctx context.Context, newBuy *domain.Buy) error {
	options := client.StartWorkflowOptions{
		ID:                    newBuy.UID,
		TaskQueue:             workflow.BuyQueue,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}

	c := *mr.c

	_, err := c.ExecuteWorkflow(ctx, options, "BuyFlow", newBuy)
	if err != nil {
		fmt.Println(err)
		return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowExecute).Add(err)
	}

	return nil
}

func (mr *Repo) BuyProduct(ctx context.Context, newBuy *domain.Buy) error {
	options := client.StartWorkflowOptions{
		ID:                    newBuy.UID,
		TaskQueue:             workflow.ProductQueue,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}

	c := *mr.c

	_, err := c.ExecuteWorkflow(ctx, options, "ProductFlow", newBuy)
	if err != nil {
		fmt.Println(err)
		return domain.NewError(buyErrorSource).SetCode(domain.ErrProductWorkflowExecute).Add(err)
	}

	return nil
}

func (mr *Repo) BuySignalPaid(ctx context.Context, buy *domain.Buy, info *domain.BuySignalPaid) error {
	c := *mr.c
	err := c.SignalWorkflow(ctx, buy.UID, buy.FlowUID, workflow.BuyPaidSignalName, info)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrSignal).Add(err)
	}
	return nil
}

func (mr *Repo) BuySignalRefund(ctx context.Context, buy *domain.Buy, info *domain.BuySignalRefund) error {
	c := *mr.c
	err := c.SignalWorkflow(ctx, buy.UID, buy.FlowUID, workflow.BuyRefundSignalName, info)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrSignal).Add(err)
	}
	return nil
}

func (mr *Repo) buyEnsureNamespaces(ctx context.Context, nc client.NamespaceClient) error {
	_ = nc.Register(ctx, &workflowservice.RegisterNamespaceRequest{
		Namespace:                        workflow.BuyQueue,
		WorkflowExecutionRetentionPeriod: durationpb.New(90 * 24 * time.Hour),
	})
	return nil
}
