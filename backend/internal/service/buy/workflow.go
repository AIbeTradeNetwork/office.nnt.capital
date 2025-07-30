package buy

import (
	"context"
	"github.com/google/uuid"
	"server/internal/config"
	"server/internal/service/auth"
	"time"

	"server/internal/domain"
)

// Create deprecated
func (s *Service) Create(ctx context.Context, buy *domain.BuyReq) (*domain.Buy, error) {
	cfg := config.Get()

	dPlan, err := s.db.PlanGetByCode(ctx, buy.PlanCode)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	dUser, err := s.db.UserGetByUID(ctx, buy.UserUID)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	defaultRefUid := auth.DefaultRefUid
	if dConf.DefaultRefUid != "" {
		defaultRefUid = dConf.DefaultRefUid
	}

	var pbOrder *domain.CreateOrderResp
	if cfg.Env == "test" {
		newUuid, _ := uuid.NewUUID()
		pbOrder = &domain.CreateOrderResp{
			BillId: newUuid,
			URL:    "",
		}
	} else {
		price := float64(dPlan.Price) / 100
		if dUser.RefUID == defaultRefUid {
			price = float64(dPlan.RetailPrice) / 100
		}
		pbOrder, err = s.pb.OrderCreate(ctx, domain.OrderCreateReq{
			Amount:       price,
			Callbacks:    domain.Callbacks{},
			Comment:      dPlan.Code,
			CurrencyCode: "USDT",
			CustomFields: []domain.CustomField{
				{
					Name:  "productType",
					Value: "plan",
				},
			},
			Customer: domain.Customer{
				Account: dUser.UID,
				Email:   dUser.Email,
				Phone:   "",
			},
			ExpirationDateTime: time.Now().UTC().Add(dPlan.PaymentWait),
			PaymentSystemID:    uuid.MustParse(cfg.PaymentGatewayPaymentSystemID),
		})
	}
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrPaymentGateway).Add(err)
	}

	dBuy := &domain.Buy{
		UID:          domain.GenUID(12),
		UserUID:      buy.UserUID,
		CreatedAt:    time.Now().UTC(),
		PlanCode:     buy.PlanCode,
		CurrencyCode: dPlan.CurrencyCode,
		Amount:       dPlan.Price,
		Cv:           dPlan.Cv,
		PayUID:       pbOrder.BillId.String(),
	}
	err = s.wf.BuyCreate(ctx, dBuy)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return dBuy, nil
}

func (s *Service) BuyProduct(ctx context.Context, buy *domain.BuyProductReq) (*domain.Buy, error) {
	dProduct, err := s.db.ProductGetByCode(ctx, buy.ProductCode)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	dUserBalance, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, buy.UserUID, dProduct.CurrencyCode)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrBalanceNotEnough).Add(err)
	}
	if dUserBalance.Amount < dProduct.Price {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrBalanceNotEnough).Add(err)
	}

	if buy.ProductCode == "safe_6th_digit" || buy.ProductCode == "safe_7th_digit" || buy.ProductCode == "safe_6th_digit_abt" || buy.ProductCode == "safe_7th_digit_abt" {
		dConf, err := s.db.ConfigGet(ctx)
		if err != nil {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrConfig).Add(err)
		}

		safeUids := []string{dConf.RefSafeUID}
		digits := 5
		if buy.ProductCode == "safe_7th_digit" || buy.ProductCode == "safe_7th_digit_abt" {
			safeUids = []string{dConf.Tier1SafeUID, dConf.Tier1CoinSafeUID}
			digits = 6
		}

		dUserSafe, err := s.db.UserSafeGetActiveBySafeUIDsAndUserUID(ctx, safeUids, buy.UserUID)
		if dUserSafe == nil {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUserSafeNotFound).Add(err)
		}
		if domain.CountHackedSafeCode(dUserSafe.Secret) != digits {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUserSafeNotOpenedEnough).Add(err)
		}
	}

	if buy.ProductCode == "boost_x50_lifetime_usd" || buy.ProductCode == "boost_x50_lifetime_abt" {
		dUserProduct, _ := s.db.UserProductGetByCodes(ctx, buy.UserUID, []string{"boost_x50_lifetime_usd", "boost_x50_lifetime_abt"})
		if dUserProduct != nil {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUserProductExists).Add(err)
		}
	}

	if buy.ProductCode == "boost_archon_abt" || buy.ProductCode == "boost_archon_usd" {
		dUserProduct, _ := s.db.UserProductGetByCodes(ctx, buy.UserUID, []string{"boost_archon_abt", "boost_archon_usd"})
		if dUserProduct != nil {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrUserProductExists).Add(err)
		}
	}

	if dProduct.Limit > 0 {
		if dProduct.Count >= dProduct.Limit {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrProductLimitReached).Add(err)
		}
	}

	err = s.db.ProductIncCount(ctx, buy.ProductCode)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrProductIncCount).Add(err)
	}

	newUid := domain.GenUID(12)
	dBuy := &domain.Buy{
		UID:          newUid,
		UserUID:      buy.UserUID,
		CreatedAt:    time.Now().UTC(),
		ProductCode:  buy.ProductCode,
		CurrencyCode: dProduct.CurrencyCode,
		Amount:       dProduct.Price,
		Cv:           dProduct.Cv,
		PayUID:       newUid,
	}
	err = s.wf.BuyProduct(ctx, dBuy)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return dBuy, nil
}

func (s *Service) BuyPlan(ctx context.Context, buy *domain.BuyReq) (*domain.Buy, error) {
	dPlan, err := s.db.PlanGetByCode(ctx, buy.PlanCode)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	dUser, err := s.db.UserGetByUID(ctx, buy.UserUID)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	defaultRefUid := auth.DefaultRefUid
	if dConf.DefaultRefUid != "" {
		defaultRefUid = dConf.DefaultRefUid
	}

	price := dPlan.Price
	if dUser.RefUID == defaultRefUid {
		price = dPlan.RetailPrice
	}

	dUserBalance, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, buy.UserUID, dPlan.CurrencyCode)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrBalanceNotEnough).Add(err)
	}
	if dUserBalance.Amount < price {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrBalanceNotEnough).Add(err)
	}

	newUid := domain.GenUID(12)
	dBuy := &domain.Buy{
		UID:          newUid,
		UserUID:      buy.UserUID,
		CreatedAt:    time.Now().UTC(),
		PlanCode:     buy.PlanCode,
		CurrencyCode: dPlan.CurrencyCode,
		Amount:       price,
		Cv:           dPlan.Cv,
		PayUID:       newUid,
	}
	err = s.wf.BuyCreate(ctx, dBuy)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return dBuy, nil
}

func (s *Service) BuyPaid(ctx context.Context, order *domain.OrderIn) (*domain.Buy, error) {
	dBuy, err := s.db.BuyGetByPayUID(ctx, order.BillID)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// TODO: add payment save and check paid amount

	info := &domain.BuySignalPaid{
		UID:    dBuy.UID,
		Amount: dBuy.Amount,
	}

	err = s.wf.BuySignalPaid(ctx, dBuy, info)
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dBuy, nil
}

func (s *Service) SignalPaid(ctx context.Context, uid string, info *domain.BuySignalPaid) error {
	dBuy, err := s.db.BuyGetByPayUID(ctx, uid)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	err = s.wf.BuySignalPaid(ctx, dBuy, info)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return nil
}

func (s *Service) SignalRefund(ctx context.Context, uid string, info *domain.BuySignalRefund) error {
	dBuy, err := s.db.BuyGetByUID(ctx, uid)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	err = s.wf.BuySignalRefund(ctx, dBuy, info)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return nil
}
