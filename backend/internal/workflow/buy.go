package workflow

import (
	"github.com/google/uuid"
	"log"
	"server/internal/config"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"server/internal/domain"
)

const (
	BuyQueue            = "buy"
	BuyPaidSignalName   = "paid"
	BuyRefundSignalName = "refund"
	BuyBuildVersion     = "v3"

	buyErrorSource = "[workflow.buy]"
)

type BuyWorkflow struct {
	buy  BuyService
	team TeamService
}

func NewBuyWorkflow(s BuyService, t TeamService) *BuyWorkflow {
	return &BuyWorkflow{s, t}
}

func (bwf *BuyWorkflow) BuyFlow(ctx workflow.Context, newBuy *domain.Buy) error {
	cfg := config.Get()

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
		return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowInit).Add(err)
	}

	// ==========================================
	// get plan
	plan := &domain.Plan{}
	err = workflow.ExecuteActivity(ctx, bwf.buy.Plan, newBuy).Get(ctx, plan)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowPlan).Add(err)
	}

	// ==========================================
	// add plan to user
	newUserPlan := &domain.UserPlan{}
	err = workflow.ExecuteActivity(ctx, bwf.buy.PlanAdd, newBuy, plan).Get(ctx, newUserPlan)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowPlanAdd).Add(err)
	}

	// ==========================================
	// add rank to user if rank is set in plan settings
	newUserRank := &domain.UserRank{}
	if plan.RankCode != "" && plan.RankPeriod > 0 {
		err = workflow.ExecuteActivity(ctx, bwf.buy.RankAddFromPlan, newBuy, plan).Get(ctx, newUserRank)
		if err != nil {
			return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowRankAddFromPlan).Add(err)
		}
	}

	// ==========================================
	// charge ref bonus to referral
	err = workflow.ExecuteActivity(ctx, bwf.team.ChargeRefBonus, newBuy).Get(ctx, newBuy)
	if err != nil {
		log.Println(err)
		return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowChargeRef).Add(err)
	}

	// ==========================================
	// charge coin ref bonus to referral
	err = workflow.ExecuteActivity(ctx, bwf.team.ChargeCoinRefBonus, newBuy).Get(ctx, newBuy)
	if err != nil {
		log.Println(err)
		return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowChargeRef).Add(err)
	}

	// ==========================================
	// check user place if exists then user is distributor
	userPlace := &domain.UserPlace{}
	err = workflow.ExecuteActivity(ctx, bwf.team.UserPlaceGet, newBuy).Get(ctx, userPlace)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowPlace).Add(err)
	}

	var upPlaces []*domain.UserPlace

	if userPlace != nil && userPlace.UserUID != "" {
		// logic for distributor's buy

		// ==========================================
		// set place row and col to buy for optimization
		err = workflow.ExecuteActivity(ctx, bwf.buy.BuySetRowCol, newBuy, userPlace).Get(ctx, newBuy)
		if err != nil {
			return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowSetRowCol).Add(err)
		}

		// ==========================================
		// set place row and col to user rank for optimization
		err = workflow.ExecuteActivity(ctx, bwf.team.UserRankSetRowCol, newUserRank, userPlace).Get(ctx, newUserRank)
		if err != nil {
			return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowSetRowCol).Add(err)
		}

		// ==========================================
		// get all places up
		err = workflow.ExecuteActivity(ctx, bwf.team.PlaceGetAllUp, userPlace).Get(ctx, &upPlaces)
		if err != nil {
			return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowPlacesUp).Add(err)
		}

	} else {
		// logic for client's buy

		// ==========================================
		// get the nearest distributor by referrals
		nearestDistPlace := &domain.UserPlace{}
		err = workflow.ExecuteActivity(ctx, bwf.team.PlaceGetRefUpByBuy, newBuy).Get(ctx, nearestDistPlace)
		if err != nil {
			return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowPlace).Add(err)
		}

		if nearestDistPlace != nil && nearestDistPlace.UserUID != "" {
			// ==========================================
			// set place row and col to buy for optimization
			err = workflow.ExecuteActivity(ctx, bwf.team.BuyClientSetRowCol, newBuy, nearestDistPlace).Get(ctx, newBuy)
			if err != nil {
				return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowSetClientRowCol).Add(err)
			}

			// ==========================================
			// calculate activity
			nextUserActivity := &domain.UserActivity{}
			err = workflow.ExecuteActivity(ctx, bwf.team.CalculateActivity, newBuy, nearestDistPlace).Get(ctx, nextUserActivity)
			if err != nil {
				return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowCalculateActivity).Add(err)
			}

			// ==========================================
			// get all places up
			err = workflow.ExecuteActivity(ctx, bwf.team.PlaceGetAllUp, nearestDistPlace).Get(ctx, &upPlaces)
			if err != nil {
				return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowPlacesUp).Add(err)
			}

			upPlaces = append([]*domain.UserPlace{
				nearestDistPlace,
			}, upPlaces...)
		}
	}

	for binLevel, upPlace := range upPlaces {

		// ==========================================
		// charge bin bonus to top
		binBonus := &domain.Transaction{}
		err = workflow.ExecuteActivity(ctx, bwf.team.ChargeBinBonus, newBuy, upPlace, binLevel).Get(ctx, binBonus)
		if err != nil {
			return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowChargeBin).Add(err)
		}

		// if bin bonus charged then charge match bonus
		if binBonus != nil && binBonus.PosAmount > 0 {
			// ==========================================
			// get places for match bonus
			var upRefPlaces []*domain.UserPlace
			err = workflow.ExecuteActivity(ctx, bwf.team.PlaceRefGetAllUp, newBuy, upPlace).Get(ctx, &upRefPlaces)
			if err != nil {
				return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowPlacesRefUp).Add(err)
			}

			upRefPlaceLevel := 0
			for _, upRefPlace := range upRefPlaces {
				// ==========================================
				// charge next user match bonus
				matchBonus := &domain.Transaction{}
				err = workflow.ExecuteActivity(ctx, bwf.team.ChargeMatchBonus, newBuy, upRefPlace, binBonus, upRefPlaceLevel).Get(ctx, matchBonus)
				if err != nil {
					return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowChargeMatch).Add(err)
				}
				if matchBonus != nil && matchBonus.Amount > 0 {
					upRefPlaceLevel++
				}
			}
		}

		// ==========================================
		// calculate rank for next user
		nextUserRank := &domain.UserRank{}
		err = workflow.ExecuteActivity(ctx, bwf.team.CalculateNextRank, newBuy, upPlace).Get(ctx, nextUserRank)
		if err != nil {
			return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowCalculateNextRank).Add(err)
		}

		// ==========================================
		// check for one time bonus for first time rank closing
		err = workflow.ExecuteActivity(ctx, bwf.team.ChargeFirstRankBonus, newBuy, upPlace).Get(ctx, newBuy)
		if err != nil {
			return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowChargeFirstRank).Add(err)
		}

		// ==========================================
		// check for one time bonus for confirmation of previous rank
		err = workflow.ExecuteActivity(ctx, bwf.team.ChargeApproveRankBonus, newBuy, upPlace).Get(ctx, newBuy)
		if err != nil {
			return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowChargeFirstRank).Add(err)
		}

		// ==========================================
		// check fast start bonus
		err = workflow.ExecuteActivity(ctx, bwf.team.ChargeFastStartBonus, newBuy, upPlace).Get(ctx, newBuy)
		if err != nil {
			return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowChargeFastStart).Add(err)
		}

		// ==========================================
		// unknown bonuses
		// TODO: activity bonus (?)
		// TODO: leader bonus (?)

	}

	// ==========================================
	// save the date of charging all bonuses if all ok
	err = workflow.ExecuteActivity(ctx, bwf.buy.Charged, newBuy, getWorkflowNow(ctx, newBuy.CreatedAt)).Get(ctx, newBuy)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrBuyWorkflowCharged).Add(err)
	}

	return nil
}

func getWorkflowNow(ctx workflow.Context, curTime time.Time) time.Time {
	cfg := config.Get()
	if cfg.Env == "test" {
		return curTime.Add(24 * time.Hour)
	}
	return workflow.Now(ctx)
}
