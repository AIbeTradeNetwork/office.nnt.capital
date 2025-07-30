package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"server/internal/config"
	"server/internal/domain"
)

const (
	// table name in DB
	userAuthTable = "user_auth"

	// errors prefix
	userAuthErrorSource = "[repository.mongodb.user_auth]"
)

type userAuthDB struct {
	Type    domain.AuthType `bson:"type"`
	UserUID string          `bson:"userUid"`
	Token   string          `bson:"token"`
	Locale  string          `bson:"locale"`
}

func (mr *Repo) UserAuthGetByUserUIDAndType(ctx context.Context, uid string, authType domain.AuthType) (*domain.UserAuth, error) {
	userAuthDb := &userAuthDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userAuthTable).
		FindOne(ctx, bson.M{"userUid": uid, "type": authType}).Decode(&userAuthDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userAuthErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userAuthErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return &domain.UserAuth{
		UserUID: userAuthDb.UserUID,
		Token:   userAuthDb.Token,
	}, nil
}

func (mr *Repo) UserAuthGetAllUserTg(ctx context.Context, uid []string, authType domain.AuthType, limit int64, skip int64) ([]*domain.UserTg, error) {
	cfg := config.Get()

	filter := bson.M{"type": authType}
	if len(uid) > 0 {
		filter["userUid"] = bson.M{"$in": uid}
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$skip": skip},
		{"$limit": limit},
		{"$lookup": bson.M{
			"from":         userTable,
			"localField":   "userUid",
			"foreignField": "uid",
			"as":           "locale",
		}},
		{"$unwind": "$locale"},
		{"$addFields": bson.M{
			"locale": "$locale.locale",
		}},
	}

	cursor, err := mr.db.Database(cfg.MongoDB).Collection(userAuthTable).
		Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)
	userAuths := make([]*domain.UserTg, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		userAuthDb := &userAuthDB{}
		err = cursor.Decode(&userAuthDb)
		if err != nil {
			return nil, domain.NewError(userAuthErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		userAuths = append(userAuths, &domain.UserTg{
			UID:    userAuthDb.UserUID,
			TgUID:  userAuthDb.Token,
			Locale: userAuthDb.Locale,
		})
	}

	return userAuths, nil
}

func (mr *Repo) UserAuthGetByTokenAndType(ctx context.Context, token string, authType domain.AuthType) (*domain.UserAuth, error) {
	userAuthDb := &userAuthDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userAuthTable).
		FindOne(ctx, bson.M{"token": token, "type": authType}).Decode(&userAuthDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userAuthErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userAuthErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return &domain.UserAuth{
		UserUID: userAuthDb.UserUID,
		Token:   userAuthDb.Token,
	}, nil
}

func (mr *Repo) UserAuthCreate(ctx context.Context, auth *domain.UserAuth) error {
	userAuthDb := &userAuthDB{
		UserUID: auth.UserUID,
		Type:    auth.Type,
		Token:   auth.Token,
	}
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userAuthTable).
		InsertOne(ctx, userAuthDb)
	if err != nil {
		return domain.NewError(userAuthErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) UserAuthUpdate(ctx context.Context, auth *domain.UserAuth) error {
	userAuthDb := &userAuthDB{
		UserUID: auth.UserUID,
		Type:    auth.Type,
		Token:   auth.Token,
	}
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userAuthTable).
		UpdateOne(ctx, bson.M{"userUid": auth.UserUID, "type": auth.Type}, bson.M{"$set": userAuthDb})
	if err != nil {
		return domain.NewError(userAuthErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) userAuthEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userAuthIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "userUid", Value: 1}, {Key: "type", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"type": 1}, Options: options.Index()},
		{Keys: bson.M{"token": 1}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection("user_auth").Indexes().CreateMany(ctx, userAuthIndexes)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
