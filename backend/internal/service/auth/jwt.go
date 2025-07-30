package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"server/internal/config"
	"server/internal/domain"
)

type Key string

// JwtCustomClaims jwt token struct
type JwtCustomClaims struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

func RefreshGenerate(userUID string) (string, error) {
	cfg := config.Get()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaims{
		UID: userUID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userUID,
			Issuer:    "api",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err := t.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", domain.NewError(authErrorSource).SetCode(domain.ErrTokenGeneration).Add(err)
	}

	return token, nil
}

func JwtGenerate(userUID string) (string, error) {
	cfg := config.Get()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaims{
		UID: userUID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userUID,
			Issuer:    "api",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err := t.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", domain.NewError(authErrorSource).SetCode(domain.ErrTokenGeneration).Add(err)
	}

	return token, nil
}

func JwtValidate(token string) (*jwt.Token, error) {
	cfg := config.Get()
	return jwt.ParseWithClaims(token, &JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTokenMethodWrong)
		}
		return []byte(cfg.JWTSecret), nil
	})
}
