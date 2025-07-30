package workflow

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"server/internal/config"
	"server/internal/transport/telegram/lng"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"server/internal/domain"
)

const (
	ProductQueue        = "product"
	ProductBuildVersion = "v4"

	productErrorSource = "[workflow.product]"
)

type ProductWorkflow struct {
	buy  BuyService
	user UserService
}

func NewProductWorkflow(s BuyService, u UserService) *ProductWorkflow {
	return &ProductWorkflow{s, u}
}

func (bwf *ProductWorkflow) ProductFlow(ctx workflow.Context, newBuy *domain.Buy) error {
	cfg := config.Get()
	l := workflow.GetLogger(ctx)

	// workflow settings
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Hour,
		MaximumAttempts:    10,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failures by default, this is just an example.
		RetryPolicy: retryPolicy,
		// VersioningIntent: temporal.VersioningIntentCompatible,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// ==========================================
	// init buy and save runId
	newBuy.FlowUID = workflow.GetInfo(ctx).WorkflowExecution.RunID
	if cfg.Env == "test" {
		runId, _ := uuid.NewUUID()
		newBuy.FlowUID = runId.String()
	}

	// ==========================================
	// pay from balance
	err := workflow.ExecuteActivity(ctx, bwf.buy.BuyPay, newBuy).Get(ctx, newBuy)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrProductWorkflowInit).Add(err)
	}

	// ==========================================
	// init paid
	err = workflow.ExecuteActivity(ctx, bwf.buy.InitPaid, newBuy).Get(ctx, newBuy)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrProductWorkflowInit).Add(err)
	}

	// ==========================================
	// get product
	product := &domain.Product{}
	err = workflow.ExecuteActivity(ctx, bwf.buy.Product, newBuy).Get(ctx, product)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrProductWorkflowProduct).Add(err)
	}

	// ==========================================
	// add product to user
	newUserProduct := &domain.UserProduct{}
	err = workflow.ExecuteActivity(ctx, bwf.buy.ProductAdd, newBuy, product).Get(ctx, newUserProduct)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrProductWorkflowPlanApply).Add(err)
	}

	// ==========================================
	// apply product to user
	err = workflow.ExecuteActivity(ctx, bwf.buy.ProductApply, newUserProduct).Get(ctx, newUserProduct)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrProductWorkflowPlanAdd).Add(err)
	}

	// ==========================================
	// charge ref bonuses
	var upUsers []*domain.User
	err = workflow.ExecuteActivity(ctx, bwf.buy.GetRefUsersUp, newBuy.UserUID, 1).Get(ctx, &upUsers)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrProductWorkflowRefUsersUp).Add(err)
	}

	for level, user := range upUsers {
		var bonusAmount int64 = 0
		err = workflow.ExecuteActivity(ctx, bwf.buy.RefUserCharge, newBuy, user, level).Get(ctx, &bonusAmount)
		if err != nil {
			return domain.NewError(productErrorSource).SetCode(domain.ErrProductWorkflowRefUserCharge).Add(err)
		}

		if user.TgID > 0 && bonusAmount > 0 {
			tgId := user.TgID
			msg := lng.Get("refBuy", user.Locale)
			if newBuy.CurrencyCode == "usd" {
				msg = fmt.Sprintf(msg, decimal.NewFromInt(bonusAmount).Div(decimal.NewFromInt(100)).StringFixed(2), "UDEX")
			} else if newBuy.CurrencyCode == "abt" {
				msg = fmt.Sprintf(msg, decimal.NewFromInt(bonusAmount).Div(decimal.NewFromInt(1000000000)).StringFixed(4), "ABT")
			} else {
				continue
			}

			// send notification about successful ref buy bonus if telegram id is not empty
			err = workflow.ExecuteActivity(ctx, bwf.user.NotifyTgSend, tgId, msg).Get(ctx, &tgId)
			if err != nil {
				l.Error("Error send notification about ref buy bonus", "error", err)
			}
		}
	}

	// ==========================================
	// save the date of charging all bonuses if all ok
	err = workflow.ExecuteActivity(ctx, bwf.buy.Charged, newBuy, getWorkflowNow(ctx, newBuy.CreatedAt)).Get(ctx, newBuy)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowCharged).Add(err)
	}

	return nil
}
