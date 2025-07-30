package temporal

import (
	"context"
	"fmt"
	"server/internal/workflow"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"server/internal/config"
	"server/internal/domain"
)

const (
	transportErrorSource = "[transport.temporal]"
)

type Worker struct {
	w *worker.Worker
}

func Connect(ctx context.Context, namespace string, queue string, version string) (*Worker, error) {
	cfg := config.Get()

	c, err := client.Dial(client.Options{
		HostPort:  cfg.TemporalURL,
		Namespace: namespace,
	})

	if err != nil {
		fmt.Println(err)
		return nil, domain.NewError(transportErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	err = nil
	switch queue {
	case workflow.BuyQueue:
		err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
			TaskQueue: workflow.BuyQueue,
			Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
				BuildID: workflow.BuyBuildVersion,
			},
		})
	case workflow.NotifyQueue:
		err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
			TaskQueue: workflow.NotifyQueue,
			Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
				BuildID: workflow.NotifyBuildVersion,
			},
		})
	case workflow.DepositQueue:
		err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
			TaskQueue: workflow.DepositQueue,
			Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
				BuildID: workflow.DepositBuildVersion,
			},
		})
	case workflow.ProductQueue:
		err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
			TaskQueue: workflow.ProductQueue,
			Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
				BuildID: workflow.ProductBuildVersion,
			},
		})
	case workflow.AutofarmQueue:
		err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
			TaskQueue: workflow.AutofarmQueue,
			Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
				BuildID: workflow.AutofarmBuildVersion,
			},
		})
	}

	if err != nil {
		fmt.Println(err)
	}

	w := worker.New(c, queue, worker.Options{
		BuildID:                 version,
		UseBuildIDForVersioning: true,
	})

	return &Worker{&w}, nil
}

func (w *Worker) Register(flow interface{}, activities ...interface{}) {
	wn := *w.w
	wn.RegisterWorkflow(flow)
	for _, act := range activities {
		wn.RegisterActivity(act)
	}
}

func (w *Worker) Run() error {
	wn := *w.w
	err := wn.Run(worker.InterruptCh())
	if err != nil {
		return domain.NewError(transportErrorSource).SetCode(domain.ErrWorkerRun).Add(err)
	}
	return nil
}
