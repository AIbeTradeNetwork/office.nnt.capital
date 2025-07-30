package auth

import (
	"context"
	"time"

	"server/internal/domain"
)

const (
	authErrorSource = "[service.auth]"

	emailPattern    = `@`
	nicknamePattern = `^[a-zA-Z0-9_-]{3,32}$`

	DefaultRefUid = "root_ABT_8"

	jwtExpiresIn     = 24 * time.Hour
	refreshExpiresIn = 7 * 24 * time.Hour
)

//go:generate mockery --dir . --name DbRepository --output ./mocks
type DbRepository interface {
	UserGetByUID(context.Context, string) (*domain.User, error)
	UserGetByEmail(context.Context, string) (*domain.User, error)
	UserGetByNickname(context.Context, string) (*domain.User, error)
	UserCreate(context.Context, *domain.User) error
	UserAuthGetByUserUIDAndType(context.Context, string, domain.AuthType) (*domain.UserAuth, error)
	UserAuthGetByTokenAndType(context.Context, string, domain.AuthType) (*domain.UserAuth, error)
	UserAuthCreate(context.Context, *domain.UserAuth) error
	ConfigGet(context.Context) (*domain.Config, error)
	ClaimGetByCode(context.Context, string) (*domain.Claim, error)
	UserClaimGetLastByClaimCodeAndTypeAndUserUID(context.Context, string, domain.UserClaimType, string) (*domain.UserClaim, error)
	UserClaimCreate(context.Context, *domain.UserClaim) error
	UserBalanceGetByUserUIDAndCurrencyCode(context.Context, string, string) (*domain.UserBalance, error)
	UserBalanceCreate(context.Context, *domain.UserBalance) error
	UserBalanceUpdate(context.Context, *domain.UserBalance) error
	UserCountByRefUID(context.Context, string) (int64, error)
	UserPlaceGetByUserUID(context.Context, string) (*domain.UserPlace, error)
	UserProductGetByCategoryAndDate(context.Context, string, string, time.Time) (*domain.UserProduct, error)
	UserProductGetLastBoostAndDate(context.Context, string, time.Time) (*domain.UserProduct, error)
	UserSafeGetActiveBySafeUIDsAndUserUID(context.Context, []string, string) (*domain.UserSafe, error)
	UserSafeCreate(context.Context, *domain.UserSafe) error
	UserSafeUpdateSecret(context.Context, string, string) error
	TaskGetByRefUID(context.Context, string) (*domain.Task, error)
	TaskIncRefCount(context.Context, string) error
}

type Service struct {
	db DbRepository
	ch chan string
}

func NewAuthService(db DbRepository, ch chan string) *Service {
	return &Service{db, ch}
}

func (s *Service) GetUserByAuth(ctx context.Context, userUID string) (*domain.User, error) {
	dUser, err := s.db.UserGetByUID(ctx, userUID)

	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrUserNotFound).Add(err)
	}

	return dUser, nil
}
