package auth

import (
	"context"
	"server/internal/domain"
	"time"
)

func (s *Service) UnlimInvite(ctx context.Context, cfg *domain.Config, userUid string) bool {
	if cfg.UnlimInvite {
		return true
	}

	refUser, err := s.db.UserGetByUID(ctx, userUid)
	if err != nil {
		return false
	}
	if refUser.UnlimInvite {
		return true
	}

	// for distributor unlimited invite
	refPlace, _ := s.db.UserPlaceGetByUserUID(ctx, userUid)
	if refPlace != nil {
		return true
	}

	userLevel, err := s.UserLevel(ctx, refUser)
	if err != nil {
		return false
	}
	refUserCount, err := s.db.UserCountByRefUID(ctx, userUid)
	if err != nil {
		return false
	}
	limit := userLevel.InviteLimit
	premium, _ := s.db.UserProductGetByCategoryAndDate(ctx, userUid, "premium", time.Now().UTC())
	if premium != nil {
		switch premium.ProductCode {
		case "premium_year_abt", "premium_year_usd":
			return true
		case "premium_month_abt", "premium_month_usd":
			limit += 150
		}
	}

	if refUserCount < limit {
		return true
	}

	return false
}
