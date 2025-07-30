package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/internal/config"
	"server/internal/domain"
)

const (
	// table name in DB
	distTable = "dist"

	// errors prefix
	distErrorSource = "[repository.mongodb.distributor]"
)

type distDB struct {
	UID          string    `bson:"uid"`
	UserUID      string    `bson:"userUid"`
	RefUID       string    `bson:"refUid"`
	CurrencyCode string    `bson:"currencyCode"`
	Amount       int64     `bson:"amount"`
	CreatedAt    time.Time `bson:"createdAt"`
	PaidAt       time.Time `bson:"paidAt"`
	PayUID       string    `bson:"payUid"`
}

func (mr *Repo) DistGetByUID(ctx context.Context, uid string) (*domain.Dist, error) {
	distDb := &distDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(distTable).
		FindOne(ctx, bson.M{"uid": uid}).Decode(&distDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(distErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(distErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.distConvertFromDB(distDb), nil
}

func (mr *Repo) DistGetByPayUID(ctx context.Context, uid string) (*domain.Dist, error) {
	distDb := &distDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(distTable).
		FindOne(ctx, bson.M{"payUid": uid}).Decode(&distDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(distErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(distErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.distConvertFromDB(distDb), nil
}

func (mr *Repo) DistCreate(ctx context.Context, dist *domain.Dist) error {
	distDb := mr.distConvertToDB(dist)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(distTable).
		InsertOne(ctx, distDb)
	if err != nil {
		return domain.NewError(distErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) DistUpdate(ctx context.Context, dist *domain.Dist) error {
	distDb := mr.distConvertToDB(dist)
	update := bson.M{"$set": distDb}
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(distTable).
		UpdateOne(ctx, bson.D{{Key: "uid", Value: distDb.UID}}, update)
	if err != nil {
		return err
	}
	return nil
}

func (mr *Repo) distConvertToDB(dist *domain.Dist) *distDB {
	return &distDB{
		UID:          dist.UID,
		UserUID:      dist.UserUID,
		RefUID:       dist.RefUID,
		CurrencyCode: dist.CurrencyCode,
		Amount:       dist.Amount,
		CreatedAt:    dist.CreatedAt,
		PaidAt:       dist.PaidAt,
		PayUID:       dist.PayUID,
	}
}

func (mr *Repo) distConvertFromDB(distDB *distDB) *domain.Dist {
	return &domain.Dist{
		UID:          distDB.UID,
		UserUID:      distDB.UserUID,
		RefUID:       distDB.RefUID,
		CreatedAt:    distDB.CreatedAt,
		PaidAt:       distDB.PaidAt,
		CurrencyCode: distDB.CurrencyCode,
		Amount:       distDB.Amount,
		PayUID:       distDB.PayUID,
	}
}

func (mr *Repo) distEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userIndexes := []mongo.IndexModel{
		{Keys: bson.M{"uid": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"payUid": 1}, Options: options.Index()},
		{Keys: bson.D{{Key: "userUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "refUid", Value: 1}}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(distTable).Indexes().CreateMany(ctx, userIndexes)
	if err != nil {
		return domain.NewError(distErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
