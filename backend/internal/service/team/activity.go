package team

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/internal/domain"
)

// UserPlaceGet get user place for workflow
func (s *Service) UserPlaceGet(ctx context.Context, newBuy *domain.Buy) (*domain.UserPlace, error) {
	userPlace, err := s.db.UserPlaceGetByUserUID(ctx, newBuy.UserUID)
	if err != nil {
		if domain.ErrorIs(err, domain.ErrNoDocuments) {
			return nil, nil
		}
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return userPlace, nil
}

func (s *Service) PlaceRefGetAllUp(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace) ([]*domain.UserPlace, error) {
	// TODO: add capacity
	refPlaces := make([]*domain.UserPlace, 0)
	row := new(big.Int).Set(zero)
	dUserPlace, _ := s.db.UserPlaceGetByUserUID(ctx, place.MatchUID)

	for dUserPlace != nil && dUserPlace.UserUID != "" && row.Cmp(place.Row) < 0 {
		// get parent user place by ref uid
		row.Add(row, one)
		refPlaces = append(refPlaces, dUserPlace)
		dUserPlace, _ = s.db.UserPlaceGetByUserUID(ctx, dUserPlace.MatchUID)
	}
	return refPlaces, nil
}

// PlaceGetAllUp get all parents tree places up
func (s *Service) PlaceGetAllUp(ctx context.Context, place *domain.UserPlace) ([]*domain.UserPlace, error) {
	// TODO: remove bson find options mongo-driver dependency from services
	filter := bson.M{}
	poses := s.getAllPositionsUp(place)
	if len(poses) == 0 {
		return nil, nil
	}

	orFilter := make([]bson.M, 0, len(poses))
	for _, pos := range poses {
		orFilter = append(orFilter, bson.M{
			"row": pos[0].String(),
			"col": pos[1].String(),
		})
	}
	if len(orFilter) == 0 {
		return nil, nil
	}
	filter["$or"] = orFilter

	opts := options.Find().SetSort(bson.D{{"row", -1}, {"col", -1}}).SetCollation(&options.Collation{
		Locale:          "en_US",
		NumericOrdering: true,
	})
	opts.SetLimit(maxRow.Int64())

	dUserPlaces, err := s.db.UserPlaceGetAll(ctx, filter, opts)
	if err != nil {
		if domain.ErrorIs(err, domain.ErrNoDocuments) {
			return nil, nil
		}
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dUserPlaces, nil
}

// ChargeBinBonus makes all charges to all levels up to net by their rank settings
func (s *Service) ChargeBinBonus(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace, level int) (*domain.Transaction, error) {
	// if buy without row and col then is error
	if len(newBuy.Col.Bits()) == 0 || len(newBuy.Col.Bits()) == 0 {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrBuyWithoutPlace)
	}

	// charge if rank exists and is not expired
	userRank, _ := s.db.UserRankGetByDate(ctx, place.UserUID, newBuy.PaidAt)
	if userRank == nil {
		return nil, nil
	}

	// get rank
	rank, err := s.db.RankGetByCode(ctx, userRank.RankCode)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// if ref bonus not set
	if rank.BinBonus == 0 {
		return nil, nil
	}

	// get max row in team
	rows := new(big.Int).Set(maxRow)
	dUserPlaceMax, err := s.PlaceGetLastRow(ctx, place, domain.UserTeamTypeUndefined, maxRow)
	if dUserPlaceMax != nil {
		rows.Sub(dUserPlaceMax.Row, place.Row)
	}

	// get cv amount of left and right legs
	sumLeft, err := s.PlaceSumAllCvToDate(ctx, place, domain.UserTeamTypeLeft, rows, newBuy.PaidAt)
	if err != nil {
		log.Println(err)
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	sumRight, err := s.PlaceSumAllCvToDate(ctx, place, domain.UserTeamTypeRight, rows, newBuy.PaidAt)
	if err != nil {
		log.Println(err)
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// find diff between legs
	buyTeamType := s.getUserTeamTypeByPlaceAndBuy(place, newBuy)
	sumDiffBefore := sumLeft - sumRight
	sumDiffAfter := sumDiffBefore
	switch buyTeamType {
	case domain.UserTeamTypeLeft:
		sumDiffAfter = sumDiffBefore + newBuy.Cv
	case domain.UserTeamTypeRight:
		sumDiffAfter = sumDiffBefore - newBuy.Cv
	default:
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	var bonusAmount int64 = 0

	if sumDiffBefore >= 0 && sumDiffAfter < sumDiffBefore {
		bonusAmount = domain.Min(sumDiffBefore, newBuy.Cv)
	} else if sumDiffBefore <= 0 && sumDiffAfter > sumDiffBefore {
		bonusAmount = domain.Min(sumDiffBefore*-1, newBuy.Cv)
	}

	if bonusAmount <= 0 {
		return nil, nil
	}

	transaction := &domain.Transaction{
		UID:         domain.GenUID(12),
		UserUID:     place.UserUID,
		FromUID:     newBuy.UserUID,
		Percent:     rank.BinBonus,
		Level:       level,
		Type:        domain.TransactionTypeBinBonus,
		Amount:      0,
		PosAmount:   domain.Percent(bonusAmount, rank.BinBonus),
		FullAmount:  0,
		Coefficient: defaultCoefficient,
		BuyUID:      newBuy.UID,
		PayoutUID:   "",
		CreatedAt:   newBuy.ApprovedAt, // TODO: time.Now().UTC(),
		ChargedAt:   newBuy.ApprovedAt, // domain.BeginningOfNextWeek(newBuy.ApprovedAt).Add(chargeDelay),
		MsgCodes:    nil,
	}

	// check week limit
	var transactionAmount int64 = 0
	if rank.BinBonusLimit > 0 {
		transactionAmount, err = s.db.TransactionSumByUserUIDAndTypeAndDates(ctx, place.UserUID, domain.TransactionTypeBinBonus, domain.BeginningOfWeek(newBuy.PaidAt), domain.EndOfWeek(newBuy.PaidAt))
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
		}
	}
	var levelCharge int64 = 0
	msgCode := domain.TransactionMsgCodeWeekLimit
	weekLimit := rank.BinBonusLimit / 4
	if weekLimit > 0 && transactionAmount < weekLimit {
		levelCharge = domain.Percent(bonusAmount, rank.BinBonus)
		msgCode = domain.TransactionMsgCodeUndefined
		if transactionAmount+levelCharge > weekLimit {
			levelCharge = weekLimit - transactionAmount
			msgCode = domain.TransactionMsgCodeWeekLimit
		}
	}
	if msgCode == domain.TransactionMsgCodeWeekLimit {
		transaction.MsgCodes = append(transaction.MsgCodes, domain.TransactionMsgCodeMonthLimit)
	}
	transaction.Amount = levelCharge
	transaction.FullAmount = levelCharge

	// and if user activity is not expired
	userActivity, _ := s.db.UserActivityGetByDate(ctx, place.UserUID, newBuy.PaidAt)
	if userActivity == nil {
		transaction.Amount = 0
		transaction.FullAmount = 0
		transaction.MsgCodes = append(transaction.MsgCodes, domain.TransactionMsgCodeActivityExpired)
	}

	// exclude zero transaction without a reason
	if transaction.Amount <= 0 && len(transaction.MsgCodes) == 0 {
		return nil, nil
	}

	// create transaction
	err = s.db.TransactionCreate(ctx, transaction)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return transaction, nil
}

func (s *Service) ChargeCoinRefBonus(ctx context.Context, newBuy *domain.Buy) (*domain.Buy, error) {
	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	dPlan, err := s.db.PlanGetByCode(ctx, newBuy.PlanCode)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	if dPlan.CoinBonus <= 0 || dConf.CoinCode == "" {
		return newBuy, nil
	}

	dUser, err := s.db.UserGetByUID(ctx, newBuy.UserUID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserNotFound).Add(err)
	}
	if dUser.RefUID == "" || dUser.RefUID == dConf.DefaultRefUid {
		return newBuy, nil
	}

	claim, err := s.db.ClaimGetByCode(ctx, dConf.CoinCode)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	userClaim := &domain.UserClaim{
		UID:          domain.GenUID(12),
		ClaimCode:    claim.Code,
		UserUID:      dUser.RefUID,
		RefUID:       dUser.UID,
		Level:        0,
		CreatedAt:    time.Now().UTC(),
		ClaimedAt:    time.UnixMilli(0),
		Amount:       dPlan.CoinBonus,
		CurrencyCode: claim.CurrencyCode,
		Precision:    claim.Precision,
		Type:         domain.UserClaimTypePlan,
	}
	err = s.db.UserClaimCreate(ctx, userClaim)
	if err != nil {
		return nil, err
	}

	bal, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, dUser.RefUID, claim.CurrencyCode)
	if err != nil {
		if !domain.ErrorIs(err, domain.ErrNoDocuments) {
			return nil, err
		}
		bal = &domain.UserBalance{
			UserUID:      dUser.RefUID,
			CurrencyCode: claim.CurrencyCode,
			Precision:    claim.Precision,
			Amount:       0,
		}
		err = s.db.UserBalanceCreate(ctx, bal)
		if err != nil {
			return nil, err
		}
	}
	bal.Amount += dPlan.CoinBonus
	err = s.db.UserBalanceUpdate(ctx, bal)
	if err != nil {
		return nil, err
	}

	err = s.auth.ChargeLineCoin(ctx, claim, dUser.RefUID, dPlan.CoinBonus, dConf.DefaultRefUid, 0, dConf.CoinRefLinePercent)
	if err != nil {
		fmt.Println(err)
	}

	return newBuy, nil
}

// ChargeRefBonus makes the charge to referral by rank or by plan settings
func (s *Service) ChargeRefBonus(ctx context.Context, newBuy *domain.Buy) (*domain.Buy, error) {
	dUser, err := s.db.UserGetByUID(ctx, newBuy.UserUID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserNotFound).Add(err)
	}

	if dUser.RefUID == "" {
		return newBuy, nil
	}

	transaction := &domain.Transaction{
		UID:         domain.GenUID(12),
		UserUID:     dUser.RefUID,
		FromUID:     newBuy.UserUID,
		Percent:     0,
		Level:       0,
		Type:        domain.TransactionTypeRefBonus,
		Amount:      0,
		PosAmount:   0,
		FullAmount:  0,
		Coefficient: defaultCoefficient,
		BuyUID:      newBuy.UID,
		PayoutUID:   "",
		CreatedAt:   newBuy.ApprovedAt, // TODO: time.Now().UTC(),
		ChargedAt:   newBuy.ApprovedAt,
		MsgCodes:    nil,
	}

	refRank, _ := s.db.UserRankGetByDate(ctx, dUser.RefUID, newBuy.PaidAt)
	refActivity, _ := s.db.UserActivityGetByDate(ctx, dUser.RefUID, newBuy.PaidAt)
	refPlace, _ := s.db.UserPlaceGetByUserUID(ctx, dUser.RefUID)
	var refCharge int64
	// if is distributor and rank exists and activity is not expired
	if refPlace != nil && refRank != nil {
		rank, err := s.db.RankGetByCode(ctx, refRank.RankCode)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		refPercent, ok := rank.RefBonus[newBuy.PlanCode]
		if ok {
			refCharge = domain.Percent(newBuy.Cv, refPercent)
		}
		transaction.PosAmount = refCharge
		transaction.Percent = refPercent
		if refActivity != nil {
			transaction.Amount = refCharge
			transaction.FullAmount = refCharge
		} else {
			transaction.MsgCodes = append(transaction.MsgCodes, domain.TransactionMsgCodeActivityExpired)
		}
		// if not then check plan
	} else {
		dConf, err := s.db.ConfigGet(ctx)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrConfig).Add(err)
		}
		var refPercent int64
		var ok bool
		if refPlace != nil {
			refPercent, ok = dConf.DistributorRefBonus[newBuy.PlanCode]
			if ok {
				refCharge = domain.Percent(newBuy.Cv, refPercent)
			}
		} else {
			refPercent, ok = dConf.ClientRefBonus[newBuy.PlanCode]
			if ok {
				refCharge = domain.Percent(newBuy.Cv, refPercent)
			}
		}
		transaction.Percent = refPercent
		transaction.PosAmount = refCharge
		transaction.Amount = refCharge
		transaction.FullAmount = refCharge
	}

	// exclude zero transaction without a reason
	if transaction.Amount <= 0 && len(transaction.MsgCodes) == 0 {
		return newBuy, nil
	}

	// create transaction
	err = s.db.TransactionCreate(ctx, transaction)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return newBuy, nil
}

// ChargeMatchBonus makes the charge to referral by rank to ref in places
func (s *Service) ChargeMatchBonus(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace, trans *domain.Transaction, level int) (*domain.Transaction, error) {
	// charge for place up to net
	userRank, _ := s.db.UserRankGetByDate(ctx, place.UserUID, newBuy.PaidAt)
	if userRank == nil {
		return nil, nil
	}

	transaction := &domain.Transaction{
		UID:         domain.GenUID(12),
		UserUID:     place.UserUID,
		FromUID:     newBuy.UserUID,
		Percent:     0,
		Level:       level,
		Type:        domain.TransactionTypeMatchBonus,
		Amount:      0,
		PosAmount:   0,
		FullAmount:  0,
		Coefficient: defaultCoefficient,
		BuyUID:      newBuy.UID,
		PayoutUID:   "",
		CreatedAt:   newBuy.ApprovedAt, // TODO: time.Now().UTC(),
		ChargedAt:   newBuy.ApprovedAt, // domain.BeginningOfNextWeek(newBuy.ApprovedAt).Add(chargeDelay),
		MsgCodes:    nil,
	}

	rank, err := s.db.RankGetByCode(ctx, userRank.RankCode)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	nextRank, _ := s.db.RankGetNext(ctx, rank)

	if len(rank.MatchBonus) > level {
		levelBonus := rank.MatchBonus[level]
		levelCharge := domain.Percent(trans.PosAmount, levelBonus)
		transaction.Percent = levelBonus
		transaction.Amount = levelCharge
		transaction.FullAmount = levelCharge
		transaction.PosAmount = levelCharge
	} else {
		if nextRank != nil && len(nextRank.MatchBonus) > level {
			transaction.PosAmount = domain.Percent(trans.PosAmount, nextRank.MatchBonus[level])
			transaction.MsgCodes = append(transaction.MsgCodes, domain.TransactionMsgCodeNotEnoughRank)
			transaction.RankCode = nextRank.Code
		}
	}

	userActivity, _ := s.db.UserActivityGetByDate(ctx, place.UserUID, newBuy.PaidAt)
	if userActivity == nil {
		transaction.Amount = 0
		transaction.FullAmount = 0
		transaction.MsgCodes = append(transaction.MsgCodes, domain.TransactionMsgCodeActivityExpired)
	}

	// exclude zero transaction without a reason
	if transaction.Amount <= 0 && len(transaction.MsgCodes) == 0 {
		return nil, nil
	}

	err = s.db.TransactionCreate(ctx, transaction)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return transaction, nil
}

// UserRankSetRowCol saves the row and col to user rank for optimization of calculations
func (s *Service) UserRankSetRowCol(ctx context.Context, userRank *domain.UserRank, place *domain.UserPlace) (*domain.UserRank, error) {
	userRank.MatchUID = place.MatchUID
	userRank.Row = place.Row
	userRank.Col = place.Col

	err := s.db.UserRankUpdate(ctx, userRank)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return userRank, nil
}

// BuyClientSetRowCol saves the row and col to client's buy for optimization of calculations
func (s *Service) BuyClientSetRowCol(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace) (*domain.Buy, error) {
	// get max row in team
	rows := new(big.Int).Set(maxRow)
	dUserPlaceMax, err := s.PlaceGetLastRow(ctx, place, domain.UserTeamTypeUndefined, maxRow)
	if dUserPlaceMax != nil {
		rows.Sub(dUserPlaceMax.Row, place.Row)
	}

	sumLeft, err := s.PlaceSumAllCvToDate(ctx, place, domain.UserTeamTypeLeft, rows, newBuy.PaidAt)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	sumRight, err := s.PlaceSumAllCvToDate(ctx, place, domain.UserTeamTypeRight, rows, newBuy.PaidAt)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	dUser, err := s.db.UserGetByUID(ctx, newBuy.UserUID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	newBuy.MatchUID = place.UserUID
	newBuy.RefUID = dUser.RefUID
	newBuy.Row = new(big.Int)
	newBuy.Row.Add(place.Row, one)
	newBuy.Col = new(big.Int)
	newBuy.Col.Mul(place.Col, two)
	if sumLeft < sumRight {
		newBuy.Col.Sub(newBuy.Col, one)
	}
	if sumLeft == sumRight {
		var q, r big.Int
		q.DivMod(place.Col, two, &r)
		if r.Int64() == 1 {
			newBuy.Col.Sub(newBuy.Col, one)
		}
	}
	newBuy.Type = domain.BuyTypeClient

	err = s.db.BuyUpdate(ctx, newBuy)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return newBuy, nil
}

// ChargeFirstRankBonus get current rank for user
func (s *Service) ChargeFirstRankBonus(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace) (*domain.Buy, error) {
	if place.UserUID == domain.SystemUserUID {
		return newBuy, nil
	}

	// check if bonus type already charged by another buy
	nextRank, _ := s.db.UserRankGetByDateAndType(ctx, place.UserUID, domain.BeginningOfNextMonth(newBuy.PaidAt), domain.UserRankTypeAuto)
	if nextRank == nil {
		return newBuy, nil
	}
	trans, _ := s.db.TransactionGetByUserUIDAndTypeAndRankCode(ctx, place.UserUID, domain.TransactionTypeFirstRankBonus, nextRank.RankCode)
	if trans != nil {
		return newBuy, nil
	}

	// compare current rank and next calculated rank
	curRank, _ := s.db.UserRankGetByDate(ctx, place.UserUID, newBuy.PaidAt)
	if curRank != nil && nextRank != nil && nextRank.Priority > curRank.Priority || curRank == nil {
		rank, err := s.db.RankGetByCode(ctx, nextRank.RankCode)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
		}

		// check previous each month if needs months for first time rank bonus
		if rank.FirstBonus.Months > 0 && rank.FirstBonus.Amount > 0 {
			for m := 0; m < rank.FirstBonus.Months; m++ {
				prevTime := time.Duration(m) * -30 * 24 * time.Hour
				prevRank, _ := s.db.UserRankGetByDateAndType(ctx, place.UserUID, newBuy.PaidAt.Add(prevTime), domain.UserRankTypeAuto)
				if prevRank == nil || curRank.Priority > prevRank.Priority {
					return newBuy, nil
				}
			}
		}

		transaction := &domain.Transaction{
			UID:         domain.GenUID(12),
			UserUID:     place.UserUID,
			FromUID:     newBuy.UserUID,
			Percent:     0,
			Level:       0,
			Type:        domain.TransactionTypeFirstRankBonus,
			RankCode:    rank.Code,
			Amount:      rank.FirstBonus.Amount,
			PosAmount:   rank.FirstBonus.Amount,
			FullAmount:  rank.FirstBonus.Amount,
			Coefficient: defaultCoefficient,
			BuyUID:      newBuy.UID,
			PayoutUID:   "",
			CreatedAt:   newBuy.ApprovedAt, // TODO: time.Now().UTC(),
			ChargedAt:   newBuy.ApprovedAt,
			MsgCodes:    nil,
		}
		err = s.db.TransactionCreate(ctx, transaction)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
		}
	}
	return newBuy, nil
}

// ChargeApproveRankBonus get current rank for user
func (s *Service) ChargeApproveRankBonus(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace) (*domain.Buy, error) {
	if place.UserUID == domain.SystemUserUID {
		return newBuy, nil
	}

	// get next rank
	nextRank, _ := s.db.UserRankGetByDateAndType(ctx, place.UserUID, domain.BeginningOfNextMonth(newBuy.PaidAt), domain.UserRankTypeAuto)
	if nextRank == nil {
		return newBuy, nil
	}

	// check if bonus type already charged by another buy
	trans, _ := s.db.TransactionGetByUserUIDAndTypeAndRankCode(ctx, place.UserUID, domain.TransactionTypeApproveRankBonus, nextRank.RankCode)
	if trans != nil {
		return newBuy, nil
	}

	// if first rank bonus not charged yet then don't check approve rank bonus
	firstRankTrans, _ := s.db.TransactionGetByUserUIDAndTypeAndRankCode(ctx, place.UserUID, domain.TransactionTypeFirstRankBonus, nextRank.RankCode)
	if firstRankTrans == nil {
		return newBuy, nil
	}

	curRank, _ := s.db.UserRankGetByDate(ctx, place.UserUID, newBuy.PaidAt)
	if curRank == nil {
		return newBuy, nil
	}

	// compare current rank and next calculated rank
	if nextRank.Priority >= curRank.Priority {
		rank, err := s.db.RankGetByCode(ctx, nextRank.RankCode)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
		}

		// check previous each month if it needs months for first time rank bonus
		prevTime := -30 * 24 * time.Hour
		if rank.ApproveBonus.Months > 0 && rank.ApproveBonus.Amount > 0 {
			for m := 0; m < rank.ApproveBonus.Months; m++ {
				prevTime = time.Duration(m) * prevTime
				prevRank, _ := s.db.UserRankGetByDateAndType(ctx, place.UserUID, newBuy.PaidAt.Add(prevTime), domain.UserRankTypeAuto)
				if prevRank == nil || curRank.Priority > prevRank.Priority {
					return newBuy, nil
				}
			}
		}
		if rank.FirstBonus.Months > 0 && rank.FirstBonus.Amount > 0 {
			for m := 0; m < rank.FirstBonus.Months; m++ {
				prevTime = time.Duration(m) * prevTime
				prevRank, _ := s.db.UserRankGetByDateAndType(ctx, place.UserUID, newBuy.PaidAt.Add(prevTime), domain.UserRankTypeAuto)
				if prevRank == nil || curRank.Priority > prevRank.Priority {
					return newBuy, nil
				}
			}
		}

		transaction := &domain.Transaction{
			UID:         domain.GenUID(12),
			UserUID:     place.UserUID,
			FromUID:     newBuy.UserUID,
			Percent:     0,
			Level:       0,
			Type:        domain.TransactionTypeApproveRankBonus,
			RankCode:    rank.Code,
			Amount:      rank.ApproveBonus.Amount,
			PosAmount:   rank.ApproveBonus.Amount,
			FullAmount:  rank.ApproveBonus.Amount,
			Coefficient: defaultCoefficient,
			BuyUID:      newBuy.UID,
			PayoutUID:   "",
			CreatedAt:   newBuy.ApprovedAt, // TODO: time.Now().UTC(),
			ChargedAt:   newBuy.ApprovedAt,
			MsgCodes:    nil,
		}
		err = s.db.TransactionCreate(ctx, transaction)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
		}
	}
	return newBuy, nil
}

// ChargeFastStartBonus check and charge fast start bonus
// initial settings:
// 1950 CV in 30 days = 200 USD
// 2600 CV in 30 days = 400 USD
// 3900 CV in 30 days = 600 USD
func (s *Service) ChargeFastStartBonus(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace) (*domain.Buy, error) {
	if place.UserUID == domain.SystemUserUID {
		return newBuy, nil
	}

	dConfig, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	if len(dConfig.FastStartBonuses) == 0 {
		return newBuy, nil
	}

	trans, _ := s.db.TransactionGetByUserUIDAndType(ctx, place.UserUID, domain.TransactionTypeFastStartBonus)
	if trans != nil {
		return newBuy, nil
	}

	if place.CreatedAt.Add(dConfig.FastStartDuration).After(newBuy.PaidAt) {
		return newBuy, nil
	}

	transaction := &domain.Transaction{
		UID:         domain.GenUID(12),
		UserUID:     place.UserUID,
		FromUID:     newBuy.UserUID,
		Percent:     0,
		Level:       0,
		Type:        domain.TransactionTypeFastStartBonus,
		Amount:      0,
		PosAmount:   0,
		FullAmount:  0,
		Coefficient: defaultCoefficient,
		BuyUID:      newBuy.UID,
		PayoutUID:   "",
		CreatedAt:   newBuy.ApprovedAt, // TODO: time.Now().UTC(),
		ChargedAt:   newBuy.ApprovedAt,
		MsgCodes:    nil,
	}

	// get max row in team
	rows := new(big.Int).Set(maxRow)
	dUserPlaceMax, err := s.PlaceGetLastRow(ctx, place, domain.UserTeamTypeUndefined, maxRow)
	if dUserPlaceMax != nil {
		rows.Sub(dUserPlaceMax.Row, place.Row)
	}

	filter := bson.M{}
	filter["paidAt"] = bson.M{"$gte": place.CreatedAt, "$lte": place.CreatedAt.Add(dConfig.FastStartDuration)}
	orFilter := s.getFilterForTeamDown(place, rows, domain.UserTeamTypeUndefined)
	if len(orFilter) == 0 {
		return newBuy, nil
	}
	filter["$or"] = orFilter
	sum, err := s.db.BuySumAllByField(ctx, filter, "cv")
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	var maxPay int64
	var fastStartBonus domain.FastStartBonus
	for _, fsb := range dConfig.FastStartBonuses {
		if sum >= fsb.MinCv && fsb.Amount > fastStartBonus.Amount {
			fastStartBonus = fsb
			transaction.Amount = fsb.Amount
		}
		if fsb.Amount > maxPay {
			maxPay = fsb.Amount
		}
	}
	transaction.PosAmount = maxPay
	if transaction.Amount < transaction.PosAmount {
		transaction.MsgCodes = append(transaction.MsgCodes, domain.TransactionMsgCodeAmountNotEnough)
	}

	if fastStartBonus.Amount == 0 && len(transaction.MsgCodes) == 0 {
		return newBuy, nil
	}

	err = s.db.TransactionCreate(ctx, transaction)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return newBuy, nil
}

func (s *Service) CalculateActivity(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace) (*domain.UserActivity, error) {
	userActivity, err := s.db.UserActivityGet(ctx, place.UserUID)
	if userActivity == nil {
		return nil, nil
	}

	dConfig, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	// activity calculation based on amount of purchase
	cvAmount := userActivity.CvAmount + newBuy.Amount
	if cvAmount >= dConfig.ActivityCvAmount {
		startAt := userActivity.EndAt
		if startAt.Before(newBuy.PaidAt) {
			startAt = newBuy.PaidAt
		}
		endAt := startAt.Add(30 * 24 * time.Hour)
		maxAt := newBuy.PaidAt.Add(dConfig.ActivityMaxDuration)
		if endAt.After(maxAt) {
			endAt = maxAt
		}
		userActivity.EndAt = endAt
		userActivity.CvAmount = 0
	} else {
		userActivity.CvAmount += newBuy.Amount
	}
	err = s.db.UserActivityUpdate(ctx, userActivity)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return userActivity, nil
}

// CalculateNextRank sum all cv in teams and calculate automatic rank for the next period
func (s *Service) CalculateNextRank(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace) (*domain.UserRank, error) {
	if place.UserUID == domain.SystemUserUID {
		return nil, nil
	}

	curRank, _ := s.db.UserRankGetByDateAndType(ctx, place.UserUID, domain.BeginningOfNextMonth(newBuy.PaidAt), domain.UserRankTypeAuto)

	// get max row in team
	rows := new(big.Int).Set(maxRow)
	dUserPlaceMax, err := s.PlaceGetLastRow(ctx, place, domain.UserTeamTypeUndefined, maxRow)
	if dUserPlaceMax != nil {
		rows.Sub(dUserPlaceMax.Row, place.Row)
	}

	sumLeftToBeginningOfMonth, err := s.PlaceSumAllCvToDate(ctx, place, domain.UserTeamTypeLeft, rows, domain.BeginningOfMonth(newBuy.PaidAt))
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	sumRightToBeginningOfMonth, err := s.PlaceSumAllCvToDate(ctx, place, domain.UserTeamTypeRight, rows, domain.BeginningOfMonth(newBuy.PaidAt))
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	sumLeftMonth, err := s.PlaceSumAllCvFromDateToDate(ctx, place, domain.UserTeamTypeLeft, rows, domain.BeginningOfMonth(newBuy.PaidAt), newBuy.PaidAt)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	sumRightMonth, err := s.PlaceSumAllCvFromDateToDate(ctx, place, domain.UserTeamTypeRight, rows, domain.BeginningOfMonth(newBuy.PaidAt), newBuy.PaidAt)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	sumDiffToBeginningOfMonth := sumLeftToBeginningOfMonth - sumRightToBeginningOfMonth

	var binValue int64
	if sumDiffToBeginningOfMonth >= 0 {
		binValue = domain.Min(sumDiffToBeginningOfMonth+sumLeftMonth, sumRightMonth)
	} else {
		binValue = domain.Min((sumDiffToBeginningOfMonth-sumRightMonth)*-1, sumLeftMonth)
	}

	minRank, err := s.db.RankGetByMinCv(ctx, binValue)
	if minRank != nil {
		if curRank != nil && curRank.Priority < minRank.Priority || curRank == nil {
			satisfied := true
			if len(minRank.TeamCondition) > 0 {
				satisfied = false
				teamCounts := make(map[string]uint64)
				for _, cond := range minRank.TeamCondition {
					switch cond.TeamType {
					case domain.UserTeamTypeLeft:
						teamCounts, err = s.PlaceCountRefTeamMonth(ctx, place, domain.UserTeamTypeLeft, rows, newBuy.PaidAt)
						if err != nil {
							return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
						}
					case domain.UserTeamTypeRight:
						teamCounts, err = s.PlaceCountRefTeamMonth(ctx, place, domain.UserTeamTypeRight, rows, newBuy.PaidAt)
						if err != nil {
							return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
						}
					default:
						satisfied = false
					}
					cnt, ok := teamCounts[cond.RankCode]
					if ok && cnt >= cond.Count {
						satisfied = true
					}
				}
			}
			if satisfied {
				newUserRank := &domain.UserRank{
					UID:      domain.GenUID(12),
					Type:     domain.UserRankTypeAuto,
					UserUID:  place.UserUID,
					MatchUID: place.MatchUID,
					Row:      place.Row,
					Col:      place.Col,
					BuyUID:   newBuy.UID,
					RankCode: minRank.Code,
					StartAt:  domain.BeginningOfNextMonth(newBuy.PaidAt),
					EndAt:    domain.EndOfNextMonth(newBuy.PaidAt),
					Priority: minRank.Priority,
				}
				err = s.db.UserRankCreate(ctx, newUserRank)
				if err != nil {
					return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
				}
				return newUserRank, nil
			}
		}
	}

	return nil, nil
}
