package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"server/internal/config"
	"server/internal/domain"
	"strconv"
	"strings"
	"time"
)

func (s *Service) Telegram(ctx context.Context, telegramAuth *domain.TelegramAuth) (*domain.AuthRes, error) {
	cfg := config.Get()
	secret := cfg.TgBotApiKey
	origin, ok := ctx.Value(Key("origin")).(string)
	if ok {
		if strings.Contains(origin, "sg.aibetrade.com") {
			secret = cfg.CgBotApiKey
		}
	}

	err := s.telegramHashCheck(telegramAuth, secret)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTelegramAuthFailed).Add(err)
	}

	userAuth, _ := s.db.UserAuthGetByTokenAndType(ctx, telegramAuth.ID, domain.AuthTypeTelegram)
	if userAuth != nil {
		// Generate encoded token
		tokenSign, err := JwtGenerate(userAuth.UserUID)
		if err != nil {
			return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTokenGeneration).Add(err)
		}

		// Generate encoded refresh
		refreshSign, err := RefreshGenerate(userAuth.UserUID)
		if err != nil {
			return nil, domain.NewError(authErrorSource).SetCode(domain.ErrRefreshGeneration).Add(err)
		}

		authRes := &domain.AuthRes{
			AuthToken:    tokenSign,
			RefreshToken: refreshSign,
		}

		return authRes, nil
	}

	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	newUID := domain.GenUID(8)

	refUid := telegramAuth.RefUid
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

	newNickname := telegramAuth.Username
	if newNickname != "" {
		checkUsername, _ := s.db.UserGetByNickname(ctx, newNickname)
		if checkUsername != nil {
			newNickname = domain.GenNickname()
		}
	} else {
		newNickname = domain.GenNickname()
	}

	tgId, err := strconv.ParseInt(telegramAuth.ID, 10, 64)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTelegramAuthFailed).Add(err)
	}

	// save user
	newUser := &domain.User{
		UID:        newUID,
		Nickname:   newNickname,
		Email:      fmt.Sprintf("%s@t.me", telegramAuth.ID),
		FirstName:  telegramAuth.FirstName,
		LastName:   telegramAuth.LastName,
		RefUID:     refUid,
		LimRefUID:  limRefUid,
		PhotoUrl:   telegramAuth.PhotoUrl,
		TgID:       tgId,
		TgUsername: telegramAuth.Username,
		CreatedAt:  time.Now().UTC(),
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
		Type:      domain.AuthTypeTelegram,
		Token:     telegramAuth.ID,
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

func (s *Service) telegramHashCheck(telegramAuth *domain.TelegramAuth, secret string) error {
	stringToHash := fmt.Sprint("auth_date=", telegramAuth.AuthDate, "\nfirst_name=", telegramAuth.FirstName, "\nid=", telegramAuth.ID)

	if telegramAuth.LastName != "" {
		stringToHash = fmt.Sprint(stringToHash, "\nlast_name=", telegramAuth.LastName)
	}

	if telegramAuth.PhotoUrl != "" {
		stringToHash = fmt.Sprint(stringToHash, "\nphoto_url=", telegramAuth.PhotoUrl)
	}

	if telegramAuth.Username != "" {
		stringToHash = fmt.Sprint(stringToHash, "\nusername=", telegramAuth.Username)
	}

	sha256hash := sha256.New()

	_, err := io.WriteString(sha256hash, secret)
	if err != nil {
		return err
	}

	hmacHash := hmac.New(sha256.New, sha256hash.Sum(nil))

	_, err = io.WriteString(hmacHash, stringToHash)
	if err != nil {
		return err
	}

	checkHash := hex.EncodeToString(hmacHash.Sum(nil))

	fmt.Println(secret)
	fmt.Println(telegramAuth)
	fmt.Println(checkHash)
	fmt.Println(telegramAuth.Hash)

	if telegramAuth.Hash != checkHash {
		return domain.NewError(authErrorSource).SetCode(domain.ErrTelegramHashCheckFailed)
	}

	return nil
}
