package auth

import (
	"context"
	"server/internal/domain"
	"time"
)

func (s *Service) ChargeRefCoins(ctx context.Context, uid string) error {
	ctx = context.Background()

	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return err
	}

	dUser, err := s.db.UserGetByUID(ctx, uid)
	if err != nil {
		return err
	}

	if uid == dConf.DefaultRefUid {
		return nil
	}

	if dConf.CoinRefBonus <= 0 || dConf.CoinCode == "" {
		return nil
	}

	claim, err := s.db.ClaimGetByCode(ctx, dConf.CoinCode)
	if err != nil {
		return err
	}

	amount := dConf.CoinRefBonus
	boost, _ := s.db.UserProductGetLastBoostAndDate(ctx, uid, time.Now().UTC())
	if boost != nil {
		switch boost.ProductCategory {
		case "boost_x3", "boost_x5", "boost_x10", "boost_x20", "boost_x50":
			amount = amount * boost.Multiplier
		}
	}

	userClaim := &domain.UserClaim{
		UID:          domain.GenUID(12),
		ClaimCode:    claim.Code,
		UserUID:      uid,
		RefUID:       dUser.RefUID,
		Level:        0,
		CreatedAt:    time.Now().UTC(),
		ClaimedAt:    time.UnixMilli(0),
		Amount:       amount,
		CurrencyCode: claim.CurrencyCode,
		Precision:    claim.Precision,
		Type:         domain.UserClaimTypeInvite,
	}
	err = s.db.UserClaimCreate(ctx, userClaim)
	if err != nil {
		return err
	}

	bal, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, uid, claim.CurrencyCode)
	if err != nil {
		if !domain.ErrorIs(err, domain.ErrNoDocuments) {
			return err
		}
		bal = &domain.UserBalance{
			UserUID:      uid,
			CurrencyCode: claim.CurrencyCode,
			Precision:    claim.Precision,
			Amount:       0,
		}
		err = s.db.UserBalanceCreate(ctx, bal)
		if err != nil {
			return err
		}
	}
	bal.Amount += amount
	err = s.db.UserBalanceUpdate(ctx, bal)
	if err != nil {
		return err
	}

	//err = s.ChargeLineCoin(ctx, claim, dUser.RefUID, dConf.CoinRefBonus, dConf.DefaultRefUid, 0, dConf.CoinRefLinePercent)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (s *Service) ChargeToRefCoins(ctx context.Context, fromUid string, toUid string) error {
	ctx = context.Background()

	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return err
	}

	dUser, err := s.db.UserGetByUID(ctx, toUid)
	if err != nil {
		return err
	}

	if dUser.RefUID == dConf.DefaultRefUid {
		return nil
	}

	if dConf.CoinToRefBonus <= 0 || dConf.CoinCode == "" {
		return nil
	}

	claim, err := s.db.ClaimGetByCode(ctx, dConf.CoinCode)
	if err != nil {
		return err
	}

	amount := dConf.CoinToRefBonus
	boost, _ := s.db.UserProductGetLastBoostAndDate(ctx, fromUid, time.Now().UTC())
	if boost != nil {
		switch boost.ProductCategory {
		case "boost_x3", "boost_x5", "boost_x10", "boost_x20", "boost_x50":
			amount = amount * boost.Multiplier
		}
	}

	userClaim := &domain.UserClaim{
		UID:          domain.GenUID(12),
		ClaimCode:    claim.Code,
		UserUID:      toUid,
		RefUID:       fromUid,
		Level:        0,
		CreatedAt:    time.Now().UTC(),
		ClaimedAt:    time.UnixMilli(0),
		Amount:       amount,
		CurrencyCode: claim.CurrencyCode,
		Precision:    claim.Precision,
		Type:         domain.UserClaimTypeRefBonus,
	}
	err = s.db.UserClaimCreate(ctx, userClaim)
	if err != nil {
		return err
	}

	bal, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, toUid, claim.CurrencyCode)
	if err != nil {
		if !domain.ErrorIs(err, domain.ErrNoDocuments) {
			return err
		}
		bal = &domain.UserBalance{
			UserUID:      toUid,
			CurrencyCode: claim.CurrencyCode,
			Precision:    claim.Precision,
			Amount:       0,
		}
		err = s.db.UserBalanceCreate(ctx, bal)
		if err != nil {
			return err
		}
	}
	bal.Amount += amount
	err = s.db.UserBalanceUpdate(ctx, bal)
	if err != nil {
		return err
	}

	//err = s.ChargeLineCoin(ctx, claim, dUser.RefUID, dConf.CoinRefBonus, dConf.DefaultRefUid, 0, dConf.CoinRefLinePercent)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (s *Service) ChargeLineCoin(ctx context.Context, claim *domain.Claim, uid string, amount int64, defaultRefUid string, count int, lines []int64) error {
	if uid == "" {
		return nil
	}

	if count >= len(lines) {
		return nil
	}

	ctx = context.Background()

	dUser, err := s.db.UserGetByUID(ctx, uid)
	if err != nil {
		return err
	}

	bonusAmount := amount * lines[count] / 10000
	if bonusAmount <= 0 {
		return nil
	}

	userClaim := &domain.UserClaim{
		UID:          domain.GenUID(12),
		ClaimCode:    claim.Code,
		UserUID:      dUser.UID,
		RefUID:       dUser.RefUID,
		Level:        count + 1,
		CreatedAt:    time.Now().UTC(),
		ClaimedAt:    time.UnixMilli(0),
		Amount:       bonusAmount,
		CurrencyCode: claim.CurrencyCode,
		Precision:    claim.Precision,
		Type:         domain.UserClaimTypeRef,
	}
	err = s.db.UserClaimCreate(ctx, userClaim)
	if err != nil {
		return err
	}

	bal, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, uid, claim.CurrencyCode)
	if err != nil {
		if !domain.ErrorIs(err, domain.ErrNoDocuments) {
			return err
		}
		bal = &domain.UserBalance{
			UserUID:      uid,
			CurrencyCode: claim.CurrencyCode,
			Precision:    claim.Precision,
			Amount:       0,
		}
		err = s.db.UserBalanceCreate(ctx, bal)
		if err != nil {
			return err
		}
	}
	bal.Amount += bonusAmount
	err = s.db.UserBalanceUpdate(ctx, bal)
	if err != nil {
		return err
	}

	if dUser.RefUID == "" || dUser.RefUID == defaultRefUid {
		return nil
	}

	count++

	return s.ChargeLineCoin(ctx, claim, dUser.RefUID, amount, defaultRefUid, count, lines)
}
