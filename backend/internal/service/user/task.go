package user

import (
	"context"
	"fmt"
	"log/slog"
	"server/internal/domain"
	"strconv"
	"time"
)

func (s *Service) TaskGetAllByUser(ctx context.Context, user *domain.User) ([]*domain.Task, error) {
	tasks, err := s.db.TaskGetAllWithCompletedByUserUIDAndLocale(ctx, user.UID, user.Locale)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return tasks, nil
}

func (s *Service) TaskApprove(ctx context.Context, user *domain.User, taskCode string) (*domain.Task, error) {
	dTask, err := s.db.TaskGetByCode(ctx, taskCode)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	if user.TgID == 0 {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrUserNotHaveTgID)
	}

	if dTask.IsApprove {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrTaskNeedAutoApprove)
	}

	err = s.ProcessTask(ctx, &domain.TaskReq{UserId: user.TgID, Code: taskCode, Action: true})
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrTaskApprove).Add(err)
	}

	return dTask, nil
}

func (s *Service) ProcessTask(ctx context.Context, taskReq *domain.TaskReq) error {
	slog.Info("Process task", "user", taskReq.UserId, "code", taskReq.Code, "action", taskReq.Action)
	dUserAuth, err := s.db.UserAuthGetByTokenAndType(ctx, fmt.Sprintf("%d", taskReq.UserId), domain.AuthTypeTelegram)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrUserNotFound).Add(err).Log()
	}

	go func() {
		ctx := context.Background()

		dUser, err := s.db.UserGetByUID(ctx, dUserAuth.UserUID)
		if err != nil {
			domain.NewError(userErrorSource).SetCode(domain.ErrUserNotFound).Add(err).Log()
			return
		}

		dTask, err := s.db.TaskGetByCode(ctx, taskReq.Code)
		if err != nil {
			domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err).Log()
			return
		}

		if taskReq.Action && !dTask.IsActive {
			domain.NewError(userErrorSource).SetCode(domain.ErrTaskInactive).Log()
			return
		}

		if !dTask.EndAt.IsZero() && time.Now().UTC().After(dTask.EndAt) {
			domain.NewError(userErrorSource).SetCode(domain.ErrTaskEnded).Log()
			return
		}

		if !dTask.StartAt.IsZero() && time.Now().UTC().Before(dTask.StartAt) {
			domain.NewError(userErrorSource).SetCode(domain.ErrTaskNotStarted).Log()
			return
		}

		if taskReq.Action {
			err = s.db.TaskIncCount(ctx, dTask.Code)
		} else {
			err = s.db.TaskDecCount(ctx, dTask.Code)
		}

		switch dTask.CurrencyCode {
		case "usd":
			transCheck, err := s.db.TransactionGetByUserUIDAndTypeAndTaskCode(ctx, dUserAuth.UserUID, domain.TransactionTypeTask, dTask.Code)
			if err != nil {
				if !domain.ErrorIs(err, domain.ErrNoDocuments) {
					domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err).Log()
					return
				}
			}
			if transCheck != nil && taskReq.Action {
				domain.NewError(userErrorSource).SetCode(domain.ErrTaskAlreadyCompleted).Log()
				return
			}

			if transCheck == nil && !taskReq.Action {
				domain.NewError(userErrorSource).SetCode(domain.ErrTaskNotCompleted).Log()
				return
			}

			if !taskReq.Action && transCheck != nil {
				err = s.db.TransactionDelete(ctx, transCheck)
				if err != nil {
					domain.NewError(userErrorSource).SetCode(domain.ErrCreate).Add(err).Log()
					return
				}
			}

			if taskReq.Action && transCheck == nil {
				trans := &domain.Transaction{
					UID:         domain.GenUID(12),
					UserUID:     dUser.UID,
					FromUID:     dUser.UID,
					Percent:     0,
					Level:       0,
					Type:        domain.TransactionTypeTask,
					RankCode:    "",
					TaskCode:    dTask.Code,
					Amount:      dTask.Amount,
					PosAmount:   dTask.Amount,
					FullAmount:  dTask.Amount,
					Coefficient: 100,
					BuyUID:      "",
					PayoutUID:   "",
					DepositUID:  "",
					CreatedAt:   time.Now().UTC(),
					ChargedAt:   time.Now().UTC(),
					MsgCodes:    nil,
				}

				err = s.db.TransactionCreate(ctx, trans)
				if err != nil {
					domain.NewError(userErrorSource).SetCode(domain.ErrCreate).Add(err).Log()
					return
				}
			}

		case "abt":
			claimCheck, err := s.db.UserClaimGetByUserUIDAndClaimCodeAndTypeAndTaskCode(ctx, dUserAuth.UserUID, "ABT", domain.UserClaimTypeTask, dTask.Code)
			if err != nil {
				if !domain.ErrorIs(err, domain.ErrNoDocuments) {
					domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err).Log()
					return
				}
			}
			if claimCheck != nil && taskReq.Action {
				domain.NewError(userErrorSource).SetCode(domain.ErrTaskAlreadyCompleted).Log()
				return
			}

			if claimCheck == nil && !taskReq.Action {
				domain.NewError(userErrorSource).SetCode(domain.ErrTaskNotCompleted).Log()
				return
			}

			amount := dTask.Amount
			boost, _ := s.db.UserProductGetLastBoostAndDate(ctx, dUserAuth.UserUID, time.Now().UTC())
			if boost != nil {
				switch boost.ProductCategory {
				case "boost_x3", "boost_x5", "boost_x10", "boost_x20", "boost_x50":
					amount = amount * boost.Multiplier
				}
			}

			if !taskReq.Action && claimCheck != nil {
				err = s.db.UserBalanceChange(ctx, dUser.UID, dTask.CurrencyCode, dTask.Precision, -claimCheck.Amount)
				if err != nil {
					domain.NewError(userErrorSource).SetCode(domain.ErrUserBalanceChange).Add(err).Log()
					return
				}

				err = s.db.UserClaimDelete(ctx, claimCheck)
				if err != nil {
					errBack := s.db.UserBalanceChange(ctx, dUser.UID, dTask.CurrencyCode, dTask.Precision, claimCheck.Amount)
					if errBack != nil {
						domain.NewError(userErrorSource).SetCode(domain.ErrUserBalanceChange).Add(errBack).Log()
						return
					}
					domain.NewError(userErrorSource).SetCode(domain.ErrCreate).Add(err).Log()
					return
				}
			}

			if taskReq.Action && claimCheck == nil {
				claim := &domain.UserClaim{
					UID:          domain.GenUID(12),
					ClaimCode:    "ABT",
					UserUID:      dUserAuth.UserUID,
					RefUID:       dUser.RefUID,
					Level:        0,
					CreatedAt:    time.Now().UTC(),
					ClaimedAt:    time.Now().UTC(),
					Amount:       amount,
					CurrencyCode: dTask.CurrencyCode,
					Precision:    dTask.Precision,
					Type:         domain.UserClaimTypeTask,
					TaskCode:     dTask.Code,
				}

				err = s.db.UserBalanceChange(ctx, dUser.UID, dTask.CurrencyCode, dTask.Precision, amount)
				if err != nil {
					domain.NewError(userErrorSource).SetCode(domain.ErrUserBalanceChange).Add(err).Log()
					return
				}

				err = s.db.UserClaimCreate(ctx, claim)
				if err != nil {
					errBack := s.db.UserBalanceChange(ctx, dUser.UID, dTask.CurrencyCode, dTask.Precision, -amount)
					if errBack != nil {
						domain.NewError(userErrorSource).SetCode(domain.ErrUserBalanceChange).Add(errBack).Log()
						return
					}
					domain.NewError(userErrorSource).SetCode(domain.ErrCreate).Add(err).Log()
					return
				}
			}

		default:
			domain.NewError(userErrorSource).SetCode(domain.ErrTaskCurrencyNotSupported).Log()
			return

		}
	}()

	return nil
}

func (s *Service) TaskDecline(ctx context.Context, taskCode string, tgIds []string) (map[string]string, error) {
	results := make(map[string]string, len(tgIds))

	ctx = context.Background()

	for _, tgIdString := range tgIds {
		tgId, err := strconv.ParseInt(tgIdString, 10, 64)
		if err != nil {
			results[tgIdString] = err.Error()
			continue
		}

		err = s.ProcessTask(ctx, &domain.TaskReq{UserId: tgId, Code: taskCode, Action: false})
		if err != nil {
			results[tgIdString] = err.Error()
		} else {
			results[tgIdString] = "ok"
		}
	}

	return results, nil
}
