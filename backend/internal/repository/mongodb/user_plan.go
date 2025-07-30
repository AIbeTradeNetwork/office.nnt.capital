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
	userPlanTable = "user_plan"

	// errors prefix
	userPlanErrorSource = "[repository.mongodb.user_plan]"
)

type userPlanDB struct {
	UID      string    `bson:"uid"`
	UserUID  string    `bson:"userUid"`
	BuyUID   string    `bson:"buyUid"`
	PlanCode string    `bson:"planCode"`
	StartAt  time.Time `bson:"startAt"`
	EndAt    time.Time `bson:"endAt"`
	Priority int64     `bson:"priority"`
}

func (mr *Repo) UserPlanGetByDate(ctx context.Context, userUid string, currentDate time.Time) (*domain.UserPlan, error) {
	userPlanDb := &userPlanDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
		"startAt": bson.M{"$lte": currentDate},
		"endAt":   bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userPlanTable).
		FindOne(ctx, filter, opts).Decode(&userPlanDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userPlanConvertFromDb(userPlanDb), nil
}

func (mr *Repo) UserPlanGetByCodeAndDate(ctx context.Context, userUid string, planCode string, currentDate time.Time) (*domain.UserPlan, error) {
	userPlanDb := &userPlanDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":  userUid,
		"planCode": planCode,
		"startAt":  bson.M{"$lte": currentDate},
		"endAt":    bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userPlanTable).
		FindOne(ctx, filter, opts).Decode(&userPlanDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userPlanConvertFromDb(userPlanDb), nil
}

func (mr *Repo) UserPlanGetLastByCodeAndDate(ctx context.Context, userUid string, planCode string, currentDate time.Time) (*domain.UserPlan, error) {
	userPlanDb := &userPlanDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":  userUid,
		"planCode": planCode,
		"endAt":    bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.D{{"endAt", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userPlanTable).
		FindOne(ctx, filter, opts).Decode(&userPlanDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userPlanConvertFromDb(userPlanDb), nil
}

func (mr *Repo) UserPlanGetByUID(ctx context.Context, uid string) (*domain.UserPlan, error) {
	userPlanDb := &userPlanDB{}
	cfg := config.Get()

	filter := bson.M{
		"uid": uid,
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userPlanTable).
		FindOne(ctx, filter, opts).Decode(&userPlanDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userPlanConvertFromDb(userPlanDb), nil
}

func (mr *Repo) UserPlanGetAllWith(ctx context.Context, filter interface{}, options *options.FindOptions) ([]*domain.UserPlan, error) {
	cfg := config.Get()

	cur, err := mr.db.Database(cfg.MongoDB).Collection(userPlanTable).
		Find(ctx, filter, options)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cur.Close(ctx)
	all := make([]*domain.UserPlan, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		dbObj := &userPlanDB{}
		err = cur.Decode(dbObj)

		if err != nil {
			return nil, err
		}

		all = append(all, mr.userPlanConvertFromDb(dbObj))
	}
	err = cur.Err()
	if err != nil {
		return nil, domain.NewError(userPlanErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return all, nil
}

func (mr *Repo) UserPlanGetAllByUserUIDAndDate(ctx context.Context, userUid string, curDate time.Time) ([]*domain.UserPlan, error) {
	filter := bson.M{
		"userUid": userUid,
		"endAt":   bson.M{"$gte": curDate},
	}
	opts := options.Find().SetSort(bson.M{"priority": -1})
	return mr.UserPlanGetAllWith(ctx, filter, opts)
}

func (mr *Repo) UserPlanCreate(ctx context.Context, plan *domain.UserPlan) error {
	userPlanDb := mr.userPlanConvertToDb(plan)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userPlanTable).
		InsertOne(ctx, userPlanDb)
	if err != nil {
		return domain.NewError(userPlanErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) UserPlanUpdate(ctx context.Context, plan *domain.UserPlan) error {
	userPlanDb := mr.userPlanConvertToDb(plan)
	update := bson.M{"$set": userPlanDb}
	cfg := config.Get()
	l := logger.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userPlanTable).
		UpdateOne(ctx, bson.D{{Key: "userUid", Value: userPlanDb.UserUID}}, update)
	if err != nil {
		l.Error(err)
		return domain.NewError(userPlanErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) userPlanConvertToDb(plan *domain.UserPlan) *userPlanDB {
	return &userPlanDB{
		UID:      plan.UID,
		UserUID:  plan.UserUID,
		BuyUID:   plan.BuyUID,
		PlanCode: plan.PlanCode,
		StartAt:  plan.StartAt,
		EndAt:    plan.EndAt,
		Priority: plan.Priority,
	}
}

func (mr *Repo) userPlanConvertFromDb(plan *userPlanDB) *domain.UserPlan {
	return &domain.UserPlan{
		UID:      plan.UID,
		UserUID:  plan.UserUID,
		BuyUID:   plan.BuyUID,
		PlanCode: plan.PlanCode,
		StartAt:  plan.StartAt,
		EndAt:    plan.EndAt,
		Priority: plan.Priority,
	}
}

func (mr *Repo) userPlanConvertFromDbList(plans []*userPlanDB) []*domain.UserPlan {
	converted := make([]*domain.UserPlan, len(plans))
	for i, plan := range plans {
		converted[i] = mr.userPlanConvertFromDb(plan)
	}
	return converted
}

func (mr *Repo) userPlanEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userPlanIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "uid", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "userUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "buyUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "priority", Value: -1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "startAt", Value: 1}, {Key: "endAt", Value: 1}}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(userPlanTable).Indexes().CreateMany(ctx, userPlanIndexes)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
