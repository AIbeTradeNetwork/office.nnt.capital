package team

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"math/big"
	"time"

	"server/internal/domain"
)

func (s *Service) DistBuy(ctx context.Context, userUid string) (*domain.Dist, error) {
	userPlaceCheck, err := s.db.UserPlaceGetByUserUID(ctx, userUid)
	if userPlaceCheck != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserPlaceExists)
	}

	dUser, err := s.db.UserGetByUID(ctx, userUid)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserNotFound).Add(err)
	}

	dConf, err := s.db.ConfigGet(ctx)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	dUserBalance, err := s.db.UserBalanceGetByUserUIDAndCurrencyCode(ctx, dUser.UID, dConf.DefaultCurrencyCode)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrBalanceNotEnough).Add(err)
	}
	if dUserBalance.Amount < dConf.DistributorPrice {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrBalanceNotEnough).Add(err)
	}

	newUid := domain.GenUID(12)
	newDist := &domain.Dist{
		UID:          newUid,
		UserUID:      dUser.UID,
		RefUID:       dUser.RefUID,
		Amount:       dConf.DistributorPrice,
		CurrencyCode: dConf.DefaultCurrencyCode,
		CreatedAt:    time.Now().UTC(),
		PaidAt:       time.Now().UTC(),
		PayUID:       newUid,
	}

	err = s.db.DistCreate(ctx, newDist)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrDistCreate).Add(err)
	}

	trans := domain.Transaction{
		UID:         domain.GenUID(12),
		UserUID:     dUser.UID,
		FromUID:     dUser.UID,
		Percent:     0,
		Level:       0,
		Type:        domain.TransactionTypeDist,
		RankCode:    "",
		Amount:      -dConf.DistributorPrice,
		PosAmount:   -dConf.DistributorPrice,
		FullAmount:  -dConf.DistributorPrice,
		Coefficient: 100,
		BuyUID:      newUid,
		PayoutUID:   "",
		DepositUID:  "",
		CreatedAt:   time.Now().UTC(),
		ChargedAt:   time.Now().UTC(),
		MsgCodes:    nil,
	}

	err = s.db.TransactionCreate(ctx, &trans)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrPayForOrder).Add(err)
	}

	dUserPlace, err := s.PlaceCreateForRef(ctx, newDist.UserUID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserPlaceCreate).Add(err)
	}

	dUserConfig := &domain.UserConfig{
		UserUID:      dUserPlace.UserUID,
		TeamType:     domain.UserTeamTypeUndefined,
		LastTeamType: s.getOppositeUserTeamTypeByPlace(dUserPlace),
		AllowSwitch:  false,
	}
	err = s.db.UserConfigCreate(ctx, dUserConfig)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserConfigCreate).Add(err)
	}

	userActivity, err := s.db.UserActivityGet(ctx, dUserPlace.UserUID)
	if userActivity != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserActivityExists).Add(err)
	}

	newUserActivity := &domain.UserActivity{
		UserUID:  newDist.UserUID,
		StartAt:  newDist.PaidAt,
		EndAt:    newDist.PaidAt.Add(30 * 24 * time.Hour),
		CvAmount: 0,
	}
	err = s.db.UserActivityCreate(ctx, newUserActivity)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserActivityCreate).Add(err)
	}

	err = s.db.DistUpdate(ctx, newDist)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	// check all ranks to set row, col and matchUid
	err = s.UserRankUpdateAllWithPlace(ctx, newDist.UserUID, dUserPlace)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return newDist, nil
}

// DistPaid Deprecated
func (s *Service) DistPaid(ctx context.Context, order *domain.OrderIn) (*domain.Dist, error) {

	// TODO: add payment save and check paid amount

	newDist, err := s.db.DistGetByPayUID(ctx, order.BillID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	newDist.PaidAt = time.Now().UTC()

	dUserPlace, err := s.PlaceCreateForRef(ctx, newDist.UserUID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserPlaceCreate).Add(err)
	}

	dUserConfig := &domain.UserConfig{
		UserUID:      dUserPlace.UserUID,
		TeamType:     domain.UserTeamTypeUndefined,
		LastTeamType: s.getOppositeUserTeamTypeByPlace(dUserPlace),
		AllowSwitch:  false,
	}
	err = s.db.UserConfigCreate(ctx, dUserConfig)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserConfigCreate).Add(err)
	}

	userActivity, err := s.db.UserActivityGet(ctx, dUserPlace.UserUID)
	if userActivity != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserActivityExists).Add(err)
	}

	newUserActivity := &domain.UserActivity{
		UserUID:  newDist.UserUID,
		StartAt:  newDist.PaidAt,
		EndAt:    newDist.PaidAt.Add(30 * 24 * time.Hour),
		CvAmount: 0,
	}
	err = s.db.UserActivityCreate(ctx, newUserActivity)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserActivityCreate).Add(err)
	}

	err = s.db.DistUpdate(ctx, newDist)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	// check all ranks to set row, col and matchUid
	err = s.UserRankUpdateAllWithPlace(ctx, newDist.UserUID, dUserPlace)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return newDist, nil
}

// PlaceCreateForRef create new place for ref user
func (s *Service) PlaceCreateForRef(ctx context.Context, refUid string) (*domain.UserPlace, error) {
	// if place for ref already exists return this place
	dRefUserPlaced, err := s.db.UserPlaceGetByUserUID(ctx, refUid)
	if dRefUserPlaced != nil {
		return dRefUserPlaced, nil
	}

	// get ref user
	dRefUser, err := s.db.UserGetByUID(ctx, refUid)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// find nearest parent user with place
	dUserPlace, err := s.PlaceGetRefUp(ctx, dRefUser)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// get parent user config for place new ref user in his team
	dUserConfig, err := s.db.UserConfigGetByUserUID(ctx, dUserPlace.UserUID)
	if err != nil || dUserConfig == nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	side := dUserConfig.TeamType

	// if parent user config placing set to Auto (Undefined)
	if side == domain.UserTeamTypeUndefined {
		if !dUserConfig.AllowSwitch {
			checkAmount, err := s.CheckAmountForTeamTypeSwitch(ctx, dUserPlace.UserUID, s.getUserTeamTypeByPlace(dUserPlace))
			if err != nil {
				return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserConfig).Add(err)
			}

			side = s.getUserTeamTypeByPlace(dUserPlace)
			if checkAmount {
				dUserConfig.AllowSwitch = true
				side = s.getOppositeUserTeamTypeByPlace(dUserPlace)
			}
		} else {
			// then place new registration to opposite side of previous placing
			if dUserConfig.LastTeamType == domain.UserTeamTypeLeft {
				side = domain.UserTeamTypeRight
			} else {
				side = domain.UserTeamTypeLeft
			}
		}
	}
	dUserConfig.LastTeamType = side

	// update user config
	err = s.db.UserConfigUpdate(ctx, dUserConfig)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserConfig).Add(err)
	}

	// get new place for ref user
	newPlace, err := s.PlaceGetNew(ctx, dUserPlace, side)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// and save this place for new ref user with his UID
	newPlace.UserUID = refUid
	err = s.db.UserPlaceCreate(ctx, newPlace)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUserPlaceCreate).Add(err)
	}

	return newPlace, nil
}

func (s *Service) CheckAmountForTeamTypeSwitch(ctx context.Context, userUid string, side domain.UserTeamType) (bool, error) {
	dConfig, err := s.db.ConfigGet(ctx)
	if err != nil {
		return false, domain.NewError(teamErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	userPlace, err := s.db.UserPlaceGetByUserUID(ctx, userUid)
	if err != nil {
		return false, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	rows := new(big.Int).Set(maxRow)
	dUserPlaceMax, err := s.PlaceGetLastRow(ctx, userPlace, side, maxRow)
	if dUserPlaceMax != nil {
		rows.Sub(dUserPlaceMax.Row, userPlace.Row)
	}

	teamUsers, err := s.PlaceGetAllDown(ctx, userPlace, side, rows)
	if err != nil {
		return false, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	for _, teamUser := range teamUsers {
		distBuyAmount, err := s.db.BuySumAllByField(ctx, bson.M{"userUid": teamUser.UserUID, "paidAt": bson.M{"$gt": time.Time{}}}, "cv")
		if err != nil {
			return false, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		clientBuyAmount, err := s.db.BuySumAllByField(ctx, bson.M{"refUid": teamUser.UserUID, "type": domain.BuyTypeClient, "paidAt": bson.M{"$gt": time.Time{}}}, "cv")
		if err != nil {
			return false, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		if distBuyAmount+clientBuyAmount >= dConfig.AllowTeamTypeSwitchAmount {
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) UserRankUpdateAllWithPlace(ctx context.Context, userUid string, place *domain.UserPlace) error {
	ranks, err := s.db.UserRankGetAllWith(ctx, bson.M{"userUid": userUid}, nil)
	if len(ranks) == 0 {
		return nil
	}
	for _, rank := range ranks {
		rank.Row = place.Row
		rank.Col = place.Col
		rank.MatchUID = place.MatchUID
		err = s.db.UserRankUpdate(ctx, rank)
		if err != nil {
			return domain.NewError(teamErrorSource).SetCode(domain.ErrUpdate).Add(err)
		}
	}
	return nil
}
