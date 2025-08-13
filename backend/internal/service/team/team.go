package team

import (
	"context"
	"math/big"
	"time"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/internal/domain"
)

const (
	teamErrorSource = "[service.team]"
)

var (
	lineCap = big.NewInt(2)
	maxRow  = big.NewInt(100)

	chargeDelay              = 12 * time.Hour
	defaultCoefficient int64 = 100
)

//go:generate mockery --dir . --name DbRepository --output ./mocks
type DbRepository interface {
	UserPlaceGetByUserUID(context.Context, string) (*domain.UserPlace, error)
	UserPlaceGetByMatchUID(context.Context, string) (*domain.UserPlace, error)
	UserPlaceGetByRowCol(context.Context, *big.Int, *big.Int) (*domain.UserPlace, error)
	UserPlaceGetAll(context.Context, interface{}, *options.FindOptions) ([]*domain.UserPlace, error)
	UserPlaceGet(context.Context, interface{}, *options.FindOneOptions) (*domain.UserPlace, error)
	UserConfigGetByUserUID(context.Context, string) (*domain.UserConfig, error)
	UserPlaceCreate(context.Context, *domain.UserPlace) error
	UserGetByUID(context.Context, string) (*domain.User, error)
	UserConfigCreate(context.Context, *domain.UserConfig) error
	UserConfigUpdate(context.Context, *domain.UserConfig) error
	TransactionCreate(context.Context, *domain.Transaction) error
	TransactionSumByUserUIDAndTypeAndDates(context.Context, string, domain.TransactionType, time.Time, time.Time) (int64, error)
	TransactionDeleteAllByBuyUIDAndType(context.Context, string, domain.TransactionType) error
	TransactionDeleteAllByUserUIDAndBuyUIDAndType(context.Context, string, string, domain.TransactionType) error
	TransactionGetByUserUIDAndType(context.Context, string, domain.TransactionType) (*domain.Transaction, error)
	TransactionGetByUserUIDAndTypeAndRankCode(context.Context, string, domain.TransactionType, string) (*domain.Transaction, error)
	RankGetByCode(context.Context, string) (*domain.Rank, error)
	RankGetNext(context.Context, *domain.Rank) (*domain.Rank, error)
	UserRankGetByDateAndType(context.Context, string, time.Time, domain.UserRankType) (*domain.UserRank, error)
	UserRankGetByDate(context.Context, string, time.Time) (*domain.UserRank, error)
	UserPlanGetByDate(context.Context, string, time.Time) (*domain.UserPlan, error)
	PlanGetByCode(context.Context, string) (*domain.Plan, error)
	BuySumAllByField(context.Context, interface{}, string) (int64, error)
	BuyUpdate(context.Context, *domain.Buy) error
	RankGetByMinCv(context.Context, int64) (*domain.Rank, error)
	UserRankUpdate(context.Context, *domain.UserRank) error
	UserRankCountAllByRankCode(context.Context, interface{}) (map[string]uint64, error)
	UserRankCreate(context.Context, *domain.UserRank) error
	UserRankGetByUserUID(context.Context, string) (*domain.UserRank, error)
	UserRankGetAllWith(context.Context, interface{}, *options.FindOptions) ([]*domain.UserRank, error)
	UserActivityCreate(context.Context, *domain.UserActivity) error
	UserActivityGetByDate(context.Context, string, time.Time) (*domain.UserActivity, error)
	UserActivityGet(context.Context, string) (*domain.UserActivity, error)
	UserActivityUpdate(context.Context, *domain.UserActivity) error
	DistGetByUID(context.Context, string) (*domain.Dist, error)
	DistCreate(context.Context, *domain.Dist) error
	DistUpdate(context.Context, *domain.Dist) error
	ConfigGet(context.Context) (*domain.Config, error)
	UserPlaceAggregateWith(context.Context, interface{}, []string, time.Time) ([]*domain.TeamUser, error)
	UserAggregateWith(context.Context, interface{}, []string, time.Time) ([]*domain.TeamUser, error)
	BuyGetAllWith(context.Context, interface{}, *options.FindOptions) ([]*domain.Buy, error)
	DistGetByPayUID(context.Context, string) (*domain.Dist, error)
	ClaimGetByCode(context.Context, string) (*domain.Claim, error)
	UserClaimGetLastByClaimCodeAndTypeAndUserUID(context.Context, string, domain.UserClaimType, string) (*domain.UserClaim, error)
	UserBalanceGetByUserUIDAndCurrencyCode(context.Context, string, string) (*domain.UserBalance, error)
	UserBalanceCreate(context.Context, *domain.UserBalance) error
	UserBalanceUpdate(context.Context, *domain.UserBalance) error
	UserClaimCreate(context.Context, *domain.UserClaim) error
	UserCountByRefUID(context.Context, string) (int64, error)
	UserGetAllWithCountByRefUID(context.Context, string, int64, int64) ([]*domain.User, error)
	UserPlaceCountByMatchUID(context.Context, string) (int64, error)
	UserPlaceDelete(context.Context, string) error
	UserProductGetAllByUserUIDAndProductCategoryAndDate(context.Context, string, string, time.Time) ([]*domain.UserProduct, error)

	// Partner Application methods
	PartnerApplicationCreate(context.Context, *domain.PartnerApplication) error
	PartnerApplicationGetByUID(context.Context, string) (*domain.PartnerApplication, error)
	PartnerApplicationGetAllByPartnerUID(context.Context, string, int64, int64) ([]*domain.PartnerApplication, error)
	PartnerApplicationGetAllByApplicantUID(context.Context, string, int64, int64) ([]*domain.PartnerApplication, error)
	PartnerApplicationUpdate(context.Context, *domain.PartnerApplication) error
	PartnerApplicationCountByPartnerUID(context.Context, string) (int64, error)
	PartnerApplicationCountByApplicantUID(context.Context, string) (int64, error)
	PartnerApplicationGetExpired(context.Context, time.Duration) ([]*domain.PartnerApplication, error)
}

type AuthService interface {
	ChargeLineCoin(context.Context, *domain.Claim, string, int64, string, int, []int64) error
}

//go:generate mockery --dir . --name PBClient --output ./mocks
type PBClient interface {
	OrderCreate(ctx context.Context, order domain.OrderCreateReq) (*domain.CreateOrderResp, error)
	OrderCancel(ctx context.Context, id uuid.UUID) (bool, error)
	OrderInfo(ctx context.Context, id uuid.UUID) (*domain.OrderInfo, error)
}

type Service struct {
	db   DbRepository
	pb   PBClient
	auth AuthService
}

func NewTeamService(db DbRepository, pb PBClient, auth AuthService) *Service {
	return &Service{db, pb, auth}
}

// PlaceGetNew get the first free place in team and instantiate new place without user uid
func (s *Service) PlaceGetNew(ctx context.Context, userPlace *domain.UserPlace, side domain.UserTeamType) (*domain.UserPlace, error) {

	dUserPlaceMax, err := s.PlaceGetLastRow(ctx, userPlace, side, maxRow)
	if dUserPlaceMax == nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	newPlaceRow := new(big.Int).Set(dUserPlaceMax.Row)
	newPlaceRow.Add(newPlaceRow, big.NewInt(1))
	newPlaceCol := new(big.Int).Set(dUserPlaceMax.Col)
	newPlaceCol.Mul(newPlaceCol, big.NewInt(2))
	if side == domain.UserTeamTypeLeft {
		newPlaceCol.Sub(newPlaceCol, big.NewInt(1))
	}

	return &domain.UserPlace{
		UserUID:   "",
		MatchUID:  userPlace.UserUID,
		Row:       newPlaceRow,
		Col:       newPlaceCol,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// PlaceGetSideLastDown get all team side places down
func (s *Service) PlaceGetSideLastDown(ctx context.Context, place *domain.UserPlace, side domain.UserTeamType, rows *big.Int) (*domain.UserPlace, error) {
	// TODO: remove bson find options mongo-driver dependency from services
	filter := bson.M{}
	orFilter := s.getFilterForSideDown(place, rows, side)
	if len(orFilter) == 0 {
		return place, nil
	}
	filter["$or"] = orFilter
	sort := options.FindOne().SetSort(bson.D{{"row", -1}}).SetCollation(&options.Collation{
		Locale:          "en_US",
		NumericOrdering: true,
	})

	dUserPlace, err := s.db.UserPlaceGet(ctx, filter, sort)
	if err != nil {
		if domain.ErrorIs(err, domain.ErrNoDocuments) {
			return place, nil
		}
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dUserPlace, nil
}

// PlaceGetLastRow get all team side places down
func (s *Service) PlaceGetLastRow(ctx context.Context, place *domain.UserPlace, side domain.UserTeamType, rows *big.Int) (*domain.UserPlace, error) {
	// TODO: remove bson find options mongo-driver dependency from services
	filter := bson.M{}
	var orFilter []bson.M
	if side == domain.UserTeamTypeUndefined {
		orFilter = s.getFilterForTeamDown(place, rows, side)
	} else {
		orFilter = s.getFilterForSideDown(place, rows, side)
	}
	if len(orFilter) == 0 {
		return place, nil
	}
	filter["$or"] = orFilter
	sort := options.FindOne().SetSort(bson.D{{"row", -1}}).SetCollation(&options.Collation{
		Locale:          "en_US",
		NumericOrdering: true,
	})

	dUserPlace, err := s.db.UserPlaceGet(ctx, filter, sort)
	if err != nil {
		if domain.ErrorIs(err, domain.ErrNoDocuments) {
			return place, nil
		}
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dUserPlace, nil
}

// PlaceGetAllDown get all team tree places down
func (s *Service) PlaceGetAllDown(ctx context.Context, place *domain.UserPlace, side domain.UserTeamType, rows *big.Int) ([]*domain.UserPlace, error) {
	// TODO: remove bson find options mongo-driver dependency from services
	filter := bson.M{}
	orFilter := s.getFilterForTeamDown(place, rows, side)
	if len(orFilter) == 0 {
		return nil, nil
	}
	filter["$or"] = orFilter
	sort := options.Find().SetSort(bson.D{{"row", 1}, {"col", 1}}).SetCollation(&options.Collation{
		Locale:          "en_US",
		NumericOrdering: true,
	})

	dUserPlaces, err := s.db.UserPlaceGetAll(ctx, filter, sort)
	if err != nil {
		if domain.ErrorIs(err, domain.ErrNoDocuments) {
			return nil, nil
		}
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return dUserPlaces, nil
}

// PlaceGetRefUp find nearest existing place up by ref UID
func (s *Service) PlaceGetRefUp(ctx context.Context, user *domain.User) (*domain.UserPlace, error) {
	// get parent user place
	// TODO: think about recursion
	dUserPlace, _ := s.db.UserPlaceGetByUserUID(ctx, user.UID)
	if dUserPlace == nil {
		// get parent user
		dUser, err := s.db.UserGetByUID(ctx, user.RefUID)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		// TODO: find a way to get rid of recursion
		return s.PlaceGetRefUp(ctx, dUser)
	}

	return dUserPlace, nil
}

// PlaceGetRefUpByBuy find nearest existing place up by ref UID
func (s *Service) PlaceGetRefUpByBuy(ctx context.Context, buy *domain.Buy) (*domain.UserPlace, error) {
	// get parent user
	dUser, err := s.db.UserGetByUID(ctx, buy.UserUID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	dRefUser, err := s.db.UserGetByUID(ctx, dUser.RefUID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return s.PlaceGetRefUp(ctx, dRefUser)
}

// PlaceGetAllRefDown find all ref places down and pluck user uid
func (s *Service) PlaceGetAllRefDown(ctx context.Context, place *domain.UserPlace, rows *big.Int) ([]*domain.UserPlace, error) {
	refPlaces := make([]*domain.UserPlace, 0)
	row := new(big.Int).Set(zero)
	dUserPlace, _ := s.db.UserPlaceGetByMatchUID(ctx, place.UserUID)

	for dUserPlace != nil && dUserPlace.UserUID != "" && row.Cmp(rows) < 0 {
		// get child user place by ref uid
		row.Add(row, one)
		refPlaces = append(refPlaces, dUserPlace)
		dUserPlace, _ = s.db.UserPlaceGetByMatchUID(ctx, dUserPlace.UserUID)
	}
	return refPlaces, nil
}

// PlaceSumAllCv sum all cv from buys in the team
func (s *Service) PlaceSumAllCv(ctx context.Context, newBuy *domain.Buy, place *domain.UserPlace, side domain.UserTeamType, rows *big.Int) (int64, error) {
	filter := bson.M{}
	filter["paidAt"] = bson.M{"$lte": newBuy.PaidAt, "$gt": time.Time{}}
	orFilter := s.getFilterForTeamDown(place, rows, side)
	if len(orFilter) == 0 {
		return 0, nil
	}
	filter["$or"] = orFilter
	sum, err := s.db.BuySumAllByField(ctx, filter, "cv")
	if err != nil {
		return 0, err
	}

	// TODO: add smart cache by date

	return sum, nil
}

// PlaceSumAllCvToDate sum all cv from buys in the team
func (s *Service) PlaceSumAllCvToDate(ctx context.Context, place *domain.UserPlace, side domain.UserTeamType, rows *big.Int, dateTo time.Time) (int64, error) {
	filter := bson.M{}
	filter["paidAt"] = bson.M{"$lt": dateTo, "$gt": time.Time{}}
	orFilter := s.getFilterForTeamDown(place, rows, side)
	if len(orFilter) == 0 {
		return 0, nil
	}
	filter["$or"] = orFilter
	sum, err := s.db.BuySumAllByField(ctx, filter, "cv")
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// PlaceSumAllCvFromDateToDate sum all cv from buys in the team
func (s *Service) PlaceSumAllCvFromDateToDate(ctx context.Context, place *domain.UserPlace, side domain.UserTeamType, rows *big.Int, dateFrom time.Time, dateTo time.Time) (int64, error) {
	filter := bson.M{}
	filter["paidAt"] = bson.M{"$lte": dateTo, "$gte": dateFrom}
	orFilter := s.getFilterForTeamDown(place, rows, side)
	if len(orFilter) == 0 {
		return 0, nil
	}
	filter["$or"] = orFilter
	sum, err := s.db.BuySumAllByField(ctx, filter, "cv")
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// PlaceCountAllTeamMonth count all ranks in the team
func (s *Service) PlaceCountAllTeamMonth(ctx context.Context, place *domain.UserPlace, side domain.UserTeamType, rows *big.Int, currentTime time.Time) (map[string]uint64, error) {
	filter := bson.M{}
	filter["startAt"] = bson.M{"$gte": domain.BeginningOfNextMonth(currentTime)}
	filter["refUid"] = place.UserUID
	orFilter := s.getFilterForTeamDown(place, rows, side)
	if len(orFilter) == 0 {
		return map[string]uint64{}, nil
	}
	filter["$or"] = orFilter
	counts, err := s.db.UserRankCountAllByRankCode(ctx, filter)
	if err != nil {
		return counts, err
	}
	return counts, nil
}

// PlaceCountRefTeamMonth count ref ranks in the team
func (s *Service) PlaceCountRefTeamMonth(ctx context.Context, place *domain.UserPlace, side domain.UserTeamType, rows *big.Int, currentTime time.Time) (map[string]uint64, error) {
	refPlaces, err := s.PlaceGetAllRefDown(ctx, place, rows)
	refPlacesUids := make([]string, 0, len(refPlaces))
	for _, pl := range refPlaces {
		refPlacesUids = append(refPlacesUids, pl.UserUID)
	}
	filter := bson.M{}
	filter["startAt"] = bson.M{"$gte": domain.BeginningOfNextMonth(currentTime)}
	filter["userUid"] = bson.M{"$in": refPlacesUids}
	orFilter := s.getFilterForTeamDown(place, rows, side)
	if len(orFilter) == 0 {
		return map[string]uint64{}, nil
	}
	filter["$or"] = orFilter
	counts, err := s.db.UserRankCountAllByRankCode(ctx, filter)
	if err != nil {
		return counts, err
	}
	return counts, nil
}
