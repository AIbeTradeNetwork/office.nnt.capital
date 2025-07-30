package user

import (
	"context"
	"server/internal/domain"
	"time"
)

func (s *Service) ProcessPrises(ctx context.Context, priseReq *domain.PriseReq) (map[string]string, error) {

	results := make(map[string]string, len(priseReq.ToUserUID))

	for _, userUid := range priseReq.ToUserUID {
		results[userUid] = "ok"
		switch priseReq.Type {
		case domain.PriseTypeUDEX:
			trans := &domain.Transaction{
				UID:         domain.GenUID(12),
				UserUID:     userUid,
				FromUID:     userUid,
				Percent:     0,
				Level:       0,
				Type:        domain.TransactionTypePrise,
				RankCode:    "",
				TaskCode:    "",
				Amount:      priseReq.Amount,
				PosAmount:   priseReq.Amount,
				FullAmount:  priseReq.Amount,
				Coefficient: 100,
				BuyUID:      "",
				PayoutUID:   "",
				DepositUID:  "",
				CreatedAt:   time.Now().UTC(),
				ChargedAt:   time.Now().UTC(),
				MsgCodes:    nil,
			}

			err := s.db.TransactionCreate(ctx, trans)
			if err != nil {
				results[userUid] = err.Error()
				continue
			}

		case domain.PriseTypePremiumMonth:
			startAt := time.Now().UTC()
			userProductEx, err := s.db.UserProductGetLastByCategoryAndDate(ctx, userUid, "premium", time.Now().UTC())
			if userProductEx != nil && userProductEx.EndAt.After(startAt) {
				startAt = userProductEx.EndAt
			}

			userProduct := &domain.UserProduct{
				UID:             domain.GenUID(12),
				UserUID:         userUid,
				BuyUID:          "",
				ProductCode:     "premium_month_udex",
				ProductCategory: "premium",
				StartAt:         startAt,
				EndAt:           startAt.Add(time.Hour * 24 * 30),
				Priority:        100,
			}

			err = s.db.UserProductCreate(ctx, userProduct)
			if err != nil {
				results[userUid] = err.Error()
				continue
			}

		case domain.PriseTypePremiumYear:
			startAt := time.Now().UTC()
			userProductEx, err := s.db.UserProductGetLastByCategoryAndDate(ctx, userUid, "premium", time.Now().UTC())
			if userProductEx != nil && userProductEx.EndAt.After(startAt) {
				startAt = userProductEx.EndAt
			}

			userProduct := &domain.UserProduct{
				UID:             domain.GenUID(12),
				UserUID:         userUid,
				BuyUID:          "",
				ProductCode:     "premium_year_udex",
				ProductCategory: "premium",
				StartAt:         startAt,
				EndAt:           startAt.Add(time.Hour * 24 * 365),
				Priority:        200,
			}

			err = s.db.UserProductCreate(ctx, userProduct)
			if err != nil {
				results[userUid] = err.Error()
				continue
			}

		case domain.PriseTypePremiumMonthUdex:
			startAt := time.Now().UTC()
			userProductEx, err := s.db.UserProductGetLastByCategoryAndDate(ctx, userUid, "premium", time.Now().UTC())
			if userProductEx != nil && userProductEx.EndAt.After(startAt) {
				startAt = userProductEx.EndAt
			}

			userProduct := &domain.UserProduct{
				UID:             domain.GenUID(12),
				UserUID:         userUid,
				BuyUID:          "",
				ProductCode:     "premium_month_usd",
				ProductCategory: "premium",
				StartAt:         startAt,
				EndAt:           startAt.Add(time.Hour * 24 * 30),
				Priority:        100,
			}

			err = s.db.UserProductCreate(ctx, userProduct)
			if err != nil {
				results[userUid] = err.Error()
				continue
			}

		case domain.PriseTypePremiumYearUdex:
			startAt := time.Now().UTC()
			userProductEx, err := s.db.UserProductGetLastByCategoryAndDate(ctx, userUid, "premium", time.Now().UTC())
			if userProductEx != nil && userProductEx.EndAt.After(startAt) {
				startAt = userProductEx.EndAt
			}

			userProduct := &domain.UserProduct{
				UID:             domain.GenUID(12),
				UserUID:         userUid,
				BuyUID:          "",
				ProductCode:     "premium_year_usd",
				ProductCategory: "premium",
				StartAt:         startAt,
				EndAt:           startAt.Add(time.Hour * 24 * 365),
				Priority:        200,
			}

			err = s.db.UserProductCreate(ctx, userProduct)
			if err != nil {
				results[userUid] = err.Error()
				continue
			}

		case domain.PriseTypeAutofarmMonth:
			startAt := time.Now().UTC()
			userProductEx, err := s.db.UserProductGetLastByCategoryAndDate(ctx, userUid, "autofarm", time.Now().UTC())
			if userProductEx != nil && userProductEx.EndAt.After(startAt) {
				startAt = userProductEx.EndAt
			}

			userProduct := &domain.UserProduct{
				UID:             domain.GenUID(12),
				UserUID:         userUid,
				BuyUID:          "",
				ProductCode:     "autofarm_month_usd",
				ProductCategory: "autofarm",
				StartAt:         startAt,
				EndAt:           startAt.Add(time.Hour * 24 * 30),
				Priority:        300,
				Multiplier:      0,
			}

			err = s.db.UserProductCreate(ctx, userProduct)
			if err != nil {
				results[userUid] = err.Error()
				continue
			}

			err = s.wf.AutofarmCreate(ctx, userProduct)
			if err != nil {
				results[userUid] = err.Error()
				continue
			}

		case domain.PriseTypeAutofarmYear:
			startAt := time.Now().UTC()
			userProductEx, err := s.db.UserProductGetLastByCategoryAndDate(ctx, userUid, "autofarm", time.Now().UTC())
			if userProductEx != nil && userProductEx.EndAt.After(startAt) {
				startAt = userProductEx.EndAt
			}

			userProduct := &domain.UserProduct{
				UID:             domain.GenUID(12),
				UserUID:         userUid,
				BuyUID:          "",
				ProductCode:     "autofarm_year_usd",
				ProductCategory: "autofarm",
				StartAt:         startAt,
				EndAt:           startAt.Add(time.Hour * 24 * 365),
				Priority:        300,
				Multiplier:      0,
			}

			err = s.db.UserProductCreate(ctx, userProduct)
			if err != nil {
				results[userUid] = err.Error()
				continue
			}

			err = s.wf.AutofarmCreate(ctx, userProduct)
			if err != nil {
				results[userUid] = err.Error()
				continue
			}

		default:
			results[userUid] = "prise type unknown"
			continue
		}
	}

	return results, nil
}
