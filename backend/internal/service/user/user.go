package user

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"server/internal/domain"
	"server/internal/service/auth"
)

const (
	userErrorSource = "[service.user]"
)

//go:generate mockery --dir . --name DbRepository --output ./mocks
type DbRepository interface {
	UserGetByUID(context.Context, string) (*domain.User, error)
	UserGetByTonWallet(context.Context, string) (*domain.User, error)
	UserUpdate(context.Context, *domain.User) error
	UserPlaceGetByUserUID(context.Context, string) (*domain.UserPlace, error)
	UserPlanGetAllByUserUIDAndDate(context.Context, string, time.Time) ([]*domain.UserPlan, error)
	UserProductGetByCategoryAndDate(context.Context, string, string, time.Time) (*domain.UserProduct, error)
	UserProductGetByComboUID(context.Context, string, string) (*domain.UserProduct, error)
	UserProductGetLastBoostAndDate(context.Context, string, time.Time) (*domain.UserProduct, error)
	UserProductGetAllByUserUIDAndDate(context.Context, string, time.Time) ([]*domain.UserProduct, error)
	UserProductGetAllByUserUIDAndProductCategoryAndDate(context.Context, string, string, time.Time) ([]*domain.UserProduct, error)
	UserProductGetLastByCategoryAndDate(context.Context, string, string, time.Time) (*domain.UserProduct, error)
	UserProductCreate(context.Context, *domain.UserProduct) error
	ProductGetByCode(context.Context, string) (*domain.Product, error)
	UserRankGetAllByUserUIDAndDate(context.Context, string, time.Time) ([]*domain.UserRank, error)
	UserRankGetByDate(context.Context, string, time.Time) (*domain.UserRank, error)
	UserPlanGetByDate(context.Context, string, time.Time) (*domain.UserPlan, error)
	UserActivityGet(context.Context, string) (*domain.UserActivity, error)
	UserConfigGetByUserUID(context.Context, string) (*domain.UserConfig, error)
	UserConfigUpdate(context.Context, *domain.UserConfig) error
	TransactionBalanceByUserUIDAndDate(context.Context, string, time.Time) (int64, error)
	TransactionGetAllByUserUID(context.Context, string, time.Time, int64, int64) ([]*domain.Transaction, error)
	TransactionUpdateAllWithUserUIDAndPayoutUIDAndDate(context.Context, string, string, time.Time) error
	TransactionCreate(context.Context, *domain.Transaction) error
	TransactionDelete(context.Context, *domain.Transaction) error
	PayoutCreate(context.Context, *domain.Payout) error
	PayoutGetAllByUserUID(context.Context, string, int64, int64) ([]*domain.Payout, error)
	ConfigGet(context.Context) (*domain.Config, error)
	ClaimGetByCode(context.Context, string) (*domain.Claim, error)
	UserClaimGetLastByClaimCodeAndTypeAndUserUID(context.Context, string, domain.UserClaimType, string) (*domain.UserClaim, error)
	UserClaimCreate(context.Context, *domain.UserClaim) error
	UserClaimUpsert(context.Context, *domain.UserClaim) error
	UserClaimUpdate(context.Context, *domain.UserClaim) error
	UserClaimDelete(context.Context, *domain.UserClaim) error
	UserBalanceGetByUserUIDAndCurrencyCode(context.Context, string, string) (*domain.UserBalance, error)
	UserBalanceCreate(context.Context, *domain.UserBalance) error
	UserBalanceUpdate(context.Context, *domain.UserBalance) error
	UserAuthGetByTokenAndType(context.Context, string, domain.AuthType) (*domain.UserAuth, error)
	UserAuthGetByUserUIDAndType(context.Context, string, domain.AuthType) (*domain.UserAuth, error)
	UserAuthGetAllUserTg(context.Context, []string, domain.AuthType, int64, int64) ([]*domain.UserTg, error)
	UserCountByRefUID(context.Context, string) (int64, error)
	NotificationGetAllByUserUID(context.Context, string) ([]*domain.Notification, error)
	NotificationCreate(context.Context, *domain.Notification) error
	NotificationUpdate(context.Context, *domain.Notification) error
	NotificationGetByUID(context.Context, string) (*domain.Notification, error)
	DepositCreate(context.Context, *domain.Deposit) error
	DepositUpdate(context.Context, *domain.Deposit) error
	ConfigUpdateTonLastLT(context.Context, uint64) error
	UserBalanceChange(context.Context, string, string, uint8, int64) error
	TaskGetByCode(context.Context, string) (*domain.Task, error)
	TransactionGetByUserUIDAndTypeAndTaskCode(context.Context, string, domain.TransactionType, string) (*domain.Transaction, error)
	TransactionGetByUserUIDAndTypeAndComboCode(context.Context, string, domain.TransactionType, string) (*domain.Transaction, error)
	TransactionGetByUserUIDAndTypeAndDepositUID(context.Context, string, domain.TransactionType, string) (*domain.Transaction, error)
	UserClaimGetByUserUIDAndClaimCodeAndTypeAndTaskCode(context.Context, string, string, domain.UserClaimType, string) (*domain.UserClaim, error)
	UserClaimGetByUserUIDAndClaimCodeAndTypeAndComboCode(context.Context, string, string, domain.UserClaimType, string) (*domain.UserClaim, error)
	UserClaimGetByUserUIDAndClaimCodeAndTypeAndPartnerCode(context.Context, string, string, domain.UserClaimType, string) (*domain.UserClaim, error)
	TaskGetAllWithCompletedByUserUIDAndLocale(context.Context, string, string) ([]*domain.Task, error)
	ComboGetByCode(context.Context, string) (*domain.Combo, error)
	ComboGetAll(context.Context) ([]*domain.Combo, error)
	ComboIncCount(context.Context, string) error
	TaskIncCount(context.Context, string) error
	TaskDecCount(context.Context, string) error
	UserSafeGetActiveBySafeUIDsAndUserUID(context.Context, []string, string) (*domain.UserSafe, error)
	SafeGetByUID(context.Context, string) (*domain.Safe, error)
	SafeUpdate(context.Context, *domain.Safe) error
	UserSafeUpdate(context.Context, *domain.UserSafe) error
	DepositGetByUID(context.Context, string) (*domain.Deposit, error)
	DepositSumByUserUID(context.Context, string) (int64, error)
}

type WfRepository interface {
	NotifyCreate(context.Context, *domain.Notification) error
	DepositCreate(context.Context, *domain.Deposit) error
	AutofarmCreate(context.Context, *domain.UserProduct) error
}

type TonRepository interface {
	TransactionGetAll(context.Context, uint32) ([]*domain.TonTransaction, error)
	GeneratePayload(context.Context) (string, error)
	CheckProof(context.Context, string) (string, error)
	TransactionListenAll(context.Context, string, uint64, chan<- *domain.TonTransaction) error
	TransactionGetByHashAndLT(context.Context, string, uint64) (*domain.TonTransaction, error)
}

type TeamService interface {
	PlaceSumAllCvToDate(context.Context, *domain.UserPlace, domain.UserTeamType, *big.Int, time.Time) (int64, error)
	PlaceGetAllDown(context.Context, *domain.UserPlace, domain.UserTeamType, *big.Int) ([]*domain.UserPlace, error)

	// Partner Application methods
	CreatePartnerApplication(context.Context, string, *domain.PartnerApplicationReq) (*domain.PartnerApplication, error)
	ProcessPartnerApplication(context.Context, string, *domain.PartnerApplicationResponseReq) (*domain.PartnerApplication, error)
	GetPartnerApplications(context.Context, string, int64, int64) ([]*domain.PartnerApplication, error)
	GetMyApplications(context.Context, string, int64, int64) ([]*domain.PartnerApplication, error)
	GetPartnerApplicationsCount(context.Context, string) (int64, error)
	GetMyApplicationsCount(context.Context, string) (int64, error)
}

type AuthService interface {
	ChargeRefCoins(context.Context, string) error
	ChargeToRefCoins(context.Context, string, string) error
	ChargeLineCoin(context.Context, *domain.Claim, string, int64, string, int, []int64) error
	UserLevel(context.Context, *domain.User) (*domain.UserLevel, error)
	UnlimInvite(context.Context, *domain.Config, string) bool
	CheckSafeCodes(context.Context, *domain.Config, string) error
}

type Service struct {
	db   DbRepository
	wf   WfRepository
	ton  TonRepository
	team TeamService
	auth AuthService
	ch   chan string
}

func NewUserService(db DbRepository, wf WfRepository, ton TonRepository, team TeamService, auth AuthService, ch chan string) *Service {
	return &Service{db, wf, ton, team, auth, ch}
}

func (s *Service) Me(ctx context.Context) (*domain.User, error) {
	dUser, ok := ctx.Value(auth.Key("user")).(*domain.User)
	if !ok {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrAccessDenied)
	}

	return dUser, nil
}

func (s *Service) AmountByUserAndSide(ctx context.Context, user *domain.User, side domain.UserTeamType) (int64, error) {
	dUserPlace, err := s.db.UserPlaceGetByUserUID(ctx, user.UID)
	if err != nil {
		return 0, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	sum, err := s.team.PlaceSumAllCvToDate(ctx, dUserPlace, side, big.NewInt(60), time.Now().UTC())
	if err != nil {
		return 0, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return sum, nil
}

func (s *Service) Ranks(ctx context.Context, user *domain.User) ([]*domain.UserRank, error) {
	return s.db.UserRankGetAllByUserUIDAndDate(ctx, user.UID, time.Now().UTC())
}

func (s *Service) Plans(ctx context.Context, user *domain.User) ([]*domain.UserPlan, error) {
	return s.db.UserPlanGetAllByUserUIDAndDate(ctx, user.UID, time.Now().UTC())
}

func (s *Service) Products(ctx context.Context, user *domain.User) ([]*domain.UserProduct, error) {
	return s.db.UserProductGetAllByUserUIDAndDate(ctx, user.UID, time.Now().UTC())
}

func (s *Service) IsPremium(ctx context.Context, user *domain.User) (time.Time, int, int, error) {
	// Получаем все активные абонементы пользователя
	subscriptions, err := s.db.UserProductGetAllByUserUIDAndProductCategoryAndDate(ctx, user.UID, "subscription", time.Now().UTC())
	if err != nil {
		return time.Time{}, 0, 0, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	if len(subscriptions) > 0 {
		lastSubscription := subscriptions[len(subscriptions)-1]
		subscriptionInfo := domain.GetSubscriptionInfo(lastSubscription.ProductCode)

		if subscriptionInfo != nil {
			// Возвращаем информацию об абонементе
			return lastSubscription.EndAt, int(subscriptionInfo.PartnersLimit), int(subscriptionInfo.CommissionRate), nil
		}
	}

	// Если нет активного абонемента, возвращаем базовый Researcher
	subscriptionInfo := domain.GetSubscriptionInfo(domain.SubscriptionResearcher)
	return time.Time{}, int(subscriptionInfo.PartnersLimit), int(subscriptionInfo.CommissionRate), nil
}

func (s *Service) IsAutofarm(ctx context.Context, user *domain.User) (time.Time, error) {
	autofarms, err := s.db.UserProductGetAllByUserUIDAndProductCategoryAndDate(ctx, user.UID, "autofarm", time.Now().UTC())
	if err != nil {
		return time.Time{}, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	if len(autofarms) > 0 {
		lastPrem := autofarms[len(autofarms)-1]
		return lastPrem.EndAt, nil
	}
	return time.Time{}, nil
}

func (s *Service) Rank(ctx context.Context, user *domain.User) (*domain.UserRank, error) {
	return s.db.UserRankGetByDate(ctx, user.UID, time.Now().UTC())
}

func (s *Service) Plan(ctx context.Context, user *domain.User) (*domain.UserPlan, error) {
	return s.db.UserPlanGetByDate(ctx, user.UID, time.Now().UTC())
}

func (s *Service) Place(ctx context.Context, user *domain.User) (*domain.UserPlace, error) {
	dUserPlace, err := s.db.UserPlaceGetByUserUID(ctx, user.UID)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dUserPlace, nil
}

func (s *Service) Activity(ctx context.Context, user *domain.User) (*domain.UserActivity, error) {
	dUserActivity, err := s.db.UserActivityGet(ctx, user.UID)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dUserActivity, nil
}

func (s *Service) Config(ctx context.Context, user *domain.User) (*domain.UserConfig, error) {
	dUserConfig, err := s.db.UserConfigGetByUserUID(ctx, user.UID)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dUserConfig, nil
}

func (s *Service) GlobalConfig(ctx context.Context) (*domain.Config, error) {
	return s.db.ConfigGet(ctx)
}

func (s *Service) Balance(ctx context.Context, user *domain.User) (int64, error) {
	balance, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, user.UID, "usd")
	if err != nil {
		if domain.ErrorIs(err, domain.ErrNoDocuments) {
			return 0, nil
		}
		return 0, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return balance.Amount, nil
}

func (s *Service) TeamCount(ctx context.Context, user *domain.User) (int64, error) {
	return s.db.UserCountByRefUID(ctx, user.UID)
}

func (s *Service) Transactions(ctx context.Context, user *domain.User, limit int64, skip int64) ([]*domain.Transaction, error) {
	transactions, err := s.db.TransactionGetAllByUserUID(ctx, user.UID, time.Now().UTC(), limit, skip)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return transactions, nil
}

func (s *Service) SetTeamType(ctx context.Context, user *domain.User, teamType domain.UserTeamType) (*domain.UserConfig, error) {
	userPlace, err := s.db.UserPlaceGetByUserUID(ctx, user.UID)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	userTeam, err := s.team.PlaceGetAllDown(ctx, userPlace, domain.UserTeamTypeUndefined, big.NewInt(1))
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	if len(userTeam) < 2 {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrTeamTooSmall).Add(err)
	}

	userConfig, err := s.db.UserConfigGetByUserUID(ctx, user.UID)
	if err != nil || userConfig == nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrTeamTooSmall).Add(err)
	}

	if !userConfig.AllowSwitch {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrTeamTooSmall).Add(err)
	}

	userConfig.TeamType = teamType
	err = s.db.UserConfigUpdate(ctx, userConfig)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return userConfig, nil
}

func (s *Service) Payout(ctx context.Context, user *domain.User, payout *domain.Payout) error {
	conf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	// TODO: implement payout methods
	if payout.MethodCode != "USDT_TON" {
		return domain.NewError(userErrorSource).SetCode(domain.ErrPayoutMethod).Add(err)
	}

	if payout.CurrencyCode != conf.DefaultCurrencyCode {
		return domain.NewError(userErrorSource).SetCode(domain.ErrPayoutCurrency).Add(err)
	}

	if payout.AccountNumber == "" {
		return domain.NewError(userErrorSource).SetCode(domain.ErrPayoutAccountNumber).Add(err)
	}

	balance, err := s.db.TransactionBalanceByUserUIDAndDate(ctx, user.UID, time.Now().UTC())
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	if balance < conf.PayoutAmountMin {
		return domain.NewError(userErrorSource).SetCode(domain.ErrPayoutAmountMin).Add(err)
	}

	payout.Amount = balance

	fee := domain.Percent(payout.Amount, conf.PayoutFeePercent)
	if fee < conf.PayoutFeeMin {
		fee = conf.PayoutFeeMin
	}
	payout.Fee = fee

	err = s.db.PayoutCreate(ctx, payout)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	err = s.db.TransactionUpdateAllWithUserUIDAndPayoutUIDAndDate(ctx, user.UID, payout.UID, payout.CreatedAt)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	payoutTransaction := &domain.Transaction{
		UID:         domain.GenUID(12),
		UserUID:     user.UID,
		FromUID:     user.UID,
		Percent:     0,
		Level:       0,
		Type:        domain.TransactionTypePayout,
		RankCode:    "",
		Amount:      -payout.Amount,
		PosAmount:   -payout.Amount,
		FullAmount:  -payout.Amount,
		Coefficient: 100,
		BuyUID:      "",
		PayoutUID:   payout.UID,
		CreatedAt:   payout.CreatedAt,
		ChargedAt:   payout.CreatedAt,
		MsgCodes:    nil,
	}

	err = s.db.TransactionCreate(ctx, payoutTransaction)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return nil
}

func (s *Service) Payouts(ctx context.Context, user *domain.User, limit int64, skip int64) ([]*domain.Payout, error) {
	payouts, err := s.db.PayoutGetAllByUserUID(ctx, user.UID, limit, skip)
	if err != nil {
		if domain.ErrorIs(err, domain.ErrNoDocuments) {
			return nil, nil
		}
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return payouts, nil
}

func (s *Service) SetRef(ctx context.Context, user *domain.User, uid string) (*domain.User, error) {
	dRefUser, err := s.db.UserGetByUID(ctx, uid)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	if user.UID == uid {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrRefSelf)
	}

	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	if user.RefUID != dConf.DefaultRefUid {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrRefExists)
	}

	if !s.auth.UnlimInvite(ctx, dConf, dRefUser.UID) {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrRefUserLimit)
	}

	dPlace, err := s.db.UserPlaceGetByUserUID(ctx, user.UID)
	if dPlace != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrUserPlaceExists).Add(err)
	}

	dUser, err := s.db.UserGetByUID(ctx, user.UID)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	dUser.RefUID = dRefUser.UID

	err = s.db.UserUpdate(ctx, dUser)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	// Charge coins
	err = s.auth.ChargeRefCoins(ctx, dUser.RefUID)
	if err != nil {
		fmt.Println(err)
	}

	// Check safe codes
	err = s.auth.CheckSafeCodes(ctx, dConf, dUser.RefUID)
	if err != nil {
		fmt.Println(err)
	}

	s.ch <- user.UID

	// Charge to ref coins
	err = s.auth.ChargeToRefCoins(ctx, dUser.RefUID, dUser.UID)
	if err != nil {
		fmt.Println(err)
	}

	return dUser, nil
}

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID string) (*domain.User, error) {
	dAuth, err := s.db.UserAuthGetByTokenAndType(ctx, telegramID, domain.AuthTypeTelegram)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	dUser, err := s.db.UserGetByUID(ctx, dAuth.UserUID)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dUser, nil
}

func (s *Service) GetTelegramIDByUserUID(ctx context.Context, uid string) (string, error) {
	dAuth, err := s.db.UserAuthGetByUserUIDAndType(ctx, uid, domain.AuthTypeTelegram)
	if err != nil {
		return "", domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dAuth.Token, nil
}

func (s *Service) GetUserByUID(ctx context.Context, uid string) (*domain.User, error) {
	dUser, err := s.db.UserGetByUID(ctx, uid)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dUser, nil
}

func (s *Service) GetUserCountByRefUID(ctx context.Context, uid string) (int64, error) {
	count, err := s.db.UserCountByRefUID(ctx, uid)
	if err != nil {
		return 0, domain.NewError(userErrorSource).SetCode(domain.ErrCount).Add(err)
	}

	return count, nil
}

func (s *Service) GetUserLastBoost(ctx context.Context, uid string) (*domain.UserProduct, error) {
	boost, err := s.db.UserProductGetLastBoostAndDate(ctx, uid, time.Now().UTC())
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrCount).Add(err)
	}

	return boost, nil
}
