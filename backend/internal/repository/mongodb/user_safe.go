package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/internal/config"
	"server/internal/domain"
	"time"
)

const (
	// table name in DB
	userSafeTable = "user_safe"

	// errors prefix
	userSafeErrorSource = "[repository.mongodb.user_safe]"
)

func (mr *Repo) UserSafeCreate(ctx context.Context, userSafe *domain.UserSafe) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userSafeTable).
		InsertOne(ctx, userSafe)
	if err != nil {
		return domain.NewError(userSafeErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) UserSafeUpdate(ctx context.Context, userSafe *domain.UserSafe) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userSafeTable).
		UpdateOne(ctx, bson.M{"uid": userSafe.UID}, bson.M{"$set": userSafe})
	if err != nil {
		return domain.NewError(userSafeErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) UserSafeUpdateSecret(ctx context.Context, uid string, secret string) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userSafeTable).
		UpdateMany(ctx, bson.M{"uid": uid}, bson.M{"$set": bson.M{"secret": secret}})
	if err != nil {
		return domain.NewError(userSafeErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) UserSafeGetByUID(ctx context.Context, uid string) (*domain.UserSafe, error) {
	userSafeDb := &domain.UserSafe{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userSafeTable).
		FindOne(ctx, bson.M{"uid": uid}).Decode(&userSafeDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userSafeErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userSafeErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return userSafeDb, nil
}

func (mr *Repo) UserSafeGetActiveBySafeUIDsAndUserUID(ctx context.Context, safeUids []string, userUid string) (*domain.UserSafe, error) {
	userSafeDb := &domain.UserSafe{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userSafeTable).
		FindOne(ctx, bson.M{"safeUid": bson.M{"$in": safeUids}, "userUid": userUid, "claimedAt": time.Time{}}).Decode(&userSafeDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userSafeErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userSafeErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return userSafeDb, nil
}

func (mr *Repo) ensureUserSafeIndexes(ctx context.Context) error {
	cfg := config.Get()
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "uid", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "userUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "safeUid", Value: 1}}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(userSafeTable).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
