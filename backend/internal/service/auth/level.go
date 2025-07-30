package auth

import (
	"context"
	"server/internal/domain"
)

func (s *Service) UserLevel(ctx context.Context, user *domain.User) (*domain.UserLevel, error) {
	var bal int64
	userBalance, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, user.UID, "UDEX")
	if userBalance != nil && err == nil {
		bal = userBalance.Amount
	}

	curLevel := domain.UserLevels[0]
	for _, l := range domain.UserLevels {
		if bal >= l.Balance && curLevel.Balance < l.Balance {
			curLevel = l
		}
	}
	return curLevel, nil
}

func (s *Service) Levels(ctx context.Context) ([]*domain.UserLevel, error) {
	return domain.UserLevels, nil
}
