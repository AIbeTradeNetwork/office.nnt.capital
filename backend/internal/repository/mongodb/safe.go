package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/internal/config"
	"server/internal/domain"
)

const (
	safeTable = "safe"

	safeErrorSource = "[repository.mongodb.safe]"
)

func (mr *Repo) SafeCreate(ctx context.Context, safe *domain.Safe) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(safeTable).
		InsertOne(ctx, safe)
	if err != nil {
		return domain.NewError(safeErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) SafeUpdate(ctx context.Context, safe *domain.Safe) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(safeTable).
		UpdateOne(ctx, bson.M{"uid": safe.UID}, bson.M{"$set": safe})
	if err != nil {
		return domain.NewError(safeErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) SafeGetByUID(ctx context.Context, uid string) (*domain.Safe, error) {
	safeDb := &domain.Safe{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(safeTable).
		FindOne(ctx, bson.M{"uid": uid}).Decode(&safeDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(safeErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(safeErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return safeDb, nil
}

func (mr *Repo) safeEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	indexes := []mongo.IndexModel{
		{Keys: bson.M{"uid": 1}, Options: options.Index().SetUnique(true)},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(safeTable).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return domain.NewError(safeErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
