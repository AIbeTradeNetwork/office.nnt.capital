package user

import (
	"context"
	"server/internal/domain"
	"time"
)

func (s *Service) GetSafe(ctx context.Context, user *domain.User) (*domain.UserSafe, error) {
	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	safeUids := []string{dConf.RefSafeUID, dConf.Tier1SafeUID, dConf.Tier2SafeUID, dConf.Tier1CoinSafeUID, dConf.Tier2CoinSafeUID}
	userSafe, err := s.db.UserSafeGetActiveBySafeUIDsAndUserUID(ctx, safeUids, user.UID)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrUserSafeNotFound).Add(err)
	}

	return userSafe, nil
}

func (s *Service) HackSafe(ctx context.Context, user *domain.User, code string) (*domain.UserSafe, error) {
	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	safeUids := []string{dConf.RefSafeUID, dConf.Tier1SafeUID, dConf.Tier2SafeUID, dConf.Tier1CoinSafeUID, dConf.Tier2CoinSafeUID}
	userSafe, err := s.db.UserSafeGetActiveBySafeUIDsAndUserUID(ctx, safeUids, user.UID)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrUserSafeNotFound).Add(err)
	}

	if userSafe.Code != code {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrSafeNotHacked).Add(err)
	}

	safe, err := s.db.SafeGetByUID(ctx, userSafe.SafeUID)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrSafeNotFound).Add(err)
	}

	variant, err := s.RandSafeVariant(ctx, safe.Variants)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrSafeVariantLimit).Add(err)
	}

	err = s.db.SafeUpdate(ctx, safe)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrSafeUpdate).Add(err)
	}

	err = s.ChargeSafeVariant(ctx, user, variant)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrSafeVariant).Add(err)
	}

	userSafe.ClaimedAt = time.Now().UTC()
	userSafe.Variant = variant

	err = s.db.UserSafeUpdate(ctx, userSafe)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrUserSafeUpdate).Add(err)
	}

	return userSafe, nil
}

func (s *Service) RandSafeVariant(ctx context.Context, safeVariants []domain.SafeVariant) (*domain.SafeVariant, error) {
	lastVariant := domain.SafeVariant{}
	lastI := 0
	for i := 0; i < len(safeVariants); i++ {
		safeVariant := safeVariants[i]
		if safeVariant.LimitCount > 0 && safeVariant.Count >= safeVariant.LimitCount {
			continue
		}
		if safeVariant.Chance > lastVariant.Chance {
			lastVariant = safeVariant
			lastI = i
		}

		rand := int64(domain.GenRandomInt(0, 10000))

		if rand < safeVariant.Chance {
			safeVariant.Count += 1
			safeVariants[i] = safeVariant

			if safeVariant.Type == domain.SafeVariantTypeVariants {
				return s.RandSafeVariant(ctx, safeVariant.Variants)
			}

			return &safeVariant, nil
		}
	}

	lastVariant.Count += 1
	safeVariants[lastI] = lastVariant
	if lastVariant.Type == domain.SafeVariantTypeVariants {
		return s.RandSafeVariant(ctx, lastVariant.Variants)
	}

	if lastVariant.LimitCount > 0 && lastVariant.Count >= lastVariant.LimitCount {
		return s.RandSafeVariant(ctx, safeVariants)
	}

	return &lastVariant, nil
}

func (s *Service) ChargeSafeVariant(ctx context.Context, user *domain.User, safeVariant *domain.SafeVariant) error {
	switch safeVariant.Type {
	case domain.SafeVariantTypeUDEX:
		claim := &domain.UserClaim{
			UID:          domain.GenUID(12),
			ClaimCode:    "UDEX",
			UserUID:      user.UID,
			RefUID:       user.RefUID,
			CreatedAt:    time.Now().UTC(),
			ClaimedAt:    time.Now().UTC(),
			Amount:       safeVariant.Amount,
			CurrencyCode: "UDEX",
			Precision:    9,
			Type:         domain.UserClaimTypeSafe,
		}

		err := s.db.UserBalanceChange(ctx, user.UID, claim.CurrencyCode, claim.Precision, claim.Amount)
		if err != nil {
			return domain.NewError(userErrorSource).SetCode(domain.ErrUserBalanceChange).Add(err)
		}

		err = s.db.UserClaimCreate(ctx, claim)
		if err != nil {
			errBack := s.db.UserBalanceChange(ctx, user.UID, claim.CurrencyCode, claim.Precision, -claim.Amount)
			if errBack != nil {
				return domain.NewError(userErrorSource).SetCode(domain.ErrUserBalanceChange).Add(errBack)
			}
			return domain.NewError(userErrorSource).SetCode(domain.ErrUserBalanceChange).Add(err)
		}

	case domain.SafeVariantTypeUsd:
		trans := &domain.Transaction{
			UID:         domain.GenUID(12),
			UserUID:     user.UID,
			FromUID:     user.UID,
			Type:        domain.TransactionTypeSafe,
			Amount:      safeVariant.Amount,
			PosAmount:   safeVariant.Amount,
			FullAmount:  safeVariant.Amount,
			Coefficient: 100,
			CreatedAt:   time.Now().UTC(),
			ChargedAt:   time.Now().UTC(),
			MsgCodes:    nil,
		}

		err := s.db.TransactionCreate(ctx, trans)
		if err != nil {
			return domain.NewError(userErrorSource).SetCode(domain.ErrSafeVariantTransaction).Add(err)
		}

	default:
		return domain.NewError(userErrorSource).SetCode(domain.ErrSafeVariantType)

	}

	return nil
}
