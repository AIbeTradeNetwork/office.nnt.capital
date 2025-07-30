package web

import (
	"context"
	"server/internal/config"
	"server/internal/domain"
)

const (
	userErrorSource = "[service.web]"
)

//go:generate mockery --dir . --name DbRepository --output ./mocks
type DbRepository interface {
	UserGetByUID(context.Context, string) (*domain.User, error)
	PlanGetAll(context.Context) ([]*domain.Plan, error)
	ProductGetAllIsActive(context.Context) ([]*domain.Product, error)
	ProductGetAllIsActiveByCategory(context.Context, string) ([]*domain.Product, error)
	RankGetAll(context.Context) ([]*domain.Rank, error)
	CurrencyGetAll(context.Context) ([]*domain.Currency, error)
	ConfigGet(context.Context) (*domain.Config, error)
}

type Service struct {
	db DbRepository
}

func NewWebService(db DbRepository) *Service {
	return &Service{db}
}

func (s *Service) Plans(ctx context.Context) ([]*domain.Plan, error) {
	plans, err := s.db.PlanGetAll(ctx)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return plans, nil
}

func (s *Service) Products(ctx context.Context, category string) ([]*domain.Product, error) {
	var err error
	var products []*domain.Product
	if category != "" {
		products, err = s.db.ProductGetAllIsActiveByCategory(ctx, category)
	} else {
		products, err = s.db.ProductGetAllIsActive(ctx)
	}

	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return products, nil
}

func (s *Service) Ranks(ctx context.Context) ([]*domain.Rank, error) {
	ranks, err := s.db.RankGetAll(ctx)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return ranks, nil
}

func (s *Service) Currencies(ctx context.Context) ([]*domain.Currency, error) {
	currencies, err := s.db.CurrencyGetAll(ctx)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return currencies, nil
}

func (s *Service) Config(ctx context.Context) (*domain.Config, error) {
	cfg := config.Get()

	conf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	conf.TonWallet = cfg.TonWallet

	return conf, nil
}
