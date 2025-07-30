package user

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"server/internal/domain"
)

func (s *Service) ClaimCreatePartner(ctx context.Context, code string, partnerCode string, userUid string, amount string) error {
	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return err
	}

	claim, err := s.db.ClaimGetByCode(ctx, code)
	if err != nil {
		return err
	}

	user, err := s.db.UserGetByUID(ctx, userUid)
	if err != nil {
		return err
	}

	claimAmount, err := decimal.NewFromString(amount)
	if err != nil {
		return err
	}

	claimAmountInt := claimAmount.Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(claim.Precision)))).IntPart()

	exUserClaim, _ := s.db.UserClaimGetByUserUIDAndClaimCodeAndTypeAndPartnerCode(ctx, userUid, code, domain.UserClaimTypePartner, partnerCode)
	if exUserClaim != nil {
		exUserClaim.Amount += claimAmountInt
		exUserClaim.ClaimedAt = time.Now().UTC()

		err = s.db.UserClaimUpdate(ctx, exUserClaim)
		if err != nil {
			return err
		}
	} else {
		userClaim := &domain.UserClaim{
			UID:          domain.GenUID(12),
			ClaimCode:    claim.Code,
			UserUID:      user.UID,
			RefUID:       user.RefUID,
			Level:        0,
			CreatedAt:    time.Now().UTC(),
			ClaimedAt:    time.Now().UTC(),
			Amount:       claimAmountInt,
			CurrencyCode: claim.CurrencyCode,
			Precision:    claim.Precision,
			Type:         domain.UserClaimTypePartner,
			PartnerCode:  partnerCode,
		}

		err = s.db.UserClaimCreate(ctx, userClaim)
		if err != nil {
			return err
		}
	}

	bal, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, userUid, claim.CurrencyCode)
	if err != nil {
		if !domain.ErrorIs(err, domain.ErrNoDocuments) {
			return err
		}
		bal = &domain.UserBalance{
			UserUID:      userUid,
			CurrencyCode: claim.CurrencyCode,
			Precision:    claim.Precision,
			Amount:       0,
		}
		err = s.db.UserBalanceCreate(ctx, bal)
		if err != nil {
			return err
		}
	}
	bal.Amount += claimAmountInt
	err = s.db.UserBalanceUpdate(ctx, bal)
	if err != nil {
		return err
	}

	go func() {
		// Charge coins
		err = s.auth.ChargeLineCoin(ctx, claim, user.RefUID, claimAmountInt, dConf.DefaultRefUid, 0, dConf.CoinRefLinePercent)
		if err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func (s *Service) ClaimGet(ctx context.Context, code string, user *domain.User) (*domain.Claim, error) {
	claim, err := s.db.ClaimGetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	userLevel, err := s.auth.UserLevel(ctx, user)
	if userLevel.ClaimAmount > 0 {
		claim.Amount = userLevel.ClaimAmount
	}

	userPremium, err := s.db.UserProductGetByCategoryAndDate(ctx, user.UID, "premium", time.Now().UTC())
	if userPremium != nil {
		switch userPremium.ProductCode {
		case "premium_month_usd", "premium_year_usd":
			claim.Amount = claim.Amount * 5
		case "premium_month_abt", "premium_year_abt":
			claim.Amount = claim.Amount * 3
		}
	}

	boost, _ := s.db.UserProductGetLastBoostAndDate(ctx, user.UID, time.Now().UTC())
	if boost != nil {
		switch boost.ProductCategory {
		case "boost_x3", "boost_x5", "boost_x10", "boost_x20", "boost_x50":
			claim.Amount = claim.Amount * boost.Multiplier
		}
	}

	return claim, nil
}

func (s *Service) UserClaimCreate(ctx context.Context, code string, userUid string) (*domain.UserClaim, error) {
	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, err
	}

	dUser, err := s.db.UserGetByUID(ctx, userUid)
	if err != nil {
		return nil, err
	}

	if dUser.RefUID == dConf.DefaultRefUid {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrClaimNotAvailable)
	}

	claim, err := s.db.ClaimGetByCode(ctx, code)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	lastClaim, err := s.db.UserClaimGetLastByClaimCodeAndTypeAndUserUID(ctx, claim.Code, domain.UserClaimTypeClaim, userUid)
	if err != nil {
		if domain.ErrorIs(err, domain.ErrNoDocuments) {
			lastClaim = &domain.UserClaim{
				UID:          domain.GenUID(12),
				ClaimCode:    claim.Code,
				UserUID:      dUser.UID,
				RefUID:       dUser.RefUID,
				Level:        0,
				CreatedAt:    time.Now().UTC(),
				ClaimedAt:    time.Now().UTC().Add(-2 * time.Hour),
				Amount:       0,
				CurrencyCode: claim.CurrencyCode,
				Precision:    claim.Precision,
				Type:         domain.UserClaimTypeClaim,
			}
		} else {
			return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
		}
	}

	userLevel, err := s.auth.UserLevel(ctx, dUser)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	if time.Since(lastClaim.ClaimedAt) < claim.MinPeriod {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrClaimMinPeriod)
	}

	period := time.Now().UTC().Sub(lastClaim.ClaimedAt)
	if period > claim.MaxPeriod {
		period = claim.MaxPeriod
	}
	claimAmount := claim.Amount
	if userLevel.ClaimAmount > 0 {
		claimAmount = userLevel.ClaimAmount
	}
	percent := decimal.NewFromInt(100).Mul(decimal.NewFromInt(int64(period))).Div(decimal.NewFromInt(int64(claim.MaxPeriod)))
	amount := decimal.NewFromInt(claimAmount).Mul(percent).Div(decimal.NewFromInt(100)).IntPart()

	premium, _ := s.db.UserProductGetByCategoryAndDate(ctx, dUser.UID, "premium", time.Now().UTC())
	if premium != nil {
		switch premium.ProductCode {
		case "premium_month_usd", "premium_year_usd":
			amount = amount * 5
		case "premium_month_abt", "premium_year_abt":
			amount = amount * 3
		}
	}

	boost, _ := s.db.UserProductGetLastBoostAndDate(ctx, dUser.UID, time.Now().UTC())
	if boost != nil {
		switch boost.ProductCategory {
		case "boost_x3", "boost_x5", "boost_x10", "boost_x20", "boost_x50":
			amount = amount * boost.Multiplier
		}
	}

	userClaim := &domain.UserClaim{
		UID:          lastClaim.UID,
		ClaimCode:    lastClaim.ClaimCode,
		UserUID:      lastClaim.UserUID,
		RefUID:       lastClaim.RefUID,
		Level:        lastClaim.Level,
		CreatedAt:    lastClaim.CreatedAt,
		ClaimedAt:    time.Now().UTC(),
		Amount:       lastClaim.Amount + amount,
		CurrencyCode: lastClaim.CurrencyCode,
		Precision:    lastClaim.Precision,
		Type:         lastClaim.Type,
	}
	err = s.db.UserClaimUpsert(ctx, userClaim)
	if err != nil {
		return nil, err
	}

	bal, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, userUid, claim.CurrencyCode)
	if err != nil {
		if !domain.ErrorIs(err, domain.ErrNoDocuments) {
			return nil, err
		}
		bal = &domain.UserBalance{
			UserUID:      userUid,
			CurrencyCode: claim.CurrencyCode,
			Precision:    claim.Precision,
			Amount:       0,
		}
		err = s.db.UserBalanceCreate(ctx, bal)
		if err != nil {
			return nil, err
		}
	}
	bal.Amount += amount
	err = s.db.UserBalanceUpdate(ctx, bal)
	if err != nil {
		return nil, err
	}

	go func() {
		// Charge coins
		err = s.auth.ChargeLineCoin(ctx, claim, dUser.RefUID, amount, dConf.DefaultRefUid, 0, dConf.CoinRefLinePercent)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// fix amount for notification
	userClaim.Amount = amount

	return userClaim, nil
}

func (s *Service) ClaimBalance(ctx context.Context, code string, userUid string) (int64, error) {
	claim, err := s.db.ClaimGetByCode(ctx, code)
	if err != nil {
		return 0, err
	}

	bal, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, userUid, claim.CurrencyCode)
	if err != nil {
		if !domain.ErrorIs(err, domain.ErrNoDocuments) {
			return 0, err
		}
		bal = &domain.UserBalance{
			UserUID:      userUid,
			CurrencyCode: claim.CurrencyCode,
			Precision:    claim.Precision,
			Amount:       0,
		}
	}

	return bal.Amount, nil
}

func (s *Service) UserClaimGetLast(ctx context.Context, code string, userUid string) (*domain.UserClaim, error) {
	userClaim, err := s.db.UserClaimGetLastByClaimCodeAndTypeAndUserUID(ctx, code, domain.UserClaimTypeClaim, userUid)
	if err != nil {
		if domain.ErrorIs(err, domain.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return userClaim, nil
}
