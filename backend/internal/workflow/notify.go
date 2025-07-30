package workflow

import (
	"github.com/google/uuid"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"server/internal/config"
	"server/internal/domain"
	"strconv"
	"time"
)

const (
	NotifyQueue        = "notify"
	NotifyBuildVersion = "v1"

	notifyErrorSource = "[workflow.notify]"
)

type NotifyWorkflow struct {
	user UserService
}

func NewNotifyWorkflow(u UserService) *NotifyWorkflow {
	return &NotifyWorkflow{u}
}

func (nwf *NotifyWorkflow) NotifyFlow(ctx workflow.Context, newNotify *domain.Notification, page int64) error {
	cfg := config.Get()

	// workflow settings
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
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

	newNotify.FlowUID = workflow.GetInfo(ctx).WorkflowExecution.RunID
	if cfg.Env == "test" {
		runId, _ := uuid.NewUUID()
		newNotify.FlowUID = runId.String()
	}

	err := workflow.ExecuteActivity(ctx, nwf.user.NotifyCreate, newNotify).Get(ctx, newNotify)
	if err != nil {
		return domain.NewError(notifyErrorSource).SetCode(domain.ErrNotifyWorkflowCreate).Add(err)
	}

	var tgIds []*domain.UserTg
	err = workflow.ExecuteActivity(ctx, nwf.user.NotifyGetAllTgID, newNotify, page).Get(ctx, &tgIds)
	if err != nil {
		return domain.NewError(notifyErrorSource).SetCode(domain.ErrNotifyWorkflowGetAllTgIDs).Add(err)
	}

	for _, tgId := range tgIds {

		lang := "en"
		if tgId.Locale == "ru" {
			lang = tgId.Locale
		}
		msg, ok := newNotify.Texts[lang]
		if !ok {
			continue
		}

		tgIdInt, err := strconv.ParseInt(tgId.TgUID, 10, 64)
		if err != nil {
			continue
		}

		var tgIdSent int64
		err = workflow.ExecuteActivity(ctx, nwf.user.NotifyTgSend, tgIdInt, msg).Get(ctx, &tgIdSent)
		if err != nil {
			return domain.NewError(notifyErrorSource).SetCode(domain.ErrNotifyWorkflowTgSend).Add(err)
		}

		//workflow.Sleep(ctx, 50*time.Millisecond) // Telegram rate limit
	}

	if len(tgIds) >= 5000 {
		page++
		return workflow.NewContinueAsNewError(ctx, nwf.NotifyFlow, newNotify, page)
	}

	newNotify.SentAt = workflow.Now(ctx)
	err = workflow.ExecuteActivity(ctx, nwf.user.NotifyUpdate, newNotify).Get(ctx, newNotify)
	if err != nil {
		return domain.NewError(notifyErrorSource).SetCode(domain.ErrNotifyWorkflowUpdate).Add(err)
	}

	return nil
}
