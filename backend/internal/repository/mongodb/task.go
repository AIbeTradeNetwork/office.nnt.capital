package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/internal/config"
	"server/internal/domain"
	"time"
)

const (
	taskTable = "task"

	taskErrorSource = "[repository.mongodb.task]"
)

type taskDB struct {
	Code         string            `bson:"code"`
	Texts        map[string]string `bson:"texts"`
	IsActive     bool              `bson:"isActive"`
	StartAt      time.Time         `bson:"startAt"`
	EndAt        time.Time         `bson:"endAt"`
	CurrencyCode string            `bson:"currencyCode"`
	Precision    uint8             `bson:"precision"`
	Amount       int64             `bson:"amount"`
	Priority     int64             `bson:"priority"`
	Link         string            `bson:"link"`
	Completed    bool              `bson:"completed"`
	Locale       string            `bson:"locale"`
	Count        int64             `bson:"count"`
	Limit        int64             `bson:"limit"`
	IsApprove    bool              `bson:"isApprove"`
	RefUID       string            `bson:"refUid"`
	RefCount     int64             `bson:"refCount"`
}

func (mr *Repo) TaskCreate(ctx context.Context, task *domain.Task) error {
	cfg := config.Get()

	taskDb := mr.taskConvertToDB(task)
	_, err := mr.db.Database(cfg.MongoDB).Collection(taskTable).InsertOne(ctx, taskDb)
	if err != nil {
		return domain.NewError(taskErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return nil
}

func (mr *Repo) TaskIncCount(ctx context.Context, code string) error {
	cfg := config.Get()

	_, err := mr.db.Database(cfg.MongoDB).Collection(taskTable).
		UpdateOne(ctx, bson.M{"code": code}, bson.M{"$inc": bson.M{"count": 1}})
	if err != nil {
		return domain.NewError(taskErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) TaskDecCount(ctx context.Context, code string) error {
	cfg := config.Get()

	_, err := mr.db.Database(cfg.MongoDB).Collection(taskTable).
		UpdateOne(ctx, bson.M{"code": code}, bson.M{"$inc": bson.M{"count": -1}})
	if err != nil {
		return domain.NewError(taskErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) TaskIncRefCount(ctx context.Context, code string) error {
	cfg := config.Get()

	_, err := mr.db.Database(cfg.MongoDB).Collection(taskTable).
		UpdateOne(ctx, bson.M{"code": code}, bson.M{"$inc": bson.M{"refCount": 1}})
	if err != nil {
		return domain.NewError(taskErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) TaskGetByCode(ctx context.Context, code string) (*domain.Task, error) {
	cfg := config.Get()

	var taskDb *taskDB
	err := mr.db.Database(cfg.MongoDB).Collection(taskTable).
		FindOne(ctx, bson.M{"code": code}).Decode(&taskDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(taskErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(taskErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return mr.taskConvertFromDB(taskDb), nil
}

func (mr *Repo) TaskGetByRefUID(ctx context.Context, uid string) (*domain.Task, error) {
	cfg := config.Get()

	var taskDb *taskDB
	err := mr.db.Database(cfg.MongoDB).Collection(taskTable).
		FindOne(ctx, bson.M{"refUid": uid}).Decode(&taskDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(taskErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(taskErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return mr.taskConvertFromDB(taskDb), nil
}

func (mr *Repo) TaskGetAllWithCompletedByUserUIDAndLocale(ctx context.Context, userUid string, locale string) ([]*domain.Task, error) {
	cfg := config.Get()

	var tasks []*domain.Task
	cursor, err := mr.db.Database(cfg.MongoDB).Collection(taskTable).Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"isActive": true,
				"locale": bson.M{
					"$in": []string{locale, ""},
				},
				"$or": []bson.M{
					{
						"$expr": bson.M{
							"$lt": []any{"$count", "$limit"},
						},
					},
					{
						"limit": 0,
					},
				},
			},
		},
		{
			"$match": bson.M{
				"$or": []bson.M{
					{
						"refUid": bson.M{
							"$ne": "",
						},
						"$or": []bson.M{
							{
								"$expr": bson.M{
									"$gte": []any{
										"$refCount",
										bson.M{
											"$multiply": []any{
												"$count",
												0.8,
											},
										},
									},
								},
							},
							{
								"count": bson.M{
									"$lte": 1000,
								},
							},
						},
					},
					{
						"refUid": "",
					},
				},
			},
		},
		{
			"$sort": bson.M{
				"priority": 1,
			},
		},
		{
			"$lookup": bson.M{
				"from":         transactionTable,
				"localField":   "code",
				"foreignField": "taskCode",
				"as":           "transactions",
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"userUid": userUid,
							"type":    domain.TransactionTypeTask,
						},
					},
					{
						"$limit": 1,
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         userClaimTable,
				"localField":   "code",
				"foreignField": "taskCode",
				"as":           "claims",
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"userUid": userUid,
							"type":    domain.UserClaimTypeTask,
						},
					},
					{
						"$limit": 1,
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"completed": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$or": []bson.M{
								{"$eq": []interface{}{bson.M{"$size": "$transactions"}, 1}},
								{"$eq": []interface{}{bson.M{"$size": "$claims"}, 1}},
							},
						},
						"then": true,
						"else": false,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, domain.NewError(taskErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	for cursor.Next(ctx) {
		taskDb := &taskDB{}
		err = cursor.Decode(&taskDb)
		if err != nil {
			return nil, domain.NewError(taskErrorSource).SetCode(domain.ErrFind).Add(err)
		}

		tasks = append(tasks, mr.taskConvertFromDB(taskDb))
	}
	return tasks, nil
}

func (mr *Repo) taskConvertFromDB(taskDB *taskDB) *domain.Task {
	return &domain.Task{
		Code:         taskDB.Code,
		Texts:        taskDB.Texts,
		IsActive:     taskDB.IsActive,
		StartAt:      taskDB.StartAt,
		EndAt:        taskDB.EndAt,
		CurrencyCode: taskDB.CurrencyCode,
		Precision:    taskDB.Precision,
		Amount:       taskDB.Amount,
		Priority:     taskDB.Priority,
		Link:         taskDB.Link,
		Completed:    taskDB.Completed,
		Locale:       taskDB.Locale,
		Count:        taskDB.Count,
		Limit:        taskDB.Limit,
		IsApprove:    taskDB.IsApprove,
		RefUID:       taskDB.RefUID,
		RefCount:     taskDB.RefCount,
	}
}

func (mr *Repo) taskConvertToDB(task *domain.Task) *taskDB {
	return &taskDB{
		Code:         task.Code,
		Texts:        task.Texts,
		IsActive:     task.IsActive,
		StartAt:      task.StartAt,
		EndAt:        task.EndAt,
		CurrencyCode: task.CurrencyCode,
		Precision:    task.Precision,
		Amount:       task.Amount,
		Priority:     task.Priority,
		Link:         task.Link,
		Locale:       task.Locale,
		Count:        task.Count,
		Limit:        task.Limit,
		IsApprove:    task.IsApprove,
		RefUID:       task.RefUID,
		RefCount:     task.RefCount,
	}
}

func (mr *Repo) taskEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()

	indexes := []mongo.IndexModel{
		{Keys: bson.M{"code": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"isActive": 1}, Options: options.Index()},
		{Keys: bson.M{"priority": 1}, Options: options.Index()},
		{Keys: bson.M{"locale": 1}, Options: options.Index()},
		{Keys: bson.M{"count": 1}, Options: options.Index()},
		{Keys: bson.M{"limit": 1}, Options: options.Index()},
		{Keys: bson.M{"refUid": 1}, Options: options.Index()},
		{Keys: bson.M{"refCount": 1}, Options: options.Index()},
		{Keys: bson.D{{Key: "startAt", Value: 1}, {Key: "endAt", Value: 1}}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(taskTable).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return domain.NewError(taskErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
