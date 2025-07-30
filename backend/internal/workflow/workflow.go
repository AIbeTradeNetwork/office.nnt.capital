package workflow

import (
	"context"
	"server/internal/domain"
	"time"
)

//go:generate mockery --dir . --name BuyService --output ./mocks
type BuyService interface {
	Init(context.Context, *domain.Buy) (*domain.Buy, error)
	BuyPay(context.Context, *domain.Buy) (*domain.Buy, error)
	InitPaid(context.Context, *domain.Buy) (*domain.Buy, error)
	InitAutofarm(context.Context, *domain.UserProduct) (*domain.UserProduct, error)
	Plan(context.Context, *domain.Buy) (*domain.Plan, error)
	Product(context.Context, *domain.Buy) (*domain.Product, error)
	Paid(context.Context, *domain.Buy, time.Time) (*domain.Buy, error)
	PlanAdd(context.Context, *domain.Buy, *domain.Plan) (*domain.UserPlan, error)
	ProductAdd(context.Context, *domain.Buy, *domain.Product) (*domain.UserProduct, error)
	ProductApply(context.Context, *domain.UserProduct) (*domain.UserProduct, error)
	RankAddFromPlan(context.Context, *domain.Buy, *domain.Plan) (*domain.UserRank, error)
	Approved(context.Context, *domain.Buy, time.Time) (*domain.Buy, error)
	BuySetRowCol(context.Context, *domain.Buy, *domain.UserPlace) (*domain.Buy, error)
	Cancelled(context.Context, *domain.Buy, time.Time) (*domain.Buy, error)
	Refund(context.Context, *domain.Buy, time.Time) (*domain.Buy, error)
	PlanEnd(context.Context, *domain.Buy, *domain.UserPlan) (*domain.UserPlan, error)
	RankEnd(context.Context, *domain.Buy, *domain.UserRank) (*domain.UserRank, error)
	Charged(context.Context, *domain.Buy, time.Time) (*domain.Buy, error)
	GetRefUsersUp(context.Context, string, int) ([]*domain.User, error)
	RefUserCharge(context.Context, *domain.Buy, *domain.User, int) (int64, error)
}

//go:generate mockery --dir . --name TeamService --output ./mocks
type TeamService interface {
	PlaceCreateForRef(context.Context, string) (*domain.UserPlace, error)
	PlaceGetAllUp(context.Context, *domain.UserPlace) ([]*domain.UserPlace, error)
	PlaceRefGetAllUp(context.Context, *domain.Buy, *domain.UserPlace) ([]*domain.UserPlace, error)
	ChargeBinBonus(context.Context, *domain.Buy, *domain.UserPlace, int) (*domain.Transaction, error)
	ChargeRefBonus(context.Context, *domain.Buy) (*domain.Buy, error)
	ChargeMatchBonus(context.Context, *domain.Buy, *domain.UserPlace, *domain.Transaction, int) (*domain.Transaction, error)
	ChargeFirstRankBonus(context.Context, *domain.Buy, *domain.UserPlace) (*domain.Buy, error)
	ChargeApproveRankBonus(context.Context, *domain.Buy, *domain.UserPlace) (*domain.Buy, error)
	ChargeFastStartBonus(context.Context, *domain.Buy, *domain.UserPlace) (*domain.Buy, error)
	CalculateNextRank(context.Context, *domain.Buy, *domain.UserPlace) (*domain.UserRank, error)
	CalculateActivity(context.Context, *domain.Buy, *domain.UserPlace) (*domain.UserActivity, error)
	UserRankSetRowCol(context.Context, *domain.UserRank, *domain.UserPlace) (*domain.UserRank, error)
	BuyClientSetRowCol(context.Context, *domain.Buy, *domain.UserPlace) (*domain.Buy, error)
	UserPlaceGet(context.Context, *domain.Buy) (*domain.UserPlace, error)
	PlaceGetRefUpByBuy(context.Context, *domain.Buy) (*domain.UserPlace, error)
	ChargeCoinRefBonus(context.Context, *domain.Buy) (*domain.Buy, error)
}

type UserService interface {
	NotifyCreate(context.Context, *domain.Notification) (*domain.Notification, error)
	NotifyGetAllTgID(context.Context, *domain.Notification, int64) ([]*domain.UserTg, error)
	NotifyUpdate(context.Context, *domain.Notification) (*domain.Notification, error)
	NotifyTgSend(context.Context, int64, string) (int64, error)
	DepositCreate(context.Context, *domain.Deposit) (*domain.Deposit, error)
	DepositTransactionCreate(context.Context, *domain.Deposit) (*domain.Deposit, error)
	UserClaimCreate(context.Context, string, string) (*domain.UserClaim, error)
	GetTelegramIDByUserUID(context.Context, string) (string, error)
	GetUserByUID(context.Context, string) (*domain.User, error)
}
