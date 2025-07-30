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
	productErrorSource = "[transport.temporal.product]"
)

func StartProduct() error {
	ctx := context.Background()

	worker, err := Connect(ctx, workflow.BuyQueue, workflow.ProductQueue, workflow.ProductBuildVersion)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrConnect).Add(err)
	}

	dbRepo, err := repository.NewDbRepo(ctx)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	wfRepo, err := repository.NewWfRepo(ctx)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	err = dbRepo.EnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	authService := auth.NewAuthService(dbRepo, nil)
	buyService := buy.NewBuyService(dbRepo, wfRepo, nil, authService)
	userService := user.NewUserService(dbRepo, nil, nil, nil, authService, nil)

	productWorkflow := workflow.NewProductWorkflow(buyService, userService)
	worker.Register(
		productWorkflow.ProductFlow,
		buyService.BuyPay,
		buyService.InitPaid,
		buyService.Product,
		buyService.ProductAdd,
		buyService.ProductApply,
		buyService.GetRefUsersUp,
		buyService.RefUserCharge,
		userService.NotifyTgSend,
		buyService.Charged,
	)

	err = worker.Run()
	if err != nil {
		return err
	}

	return nil
}
