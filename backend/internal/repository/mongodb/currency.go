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
	// table name in DB
	currencyTable = "currency"

	// errors prefix
	currencyErrorSource = "[repository.mongodb.currency]"
)

type currencyDB struct {
	Code      string `bson:"code"`
	Precision uint8  `bson:"precision"`
}

func (mr *Repo) CurrencyGetAll(ctx context.Context) ([]*domain.Currency, error) {
	cfg := config.Get()
	currencies := make([]*domain.Currency, 0)
	cur, err := mr.db.Database(cfg.MongoDB).Collection(currencyTable).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var currency currencyDB
		err = cur.Decode(&currency)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, mr.currencyConvertFromDb(currency))
	}
	return currencies, nil
}

func (mr *Repo) CurrencyGetByCode(ctx context.Context, code string) (*domain.Currency, error) {
	cfg := config.Get()
	var currency currencyDB
	err := mr.db.Database(cfg.MongoDB).Collection(currencyTable).FindOne(ctx, bson.M{"code": code}).Decode(&currency)
	if err != nil {
		return nil, err
	}
	return mr.currencyConvertFromDb(currency), nil
}

func (mr *Repo) CurrencyCreate(ctx context.Context, currency *domain.Currency) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(currencyTable).InsertOne(ctx, currency)
	if err != nil {
		return errors.Wrap(err, currencyErrorSource)
	}
	return nil
}

func (mr *Repo) currencyConvertFromDb(currency currencyDB) *domain.Currency {
	return &domain.Currency{
		Code:      currency.Code,
		Precision: currency.Precision,
	}
}

func (mr *Repo) currencyEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()

	_, err := mr.db.Database(cfg.MongoDB).Collection(currencyTable).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"code", 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return domain.NewError(currencyErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	return nil
}
