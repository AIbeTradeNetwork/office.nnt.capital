package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/internal/config"
	"server/internal/domain"
)

const (
	// table name in DB
	userRankTable = "user_rank"

	// errors prefix
	userRankErrorSource = "[repository.mongodb.user_rank]"
)

type userRankDB struct {
	UID      string              `bson:"uid"`
	Type     domain.UserRankType `bson:"type"`
	UserUID  string              `bson:"userUid"`
	MatchUID string              `bson:"matchUid"`
	Row      *BigInt             `bson:"row"`
	Col      *BigInt             `bson:"col"`
	BuyUID   string              `bson:"buyUid"`
	RankCode string              `bson:"rankCode"`
	StartAt  time.Time           `bson:"startAt"`
	EndAt    time.Time           `bson:"endAt"`
	Priority int64               `bson:"priority"`
}

func (mr *Repo) UserRankGetByDate(ctx context.Context, userUid string, currentDate time.Time) (*domain.UserRank, error) {
	userRankDb := &userRankDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
		"startAt": bson.M{"$lte": currentDate},
		"endAt":   bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.M{"priority": -1})

	err := mr.db.Database(cfg.MongoDB).Collection(userRankTable).
		FindOne(ctx, filter, opts).Decode(&userRankDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userRankConvertFromDb(userRankDb), nil
}

func (mr *Repo) UserRankGetAllWith(ctx context.Context, filter interface{}, options *options.FindOptions) ([]*domain.UserRank, error) {
	cfg := config.Get()

	cur, err := mr.db.Database(cfg.MongoDB).Collection(userRankTable).
		Find(ctx, filter, options)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cur.Close(ctx)
	all := make([]*domain.UserRank, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		dbObj := &userRankDB{}
		err = cur.Decode(dbObj)

		if err != nil {
			return nil, err
		}

		all = append(all, mr.userRankConvertFromDb(dbObj))
	}
	err = cur.Err()
	if err != nil {
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return all, nil
}

func (mr *Repo) UserRankGetAllByUserUIDAndDate(ctx context.Context, userUid string, curDate time.Time) ([]*domain.UserRank, error) {
	filter := bson.M{
		"userUid": userUid,
		"endAt":   bson.M{"$gte": curDate},
	}
	opts := options.Find().SetSort(bson.M{"priority": -1})
	return mr.UserRankGetAllWith(ctx, filter, opts)
}

func (mr *Repo) UserRankGetByDateAndType(ctx context.Context, userUid string, currentDate time.Time, userRankType domain.UserRankType) (*domain.UserRank, error) {
	userRankDb := &userRankDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
		"type":    userRankType,
		"startAt": bson.M{"$lte": currentDate},
		"endAt":   bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.M{"priority": -1})

	err := mr.db.Database(cfg.MongoDB).Collection(userRankTable).
		FindOne(ctx, filter, opts).Decode(&userRankDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userRankConvertFromDb(userRankDb), nil
}

func (mr *Repo) UserRankCreate(ctx context.Context, rank *domain.UserRank) error {
	userRankDb := mr.userRankConvertToDb(rank)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userRankTable).
		InsertOne(ctx, userRankDb)
	if err != nil {
		return domain.NewError(userRankErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) UserRankGetByUserUID(ctx context.Context, userUid string) (*domain.UserRank, error) {
	userRankDb := &userRankDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
	}
	opts := options.FindOne().SetSort(bson.M{"priority": -1})

	err := mr.db.Database(cfg.MongoDB).Collection(userRankTable).
		FindOne(ctx, filter, opts).Decode(&userRankDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userRankConvertFromDb(userRankDb), nil
}

func (mr *Repo) UserRankGetByCodeAndDate(ctx context.Context, userUid string, rankCode string, currentDate time.Time) (*domain.UserRank, error) {
	userRankDb := &userRankDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":  userUid,
		"rankCode": rankCode,
		"startAt":  bson.M{"$lte": currentDate},
		"endAt":    bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.M{"priority": -1})

	err := mr.db.Database(cfg.MongoDB).Collection(userRankTable).
		FindOne(ctx, filter, opts).Decode(&userRankDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userRankConvertFromDb(userRankDb), nil
}

func (mr *Repo) UserRankGetByUID(ctx context.Context, uid string) (*domain.UserRank, error) {
	userRankDb := &userRankDB{}
	cfg := config.Get()

	filter := bson.M{
		"uid": uid,
	}
	opts := options.FindOne().SetSort(bson.M{"priority": -1})

	err := mr.db.Database(cfg.MongoDB).Collection(userRankTable).
		FindOne(ctx, filter, opts).Decode(&userRankDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userRankConvertFromDb(userRankDb), nil
}

func (mr *Repo) UserRankUpdate(ctx context.Context, rank *domain.UserRank) error {
	userRankDb := mr.userRankConvertToDb(rank)
	update := bson.M{"$set": userRankDb}
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userRankTable).
		UpdateOne(ctx, bson.M{"uid": userRankDb.UID}, update)
	if err != nil {
		return domain.NewError(userRankErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) UserRankCountAllByRankCode(ctx context.Context, filter interface{}) (map[string]uint64, error) {
	cfg := config.Get()
	rankMap := make(map[string]uint64)
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": "$rankCode", "count": bson.M{"$sum": 1}}},
	}
	type Result struct {
		ID    string `bson:"_id"`
		Count uint64 `bson:"count"`
	}
	opts := options.Aggregate().SetCollation(&options.Collation{Locale: "en_US", NumericOrdering: true})
	cur, err := mr.db.Database(cfg.MongoDB).Collection(userRankTable).
		Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var res Result
	for cur.Next(ctx) {
		err = cur.Decode(&res)
		if err != nil {
			return rankMap, err
		}
		rankMap[res.ID] = res.Count
	}

	return rankMap, nil
}

func (mr *Repo) userRankConvertToDb(rank *domain.UserRank) *userRankDB {
	return &userRankDB{
		UID:      rank.UID,
		Type:     rank.Type,
		UserUID:  rank.UserUID,
		MatchUID: rank.MatchUID,
		Row:      NewBigInt(rank.Row),
		Col:      NewBigInt(rank.Col),
		BuyUID:   rank.BuyUID,
		RankCode: rank.RankCode,
		StartAt:  rank.StartAt,
		EndAt:    rank.EndAt,
		Priority: rank.Priority,
	}
}

func (mr *Repo) userRankConvertFromDb(rank *userRankDB) *domain.UserRank {
	return &domain.UserRank{
		UID:      rank.UID,
		Type:     rank.Type,
		UserUID:  rank.UserUID,
		MatchUID: rank.MatchUID,
		Row:      rank.Row.Int(),
		Col:      rank.Col.Int(),
		BuyUID:   rank.BuyUID,
		RankCode: rank.RankCode,
		StartAt:  rank.StartAt,
		EndAt:    rank.EndAt,
		Priority: rank.Priority,
	}
}

func (mr *Repo) userRankConvertFromDbList(ranks []*userRankDB) []*domain.UserRank {
	result := make([]*domain.UserRank, len(ranks))
	for i, rank := range ranks {
		result[i] = mr.userRankConvertFromDb(rank)
	}
	return result
}

func (mr *Repo) userRankEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userPlanIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "uid", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "userUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "type", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "matchUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "buyUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "row", Value: 1}, {Key: "col", Value: 1}}, Options: options.Index().SetCollation(&options.Collation{
			Locale:          "en_US",
			NumericOrdering: true,
		})},
		{Keys: bson.D{{Key: "priority", Value: -1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "startAt", Value: 1}, {Key: "endAt", Value: 1}}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(userRankTable).Indexes().CreateMany(ctx, userPlanIndexes)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
