package server

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"server/internal/config"
	"server/internal/domain"
	"server/internal/provider/payment"
	"server/internal/repository"
	"server/internal/service/auth"
	"server/internal/service/buy"
	"server/internal/service/team"
	"server/internal/service/user"
	"server/internal/service/web"
	"server/internal/transport/graphql/generated"
	"server/internal/transport/graphql/resolver"
	"server/internal/transport/telegram"
)

const (
	serverErrorSource = "[transport.graphql]"
)

func Start() error {
	ctx := context.Background()

	cfg := config.Get()

	//init providers
	pbProvider := payment.New(payment.Config{
		URL:   cfg.PaymentGatewayURL,
		Login: cfg.PaymentGatewayLogin,
		Key:   cfg.PaymentGatewayKey,
	})

	// init repositories
	// persistent db
	dbRepo, err := repository.NewDbRepo(ctx)
	if err != nil {
		return domain.NewError(serverErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	err = dbRepo.EnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(serverErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// workflow
	wfRepo, err := repository.NewWfRepo(ctx)
	if err != nil {
		return domain.NewError(serverErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	tonRepo, err := repository.NewTonRepo(ctx)
	if err != nil {
		return domain.NewError(serverErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// channel for user registration notifications to telegram
	tgUserRegCh := make(chan string, 100)
	defer close(tgUserRegCh)

	// init services
	authService := auth.NewAuthService(dbRepo, tgUserRegCh)
	buyService := buy.NewBuyService(dbRepo, wfRepo, pbProvider, authService)
	teamService := team.NewTeamService(dbRepo, pbProvider, authService)
	webService := web.NewWebService(dbRepo)
	userService := user.NewUserService(dbRepo, wfRepo, tonRepo, teamService, authService, tgUserRegCh)

	// telegram transport
	tgTransport := telegram.NewTransport(userService, tgUserRegCh)
	go tgTransport.ListenChanel(ctx)
	defer tgTransport.Close()

	// router
	router := chi.NewRouter()
	router.Use(authMiddleware(authService))
	router.Use(originMiddleware())
	//router.Use(middleware.Logger)

	// graphql interface
	g := generated.Config{Resolvers: resolver.NewResolver(authService, userService, buyService, teamService, webService)}
	g.Directives.HasAuth = hasAuthDirective
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(g))
	srv.SetErrorPresenter(errorPresenter)

	// routes
	router.Handle("/", playground.Handler("GraphQL AiBeTrade", "/query"))
	router.Handle("/query", srv)

	router.Post("/hook/task", func(w http.ResponseWriter, r *http.Request) {
		//if err := checkRequestSign(r); err != nil {
		//	w.WriteHeader(http.StatusBadRequest)
		//	_, _ = w.Write([]byte(err.Error()))
		//	return
		//}

		var taskReq *domain.TaskReq

		err := json.NewDecoder(r.Body).Decode(&taskReq)
		if err != nil {
			//w.WriteHeader(http.StatusBadRequest)
			//_, _ = w.Write([]byte(err.Error()))
			return
		}

		slog.Info("task", "req", taskReq)

		err = userService.ProcessTask(ctx, taskReq)
		if err != nil {
			slog.Error("task claim failed", "err", err)
			//w.WriteHeader(http.StatusBadRequest)
			//_, _ = w.Write([]byte(err.Error()))
			return
		}
		slog.Error("task claim success", "req", taskReq)
	})

	router.Post("/hook/claim", func(w http.ResponseWriter, r *http.Request) {
		if err := checkRequestSign(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var claimReq *domain.ClaimReq

		err = json.NewDecoder(r.Body).Decode(&claimReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		slog.Info("partner claim", "req", claimReq)

		err = userService.ClaimCreatePartner(ctx, claimReq.Code, claimReq.PartnerCode, claimReq.UserUID, claimReq.Amount)
		if err != nil {
			slog.Error("task claim failed", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		slog.Info("task claim success", "req", claimReq)
	})

	router.HandleFunc("/hook/payment_gateway", func(w http.ResponseWriter, r *http.Request) {
		_ = pbProvider.Handle(w, r, func(ctx context.Context, order *domain.OrderIn) error {

			if order.Status == "PAID" {
				productType := ""
				for _, f := range order.CustomFields {
					if f.Name == "productType" {
						productType = f.Value
					}
				}
				switch productType {
				case "distributor":
					_, err = teamService.DistPaid(ctx, order)
					if err != nil {
						return err
					}
				case "plan":
					_, err = buyService.BuyPaid(ctx, order)
					if err != nil {
						return err
					}
				default:
					return domain.NewError(serverErrorSource).SetCode(domain.ErrUnexpectedProductType)
				}
			}

			return nil
		})

		//log error
	})

	// telegram webhook
	if cfg.Env == "test" {
		go tgTransport.TestRouter(ctx)
	} else {
		webhook, err := tgTransport.InitWebhook(ctx)
		if err != nil {
			return domain.NewError(serverErrorSource).SetCode(domain.ErrProviderInit).Add(err)
		}
		router.Post("/hook/tg", webhook.ServeHTTP)
	}

	// start api server
	log.Printf("connect to http://%s/ for GraphQL playground", cfg.HTTPAddr)
	err = http.ListenAndServe(cfg.HTTPAddr, router)
	if err != nil {
		return domain.NewError(serverErrorSource).SetCode(domain.ErrTransportInit).Add(err)
	}
	return nil
}

func hasAuthDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	_, ok := ctx.Value(auth.Key("user")).(*domain.User)
	if !ok {
		return nil, domain.NewError(serverErrorSource).SetCode(domain.ErrAccessDenied)
	}
	return next(ctx)
}

func authMiddleware(authService *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")

			if authToken == "" {
				next.ServeHTTP(w, r)
				return
			}

			bearer := "Bearer "
			authToken = authToken[len(bearer):]

			validate, err := auth.JwtValidate(authToken)
			if err != nil || !validate.Valid {
				next.ServeHTTP(w, r)
				return
			}

			customClaims, ok := validate.Claims.(*auth.JwtCustomClaims)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			dUser, err := authService.GetUserByAuth(context.Background(), customClaims.UID)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), auth.Key("user"), dUser)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func originMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			ctx := context.WithValue(r.Context(), auth.Key("origin"), origin)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func checkRequestSign(req *http.Request) error {
	sign := req.Header.Get("X-Api-Signature-Sha256")
	if sign == "" {
		return domain.NewError(serverErrorSource).SetCode(domain.ErrAccessDenied)
	}

	mac := hmac.New(sha256.New, []byte("7V4iq4EkxWkcEVC"))

	body, err := io.ReadAll(req.Body)
	defer func() {
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}()

	if err != nil {
		return domain.NewError(serverErrorSource).SetCode(domain.ErrBadRequest).Add(err)
	}

	_, err = mac.Write(body)
	if err != nil {
		return domain.NewError(serverErrorSource).SetCode(domain.ErrBadRequest).Add(err)
	}

	signFromUser := hex.EncodeToString(mac.Sum(nil))
	if sign != signFromUser {
		return domain.NewError(serverErrorSource).SetCode(domain.ErrAccessDenied)
	}

	return nil
}

func errorPresenter(ctx context.Context, e error) *gqlerror.Error {
	var gqlErr *gqlerror.Error
	if errors.As(e, &gqlErr) {
		e = gqlErr.Err
	}

	err := domain.AsError(e)
	if err == nil {
		err = domain.NewError(serverErrorSource).SetCode(domain.ErrServer).Add(e)
	}

	return &gqlerror.Error{
		Err:     err,
		Message: err.Code,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"message": err.Code,
			"source":  err.Source,
			"errors":  err.Errors,
			"native":  err.Native,
		},
	}
}
