package temporal

import (
	"context"
	"server/internal/domain"
	"server/internal/repository"
	"server/internal/service/auth"
	"server/internal/service/buy"
	"server/internal/service/user"
	"server/internal/workflow"
)

const (
	autofarmErrorSource = "[transport.temporal.autofarm]"
)

func StartAutofarm() error {
	ctx := context.Background()

	worker, err := Connect(ctx, workflow.BuyQueue, workflow.AutofarmQueue, workflow.AutofarmBuildVersion)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrConnect).Add(err)
	}

	dbRepo, err := repository.NewDbRepo(ctx)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	err = dbRepo.EnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	buyService := buy.NewBuyService(dbRepo, nil, nil, nil)
	authService := auth.NewAuthService(dbRepo, nil)
	userService := user.NewUserService(dbRepo, nil, nil, nil, authService, nil)

	autofarmWorkflow := workflow.NewAutofarmWorkflow(buyService, userService)
	worker.Register(
		autofarmWorkflow.AutofarmFlow,
		buyService.InitAutofarm,
		userService.GetTelegramIDByUserUID,
		userService.GetUserByUID,
		userService.UserClaimCreate,
		userService.NotifyTgSend,
	)

	err = worker.Run()
	if err != nil {
		return err
	}

	return nil
}
