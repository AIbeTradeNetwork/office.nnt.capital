package user

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xssnick/tonutils-go/address"
	"log/slog"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/provider/tonapi"
	"server/internal/service/auth"
	"time"
)

func (s *Service) GetTonPayload(ctx context.Context) (string, error) {
	return s.ton.GeneratePayload(ctx)
}

func (s *Service) SetTonWallet(ctx context.Context, walletReq string) error {
	dUser, ok := ctx.Value(auth.Key("user")).(*domain.User)
	if !ok {
		return domain.NewError(userErrorSource).SetCode(domain.ErrAccessDenied)
	}

	if dUser.TonWallet != "" {
		return domain.NewError(userErrorSource).SetCode(domain.ErrWalletExists)
	}

	addr, err := s.ton.CheckProof(ctx, walletReq)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrWalletCheck).Add(err)
	}

	checkUser, err := s.db.UserGetByTonWallet(ctx, addr)
	if checkUser != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrWalletExists)
	}

	dUser, err = s.db.UserGetByUID(ctx, dUser.UID)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	parsed, err := address.ParseRawAddr(addr)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrWalletParse).Add(err)
	}

	dUser.TonWallet = parsed.String()
	err = s.db.UserUpdate(ctx, dUser)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return nil
}

func (s *Service) ListenTonDeposits(ctx context.Context) error {
	cfg := config.Get()

	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	depChan := make(chan *domain.TonTransaction)
	go func() {
		slog.Info("Listen for transactions from...", "lastLT", dConf.TonLastLT)
		err = s.ton.TransactionListenAll(ctx, cfg.TonWallet, dConf.TonLastLT, depChan)
		if err != nil {
			slog.Error(fmt.Errorf("transaction listener stopped: %w", err).Error())
		}
	}()

	for tx := range depChan {
		if tx == nil || tx.Amount == 0 || tx.Sender == cfg.TonWallet {
			continue
		}
		slog.Info("Deposit: ", "hash", tx.TxUID)
		userUid := ""
		txUser, err := s.db.UserGetByTonWallet(ctx, tx.Sender)
		if err != nil {
			slog.Error(fmt.Errorf("user not found: %w", err).Error())
		}
		if txUser != nil {
			userUid = txUser.UID
		}

		dep := &domain.Deposit{
			UID:           tx.TxUID,
			UserUID:       userUid,
			Amount:        tx.Amount,
			Fee:           tx.Fee,
			CurrencyCode:  tx.CurrencyCode,
			Precision:     9,
			MethodCode:    "ton",
			AccountNumber: tx.Sender,
			AccountName:   "",
			CreatedAt:     time.Now().UTC(),
			ApprovedAt:    time.Now().UTC(),
			ChargedAt:     time.Time{},
			CancelledAt:   time.Time{},
			Comment:       tx.Comment,
			TxUID:         tx.TxUID,
			TxLT:          tx.TxLT,
		}

		err = s.wf.DepositCreate(ctx, dep)
		if err != nil {
			fmt.Println(fmt.Errorf("deposit create failed: %w", err))
		}

		err = s.db.ConfigUpdateTonLastLT(ctx, tx.TxLT)
		if err != nil {
			fmt.Println(fmt.Errorf("config update ton last LT failed: %w", err))
		}
	}

	return nil
}

func (s *Service) DepositCreate(ctx context.Context, dep *domain.Deposit) (*domain.Deposit, error) {
	dep.CreatedAt = time.Now().UTC()
	dep.ApprovedAt = time.Now().UTC()
	dep.ChargedAt = time.Now().UTC()
	rate, err := tonapi.GetTonRate()
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrTonRate).Add(err)
	}
	dep.Rate = rate
	err = s.db.DepositCreate(ctx, dep)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrDepositCreate).Add(err)
	}

	return dep, nil
}

func (s *Service) DepositTransactionCreate(ctx context.Context, dep *domain.Deposit) (*domain.Deposit, error) {
	if dep.UserUID == "" {
		return dep, nil
	}

	exTx, err := s.db.TransactionGetByUserUIDAndTypeAndDepositUID(ctx, dep.UserUID, domain.TransactionTypeDeposit, dep.UID)
	if exTx != nil {
		return dep, nil
	}

	amount := dep.Rate.Mul(
		decimal.NewFromInt(dep.Amount).Div(
			decimal.NewFromInt(10).
				Pow(decimal.NewFromInt(int64(dep.Precision))),
		),
	).RoundFloor(2).Mul(decimal.NewFromInt(100)).IntPart()

	trx := &domain.Transaction{
		UID:         domain.GenUID(12),
		UserUID:     dep.UserUID,
		FromUID:     dep.UserUID,
		Percent:     0,
		Level:       0,
		Type:        domain.TransactionTypeDeposit,
		RankCode:    "",
		Amount:      amount,
		PosAmount:   amount,
		FullAmount:  amount,
		Coefficient: 100,
		BuyUID:      "",
		PayoutUID:   "",
		DepositUID:  dep.UID,
		CreatedAt:   dep.CreatedAt,
		ChargedAt:   dep.ChargedAt,
		MsgCodes:    nil,
	}

	err = s.db.TransactionCreate(ctx, trx)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrDepositTransactionCreate).Add(err)
	}

	return dep, nil
}
