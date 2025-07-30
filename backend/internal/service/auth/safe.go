package auth

import (
	"context"
	"server/internal/domain"
	"time"
)

func (s *Service) CheckSafeCodes(ctx context.Context, conf *domain.Config, userUid string) error {
	ctx = context.Background()

	if conf.RefSafeUID == "" {
		return nil
	}

	safeUids := []string{conf.RefSafeUID, conf.Tier1SafeUID, conf.Tier2SafeUID, conf.Tier1CoinSafeUID, conf.Tier2CoinSafeUID}
	safeCode, err := s.db.UserSafeGetActiveBySafeUIDsAndUserUID(ctx, safeUids, userUid)
	if err != nil {
		if domain.ErrorIs(err, domain.ErrNoDocuments) {
			safeCode, err = s.GenerateNewSafeCode(ctx, conf.RefSafeUID, userUid)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if domain.CountHackedSafeCode(safeCode.Secret) >= 5 {
		return nil
	}

	safeCode.Secret = domain.HackSafeCode(safeCode.Code, safeCode.Secret)

	err = s.db.UserSafeUpdateSecret(ctx, safeCode.UID, safeCode.Secret)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GenerateNewSafeCode(ctx context.Context, safeUid string, userUid string) (*domain.UserSafe, error) {
	safeCode := &domain.UserSafe{
		UID:       domain.GenUID(12),
		UserUID:   userUid,
		SafeUID:   safeUid,
		Code:      domain.GenSafeCode(8),
		Secret:    domain.GenSafeSecret(8),
		CreatedAt: time.Now().UTC(),
		ClaimedAt: time.Time{},
		Variant:   nil,
	}

	err := s.db.UserSafeCreate(ctx, safeCode)
	if err != nil {
		return nil, err
	}
	return safeCode, nil
}
