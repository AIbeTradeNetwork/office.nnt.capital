package buy

import (
	"context"
	"time"

	"github.com/google/uuid"

	"server/internal/domain"
)

const (
	buyErrorSource = "[service.buy]"
)

//go:generate mockery --dir . --name DbRepository --output ./mocks
type DbRepository interface {
	BuyGetByUID(context.Context, string) (*domain.Buy, error)
	BuyGetByPayUID(context.Context, string) (*domain.Buy, error)
	BuyCreate(context.Context, *domain.Buy) error
	BuyUpdate(context.Context, *domain.Buy) error
	PlanGetByCode(context.Context, string) (*domain.Plan, error)
	ProductGetByCode(context.Context, string) (*domain.Product, error)
	ProductIncCount(context.Context, string) error
	RankGetByCode(context.Context, string) (*domain.Rank, error)
	UserPlanCreate(context.Context, *domain.UserPlan) error
	UserProductCreate(context.Context, *domain.UserProduct) error
	UserProductUpdate(context.Context, *domain.UserProduct) error
	UserPlanUpdate(context.Context, *domain.UserPlan) error
	UserPlanGetByDate(context.Context, string, time.Time) (*domain.UserPlan, error)
	UserPlanGetByCodeAndDate(context.Context, string, string, time.Time) (*domain.UserPlan, error)
	UserPlanGetLastByCodeAndDate(context.Context, string, string, time.Time) (*domain.UserPlan, error)
	UserProductGetByCodeAndDate(context.Context, string, string, time.Time) (*domain.UserProduct, error)
	UserProductGetByCodes(context.Context, string, []string) (*domain.UserProduct, error)
	UserProductGetByCategoryAndDate(context.Context, string, string, time.Time) (*domain.UserProduct, error)
	UserProductGetLastByCategoryAndDate(context.Context, string, string, time.Time) (*domain.UserProduct, error)
	UserPlanGetByUID(context.Context, string) (*domain.UserPlan, error)
	UserRankGetByCodeAndDate(context.Context, string, string, time.Time) (*domain.UserRank, error)
	UserRankCreate(context.Context, *domain.UserRank) error
	UserRankGetByUID(context.Context, string) (*domain.UserRank, error)
	UserRankUpdate(context.Context, *domain.UserRank) error
	UserActivityGet(context.Context, string) (*domain.UserActivity, error)
	UserPlaceGetByUserUID(context.Context, string) (*domain.UserPlace, error)
	UserGetByUID(context.Context, string) (*domain.User, error)
	ConfigGet(context.Context) (*domain.Config, error)
	UserBalanceGetByUserUIDAndCurrencyCode(context.Context, string, string) (*domain.UserBalance, error)
	TransactionCreate(context.Context, *domain.Transaction) error
	UserClaimCreate(context.Context, *domain.UserClaim) error
	UserBalanceChange(context.Context, string, string, uint8, int64) error
	UserGetAllUpByUID(context.Context, string, int) ([]*domain.User, error)
	UserSafeGetActiveBySafeUIDsAndUserUID(context.Context, []string, string) (*domain.UserSafe, error)
	UserSafeUpdate(context.Context, *domain.UserSafe) error
}

//go:generate mockery --dir . --name WfRepository --output ./mocks
type WfRepository interface {
	BuyCreate(context.Context, *domain.Buy) error
	BuyProduct(context.Context, *domain.Buy) error
	BuySignalPaid(context.Context, *domain.Buy, *domain.BuySignalPaid) error
	BuySignalRefund(context.Context, *domain.Buy, *domain.BuySignalRefund) error
	AutofarmCreate(context.Context, *domain.UserProduct) error
}

//go:generate mockery --dir . --name PBClient --output ./mocks
type PBClient interface {
	OrderCreate(ctx context.Context, order domain.OrderCreateReq) (*domain.CreateOrderResp, error)
	OrderCancel(ctx context.Context, id uuid.UUID) (bool, error)
	OrderInfo(ctx context.Context, id uuid.UUID) (*domain.OrderInfo, error)
}

type AuthService interface {
	UserLevel(context.Context, *domain.User) (*domain.UserLevel, error)
}

type Service struct {
	db   DbRepository
	wf   WfRepository
	pb   PBClient
	auth AuthService
}

func NewBuyService(db DbRepository, wf WfRepository, pb PBClient, auth AuthService) *Service {
	return &Service{db, wf, pb, auth}
}
