package user

import (
	"context"
	"log/slog"
	"server/internal/domain"
	"time"
)

func (s *Service) ComboList(ctx context.Context) ([]*domain.Combo, error) {
	return s.db.ComboGetAll(ctx)
}

func (s *Service) Combo(ctx context.Context, user *domain.User, code string) (*domain.Combo, error) {
	slog.Info("Combo request", "user", user.UID, "code", code)

	combo, err := s.db.ComboGetByCode(ctx, code)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrComboNotFound).Add(err)
	}

	if combo.IsActive == false {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrComboNotFound)
	}

	if combo.StartAt.After(time.Now()) {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrComboNotFound)
	}

	if combo.EndAt.Before(time.Now()) {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrComboNotFound)
	}

	if combo.Amount > 0 {
		switch combo.CurrencyCode {
		case "usd":
			exReward, err := s.db.TransactionGetByUserUIDAndTypeAndComboCode(ctx, user.UID, domain.TransactionTypeCombo, combo.UID)
			if exReward != nil {
				return nil, domain.NewError(userErrorSource).SetCode(domain.ErrComboAlreadyUsed).Add(err)
			}
		case "UDEX":
			exReward, err := s.db.UserClaimGetByUserUIDAndClaimCodeAndTypeAndComboCode(ctx, user.UID, "UDEX", domain.UserClaimTypeCombo, combo.UID)
			if exReward != nil {
				return nil, domain.NewError(userErrorSource).SetCode(domain.ErrComboAlreadyUsed).Add(err)
			}
		}
	}

	if combo.PriseCode != "" {
		exPrise, err := s.db.UserProductGetByComboUID(ctx, user.UID, combo.UID)
		if exPrise != nil {
			return nil, domain.NewError(userErrorSource).SetCode(domain.ErrComboAlreadyUsed).Add(err)
		}
	}

	if combo.Count >= combo.Limit {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrComboLimitReached).Add(err)
	}

	slog.Info("Combo hacked", "user", user.UID, "combo", combo)

	err = s.db.ComboIncCount(ctx, code)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrComboIncrement).Add(err)
	}

	if combo.PriseCode != "" {
		_, err := s.ProductAdd(ctx, user.UID, combo.UID, combo.PriseCode)
		if err != nil {
			slog.Error("Combo product add error", "user", user.UID, "combo", combo, "error", err)
		}
	}

	if combo.Amount > 0 {
		switch combo.CurrencyCode {
		case "usd":
			trans := &domain.Transaction{
				UID:         domain.GenUID(12),
				UserUID:     user.UID,
				FromUID:     user.UID,
				Type:        domain.TransactionTypeCombo,
				ComboUID:    combo.UID,
				Amount:      combo.Amount,
				PosAmount:   combo.Amount,
				FullAmount:  combo.Amount,
				Coefficient: 100,
				CreatedAt:   time.Now().UTC(),
				ChargedAt:   time.Now().UTC(),
				MsgCodes:    nil,
			}

			err := s.db.TransactionCreate(ctx, trans)
			if err != nil {
				slog.Error("Combo transaction create error", "user", user.UID, "combo", combo, "error", err)
			}

		case "UDEX":
			claim := &domain.UserClaim{
				UID:          domain.GenUID(12),
				ClaimCode:    "UDEX",
				UserUID:      user.UID,
				RefUID:       user.RefUID,
				Level:        0,
				CreatedAt:    time.Now().UTC(),
				ClaimedAt:    time.Now().UTC(),
				Amount:       combo.Amount,
				CurrencyCode: "UDEX",
				Precision:    9, // если нужно, поменяйте precision
				Type:         domain.UserClaimTypeCombo,
				ComboUID:     combo.UID,
			}

			err = s.db.UserBalanceChange(ctx, user.UID, claim.CurrencyCode, claim.Precision, claim.Amount)
			if err != nil {
				slog.Error("Combo balance change error", "user", user.UID, "combo", combo, "error", err)
			}

			err = s.db.UserClaimCreate(ctx, claim)
			if err != nil {
				slog.Error("Combo claim create error", "user", user.UID, "combo", combo, "error", err)
			}
		}
	}

	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrComboReward).Add(err)
	}

	return combo, nil
}

func (s *Service) ProductAdd(ctx context.Context, userUid string, comboUid string, productCode string) (*domain.UserProduct, error) {
	product, err := s.db.ProductGetByCode(ctx, productCode)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrProductNotFound).Add(err)
	}

	startDate := time.Now().UTC()
	exProduct, _ := s.db.UserProductGetLastByCategoryAndDate(ctx, userUid, product.Category, startDate)
	// if current plan with this code exists then set the start date as the end date of current plan
	if exProduct != nil && exProduct.EndAt.After(startDate) {
		startDate = exProduct.EndAt
	}
	endDate := startDate.Add(product.Period)

	userProduct := &domain.UserProduct{
		UID:             domain.GenUID(12),
		UserUID:         userUid,
		ProductCode:     product.Code,
		ProductCategory: product.Category,
		StartAt:         startDate,
		EndAt:           endDate,
		Priority:        product.Priority,
		Multiplier:      product.Multiplier,
		ComboUID:        comboUid,
	}

	err = s.db.UserProductCreate(ctx, userProduct)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return userProduct, nil
}
