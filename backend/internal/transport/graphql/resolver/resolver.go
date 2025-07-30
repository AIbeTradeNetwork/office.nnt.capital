package resolver

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"math/big"
	"server/internal/domain"
	"time"
)

const (
	resolverErrorSource = "[graphql.resolver]"
)

type AuthService interface {
	Login(context.Context, *domain.LoginReq) (*domain.AuthRes, error)
	Register(context.Context, *domain.RegisterReq) (*domain.AuthRes, error)
	Refresh(context.Context, *domain.RefreshReq) (*domain.AuthRes, error)
	Telegram(context.Context, *domain.TelegramAuth) (*domain.AuthRes, error)
	TelegramApp(context.Context, string) (*domain.AuthRes, error)
	UserLevel(context.Context, *domain.User) (*domain.UserLevel, error)
	Levels(context.Context) ([]*domain.UserLevel, error)
}

type UserService interface {
	Me(context.Context) (*domain.User, error)
	IsPremium(context.Context, *domain.User) (time.Time, int, int, error)
	IsAutofarm(context.Context, *domain.User) (time.Time, error)
	AmountByUserAndSide(context.Context, *domain.User, domain.UserTeamType) (int64, error)
	Ranks(context.Context, *domain.User) ([]*domain.UserRank, error)
	Plans(context.Context, *domain.User) ([]*domain.UserPlan, error)
	Products(context.Context, *domain.User) ([]*domain.UserProduct, error)
	Rank(context.Context, *domain.User) (*domain.UserRank, error)
	Plan(context.Context, *domain.User) (*domain.UserPlan, error)
	Place(context.Context, *domain.User) (*domain.UserPlace, error)
	Activity(context.Context, *domain.User) (*domain.UserActivity, error)
	Config(context.Context, *domain.User) (*domain.UserConfig, error)
	Balance(context.Context, *domain.User) (int64, error)
	Transactions(context.Context, *domain.User, int64, int64) ([]*domain.Transaction, error)
	SetTeamType(context.Context, *domain.User, domain.UserTeamType) (*domain.UserConfig, error)
	Payout(context.Context, *domain.User, *domain.Payout) error
	Payouts(context.Context, *domain.User, int64, int64) ([]*domain.Payout, error)
	SetRef(context.Context, *domain.User, string) (*domain.User, error)
	ClaimGet(context.Context, string, *domain.User) (*domain.Claim, error)
	UserClaimCreate(context.Context, string, string) (*domain.UserClaim, error)
	ClaimBalance(context.Context, string, string) (int64, error)
	UserClaimGetLast(context.Context, string, string) (*domain.UserClaim, error)
	Notifications(context.Context, *domain.User) ([]*domain.Notification, error)
	Notify(context.Context, *domain.Notification) error
	TeamCount(context.Context, *domain.User) (int64, error)
	GetTonPayload(context.Context) (string, error)
	SetTonWallet(context.Context, string) error
	ProcessPrises(context.Context, *domain.PriseReq) (map[string]string, error)
	TaskGetAllByUser(context.Context, *domain.User) ([]*domain.Task, error)
	GetUserLastBoost(context.Context, string) (*domain.UserProduct, error)
	ComboList(context.Context) ([]*domain.Combo, error)
	Combo(context.Context, *domain.User, string) (*domain.Combo, error)
	GetSafe(context.Context, *domain.User) (*domain.UserSafe, error)
	HackSafe(context.Context, *domain.User, string) (*domain.UserSafe, error)
	TaskApprove(context.Context, *domain.User, string) (*domain.Task, error)
	TaskDecline(context.Context, string, []string) (map[string]string, error)
	CheckDeposit(context.Context, string, uint64) error
	DepositTotalByUser(ctx context.Context, userUid string) (int64, error)
}

type BuyService interface {
	Create(context.Context, *domain.BuyReq) (*domain.Buy, error)
	BuyProduct(context.Context, *domain.BuyProductReq) (*domain.Buy, error)
	BuyPlan(context.Context, *domain.BuyReq) (*domain.Buy, error)
	SignalPaid(context.Context, string, *domain.BuySignalPaid) error
	SignalRefund(context.Context, string, *domain.BuySignalRefund) error
}

type TeamService interface {
	PlaceGetNew(context.Context, *domain.UserPlace, domain.UserTeamType) (*domain.UserPlace, error)
	PlaceGetAllUp(context.Context, *domain.UserPlace) ([]*domain.UserPlace, error)
	PlaceGetAllDown(context.Context, *domain.UserPlace, domain.UserTeamType, *big.Int) ([]*domain.UserPlace, error)
	DistBuy(context.Context, string) (*domain.Dist, error)
	PlaceSumAllCvToDate(context.Context, *domain.UserPlace, domain.UserTeamType, *big.Int, time.Time) (int64, error)
	TeamUserGetAllBin(context.Context, string, *big.Int) ([]*domain.TeamUser, error)
	TeamUserGetAllRef(context.Context, string) ([]*domain.TeamUser, error)
	TeamUserGetAllMatch(context.Context, string) ([]*domain.TeamUser, error)
	TeamBuyGetAll(context.Context, string, int64, int64) ([]*domain.Buy, error)
	FriendsGetCountByUserUID(context.Context, string) (int64, error)
	FriendsGetAllByUserUID(context.Context, string, int64, int64) ([]*domain.User, error)
	PartnersGetCountByUserUID(context.Context, string) (int64, error)
	AddPartner(context.Context, string, string) error
	RemovePartner(context.Context, string, string) error
}

type WebService interface {
	Plans(context.Context) ([]*domain.Plan, error)
	Products(context.Context, string) ([]*domain.Product, error)
	Ranks(context.Context) ([]*domain.Rank, error)
	Currencies(context.Context) ([]*domain.Currency, error)
	Config(context.Context) (*domain.Config, error)
}

type Resolver struct {
	auth AuthService
	user UserService
	buy  BuyService
	team TeamService
	web  WebService
}

func NewResolver(authService AuthService, userService UserService, buyService BuyService, teamService TeamService, webService WebService) *Resolver {
	return &Resolver{authService, userService, buyService, teamService, webService}
}
