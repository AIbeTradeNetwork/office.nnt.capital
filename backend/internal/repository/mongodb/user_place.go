package mongodb

import (
	"context"
	"log"
	"math/big"
	"server/internal/config"
	"server/internal/domain"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// table name in DB
	userPlaceTable = "user_place"

	// errors prefix
	userPlaceErrorSource = "[repository.mongodb.user_place]"
)

type userPlaceDB struct {
	UserUID   string    `bson:"userUid"`
	MatchUID  string    `bson:"matchUid"`
	Row       *BigInt   `bson:"row"`
	Col       *BigInt   `bson:"col"`
	CreatedAt time.Time `bson:"createdAt"`
}

func (mr *Repo) UserPlaceGetByUserUID(ctx context.Context, userUid string) (*domain.UserPlace, error) {
	userPlaceDb := &userPlaceDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).
		FindOne(ctx, bson.M{"userUid": userUid}).Decode(&userPlaceDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userPlaceConvertFromDb(userPlaceDb), nil
}

func (mr *Repo) UserPlaceGetByMatchUID(ctx context.Context, matchUid string) (*domain.UserPlace, error) {
	userPlaceDb := &userPlaceDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).
		FindOne(ctx, bson.M{"matchUid": matchUid}).Decode(&userPlaceDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userPlaceConvertFromDb(userPlaceDb), nil
}

func (mr *Repo) UserPlaceGetByRowCol(ctx context.Context, row *big.Int, col *big.Int) (*domain.UserPlace, error) {
	userPlaceDb := &userPlaceDB{}
	cfg := config.Get()
	opts := options.FindOne().SetCollation(&options.Collation{Locale: "en_US", NumericOrdering: true})
	err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).
		FindOne(ctx, bson.M{"row": row.String(), "col": col.String()}, opts).Decode(&userPlaceDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userPlaceConvertFromDb(userPlaceDb), nil
}

func (mr *Repo) UserPlaceGet(ctx context.Context, filter interface{}, options *options.FindOneOptions) (*domain.UserPlace, error) {
	userPlaceDb := &userPlaceDB{}
	cfg := config.Get()

	err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).
		FindOne(ctx, filter, options).Decode(&userPlaceDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userPlaceConvertFromDb(userPlaceDb), nil
}

func (mr *Repo) UserPlaceGetAll(ctx context.Context, filter interface{}, options *options.FindOptions) ([]*domain.UserPlace, error) {
	cfg := config.Get()

	cur, err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).
		Find(ctx, filter, options)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cur.Close(ctx)
	all := make([]*domain.UserPlace, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		dbObj := &userPlaceDB{}
		err = cur.Decode(dbObj)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		all = append(all, mr.userPlaceConvertFromDb(dbObj))
	}
	err = cur.Err()
	if err != nil {
		return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return all, nil
}

func (mr *Repo) UserPlaceGetAllWith(ctx context.Context, filter interface{}, options *options.FindOptions) ([]*domain.UserPlace, error) {
	cfg := config.Get()

	cur, err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).
		Find(ctx, filter, options)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cur.Close(ctx)
	all := make([]*domain.UserPlace, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		dbObj := &userPlaceDB{}
		err = cur.Decode(dbObj)

		if err != nil {
			return nil, err
		}

		all = append(all, mr.userPlaceConvertFromDb(dbObj))
	}
	err = cur.Err()
	if err != nil {
		return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return all, nil
}

func (mr *Repo) UserPlaceAggregateWith(ctx context.Context, filter interface{}, joins []string, curDate time.Time) ([]*domain.TeamUser, error) {
	cfg := config.Get()

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$replaceWith": bson.M{
				"place": "$$ROOT",
			},
		},
	}

	for _, join := range joins {
		switch join {
		case "user":
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         userTable,
					"localField":   "place.userUid",
					"foreignField": "uid",
					"as":           "user",
				},
			})
			pipeline = append(pipeline, bson.M{
				"$unwind": bson.M{
					"path":                       "$user",
					"preserveNullAndEmptyArrays": true,
				},
			})
		case "plans":
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         userPlanTable,
					"localField":   "place.userUid",
					"foreignField": "userUid",
					"as":           "plans",
					"pipeline": []bson.M{
						{
							"$match": bson.M{
								"startAt": bson.M{
									"$lte": curDate,
								},
								"endAt": bson.M{
									"$gte": curDate,
								},
							},
						},
						{
							"$sort": bson.M{
								"priority": -1,
							},
						},
					},
				},
			})
		case "ranks":
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         userRankTable,
					"localField":   "place.userUid",
					"foreignField": "userUid",
					"as":           "ranks",
					"pipeline": []bson.M{
						{
							"$match": bson.M{
								"startAt": bson.M{
									"$lte": curDate,
								},
								"endAt": bson.M{
									"$gte": curDate,
								},
							},
						},
						{
							"$sort": bson.M{
								"priority": -1,
							},
						},
					},
				},
			})
		case "activity":
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         userActivityTable,
					"localField":   "place.userUid",
					"foreignField": "userUid",
					"as":           "activity",
				},
			})
			pipeline = append(pipeline, bson.M{
				"$unwind": bson.M{
					"path":                       "$activity",
					"preserveNullAndEmptyArrays": true,
				},
			})
		case "team_count":
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         userPlaceTable,
					"localField":   "place.userUid",
					"foreignField": "matchUid",
					"as":           "team",
				},
			})
			pipeline = append(pipeline, bson.M{
				"$addFields": bson.M{
					"teamCount": bson.M{
						"$size": "$team",
					},
				},
			})
		default:
			return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind)
		}
	}

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"place.row": 1,
			"place.col": 1,
		},
	})

	var res []*domain.TeamUser
	opts := options.Aggregate().SetCollation(&options.Collation{Locale: "en_US", NumericOrdering: true})
	cur, err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).
		Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var r userAggregatedDB
		err = cur.Decode(&r)
		if err != nil {
			return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		res = append(res, mr.userAggregatedConvertFromDb(&r))
	}

	return res, nil
}

func (mr *Repo) UserPlaceCreate(ctx context.Context, place *domain.UserPlace) error {
	userPlaceDb := mr.userPlaceConvertToDb(place)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).
		InsertOne(ctx, userPlaceDb)
	if err != nil {
		return domain.NewError(userPlaceErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) UserPlaceCountByMatchUID(ctx context.Context, matchUID string) (int64, error) {
	cfg := config.Get()
	count, err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).
		CountDocuments(ctx, bson.M{"matchUid": matchUID})
	if err != nil {
		return 0, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrCount).Add(err)
	}
	return count, nil
}

func (mr *Repo) UserPlaceDelete(ctx context.Context, userUid string) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).
		DeleteOne(ctx, bson.M{"userUid": userUid})
	if err != nil {
		return domain.NewError(userPlaceErrorSource).SetCode(domain.ErrDelete).Add(err)
	}
	return nil
}

func (mr *Repo) userPlaceConvertToDb(place *domain.UserPlace) *userPlaceDB {
	return &userPlaceDB{
		UserUID:   place.UserUID,
		MatchUID:  place.MatchUID,
		Row:       NewBigInt(place.Row),
		Col:       NewBigInt(place.Col),
		CreatedAt: place.CreatedAt,
	}
}

func (mr *Repo) userPlaceConvertFromDb(dbObj *userPlaceDB) *domain.UserPlace {
	if dbObj == nil {
		return nil
	}
	return &domain.UserPlace{
		UserUID:   dbObj.UserUID,
		MatchUID:  dbObj.MatchUID,
		Row:       dbObj.Row.Int(),
		Col:       dbObj.Col.Int(),
		CreatedAt: dbObj.CreatedAt,
	}
}

func (mr *Repo) userPlaceEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userPlaceIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "row", Value: 1}, {Key: "col", Value: 1}}, Options: options.Index().SetUnique(true).SetCollation(&options.Collation{
			Locale:          "en_US",
			NumericOrdering: true,
		})},
		{Keys: bson.M{"userUid": 1}, Options: options.Index()},
		{Keys: bson.M{"refUid": 1}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(userPlaceTable).Indexes().CreateMany(ctx, userPlaceIndexes)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
