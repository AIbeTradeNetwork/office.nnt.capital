package buy

import (
	"context"
	"server/internal/domain"
	"time"
)

func (s *Service) GetRefUsersUp(ctx context.Context, userUid string, depth int) ([]*domain.User, error) {
	return s.db.UserGetAllUpByUID(ctx, userUid, depth)
}

func (s *Service) RefUserCharge(ctx context.Context, newBuy *domain.Buy, user *domain.User, level int) (int64, error) {
	userLevel, err := s.auth.UserLevel(ctx, user)
	if err != nil {
		return 0, domain.NewError(buyErrorSource).SetCode(domain.ErrUserLevel).Add(err)
	}

	var amount int64 = 0

	if len(userLevel.LevelBonus) > level {
		amount = domain.Percent(newBuy.Cv, userLevel.LevelBonus[level])
	}

	if amount == 0 {
		return 0, nil
	}

	switch newBuy.CurrencyCode {
	case "usd":
		trans := &domain.Transaction{
			UID:         domain.GenUID(12),
			UserUID:     user.UID,
			FromUID:     newBuy.UserUID,
			Percent:     userLevel.LevelBonus[level],
			Level:       level,
			Type:        domain.TransactionTypeRefBuy,
			RankCode:    "",
			TaskCode:    "",
			Amount:      amount,
			PosAmount:   amount,
			FullAmount:  amount,
			Coefficient: 100,
			BuyUID:      newBuy.UID,
			PayoutUID:   "",
			DepositUID:  "",
			CreatedAt:   time.Now().UTC(),
			ChargedAt:   time.Now().UTC(),
			MsgCodes:    nil,
		}

		err = s.db.TransactionCreate(ctx, trans)
		if err != nil {
			return 0, domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
		}
	case "abt":
		claim := &domain.UserClaim{
			UID:          domain.GenUID(12),
			ClaimCode:    "ABT",
			UserUID:      user.UID,
			RefUID:       newBuy.UserUID,
			Level:        level,
			CreatedAt:    time.Now().UTC(),
			ClaimedAt:    time.Now().UTC(),
			Amount:       amount,
			CurrencyCode: "abt",
			Precision:    9,
			Type:         domain.UserClaimTypeRefBuy,
			TaskCode:     "",
		}

		err = s.db.UserBalanceChange(ctx, user.UID, claim.CurrencyCode, claim.Precision, claim.Amount)
		if err != nil {
			return 0, domain.NewError(buyErrorSource).SetCode(domain.ErrUserBalanceChange).Add(err)
		}

		err = s.db.UserClaimCreate(ctx, claim)
		if err != nil {
			errBack := s.db.UserBalanceChange(ctx, user.UID, claim.CurrencyCode, claim.Precision, -claim.Amount)
			if errBack != nil {
				return 0, domain.NewError(buyErrorSource).SetCode(domain.ErrUserBalanceRollback).Add(err)
			}
			return 0, domain.NewError(buyErrorSource).SetCode(domain.ErrUserBalanceChange).Add(err)
		}
	}

	return amount, nil
}
