package telegram

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/shopspring/decimal"
	"log"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/provider/telegram"
	"server/internal/transport/telegram/lng"
	"strconv"
)

type UserService interface {
	GetUserByTelegramID(ctx context.Context, telegramID string) (*domain.User, error)
	GetTelegramIDByUserUID(ctx context.Context, uid string) (string, error)
	GetUserByUID(ctx context.Context, uid string) (*domain.User, error)
	GetUserCountByRefUID(ctx context.Context, uid string) (int64, error)
	GetUserLastBoost(ctx context.Context, uid string) (*domain.UserProduct, error)
	GlobalConfig(ctx context.Context) (*domain.Config, error)
}

type Transport struct {
	client *tg.Client
	router *tgb.Router
	ch     chan string
	user   UserService
}

func NewTransport(userService UserService, ch chan string) *Transport {
	trans := &Transport{
		client: telegram.NewClient(),
		ch:     ch,
		user:   userService,
	}
	trans.router = trans.InitRouter()
	return trans
}

type Logger struct{}

func (l *Logger) Printf(format string, args ...any) {
	log.Printf(format, args...)
}

func (t *Transport) InitWebhook(ctx context.Context) (*tgb.Webhook, error) {
	cfg := config.Get()
	l := &Logger{}

	webhook := tgb.NewWebhook(t.router, t.client, "https://office.nnt.capital/hook/tg",
		tgb.WithDropPendingUpdates(true),
		tgb.WithWebhookSecuritySubnets(),
		tgb.WithWebhookLogger(l),
	)

	if cfg.Env != "test" {
		if err := webhook.Setup(ctx); err != nil {
			return nil, err
		}
	}

	return webhook, nil
}

func (t *Transport) InitRouter() *tgb.Router {
	return tgb.NewRouter().
		// handles /start and /help
		Message(func(ctx context.Context, msg *tgb.MessageUpdate) error {

			lg := "en"
			user, err := t.user.GetUserByTelegramID(ctx, msg.From.PeerID())
			if err != nil {
				fmt.Println(err)
			}

			if user != nil && user.Locale != "" {
				lg = user.Locale
			}

			return msg.Answer(
				tg.HTML.Text(
					tg.HTML.Bold(lng.Get("welcome", lg)),
					lng.Get("description", lg),
					lng.Get("coins", lg),
					lng.Get("about", lg),
					lng.Get("ecosystem", lg),
					lng.Get("harvest", lg),
					lng.Get("invite", lg),
					lng.Get("exchange", lg),
					lng.Get("earn", lg),
					lng.Get("post", lg),
				),
			).ReplyMarkup(tg.NewInlineKeyboardMarkup(
				tg.NewButtonColumn(
					tg.NewInlineKeyboardButtonWebApp(lng.Get("openApp", lg), tg.WebAppInfo{
						URL: "https://office.nnt.capital/",
					}),
					tg.NewInlineKeyboardButtonURL(lng.Get("ourCommunity", lg), lng.Get("link", lg)),
				)...,
			)).ParseMode(tg.HTML).DoVoid(ctx)
		}, tgb.Command("start", tgb.WithCommandAlias("help")))
}

func (t *Transport) ListenChanel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case userUid := <-t.ch:
			dConf, err := t.user.GlobalConfig(ctx)
			if err != nil {
				fmt.Println(err)
				break
			}

			dUser, err := t.user.GetUserByUID(ctx, userUid)
			if err != nil {
				fmt.Println(err)
				break
			}

			if dUser.RefUID == "root_ABT_8" {
				fmt.Println("root user found")
				break
			}

			dUserRef, err := t.user.GetUserByUID(ctx, dUser.RefUID)
			if err != nil {
				fmt.Println(err)
				break
			}

			id, err := t.user.GetTelegramIDByUserUID(ctx, dUserRef.UID)
			if err != nil {
				fmt.Println(err)
				break
			}
			if id == "" {
				fmt.Println("TgID Not found for " + dUserRef.UID)
				break
			}

			idInt, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println(err)
				break
			}

			count, err := t.user.GetUserCountByRefUID(ctx, dUserRef.UID)
			if err != nil {
				fmt.Println(err)
				break
			}

			amount := dConf.CoinRefBonus
			boost, _ := t.user.GetUserLastBoost(ctx, dUserRef.UID)
			if boost != nil {
				switch boost.ProductCategory {
				case "boost_x3", "boost_x5", "boost_x10", "boost_x20", "boost_x50":
					amount = amount * boost.Multiplier
				}
			}

			msg := fmt.Sprintf("%s %s @%s %s\n%s\n%s <b>%d</b>!",
				lng.Get("register", dUserRef.Locale),
				dUser.FirstName,
				dUser.Nickname,
				dUser.LastName,
				fmt.Sprintf(lng.Get("inviteReward", dUserRef.Locale), decimal.NewFromInt(amount).Div(decimal.NewFromInt(1000000000)).StringFixed(0)),
				lng.Get("invited", dUserRef.Locale),
				count,
			)

			fmt.Println(msg)

			_, err = t.client.SendMessage(tg.UserID(idInt), msg).
				ParseMode(tg.HTML). // optional passed like this
				Do(ctx)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (t *Transport) TestRouter(ctx context.Context) error {
	return tgb.NewPoller(
		t.router,
		t.client,
		tgb.WithPollerAllowedUpdates(
			tg.UpdateTypeMessage,
			tg.UpdateTypeMessageReaction,
		),
	).Run(ctx)
}

func (t *Transport) Close() {
	t.client.Close()
}
