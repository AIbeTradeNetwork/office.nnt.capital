package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"server/internal/config"
	"server/internal/domain"
)

const (
	// table name in DB
	userTable = "user"

	// errors prefix
	userErrorSource = "[repository.mongodb.user]"
)

type userDB struct {
	UID         string    `bson:"uid"`
	RefUID      string    `bson:"refUid"`
	LimRefUID   string    `bson:"limRefUid"`
	Nickname    string    `bson:"nickname"`
	FirstName   string    `bson:"firstName"`
	LastName    string    `bson:"lastName"`
	Email       string    `bson:"email"`
	PhotoUrl    string    `bson:"photoUrl"`
	Locale      string    `bson:"locale"`
	TgID        int64     `bson:"tgId"`
	TgUsername  string    `bson:"tgUsername"`
	CreatedAt   time.Time `bson:"createdAt"`
	UnlimInvite bool      `bson:"unlimInvite"`
	TonWallet   string    `bson:"tonWallet"`
	TeamCount   int64     `bson:"teamCount"`
}

type userAggregatedDB struct {
	User      *userDB         `bson:"user"`
	Place     *userPlaceDB    `bson:"place"`
	Plans     []*userPlanDB   `bson:"plans"`
	Ranks     []*userRankDB   `bson:"ranks"`
	Activity  *userActivityDB `bson:"activity"`
	TeamCount int64           `bson:"teamCount"`
}

func (mr *Repo) UserGetByUID(ctx context.Context, uid string) (*domain.User, error) {
	userDb := &userDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userTable).
		FindOne(ctx, bson.M{"uid": uid}).Decode(&userDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userConvertFromDb(userDb), nil
}

func (mr *Repo) UserGetByTonWallet(ctx context.Context, addr string) (*domain.User, error) {
	userDb := &userDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userTable).
		FindOne(ctx, bson.M{"tonWallet": addr}).Decode(&userDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userConvertFromDb(userDb), nil
}

func (mr *Repo) UserGetByEmail(ctx context.Context, email string) (*domain.User, error) {
	userDb := &userDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userTable).
		FindOne(ctx, bson.M{"email": email}).Decode(&userDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userConvertFromDb(userDb), nil
}

func (mr *Repo) UserGetByNickname(ctx context.Context, nickname string) (*domain.User, error) {
	userDb := &userDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userTable).
		FindOne(ctx, bson.M{"nickname": nickname}).Decode(&userDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userConvertFromDb(userDb), nil
}

func (mr *Repo) UserCountByRefUID(ctx context.Context, refUID string) (int64, error) {
	cfg := config.Get()
	count, err := mr.db.Database(cfg.MongoDB).Collection(userTable).
		CountDocuments(ctx, bson.M{"refUid": refUID})
	if err != nil {
		return 0, domain.NewError(userErrorSource).SetCode(domain.ErrCount).Add(err)
	}
	return count, nil
}

func (mr *Repo) UserGetAllWithCountByRefUID(ctx context.Context, uid string, limit int64, skip int64) ([]*domain.User, error) {
	cfg := config.Get()

	var users []*domain.User
	cur, err := mr.db.Database(cfg.MongoDB).Collection(userTable).
		Aggregate(ctx, []bson.M{
			{
				"$match": bson.M{"refUid": uid},
			},
			{
				"$sort": bson.M{"createdAt": -1},
			},
			{
				"$skip": skip,
			},
			{
				"$limit": limit,
			},
			{
				"$lookup": bson.M{
					"from":         userTable,
					"localField":   "uid",
					"foreignField": "refUid",
					"as":           "teamCount",
					"pipeline": []bson.M{
						{
							"$count": "teamCount",
						},
					},
				},
			},
			{
				"$unwind": bson.M{
					"path":                       "$teamCount",
					"preserveNullAndEmptyArrays": true,
				},
			},
			{
				"$addFields": bson.M{
					"teamCount": "$teamCount.teamCount",
				},
			},
		})

	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user userDB
		err := cur.Decode(&user)
		if err != nil {
			return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		users = append(users, mr.userConvertFromDb(&user))
	}
	return users, nil
}

func (mr *Repo) UserGetAllUpByUID(ctx context.Context, uid string, depth int) ([]*domain.User, error) {
	cfg := config.Get()

	type userWithDepth struct {
		*userDB
		Refs []*userDB `bson:"refs"`
	}

	var users []*userWithDepth
	cur, err := mr.db.Database(cfg.MongoDB).Collection(userTable).
		Aggregate(ctx, []bson.M{
			{
				"$match": bson.M{"uid": uid},
			},
			{
				"$limit": 1,
			},
			{
				"$graphLookup": bson.M{
					"from":             userTable,
					"startWith":        "$refUid",
					"connectFromField": "refUid",
					"connectToField":   "uid",
					"as":               "refs",
					"depthField":       "depth",
					"maxDepth":         depth,
				},
			},
		})

	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	err = cur.All(ctx, &users)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userConvertFromDbList(users[0].Refs), nil
}

func (mr *Repo) UserAggregateWith(ctx context.Context, filter interface{}, joins []string, curDate time.Time) ([]*domain.TeamUser, error) {
	cfg := config.Get()

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$replaceWith": bson.M{
				"user": "$$ROOT",
			},
		},
	}

	for _, join := range joins {
		switch join {
		case "place":
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         userPlaceTable,
					"localField":   "user.uid",
					"foreignField": "userUid",
					"as":           "place",
				},
			})
			pipeline = append(pipeline, bson.M{
				"$unwind": bson.M{
					"path":                       "$place",
					"preserveNullAndEmptyArrays": true,
				},
			})
		case "plans":
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         userPlanTable,
					"localField":   "user.uid",
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
					"localField":   "user.uid",
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
					"localField":   "user.uid",
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
					"from":         userTable,
					"localField":   "user.uid",
					"foreignField": "refUid",
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
			return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind)
		}
	}

	var res []*userAggregatedDB
	opts := options.Aggregate().SetCollation(&options.Collation{Locale: "en_US", NumericOrdering: true})
	cur, err := mr.db.Database(cfg.MongoDB).Collection(userTable).
		Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, &res)
	if err != nil {
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return mr.userAggregatedConvertFromDbList(res), nil
}

func (mr *Repo) UserCreate(ctx context.Context, user *domain.User) error {
	userDb := mr.userConvertToDb(user)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userTable).
		InsertOne(ctx, userDb)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) UserUpdate(ctx context.Context, user *domain.User) error {
	userDb := mr.userConvertToDb(user)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userTable).
		UpdateOne(ctx, bson.M{"uid": user.UID}, bson.M{"$set": userDb})
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	//if result.ModifiedCount == 0 {
	//	return domain.NewError(userErrorSource).SetCode(domain.ErrUpdate)
	//}
	return nil
}

func (mr *Repo) userConvertFromDb(userDb *userDB) *domain.User {
	return &domain.User{
		UID:         userDb.UID,
		RefUID:      userDb.RefUID,
		LimRefUID:   userDb.LimRefUID,
		Nickname:    userDb.Nickname,
		FirstName:   userDb.FirstName,
		LastName:    userDb.LastName,
		Email:       userDb.Email,
		PhotoUrl:    userDb.PhotoUrl,
		Locale:      userDb.Locale,
		TgID:        userDb.TgID,
		TgUsername:  userDb.TgUsername,
		CreatedAt:   userDb.CreatedAt,
		UnlimInvite: userDb.UnlimInvite,
		TonWallet:   userDb.TonWallet,
		TeamCount:   userDb.TeamCount,
	}
}

func (mr *Repo) userConvertToDb(user *domain.User) *userDB {
	return &userDB{
		UID:         user.UID,
		RefUID:      user.RefUID,
		LimRefUID:   user.LimRefUID,
		Nickname:    user.Nickname,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhotoUrl:    user.PhotoUrl,
		Locale:      user.Locale,
		TgID:        user.TgID,
		TgUsername:  user.TgUsername,
		CreatedAt:   user.CreatedAt,
		UnlimInvite: user.UnlimInvite,
		TonWallet:   user.TonWallet,
	}
}

func (mr *Repo) userAggregatedConvertFromDb(dbObj *userAggregatedDB) *domain.TeamUser {
	return &domain.TeamUser{
		User:      mr.userConvertFromDb(dbObj.User),
		Place:     mr.userPlaceConvertFromDb(dbObj.Place),
		Ranks:     mr.userRankConvertFromDbList(dbObj.Ranks),
		Plans:     mr.userPlanConvertFromDbList(dbObj.Plans),
		Activity:  mr.userActivityConvertFromDb(dbObj.Activity),
		TeamCount: dbObj.TeamCount,
	}
}

func (mr *Repo) userAggregatedConvertFromDbList(dbObjs []*userAggregatedDB) []*domain.TeamUser {
	res := make([]*domain.TeamUser, 0, len(dbObjs))
	for _, dbObj := range dbObjs {
		res = append(res, mr.userAggregatedConvertFromDb(dbObj))
	}
	return res
}

func (mr *Repo) userConvertFromDbList(dbObjs []*userDB) []*domain.User {
	res := make([]*domain.User, 0, len(dbObjs))
	for _, dbObj := range dbObjs {
		res = append(res, mr.userConvertFromDb(dbObj))
	}
	return res
}

func (mr *Repo) userEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userIndexes := []mongo.IndexModel{
		{Keys: bson.M{"uid": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"refUid": 1}, Options: options.Index()},
		{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"nickname": 1}, Options: options.Index().SetUnique(true)},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(userTable).Indexes().CreateMany(ctx, userIndexes)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
