package ton

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"log/slog"
	"server/internal/domain"
	"time"
)

func (r *Repo) TransactionGetAll(ctx context.Context, limit uint32) ([]*domain.TonTransaction, error) {
	master, err := r.clientTonUtils.CurrentMasterchainInfo(ctx) // we fetch block just to trigger chain proof check
	if err != nil {
		return nil, err
	}

	res, err := r.clientTonUtils.WaitForBlock(master.SeqNo).GetAccount(ctx, master, r.addrTonUtils)
	if err != nil {
		return nil, err
	}

	// take last tx info from account info
	lastHash := res.LastTxHash
	lastLt := res.LastTxLT

	list, err := r.clientTonUtils.ListTransactions(ctx, r.addrTonUtils, limit, lastLt, lastHash)
	if err != nil {
		return nil, err
	}

	transactions := make([]*domain.TonTransaction, 0, len(list))
	for _, tx := range list {
		fmt.Println(tx.String())

		if tx.IO.In == nil {
			continue
		}

		transactions = append(transactions, r.transactionConvertFromTon(tx))
	}

	return transactions, nil
}

func (r *Repo) TransactionListenAll(ctx context.Context, dest string, lastTxLT uint64, txChan chan<- *domain.TonTransaction) error {
	defer func() {
		slog.Info("Close txChan")
		close(txChan)
	}()
	txs := make(chan *tlb.Transaction)

	go r.SubscribeOnTransactions(ctx, r.addrTonUtils, lastTxLT, txs)

	slog.Info("Start listening for transactions")

	for tx := range txs {
		slog.Info("Transaction received:", "tx", tx.String())
		if tx.IO.In == nil {
			slog.Info("Transaction is not deposit")
			continue
		}
		txChan <- r.transactionConvertFromTon(tx)
	}
	return nil
}

func (r *Repo) TransactionGetByHashAndLT(ctx context.Context, hash string, lt uint64) (*domain.TonTransaction, error) {
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	txs, err := r.clientTonUtils.ListTransactions(ctx, r.addrTonUtils, 1, lt, hashBytes)
	if err != nil {
		return nil, err
	}
	if len(txs) == 0 {
		return nil, fmt.Errorf("transaction not found")
	}
	for _, tx := range txs {
		if hex.EncodeToString(tx.Hash) == hash && tx.LT == lt {
			return r.transactionConvertFromTon(tx), nil
		}
	}
	return nil, fmt.Errorf("transaction not found")
}

func (r *Repo) SubscribeOnTransactions(workerCtx context.Context, addr *address.Address, lastProcessedLT uint64, channel chan<- *tlb.Transaction) {
	defer func() {
		close(channel)
	}()

	wait := 0 * time.Second
	for {
		select {
		case <-workerCtx.Done():
			return
		case <-time.After(wait):
		}
		wait = 3 * time.Second

		ctx, cancel := context.WithTimeout(workerCtx, 30*time.Second)
		master, err := r.clientTonUtils.CurrentMasterchainInfo(ctx)
		cancel()
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		ctx, cancel = context.WithTimeout(workerCtx, 30*time.Second)
		acc, err := r.clientTonUtils.GetAccount(ctx, master, addr)
		cancel()
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		if !acc.IsActive || acc.LastTxLT == 0 {
			// no transactions
			slog.Error("Account is not active or has no transactions")
			continue
		}

		if lastProcessedLT == acc.LastTxLT {
			// already processed all
			slog.Error("Last processed lt is equal to current lt")
			continue
		}

		var transactions []*tlb.Transaction
		lastHash, lastLT := acc.LastTxHash, acc.LastTxLT

		waitList := 0 * time.Second
	list:
		for {
			select {
			case <-workerCtx.Done():
				return
			case <-time.After(waitList):
			}

			ctx, cancel = context.WithTimeout(workerCtx, 10*time.Second)
			res, err := r.clientTonUtils.ListTransactions(ctx, addr, 10, lastLT, lastHash)
			cancel()
			if err != nil {
				slog.Error(err.Error())
				if lsErr, ok := err.(ton.LSError); ok && lsErr.Code == -400 {
					// lt not in db error
					return
				}
				waitList = 3 * time.Second
				continue
			}

			if len(res) == 0 {
				break
			}

			// reverse slice
			for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
				res[i], res[j] = res[j], res[i]
			}

			for i, tx := range res {
				if tx.LT <= lastProcessedLT {
					transactions = append(transactions, res[:i]...)
					break list
				}
			}

			lastLT, lastHash = res[len(res)-1].PrevTxLT, res[len(res)-1].PrevTxHash
			transactions = append(transactions, res...)
			waitList = 0 * time.Second
		}

		if len(transactions) > 0 {
			lastProcessedLT = transactions[0].LT // mark last transaction as known to not trigger twice

			// reverse slice to send in correct time order (from old to new)
			for i, j := 0, len(transactions)-1; i < j; i, j = i+1, j-1 {
				transactions[i], transactions[j] = transactions[j], transactions[i]
			}

			for _, tx := range transactions {
				channel <- tx
			}

			wait = 0 * time.Second
		}
	}
}

func (r *Repo) transactionConvertFromTon(t *tlb.Transaction) *domain.TonTransaction {
	if t.IO.In != nil && t.IO.In.MsgType == tlb.MsgTypeInternal {
		return &domain.TonTransaction{
			UID:          domain.GenUID(12),
			TxUID:        hex.EncodeToString(t.Hash),
			TxLT:         t.LT,
			Addr:         t.IO.In.AsInternal().DestAddr().String(),
			Amount:       t.IO.In.AsInternal().Amount.Nano().Int64(),
			Fee:          t.TotalFees.Coins.Nano().Int64(),
			Sender:       t.IO.In.AsInternal().SenderAddr().String(),
			CurrencyCode: "ton",
			Precision:    9,
			Comment:      t.IO.In.AsInternal().Comment(),
			CreatedAt:    time.Unix(int64(t.Now), 0),
		}
	}
	return nil
}
