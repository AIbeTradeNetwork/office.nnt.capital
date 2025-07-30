package user

import (
	"context"
	"fmt"
	"server/internal/domain"
	"time"
)

func (s *Service) GetLastTonDeposits(ctx context.Context, limit uint32) ([]*domain.TonTransaction, error) {
	txs, err := s.ton.TransactionGetAll(ctx, limit)
	if err != nil {
		return nil, err
	}

	return txs, nil
}

func (s *Service) CreateTonDeposit(ctx context.Context, tx *domain.TonTransaction) (*domain.Deposit, error) {
	userUid := ""
	user, err := s.db.UserGetByTonWallet(ctx, tx.Sender)
	if err != nil && !domain.ErrorIs(err, domain.ErrNoDocuments) {
		return nil, err
	}
	if user != nil {
		userUid = user.UID
	}
	dep := &domain.Deposit{
		UID:           domain.GenUID(12),
		UserUID:       userUid,
		Amount:        tx.Amount,
		Fee:           tx.Fee,
		CurrencyCode:  tx.CurrencyCode,
		MethodCode:    "ton",
		AccountNumber: tx.Sender,
		AccountName:   "",
		CreatedAt:     time.Now().UTC(),
		ApprovedAt:    time.Time{},
		ChargedAt:     time.Time{},
		CancelledAt:   time.Time{},
		Comment:       tx.Comment,
		TxUID:         tx.TxUID,
	}
	err = s.db.DepositCreate(ctx, dep)
	if err != nil {
		return nil, err
	}

	return dep, nil
}

func (s *Service) CheckDeposit(ctx context.Context, hash string, lt uint64) error {
	tx, err := s.ton.TransactionGetByHashAndLT(ctx, hash, lt)
	if tx == nil {
		return err
	}

	user, err := s.db.UserGetByTonWallet(ctx, tx.Sender)
	if err != nil {
		return err
	}

	dep, _ := s.db.DepositGetByUID(ctx, hash)
	if dep == nil {
		dep = &domain.Deposit{
			UID:           tx.TxUID,
			UserUID:       user.UID,
			Amount:        tx.Amount,
			Fee:           tx.Fee,
			CurrencyCode:  tx.CurrencyCode,
			MethodCode:    "ton",
			AccountNumber: tx.Sender,
			AccountName:   "",
			CreatedAt:     time.Now().UTC(),
			ApprovedAt:    time.Time{},
			ChargedAt:     time.Time{},
			CancelledAt:   time.Time{},
			Comment:       tx.Comment,
			TxUID:         tx.TxUID,
		}

		dep, err = s.DepositCreate(ctx, dep)
		if err != nil {
			return err
		}
	} else {
		if dep.UserUID != "" {
			return fmt.Errorf("deposit already linked to user %s", dep.UserUID)
		}
		dep.ChargedAt = time.Now().UTC()
		dep.UserUID = user.UID
		err = s.db.DepositUpdate(ctx, dep)
	}

	dep, err = s.DepositTransactionCreate(ctx, dep)
	if err != nil {
		return err
	}

	return nil
}

// DepositTotalByUser возвращает сумму всех депозитов пользователя
func (s *Service) DepositTotalByUser(ctx context.Context, userUid string) (int64, error) {
    return s.db.DepositSumByUserUID(ctx, userUid)
}
