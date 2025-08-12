package workflow

import (
    "github.com/google/uuid"
    "server/internal/config"
    "time"

    "go.temporal.io/sdk/temporal"
    "go.temporal.io/sdk/workflow"
)

const (
    PartnerApplicationQueue        = "partner-application"
    PartnerApplicationBuildVersion = "v1"
)

type PartnerApplicationWorkflow struct {
    team TeamService
}

func NewPartnerApplicationWorkflow(t TeamService) *PartnerApplicationWorkflow {
    return &PartnerApplicationWorkflow{team: t}
}

// PartnerApplicationFlow waits 4 hours and processes the application automatically
func (pwf *PartnerApplicationWorkflow) PartnerApplicationFlow(ctx workflow.Context, applicationUID string) error {
    cfg := config.Get()

    // workflow settings
    retryPolicy := &temporal.RetryPolicy{
        InitialInterval:    time.Second,
        BackoffCoefficient: 2.0,
        MaximumInterval:    time.Minute,
        MaximumAttempts:    3,
    }
    options := workflow.ActivityOptions{
        StartToCloseTimeout: time.Minute * 5,
        RetryPolicy:         retryPolicy,
    }

    ctx = workflow.WithActivityOptions(ctx, options)

    // for deterministic tests
    flowUID := workflow.GetInfo(ctx).WorkflowExecution.RunID
    if cfg.Env == "test" {
        runId, _ := uuid.NewUUID()
        flowUID = runId.String()
    }
    _ = flowUID // reserved for future usage

    // Wait for 4 hours using Temporal timer
    err := workflow.Sleep(ctx, 4*time.Hour)
    if err != nil {
        return err
    }

    // After timer, call activity to process the single application by UID
    var actErr error
    err = workflow.ExecuteActivity(ctx, pwf.team.ProcessExpiredApplicationByUID, applicationUID).Get(ctx, &actErr)
    if err != nil {
        return err
    }
    if actErr != nil {
        return actErr
    }

    return nil
}


