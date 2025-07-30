package team

import (
	"context"
	"errors"
	"math/big"
	"server/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Service) TeamUserGetAllBin(ctx context.Context, userUid string, rows *big.Int) ([]*domain.TeamUser, error) {
	userPlace, err := s.db.UserPlaceGetByUserUID(ctx, userUid)
	if userPlace == nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	filter := bson.M{}
	orFilter := s.getFilterForTeamDown(userPlace, rows, domain.UserTeamTypeUndefined)
	if len(orFilter) == 0 {
		return nil, nil
	}
	filter["$or"] = orFilter

	teamUsers, err := s.db.UserPlaceAggregateWith(ctx, filter, []string{"user", "ranks", "plans", "activity"}, time.Now())
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return teamUsers, nil
}

func (s *Service) TeamUserGetAllRef(ctx context.Context, userUid string) ([]*domain.TeamUser, error) {
	filter := bson.M{}
	filter["refUid"] = userUid

	teamUsers, err := s.db.UserAggregateWith(ctx, filter, []string{"team_count"}, time.Now())
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return teamUsers, nil
}

func (s *Service) TeamUserGetAllMatch(ctx context.Context, userUid string) ([]*domain.TeamUser, error) {
	userPlace, err := s.db.UserPlaceGetByUserUID(ctx, userUid)
	if userPlace == nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	filter := bson.M{}
	filter["matchUid"] = userUid

	teamUsers, err := s.db.UserPlaceAggregateWith(ctx, filter, []string{"user", "ranks", "plans", "activity", "team_count"}, time.Now())
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return teamUsers, nil
}

func (s *Service) TeamBuyGetAll(ctx context.Context, userUid string, limit int64, skip int64) ([]*domain.Buy, error) {
	userPlace, _ := s.db.UserPlaceGetByUserUID(ctx, userUid)
	if userPlace == nil {
		filter := bson.M{"paidAt": bson.M{"$gt": time.Time{}}}
		filter["$or"] = []bson.M{
			{"userUid": userUid},
			{"refUid": userUid},
			{"matchUid": userUid},
		}
		opts := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetLimit(limit).SetSkip(skip)

		buys, err := s.db.BuyGetAllWith(ctx, filter, opts)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
		}

		return buys, nil
	}

	rows := new(big.Int).Set(maxRow)
	dUserPlaceMax, err := s.PlaceGetLastRow(ctx, userPlace, domain.UserTeamTypeUndefined, maxRow)
	if dUserPlaceMax != nil {
		rows.Sub(dUserPlaceMax.Row, userPlace.Row)
	}

	filter := bson.M{"paidAt": bson.M{"$gt": time.Time{}}}
	orFilter := s.getFilterForTeamDown(userPlace, rows, domain.UserTeamTypeUndefined)
	orFilter = append(orFilter, bson.M{"userUid": userUid})
	orFilter = append(orFilter, bson.M{"refUid": userUid})
	orFilter = append(orFilter, bson.M{"matchUid": userUid})
	filter["$or"] = orFilter

	opts := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetLimit(limit).SetSkip(skip)

	buys, err := s.db.BuyGetAllWith(ctx, filter, opts)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return buys, nil
}

func (s *Service) FriendsGetAllByUserUID(ctx context.Context, userUid string, limit int64, skip int64) ([]*domain.User, error) {
	users, err := s.db.UserGetAllWithCountByRefUID(ctx, userUid, limit, skip)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return users, nil
}

func (s *Service) FriendsGetCountByUserUID(ctx context.Context, userUid string) (int64, error) {
	return s.db.UserCountByRefUID(ctx, userUid)
}

func (s *Service) PartnersGetCountByUserUID(ctx context.Context, userUid string) (int64, error) {
	return s.db.UserPlaceCountByMatchUID(ctx, userUid)
}

func (s *Service) AddPartner(ctx context.Context, userUid string, partnerUid string) error {
	// Проверяем, что пользователь существует
	_, err := s.db.UserGetByUID(ctx, userUid)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// Проверяем, что партнёр существует
	partner, err := s.db.UserGetByUID(ctx, partnerUid)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// Проверяем, что партнёр является клиентом пользователя
	if partner.RefUID != userUid {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("partner must be a client of the user"))
	}

	// Проверяем, что у пользователя есть место в бинарной структуре
	userPlace, err := s.db.UserPlaceGetByUserUID(ctx, userUid)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// Проверяем, что партнёр ещё не назначен
	existingPartnerPlace, err := s.db.UserPlaceGetByUserUID(ctx, partnerUid)
	if err == nil && existingPartnerPlace != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("partner already has a place in binary structure"))
	}

	// Получаем лимит партнёров из абонемента пользователя
	subscriptions, err := s.db.UserProductGetAllByUserUIDAndProductCategoryAndDate(ctx, userUid, "subscription", time.Now().UTC())
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	var partnersLimit int64 = 1 // По умолчанию Researcher (1 партнёр)
	if len(subscriptions) > 0 {
		lastSubscription := subscriptions[len(subscriptions)-1]
		subscriptionInfo := domain.GetSubscriptionInfo(lastSubscription.ProductCode)
		if subscriptionInfo != nil {
			partnersLimit = subscriptionInfo.PartnersLimit
		}
	} else {
		// Если нет активного абонемента, используем базовый Researcher
		subscriptionInfo := domain.GetSubscriptionInfo(domain.SubscriptionResearcher)
		partnersLimit = subscriptionInfo.PartnersLimit
	}

	// Проверяем количество партнёров пользователя
	partnersCount, err := s.db.UserPlaceCountByMatchUID(ctx, userUid)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrCount).Add(err)
	}

	// Проверяем лимит партнёров
	if partnersCount >= partnersLimit {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("maximum partners limit reached for current subscription"))
	}

	// Создаём место для партнёра
	newPlace, err := s.PlaceGetNew(ctx, userPlace, domain.UserTeamTypeUndefined)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	newPlace.UserUID = partnerUid
	newPlace.MatchUID = userUid

	err = s.db.UserPlaceCreate(ctx, newPlace)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return nil
}

func (s *Service) RemovePartner(ctx context.Context, userUid string, partnerUid string) error {
	// Проверяем, что пользователь существует
	_, err := s.db.UserGetByUID(ctx, userUid)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// Проверяем, что партнёр существует
	_, err = s.db.UserGetByUID(ctx, partnerUid)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// Получаем место партнёра
	partnerPlace, err := s.db.UserPlaceGetByUserUID(ctx, partnerUid)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// Проверяем, что партнёр действительно партнёр пользователя
	if partnerPlace.MatchUID != userUid {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("partner is not assigned to this user"))
	}

	// Проверяем, что у партнёра нет своих партнёров
	partnerPartnersCount, err := s.db.UserPlaceCountByMatchUID(ctx, partnerUid)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrCount).Add(err)
	}

	if partnerPartnersCount > 0 {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("cannot remove partner who has their own partners"))
	}

	// Удаляем место партнёра
	err = s.db.UserPlaceDelete(ctx, partnerUid)
	if err != nil {
		return domain.NewError(teamErrorSource).SetCode(domain.ErrDelete).Add(err)
	}

	return nil
}
