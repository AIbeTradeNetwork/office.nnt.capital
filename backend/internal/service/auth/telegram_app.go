package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hetiansu5/urlquery"
	"net/url"
	"server/internal/config"
	"server/internal/domain"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	telegramAuthKey = "WebAppData"
)

type WebAppInitData struct {
	QueryID      string     `query:"query_id"`
	UserJson     string     `query:"user"`
	User         WebAppUser `query:"-"`
	ReceiverJson string     `query:"receiver"`
	Receiver     WebAppUser `query:"-"`
	ChatJson     string     `query:"chat"`
	Chat         WebAppChat `query:"-"`
	ChatType     string     `query:"chat_type"`
	ChatInstance string     `query:"chat_instance"`
	StartParam   string     `query:"start_param"`
	CanSendAfter string     `query:"can_send_after"`
	AuthDate     string     `query:"auth_date"`
	Hash         string     `query:"hash"`
}

type WebAppUser struct {
	ID                    int    `json:"id"`
	IsBot                 bool   `json:"is_bot"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Username              string `json:"username"`
	LanguageCode          string `json:"language_code"`
	IsPremium             bool   `json:"is_premium"`
	AddedToAttachmentMenu bool   `json:"added_to_attachment_menu"`
	AllowsWriteToPm       bool   `json:"allows_write_to_pm"`
	PhotoUrl              string `json:"photo_url"`
}

type WebAppChat struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Username string `json:"username"`
	PhotoUrl string `json:"photo_url"`
}

func (s *Service) TelegramApp(ctx context.Context, telegramAppAuth string) (*domain.AuthRes, error) {

	fmt.Println(telegramAppAuth)

	cfg := config.Get()

	secret := cfg.TgBotApiKey
	origin, ok := ctx.Value(Key("origin")).(string)
	if ok {
		if strings.Contains(origin, "sg.aibetrade.com") {
			secret = cfg.CgBotApiKey
		}
	}

	err := s.telegramAppHashCheck(telegramAppAuth, secret)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTelegramAuthFailed).Add(err).Log()
	}

	webAppData, err := s.telegramAppParse(telegramAppAuth)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTelegramDataDecodeFailed).Add(err).Log()
	}

	fmt.Println(webAppData)

	if webAppData.User.ID == 0 {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTelegramDataUserEmpty).Log()
	}

	userId := strconv.Itoa(webAppData.User.ID)

	userAuth, _ := s.db.UserAuthGetByTokenAndType(ctx, userId, domain.AuthTypeTelegram)
	if userAuth != nil {
		// Generate encoded token
		tokenSign, err := JwtGenerate(userAuth.UserUID)
		if err != nil {
			return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTokenGeneration).Add(err).Log()
		}

		// Generate encoded refresh
		refreshSign, err := RefreshGenerate(userAuth.UserUID)
		if err != nil {
			return nil, domain.NewError(authErrorSource).SetCode(domain.ErrRefreshGeneration).Add(err).Log()
		}

		authRes := &domain.AuthRes{
			AuthToken:    tokenSign,
			RefreshToken: refreshSign,
		}

		return authRes, nil
	}

	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrConfig).Add(err).Log()
	}

	newUID := domain.GenUID(8)

	refUid := webAppData.StartParam
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

	newNickname := webAppData.User.Username
	if newNickname != "" {
		checkUsername, _ := s.db.UserGetByNickname(ctx, newNickname)
		if checkUsername != nil {
			newNickname = domain.GenNickname()
		}
	} else {
		newNickname = domain.GenNickname()
	}

	// save user
	newUser := &domain.User{
		UID:        newUID,
		Nickname:   newNickname,
		Email:      fmt.Sprintf("%s@t.me", userId),
		FirstName:  webAppData.User.FirstName,
		LastName:   webAppData.User.LastName,
		RefUID:     refUid,
		LimRefUID:  limRefUid,
		PhotoUrl:   webAppData.User.PhotoUrl,
		Locale:     webAppData.User.LanguageCode,
		TgID:       int64(webAppData.User.ID),
		TgUsername: webAppData.User.Username,
		CreatedAt:  time.Now().UTC(),
	}
	err = s.db.UserCreate(ctx, newUser)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrUserCreate).Add(err).Log()
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
		Token:     userId,
		CreatedAt: time.Now().UTC(),
	}
	err = s.db.UserAuthCreate(ctx, newAuth)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrUserAuthCreate).Add(err).Log()
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
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrTokenGeneration).Add(err).Log()
	}

	// Generate encoded refresh
	refreshSign, err := RefreshGenerate(newUser.UID)
	if err != nil {
		return nil, domain.NewError(authErrorSource).SetCode(domain.ErrRefreshGeneration).Add(err).Log()
	}

	authRes := &domain.AuthRes{
		AuthToken:    tokenSign,
		RefreshToken: refreshSign,
	}

	fmt.Println("user registered: ", newUser)

	return authRes, nil
}

func (s *Service) telegramAppParse(telegramAppAuth string) (*WebAppInitData, error) {
	webAppData := WebAppInitData{}
	err := urlquery.Unmarshal([]byte(telegramAppAuth), &webAppData)
	if err != nil {
		return nil, err
	}
	if webAppData.UserJson != "" {
		err = json.Unmarshal([]byte(webAppData.UserJson), &webAppData.User)
		if err != nil {
			return nil, err
		}
	}
	if webAppData.ReceiverJson != "" {
		err = json.Unmarshal([]byte(webAppData.ReceiverJson), &webAppData.Receiver)
		if err != nil {
			return nil, err
		}
	}
	if webAppData.ChatJson != "" {
		err = json.Unmarshal([]byte(webAppData.ChatJson), &webAppData.Chat)
		if err != nil {
			return nil, err
		}
	}
	return &webAppData, err

}

func (s *Service) telegramAppHashCheck(telegramAppAuth string, secret string) error {
	params, _ := url.ParseQuery(telegramAppAuth)
	strs := []string{}
	var hash = ""
	for k, v := range params {
		if k == "hash" {
			hash = v[0]
			continue
		}
		strs = append(strs, k+"="+v[0])
	}

	sort.Strings(strs)

	var imploded = ""
	for _, s := range strs {
		if imploded != "" {
			imploded += "\n"
		}
		imploded += s
	}

	// Define your data_check_string and hash here
	dataCheckString := []byte(imploded)

	// Define your secret key
	secretKey := []byte(secret)

	// Calculate HMAC_SHA256 of secret key
	hmacSecret := hmac.New(sha256.New, []byte(telegramAuthKey))
	hmacSecret.Write(secretKey)
	secretKeyHash := hmacSecret.Sum(nil)

	// Calculate HMAC_SHA256 of data_check_string using secret key
	hmacData := hmac.New(sha256.New, secretKeyHash)
	hmacData.Write(dataCheckString)
	dataHash := hmacData.Sum(nil)

	// Convert dataHash to hex string
	hexDataHash := hex.EncodeToString(dataHash)

	fmt.Println(secret)
	fmt.Println(imploded)
	fmt.Println(hash)
	fmt.Println(hexDataHash)

	// Compare hash with hexDataHash
	if hexDataHash != hash {
		return domain.NewError(authErrorSource).SetCode(domain.ErrTelegramHashCheckFailed)
	}

	return nil
}
