package temporal

import (
	"context"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/provider/payment"
	"server/internal/repository"
	"server/internal/service/auth"
	"server/internal/service/buy"
	"server/internal/service/team"
	"server/internal/workflow"
)

const (
	workerErrorSource = "[transport.temporal.worker]"
)

func Start() error {
	ctx := context.Background()

	cfg := config.Get()

	//init providers
	pbProvider := payment.New(payment.Config{
		URL:   cfg.PaymentGatewayURL,
		Login: cfg.PaymentGatewayLogin,
		Key:   cfg.PaymentGatewayKey,
	})

    worker, err := Connect(ctx, workflow.BuyQueue, workflow.BuyQueue, workflow.BuyBuildVersion)
	if err != nil {
		return domain.NewError(workerErrorSource).SetCode(domain.ErrConnect).Add(err)
	}
    // separate worker for partner applications waiting flow (same namespace, different task queue)
    partnerWorker, err := Connect(ctx, workflow.BuyQueue, workflow.PartnerApplicationQueue, workflow.PartnerApplicationBuildVersion)
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

	buyService := buy.NewBuyService(dbRepo, nil, pbProvider, nil)
	authService := auth.NewAuthService(dbRepo, nil)
	teamService := team.NewTeamService(dbRepo, pbProvider, authService)

    buyWorkflow := workflow.NewBuyWorkflow(buyService, teamService)
    partnerApplicationWorkflow := workflow.NewPartnerApplicationWorkflow(teamService)
    worker.Register(
		buyWorkflow.BuyFlow,
		buyService.BuyPay,
		buyService.InitPaid,
		buyService.Plan,
		buyService.Paid,
		buyService.PlanAdd,
		buyService.RankAddFromPlan,
		teamService.UserPlaceGet,
		teamService.ChargeRefBonus,
		teamService.ChargeCoinRefBonus,
		buyService.BuySetRowCol,
		teamService.PlaceGetAllUp,
		teamService.PlaceGetRefUpByBuy,
		teamService.BuyClientSetRowCol,
		teamService.UserRankSetRowCol,
		teamService.CalculateActivity,
		teamService.PlaceRefGetAllUp,
		teamService.ChargeBinBonus,
		teamService.ChargeMatchBonus,
		teamService.CalculateNextRank,
		teamService.ChargeFirstRankBonus,
		teamService.ChargeApproveRankBonus,
		teamService.ChargeFastStartBonus,
		buyService.Charged,
	)

    // register partner application workflow and its single activity on partner worker
    partnerWorker.Register(
        partnerApplicationWorkflow.PartnerApplicationFlow,
        teamService.ProcessExpiredApplicationByUID,
    )

    // run partner worker in background and main worker in foreground
    go func() { _ = partnerWorker.Run() }()
    if err := worker.Run(); err != nil {
        return err
    }
    return nil
}
