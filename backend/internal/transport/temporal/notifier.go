package temporal

import (
	"context"
	"server/internal/domain"
	"server/internal/repository"
	"server/internal/service/user"
	"server/internal/workflow"
)

const (
	notifierErrorSource = "[transport.temporal.notifier]"
)

func StartNotifier() error {
	ctx := context.Background()

	worker, err := Connect(ctx, workflow.BuyQueue, workflow.NotifyQueue, workflow.NotifyBuildVersion)
	if err != nil {
		return domain.NewError(notifierErrorSource).SetCode(domain.ErrConnect).Add(err)
	}

	dbRepo, err := repository.NewDbRepo(ctx)
	if err != nil {
		return domain.NewError(notifierErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	userService := user.NewUserService(dbRepo, nil, nil, nil, nil, nil)

	notifyWorkflow := workflow.NewNotifyWorkflow(userService)
	worker.Register(
		notifyWorkflow.NotifyFlow,
		userService.NotifyCreate,
		userService.NotifyGetAllTgID,
		userService.NotifyUpdate,
		userService.NotifyTgSend,
	)

	err = worker.Run()
	if err != nil {
		return err
	}

	return nil
}
