package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/internal/config"
	"server/internal/domain"
	"time"
)

const (
	// table name in DB
	claimTable = "claim"

	// errors prefix
	claimErrorSource = "[repository.mongodb.claim]"
)

type claimDB struct {
	Code         string        `bson:"code"`
	MaxPeriod    time.Duration `bson:"maxPeriod"`
	MinPeriod    time.Duration `bson:"minPeriod"`
	Amount       int64         `bson:"amount"`
	CurrencyCode string        `bson:"currencyCode"`
	Precision    uint8         `bson:"precision"`
}

func (mr *Repo) ClaimGetByCode(ctx context.Context, code string) (*domain.Claim, error) {
	cfg := config.Get()

	var claimDb *claimDB

	err := mr.db.Database(cfg.MongoDB).Collection(claimTable).
		FindOne(ctx, bson.M{"code": code}).Decode(&claimDb)
	if err != nil {
		return nil, domain.NewError(claimErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return mr.claimConvertFromDB(claimDb), nil
}

func (mr *Repo) ClaimCreate(ctx context.Context, claim *domain.Claim) error {
	cfg := config.Get()

	claimDb := mr.claimConvertToDB(claim)

	_, err := mr.db.Database(cfg.MongoDB).Collection(claimTable).InsertOne(ctx, claimDb)
	if err != nil {
		return domain.NewError(claimErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return nil
}

func (mr *Repo) claimConvertToDB(claim *domain.Claim) *claimDB {
	return &claimDB{
		Code:         claim.Code,
		MaxPeriod:    claim.MaxPeriod,
		MinPeriod:    claim.MinPeriod,
		Amount:       claim.Amount,
		CurrencyCode: claim.CurrencyCode,
		Precision:    claim.Precision,
	}
}

func (mr *Repo) claimConvertFromDB(claimDB *claimDB) *domain.Claim {
	return &domain.Claim{
		Code:         claimDB.Code,
		MaxPeriod:    claimDB.MaxPeriod,
		MinPeriod:    claimDB.MinPeriod,
		Amount:       claimDB.Amount,
		CurrencyCode: claimDB.CurrencyCode,
		Precision:    claimDB.Precision,
	}
}

func (mr *Repo) claimEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()

	_, err := mr.db.Database(cfg.MongoDB).Collection(claimTable).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"code", 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return domain.NewError(claimErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
