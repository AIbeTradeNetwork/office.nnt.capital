package auth

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"server/internal/domain"
)

func (s *Service) Register(ctx context.Context, req *domain.RegisterReq) (*domain.AuthRes, error) {
	vErr := s.validateRegisterReq(ctx, req)
	if vErr != nil {
		return nil, vErr
	}

	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	// encode password
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 2)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrPasswordEncode).Add(err)
	}

	newUID := domain.GenUID(8)

	refUid := req.RefUid
	limRefUid := ""
	notify := true
	if refUid == "" {
		if dConf.DefaultRefUid != "" {
			refUid = dConf.DefaultRefUid
		} else {
			refUid = DefaultRefUid
		}
	} else {
		if !s.UnlimInvite(ctx, dConf, refUid) {
			notify = false
			limRefUid = refUid
			if dConf.DefaultRefUid != "" {
				refUid = dConf.DefaultRefUid
			} else {
				refUid = DefaultRefUid
			}
		}
	}

	newNickname := req.Nickname
	if newNickname == "" {
		newNickname = domain.GenNickname()
	}

	// save user
	newUser := &domain.User{
		UID:       newUID,
		Nickname:  newNickname,
		Email:     strings.ToLower(req.Email),
		RefUID:    refUid,
		LimRefUID: limRefUid,
		Locale:    req.Locale,
		CreatedAt: time.Now().UTC(),
	}
	err = s.db.UserCreate(ctx, newUser)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrUserCreate).Add(err)
	}

	go func() {
		// Charge coins
		err = s.ChargeRefCoins(ctx, newUser.RefUID)
		if err != nil {
			fmt.Println(err)
		}
	}()

	newAuth := &domain.UserAuth{
		UserUID:   newUID,
		Type:      domain.AuthTypePassword,
		Token:     string(passwordBytes),
		CreatedAt: time.Now().UTC(),
	}
	err = s.db.UserAuthCreate(ctx, newAuth)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrUserAuthCreate).Add(err)
	}

	// send notification to telegram user
	if notify {
		s.ch <- newUID
		go func() {
			// Charge to ref coins
			err = s.ChargeToRefCoins(ctx, newUser.RefUID, newUser.UID)
			if err != nil {
				fmt.Println(err)
			}
		}()
		go func() {
			// Check safe codes
			err = s.CheckSafeCodes(ctx, dConf, newUser.RefUID)
			if err != nil {
				fmt.Println(err)
			}
		}()
		go func() {
			// Check task refs
			err = s.CheckRefTasks(ctx, newUser.RefUID)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	// Generate encoded token
	tokenSign, err := JwtGenerate(newUser.UID)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTokenGeneration).Add(err)
	}

	// Generate encoded refresh
	refreshSign, err := RefreshGenerate(newUser.UID)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrRefreshGeneration).Add(err)
	}

	authRes := &domain.AuthRes{
		AuthToken:    tokenSign,
		RefreshToken: refreshSign,
	}

	return authRes, nil
}

func (s *Service) validateRegisterReq(ctx context.Context, req *domain.RegisterReq) error {
	err := domain.NewError(authErrorSource)

	// required fields and validation
	if req.Email == "" {
		err = err.Add(domain.Error{Code: domain.ErrEmailEmpty, Source: "email"})
	}

	isEmail, e := regexp.MatchString(emailPattern, req.Email)
	if !isEmail || e != nil {
		err = err.Add(domain.Error{Code: domain.ErrEmailWrong, Source: "email"})
	}

	if req.Nickname != "" {
		isNickname, e := regexp.MatchString(nicknamePattern, req.Nickname)
		if !isNickname || e != nil {
			err = err.Add(domain.Error{Code: domain.ErrNicknameWrong, Source: "nickname"})
		}
	}

	if req.Password == "" {
		err = err.Add(domain.Error{Code: domain.ErrPasswordEmpty, Source: "password"})
	}

	if len(req.Password) < 6 {
		err = err.Add(domain.Error{Code: domain.ErrPasswordTooShort, Source: "password"})
	}

	if req.Password != req.RePassword {
		err = err.Add(domain.Error{Code: domain.ErrRePasswordWrong, Source: "password"})
	}

	if len(err.Errors) > 0 {
		err.Code = domain.ErrValidation
		return err
	}

	return nil
}
