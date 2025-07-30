package auth

import (
	"context"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"server/internal/domain"
)

func (s *Service) Login(ctx context.Context, req *domain.LoginReq) (*domain.AuthRes, error) {
	var dUser *domain.User
	var err error

	isEmail, err := regexp.MatchString(emailPattern, req.Login)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrEmailWrong).Add(err)
	}

	// find user
	if isEmail {
		dUser, err = s.db.UserGetByEmail(ctx, strings.ToLower(req.Login))
	} else {
		dUser, err = s.db.UserGetByNickname(ctx, req.Login)
	}
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrUserNotFound).Add(err)
	}

	// and password
	dUserAuth, err := s.db.UserAuthGetByUserUIDAndType(ctx, dUser.UID, domain.AuthTypePassword)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrUserAuthNotFound).Add(err)
	}

	// and compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(dUserAuth.Token), []byte(req.Password))
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrPasswordWrong).Add(err)
	}

	// Generate encoded token
	tokenSign, err := JwtGenerate(dUser.UID)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTokenGeneration).Add(err)
	}

	// Generate encoded refresh
	refreshSign, err := RefreshGenerate(dUser.UID)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrRefreshGeneration).Add(err)
	}

	dAuthRes := &domain.AuthRes{
		AuthToken:    tokenSign,
		RefreshToken: refreshSign,
	}

	return dAuthRes, nil
}
