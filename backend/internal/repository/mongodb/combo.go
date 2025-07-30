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
	comboTable = "combo"

	// errors prefix
	comboErrorSource = "[repository.mongodb.combo]"
)

type comboDB struct {
	UID          string    `bson:"uid"`
	Code         string    `bson:"code"`
	Name         string    `bson:"name"`
	CurrencyCode string    `bson:"currencyCode"`
	Precision    uint8     `bson:"precision"`
	Amount       int64     `bson:"amount"`
	PriseCode    string    `bson:"priseCode"`
	Limit        int       `bson:"limit"`
	Count        int       `bson:"count"`
	StartAt      time.Time `bson:"startAt"`
	EndAt        time.Time `bson:"endAt"`
	IsActive     bool      `bson:"isActive"`
}

func (mr *Repo) ComboGetByCode(ctx context.Context, code string) (*domain.Combo, error) {
	cfg := config.Get()

	var comboDb *comboDB

	err := mr.db.Database(cfg.MongoDB).Collection(comboTable).
		FindOne(ctx, bson.M{"code": code}).Decode(&comboDb)
	if err != nil {
		return nil, domain.NewError(comboErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return mr.comboConvertFromDB(comboDb), nil
}

func (mr *Repo) ComboGetAll(ctx context.Context) ([]*domain.Combo, error) {
	var combos []*domain.Combo
	cfg := config.Get()
	cursor, err := mr.db.Database(cfg.MongoDB).Collection(comboTable).
		Find(ctx, bson.M{
			"isActive": true,
			"startAt":  bson.M{"$lte": time.Now().UTC()},
			"endAt":    bson.M{"$gte": time.Now().UTC()},
		})
	if err != nil {
		return nil, domain.NewError(comboErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		comboDb := &comboDB{}
		err = cursor.Decode(&comboDb)
		if err != nil {
			return nil, domain.NewError(comboErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		combos = append(combos, mr.comboConvertFromDB(comboDb))
	}
	return combos, nil
}

func (mr *Repo) ComboIncCount(ctx context.Context, code string) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(comboTable).
		UpdateOne(ctx, bson.M{"code": code}, bson.M{"$inc": bson.M{"count": 1}})
	if err != nil {
		return domain.NewError(comboErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) ComboCreate(ctx context.Context, combo *domain.Combo) error {
	cfg := config.Get()

	comboDb := mr.comboConvertToDB(combo)

	_, err := mr.db.Database(cfg.MongoDB).Collection(comboTable).InsertOne(ctx, comboDb)
	if err != nil {
		return domain.NewError(comboErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return nil
}

func (mr *Repo) comboConvertToDB(claim *domain.Combo) *comboDB {
	return &comboDB{
		UID:          claim.UID,
		Code:         claim.Code,
		Name:         claim.Name,
		CurrencyCode: claim.CurrencyCode,
		Precision:    claim.Precision,
		Amount:       claim.Amount,
		PriseCode:    claim.PriseCode,
		Limit:        claim.Limit,
		Count:        claim.Count,
		StartAt:      claim.StartAt,
		EndAt:        claim.EndAt,
		IsActive:     claim.IsActive,
	}
}

func (mr *Repo) comboConvertFromDB(comboDB *comboDB) *domain.Combo {
	return &domain.Combo{
		UID:          comboDB.UID,
		Code:         comboDB.Code,
		Name:         comboDB.Name,
		CurrencyCode: comboDB.CurrencyCode,
		Precision:    comboDB.Precision,
		Amount:       comboDB.Amount,
		PriseCode:    comboDB.PriseCode,
		Limit:        comboDB.Limit,
		Count:        comboDB.Count,
		StartAt:      comboDB.StartAt,
		EndAt:        comboDB.EndAt,
		IsActive:     comboDB.IsActive,
	}
}

func (mr *Repo) comboEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()

	_, err := mr.db.Database(cfg.MongoDB).Collection(comboTable).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"uid", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"code", 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return domain.NewError(comboErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
