package seed

import (
	"context"

	"server/internal/domain"
	"server/internal/service/seed/data"
)

const (
	seedErrorSource = "[service.seed]"
)

type DbRepository interface {
	ConfigCreate(context.Context, *domain.Config) error
	ConfigGet(context.Context) (*domain.Config, error)
	CurrencyCreate(context.Context, *domain.Currency) error
	PlanCreate(context.Context, *domain.Plan) error
	RankCreate(context.Context, *domain.Rank) error
	UserCreate(context.Context, *domain.User) error
	UserAuthCreate(context.Context, *domain.UserAuth) error
	UserPlaceCreate(context.Context, *domain.UserPlace) error
	UserPlanCreate(context.Context, *domain.UserPlan) error
	UserRankCreate(context.Context, *domain.UserRank) error
	UserActivityCreate(context.Context, *domain.UserActivity) error
	UserConfigCreate(context.Context, *domain.UserConfig) error
	ClaimCreate(context.Context, *domain.Claim) error
	ProductCreate(context.Context, *domain.Product) error
	TaskCreate(context.Context, *domain.Task) error
	SafeCreate(context.Context, *domain.Safe) error
}

type Service struct {
	db DbRepository
}

func NewSeedService(db DbRepository) *Service {
	return &Service{db}
}

func (s *Service) Seed(ctx context.Context) error {
	curConf, err := s.db.ConfigGet(ctx)
	if curConf != nil && !curConf.InitializedAt.IsZero() {
		return domain.NewError(seedErrorSource).SetCode(domain.ErrAlreadyInitialized).Add(err)
	}

	conf := data.Config()
	s.db.ConfigCreate(ctx, &conf)

	currencies := data.Currencies()
	for _, c := range currencies {
		s.db.CurrencyCreate(ctx, &c)
	}

	plans := data.Plans()
	for _, p := range plans {
		s.db.PlanCreate(ctx, &p)
	}

	ranks := data.Ranks()
	for _, r := range ranks {
		s.db.RankCreate(ctx, &r)
	}

	users := data.Users()
	for _, u := range users {
		s.db.UserCreate(ctx, &u)
	}

	userAuths := data.UserAuths()
	for _, ua := range userAuths {
		s.db.UserAuthCreate(ctx, &ua)
	}

	userPlaces := data.UserPlaces()
	for _, up := range userPlaces {
		s.db.UserPlaceCreate(ctx, &up)
	}

	userPlans := data.UserPlans()
	for _, up := range userPlans {
		s.db.UserPlanCreate(ctx, &up)
	}

	userRanks := data.UserRanks()
	for _, ur := range userRanks {
		s.db.UserRankCreate(ctx, &ur)
	}

	userActivities := data.UserActivities()
	for _, ua := range userActivities {
		s.db.UserActivityCreate(ctx, &ua)
	}

	userConfigs := data.UserConfigs()
	for _, uc := range userConfigs {
		s.db.UserConfigCreate(ctx, &uc)
	}

	claims := data.Claims()
	for _, c := range claims {
		s.db.ClaimCreate(ctx, &c)
	}

	products := data.Products()
	for _, p := range products {
		s.db.ProductCreate(ctx, &p)
	}

	tasks := data.Tasks()
	for _, t := range tasks {
		s.db.TaskCreate(ctx, &t)
	}

	safe := data.Safes()
	for _, sf := range safe {
		s.db.SafeCreate(ctx, &sf)
	}

	return nil
}

func (s *Service) SeedSafe(ctx context.Context) error {
	safe := data.Safes()
	for _, sf := range safe {
		s.db.SafeCreate(ctx, &sf)
	}

	return nil
}
