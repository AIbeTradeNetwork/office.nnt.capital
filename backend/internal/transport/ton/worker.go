package ton

import (
	"context"
	"server/internal/domain"
	"server/internal/repository"
	"server/internal/service/user"
	"server/internal/transport/temporal"
	"server/internal/workflow"
)

const (
	workerErrorSource = "[transport.ton.worker]"
)

func StartWorker() error {
	ctx := context.Background()

	worker, err := temporal.Connect(ctx, workflow.BuyQueue, workflow.DepositQueue, workflow.DepositBuildVersion)
	if err != nil {
		return domain.NewError(workerErrorSource).SetCode(domain.ErrConnect).Add(err)
	}

	dbRepo, err := repository.NewDbRepo(ctx)
	if err != nil {
		return domain.NewError(workerErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	err = dbRepo.EnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(workerErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	userService := user.NewUserService(dbRepo, nil, nil, nil, nil, nil)

	depWorkflow := workflow.NewDepositWorkflow(userService)
	worker.Register(
		depWorkflow.DepositFlow,
		userService.DepositCreate,
		userService.DepositTransactionCreate,
	)

	err = worker.Run()
	if err != nil {
		return err
	}

	return nil
}
