package workflow

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"server/internal/config"
	"server/internal/transport/telegram/lng"
	"strconv"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"server/internal/domain"
)

const (
	AutofarmQueue        = "autofarm"
	AutofarmBuildVersion = "v1"

	autofarmErrorSource = "[workflow.autofarm]"
)

type AutofarmWorkflow struct {
	buy  BuyService
	user UserService
}

func NewAutofarmWorkflow(b BuyService, u UserService) *AutofarmWorkflow {
	return &AutofarmWorkflow{b, u}
}

func (awf *AutofarmWorkflow) AutofarmFlow(ctx workflow.Context, userProduct *domain.UserProduct) error {
	cfg := config.Get()
	l := workflow.GetLogger(ctx)

	// workflow settings
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Hour,
		MaximumAttempts:    3,
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
	// init autofarm and save runId
	userProduct.FlowUID = workflow.GetInfo(ctx).WorkflowExecution.RunID
	if cfg.Env == "test" {
		runId, _ := uuid.NewUUID()
		userProduct.FlowUID = runId.String()
	}

	// ==========================================
	// if startAt is in the future, sleep until it is in the future
	if userProduct.StartAt.After(workflow.Now(ctx)) {
		l.Info(fmt.Sprintf("Sleeping until %s", userProduct.StartAt))
		err := workflow.Sleep(ctx, userProduct.StartAt.Sub(workflow.Now(ctx)))
		if err != nil {
			return domain.NewError(autofarmErrorSource).SetCode(domain.ErrAutofarmWorkflowSleep).Add(err)
		}
	}

	// ==========================================
	// init autofarm
	l.Info("Init autofarm")
	err := workflow.ExecuteActivity(ctx, awf.buy.InitAutofarm, userProduct).Get(ctx, userProduct)
	if err != nil {
		l.Error("Error init autofarm", "error", err)
		return domain.NewError(productErrorSource).SetCode(domain.ErrAutofarmWorkflowInit).Add(err)
	}

	var tgIdString string
	var tgId int64
	err = workflow.ExecuteActivity(ctx, awf.user.GetTelegramIDByUserUID, userProduct.UserUID).Get(ctx, &tgIdString)
	if err == nil {
		tgId, _ = strconv.ParseInt(tgIdString, 10, 64)
	}
	l.Info(fmt.Sprintf("Telegram id: %d", tgId))

	dUser := &domain.User{}
	err = workflow.ExecuteActivity(ctx, awf.user.GetUserByUID, userProduct.UserUID).Get(ctx, dUser)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrAutofarmWorkflowGetUser).Add(err)
	}
	l.Info(fmt.Sprintf("User locale: %s", dUser.Locale))

	for workflow.Now(ctx).Before(userProduct.EndAt) {
		l.Info(fmt.Sprintf("Autofarm runId: %s", userProduct.FlowUID))
		// ==========================================
		// claim tokens
		userClaim := &domain.UserClaim{}
		err = workflow.ExecuteActivity(ctx, awf.user.UserClaimCreate, "ABT", userProduct.UserUID).Get(ctx, userClaim)
		if err != nil {
			l.Error("Error claim tokens", "error", err)
		}
		if err == nil && tgId > 0 {
			l.Info("Send notification about successful autoclaim")
			msg := lng.Get("autoclaimSuccess", dUser.Locale)
			msg = fmt.Sprintf(msg, decimal.NewFromInt(userClaim.Amount).Div(decimal.NewFromInt(1000000000)).StringFixed(4))
			// send notification about successful autoclaim if telegram id is not empty
			err = workflow.ExecuteActivity(ctx, awf.user.NotifyTgSend, tgId, msg).Get(ctx, &tgId)
			if err != nil {
				l.Error("Error send notification about successful autoclaim", "error", err)
			}
		}

		l.Info("Sleeping until next claim in 8 hrs")
		err = workflow.Sleep(ctx, time.Hour*8)
		if err != nil {
			return domain.NewError(autofarmErrorSource).SetCode(domain.ErrAutofarmWorkflowSleep).Add(err)
		}
	}

	// ==========================================
	// send notification about autoclaim ends if telegram id is not empty
	if tgId > 0 {
		msg := lng.Get("autoclaimEnd", dUser.Locale)
		_ = workflow.ExecuteActivity(ctx, awf.user.NotifyTgSend, tgId, msg).Get(ctx, &tgId)
	}

	return nil
}
