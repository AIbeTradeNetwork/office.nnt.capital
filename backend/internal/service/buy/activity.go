package buy

import (
	"context"
	"math/big"
	"time"

	"server/internal/domain"
)

// Init to save new buy to database or return existing one
func (s *Service) Init(ctx context.Context, newBuy *domain.Buy) (*domain.Buy, error) {
	exBuy, _ := s.db.BuyGetByUID(ctx, newBuy.UID)
	if exBuy != nil {
		return exBuy, nil
	}
	err := s.db.BuyCreate(ctx, newBuy)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return newBuy, nil
}

// BuyPay to save new buy to database or return existing one
func (s *Service) BuyPay(ctx context.Context, newBuy *domain.Buy) (*domain.Buy, error) {
	switch newBuy.CurrencyCode {
	case "usd":
		trans := domain.Transaction{
			UID:         domain.GenUID(12),
			UserUID:     newBuy.UserUID,
			FromUID:     newBuy.UserUID,
			Percent:     0,
			Level:       0,
			Type:        domain.TransactionTypeBuy,
			RankCode:    "",
			Amount:      -newBuy.Amount,
			PosAmount:   -newBuy.Amount,
			FullAmount:  -newBuy.Amount,
			Coefficient: 100,
			BuyUID:      newBuy.UID,
			PayoutUID:   "",
			DepositUID:  "",
			CreatedAt:   time.Now().UTC(),
			ChargedAt:   time.Now().UTC(),
			MsgCodes:    nil,
		}

		err := s.db.TransactionCreate(ctx, &trans)
		if err != nil {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrPayForOrder).Add(err)
		}
	case "abt":
		err := s.db.UserBalanceChange(ctx, newBuy.UserUID, "abt", 9, -newBuy.Amount)
		if err != nil {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUserBalanceChange).Add(err)
		}

		claim := domain.UserClaim{
			UID:          domain.GenUID(12),
			ClaimCode:    "ABT",
			UserUID:      newBuy.UserUID,
			RefUID:       newBuy.RefUID,
			Level:        0,
			CreatedAt:    time.Now().UTC(),
			ClaimedAt:    time.Now().UTC(),
			Amount:       -newBuy.Amount,
			CurrencyCode: newBuy.CurrencyCode,
			Precision:    9,
			Type:         domain.UserClaimTypeBuy,
		}

		err = s.db.UserClaimCreate(ctx, &claim)
		if err != nil {
			errBack := s.db.UserBalanceChange(ctx, newBuy.UserUID, "abt", 9, newBuy.Amount)
			if errBack != nil {
				return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUserBalanceChange).Add(errBack)
			}
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrPayForOrder).Add(err)
		}
	default:
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrCurrencyNotSupported)

	}

	return newBuy, nil
}

// InitPaid to save new buy to database or return existing one
func (s *Service) InitPaid(ctx context.Context, newBuy *domain.Buy) (*domain.Buy, error) {
	exBuy, _ := s.db.BuyGetByUID(ctx, newBuy.UID)
	if exBuy != nil {
		return exBuy, nil
	}
	newBuy.PaidAt = time.Now().UTC()
	newBuy.ApprovedAt = time.Now().UTC()
	err := s.db.BuyCreate(ctx, newBuy)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return newBuy, nil
}

// InitAutofarm to save new buy to database or return existing one
func (s *Service) InitAutofarm(ctx context.Context, userProduct *domain.UserProduct) (*domain.UserProduct, error) {
	err := s.db.UserProductUpdate(ctx, userProduct)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return userProduct, nil
}

// Plan get bought plan for workflow
func (s *Service) Plan(ctx context.Context, newBuy *domain.Buy) (*domain.Plan, error) {
	plan, err := s.db.PlanGetByCode(ctx, newBuy.PlanCode)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return plan, nil
}

// Product get bought plan for workflow
func (s *Service) Product(ctx context.Context, newBuy *domain.Buy) (*domain.Product, error) {
	product, err := s.db.ProductGetByCode(ctx, newBuy.ProductCode)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return product, nil
}

// Paid saves the date when status paid was received
func (s *Service) Paid(ctx context.Context, newBuy *domain.Buy, curTime time.Time) (*domain.Buy, error) {
	newBuy.PaidAt = curTime

	err := s.db.BuyUpdate(ctx, newBuy)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return newBuy, nil
}

// PlanAdd add the bought plan
func (s *Service) PlanAdd(ctx context.Context, newBuy *domain.Buy, plan *domain.Plan) (*domain.UserPlan, error) {
	startDate := newBuy.PaidAt
	exPlan, _ := s.db.UserPlanGetLastByCodeAndDate(ctx, newBuy.UserUID, plan.Code, newBuy.PaidAt)
	// if current plan with this code exists then set the start date as the end date of current plan
	if exPlan != nil && exPlan.EndAt.After(newBuy.PaidAt) {
		startDate = exPlan.EndAt
	}
	endDate := startDate.Add(plan.Period)

	userPlan := &domain.UserPlan{
		UID:      domain.GenUID(12),
		UserUID:  newBuy.UserUID,
		BuyUID:   newBuy.UID,
		PlanCode: newBuy.PlanCode,
		StartAt:  startDate,
		EndAt:    endDate,
		Priority: plan.Priority,
	}

	err := s.db.UserPlanCreate(ctx, userPlan)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return userPlan, nil
}

// ProductAdd add the bought plan
func (s *Service) ProductAdd(ctx context.Context, newBuy *domain.Buy, product *domain.Product) (*domain.UserProduct, error) {
	startDate := newBuy.PaidAt
	exProduct, _ := s.db.UserProductGetLastByCategoryAndDate(ctx, newBuy.UserUID, product.Category, newBuy.PaidAt)
	// if current plan with this code exists then set the start date as the end date of current plan
	if exProduct != nil && exProduct.EndAt.After(newBuy.PaidAt) {
		startDate = exProduct.EndAt
	}
	endDate := startDate.Add(product.Period)

	userProduct := &domain.UserProduct{
		UID:             domain.GenUID(12),
		UserUID:         newBuy.UserUID,
		BuyUID:          newBuy.UID,
		ProductCode:     product.Code,
		ProductCategory: product.Category,
		StartAt:         startDate,
		EndAt:           endDate,
		Priority:        product.Priority,
		Multiplier:      product.Multiplier,
	}

	err := s.db.UserProductCreate(ctx, userProduct)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return userProduct, nil
}

func (s *Service) ProductApply(ctx context.Context, userProduct *domain.UserProduct) (*domain.UserProduct, error) {
	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrConfig).Add(err)
	}
	switch userProduct.ProductCategory {
	case "autofarm":
		err = s.wf.AutofarmCreate(ctx, userProduct)
		if err != nil {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrAutofarmStart).Add(err)
		}

	case "boost_archon":
		archonLevel := domain.GetLevelByCode("archon")

		if archonLevel != nil {
			diff := archonLevel.Balance

			err = s.db.UserBalanceChange(ctx, userProduct.UserUID, "abt", 9, diff)
			if err != nil {
				return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUserBalanceChange).Add(err)
			}

			claim := domain.UserClaim{
				UID:          domain.GenUID(12),
				ClaimCode:    "ABT",
				UserUID:      userProduct.UserUID,
				RefUID:       userProduct.UserUID,
				Level:        0,
				CreatedAt:    time.Now().UTC(),
				ClaimedAt:    time.Now().UTC(),
				Amount:       diff,
				CurrencyCode: "abt",
				Precision:    9,
				Type:         domain.UserClaimTypeBoost,
			}

			err = s.db.UserClaimCreate(ctx, &claim)
			if err != nil {
				errBack := s.db.UserBalanceChange(ctx, userProduct.UserUID, "abt", 9, -diff)
				if errBack != nil {
					return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUserBalanceChange).Add(errBack)
				}
				return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrPayForOrder).Add(err)
			}
		}

	case "safe_6th_digit", "safe_7th_digit":
		dUserSafe, err := s.db.UserSafeGetActiveBySafeUIDsAndUserUID(
			ctx,
			[]string{dConf.RefSafeUID, dConf.Tier1SafeUID, dConf.Tier1CoinSafeUID},
			userProduct.UserUID,
		)
		if dUserSafe == nil {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUserSafeNotFound).Add(err)
		}

		countHacked := domain.CountHackedSafeCode(dUserSafe.Secret)
		if countHacked == 5 || countHacked == 6 {
			dUserSafe.Secret = domain.HackSafeCode(dUserSafe.Code, dUserSafe.Secret)
			if countHacked == 5 {
				if userProduct.ProductCode == "safe_6th_digit_abt" {
					dUserSafe.SafeUID = dConf.Tier1CoinSafeUID
				} else {
					dUserSafe.SafeUID = dConf.Tier1SafeUID
				}
			}
			if countHacked == 6 {
				if userProduct.ProductCode == "safe_7th_digit_abt" {
					dUserSafe.SafeUID = dConf.Tier2CoinSafeUID
				} else {
					dUserSafe.SafeUID = dConf.Tier2SafeUID
				}
			}

			err = s.db.UserSafeUpdate(ctx, dUserSafe)
			if err != nil {
				return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUserSafeUpdate).Add(err)
			}
		}

	default:
		return userProduct, nil
	}
	return userProduct, nil
}

// RankAddFromPlan add rank if bought plan includes rank
func (s *Service) RankAddFromPlan(ctx context.Context, newBuy *domain.Buy, plan *domain.Plan) (*domain.UserRank, error) {
	if plan.RankCode == "" {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrPlanRankEmpty)
	}

	if plan.RankPeriod == 0 {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrPlanRankPeriodEmpty)
	}

	rank, err := s.db.RankGetByCode(ctx, plan.RankCode)
	if rank == nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrPlanRankNotFound)
	}

	startDate := newBuy.PaidAt
	exRank, _ := s.db.UserRankGetByCodeAndDate(ctx, newBuy.UserUID, plan.RankCode, newBuy.PaidAt)
	// if current rank with this code exists then set the start date as the end date of current plan
	if exRank != nil && exRank.EndAt.After(newBuy.PaidAt) {
		startDate = exRank.EndAt
	}
	// End date from plan settings
	endDate := startDate.Add(plan.RankPeriod)

	var row *big.Int
	var col *big.Int
	var matchUid string
	userPlace, _ := s.db.UserPlaceGetByUserUID(ctx, newBuy.UserUID)
	if userPlace != nil {
		row = new(big.Int).Set(userPlace.Row)
		col = new(big.Int).Set(userPlace.Col)
		matchUid = userPlace.MatchUID
	}

	userRank := &domain.UserRank{
		UID:      domain.GenUID(12),
		Type:     domain.UserRankTypeBuy,
		UserUID:  newBuy.UserUID,
		MatchUID: matchUid,
		Row:      row,
		Col:      col,
		BuyUID:   newBuy.UID,
		RankCode: plan.RankCode,
		StartAt:  startDate,
		EndAt:    endDate,
		Priority: rank.Priority,
	}

	err = s.db.UserRankCreate(ctx, userRank)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return userRank, nil
}

func (s *Service) PlanEnd(ctx context.Context, newBuy *domain.Buy, userPlan *domain.UserPlan) (*domain.UserPlan, error) {
	exPlan, err := s.db.UserPlanGetByUID(ctx, userPlan.UID)
	if exPlan != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	exPlan.EndAt = time.Now().UTC()
	if !newBuy.RefundedAt.IsZero() {
		exPlan.EndAt = newBuy.RefundedAt
	}

	err = s.db.UserPlanUpdate(ctx, exPlan)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return exPlan, nil
}

func (s *Service) RankEnd(ctx context.Context, newBuy *domain.Buy, userRank *domain.UserRank) (*domain.UserRank, error) {
	exRank, err := s.db.UserRankGetByUID(ctx, userRank.UID)
	if exRank != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	exRank.EndAt = time.Now().UTC()
	if !newBuy.RefundedAt.IsZero() {
		exRank.EndAt = newBuy.RefundedAt
	}

	err = s.db.UserRankUpdate(ctx, exRank)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return exRank, nil
}

// Refund charge back the amount of the plan and refund the amount to user balance
func (s *Service) Refund(ctx context.Context, newBuy *domain.Buy, curTime time.Time) (*domain.Buy, error) {
	newBuy.RefundedAt = curTime

	// TODO: refund transaction

	err := s.db.BuyUpdate(ctx, newBuy)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return newBuy, nil
}

// Approved saves the date when the buy was approved after charge waiting time
func (s *Service) Approved(ctx context.Context, newBuy *domain.Buy, curTime time.Time) (*domain.Buy, error) {
	newBuy.ApprovedAt = curTime

	err := s.db.BuyUpdate(ctx, newBuy)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return newBuy, nil
}

// BuySetRowCol saves the row and col to buy for optimization of calculations
func (s *Service) BuySetRowCol(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace) (*domain.Buy, error) {
	dUser, err := s.db.UserGetByUID(ctx, newBuy.UserUID)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	newBuy.RefUID = dUser.RefUID
	newBuy.MatchUID = place.MatchUID
	newBuy.Row = place.Row
	newBuy.Col = place.Col
	newBuy.Type = domain.BuyTypeDistributor

	err = s.db.BuyUpdate(ctx, newBuy)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return newBuy, nil
}

// Charged saves the date when all charges made
func (s *Service) Charged(ctx context.Context, newBuy *domain.Buy, curTime time.Time) (*domain.Buy, error) {
	newBuy.ChargedAt = curTime

	err := s.db.BuyUpdate(ctx, newBuy)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return newBuy, nil
}

// Cancelled saves the date of buy cancel if paid stage is not reached in time
func (s *Service) Cancelled(ctx context.Context, newBuy *domain.Buy, curTime time.Time) (*domain.Buy, error) {
	newBuy.CancelledAt = curTime

	err := s.db.BuyUpdate(ctx, newBuy)
	if err != nil {
		return nil, err
	}

	return newBuy, nil
}
