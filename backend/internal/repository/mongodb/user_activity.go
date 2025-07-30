package mongodb

import (
	"context"
	"errors"
	"server/internal/provider/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/internal/config"
	"server/internal/domain"
)

const (
	// table name in DB
	userActivityTable = "user_activity"

	// errors prefix
	userActivityErrorSource = "[repository.mongodb.user_activity]"
)

type userActivityDB struct {
	UserUID  string    `bson:"userUid"`
	StartAt  time.Time `bson:"startAt"`
	EndAt    time.Time `bson:"endAt"`
	CvAmount int64     `bson:"cvAmount"`
}

func (mr *Repo) UserActivityGetByDate(ctx context.Context, userUid string, currentDate time.Time) (*domain.UserActivity, error) {
	userActivityDb := &userActivityDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
		"startAt": bson.M{"$lte": currentDate},
		"endAt":   bson.M{"$gte": currentDate},
	}

	err := mr.db.Database(cfg.MongoDB).Collection(userActivityTable).
		FindOne(ctx, filter).Decode(&userActivityDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userActivityErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userActivityErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userActivityConvertFromDb(userActivityDb), nil
}

func (mr *Repo) UserActivityGet(ctx context.Context, userUid string) (*domain.UserActivity, error) {
	userActivityDb := &userActivityDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
	}

	err := mr.db.Database(cfg.MongoDB).Collection(userActivityTable).
		FindOne(ctx, filter).Decode(&userActivityDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userActivityErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userActivityErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userActivityConvertFromDb(userActivityDb), nil
}

func (mr *Repo) UserActivityCreate(ctx context.Context, act *domain.UserActivity) error {
	userActivityDb := mr.userActivityConvertToDb(act)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userActivityTable).
		InsertOne(ctx, userActivityDb)
	if err != nil {
		return domain.NewError(userActivityErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) UserActivityUpdate(ctx context.Context, act *domain.UserActivity) error {
	userActivityDb := mr.userActivityConvertToDb(act)
	update := bson.M{"$set": userActivityDb}
	cfg := config.Get()
	l := logger.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userActivityTable).
		UpdateOne(ctx, bson.M{"userUid": userActivityDb.UserUID}, update)
	if err != nil {
		l.Error(err)
		return domain.NewError(userActivityErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) userActivityConvertFromDb(dbObj *userActivityDB) *domain.UserActivity {
	if dbObj == nil {
		return nil
	}
	return &domain.UserActivity{
		UserUID:  dbObj.UserUID,
		StartAt:  dbObj.StartAt,
		EndAt:    dbObj.EndAt,
		CvAmount: dbObj.CvAmount,
	}
}

func (mr *Repo) userActivityConvertToDb(act *domain.UserActivity) *userActivityDB {
	return &userActivityDB{
		UserUID:  act.UserUID,
		StartAt:  act.StartAt,
		EndAt:    act.EndAt,
		CvAmount: act.CvAmount,
	}
}

func (mr *Repo) userActivityEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userPlanIndexes := []mongo.IndexModel{
		{Keys: bson.M{"userUid": 1}, Options: options.Index()},
		{Keys: bson.D{{Key: "startAt", Value: 1}, {Key: "endAt", Value: 1}}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(userActivityTable).Indexes().CreateMany(ctx, userPlanIndexes)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
