package auth

import (
	"context"
	"server/internal/domain"
)

func (s *Service) Refresh(ctx context.Context, req *domain.RefreshReq) (*domain.AuthRes, error) {
	var dUser *domain.User
	var err error

	validate, err := JwtValidate(req.RefreshToken)
	if err != nil || !validate.Valid {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTokenGeneration).Add(err)
	}

	customClaims, ok := validate.Claims.(*JwtCustomClaims)
	if !ok {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTokenGeneration).Add(err)
	}

	dUser, err = s.db.UserGetByUID(ctx, customClaims.UID)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTokenGeneration).Add(err)
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
