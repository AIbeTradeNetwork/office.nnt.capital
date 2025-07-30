package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"server/internal/config"
	"server/internal/domain"
)

const (
	// table name in DB
	productTable = "product"

	// errors prefix
	productErrorSource = "[repository.mongodb.product]"
)

type productDB struct {
	Code         string        `bson:"code"`
	Name         string        `bson:"name"`
	Description  string        `bson:"description"`
	Category     string        `bson:"category"`
	Period       time.Duration `bson:"period"`
	Price        int64         `bson:"price"`
	RetailPrice  int64         `bson:"retailPrice"`
	Cv           int64         `bson:"cv"`
	CurrencyCode string        `bson:"currencyCode"`
	Precision    uint8         `bson:"precision"`
	IsActive     bool          `bson:"isActive"`
	Priority     int64         `bson:"priority"`
	Multiplier   int64         `bson:"multiplier"`
	Limit        int64         `bson:"limit"`
	Count        int64         `bson:"count"`
	Sort         int64         `bson:"sort"`
}

func (mr *Repo) ProductGetByCode(ctx context.Context, code string) (*domain.Product, error) {
	productDb := &productDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(productTable).
		FindOne(ctx, bson.M{"code": code}).Decode(&productDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(productErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(productErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.productConvertFromDb(productDb), nil
}

func (mr *Repo) ProductGetAll(ctx context.Context) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)
	cfg := config.Get()
	cursor, err := mr.db.Database(cfg.MongoDB).Collection(productTable).Find(ctx, bson.M{})
	if err != nil {
		return nil, domain.NewError(productErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		productDb := &productDB{}
		err = cursor.Decode(&productDb)
		if err != nil {
			return nil, domain.NewError(productErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		products = append(products, mr.productConvertFromDb(productDb))
	}
	return products, nil
}

func (mr *Repo) ProductGetAllIsActive(ctx context.Context) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)
	cfg := config.Get()
	opts := options.Find().SetSort(bson.M{"sort": 1})
	cursor, err := mr.db.Database(cfg.MongoDB).Collection(productTable).Find(ctx, bson.M{"isActive": true}, opts)
	if err != nil {
		return nil, domain.NewError(productErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		productDb := &productDB{}
		err = cursor.Decode(&productDb)
		if err != nil {
			return nil, domain.NewError(productErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		products = append(products, mr.productConvertFromDb(productDb))
	}
	return products, nil
}

func (mr *Repo) ProductGetAllIsActiveByCategory(ctx context.Context, category string) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)
	cfg := config.Get()
	cursor, err := mr.db.Database(cfg.MongoDB).Collection(productTable).Find(ctx, bson.M{"isActive": true, "category": category})
	if err != nil {
		return nil, domain.NewError(productErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		productDb := &productDB{}
		err = cursor.Decode(&productDb)
		if err != nil {
			return nil, domain.NewError(productErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		products = append(products, mr.productConvertFromDb(productDb))
	}
	return products, nil
}

func (mr *Repo) ProductIncCount(ctx context.Context, code string) error {
	cfg := config.Get()

	_, err := mr.db.Database(cfg.MongoDB).Collection(productTable).
		UpdateOne(ctx, bson.M{"code": code}, bson.M{"$inc": bson.M{"count": 1}})
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) ProductCreate(ctx context.Context, product *domain.Product) error {
	productDb := mr.productConvertToDb(product)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(productTable).
		InsertOne(ctx, productDb)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) productConvertToDb(product *domain.Product) *productDB {
	return &productDB{
		Code:         product.Code,
		Name:         product.Name,
		Description:  product.Description,
		Category:     product.Category,
		Period:       product.Period,
		Price:        product.Price,
		RetailPrice:  product.RetailPrice,
		Cv:           product.Cv,
		CurrencyCode: product.CurrencyCode,
		Precision:    product.Precision,
		IsActive:     product.IsActive,
		Priority:     product.Priority,
		Multiplier:   product.Multiplier,
		Limit:        product.Limit,
		Count:        product.Count,
		Sort:         product.Sort,
	}
}

func (mr *Repo) productConvertFromDb(product *productDB) *domain.Product {
	return &domain.Product{
		Code:         product.Code,
		Name:         product.Name,
		Description:  product.Description,
		Category:     product.Category,
		Period:       product.Period,
		CurrencyCode: product.CurrencyCode,
		Precision:    product.Precision,
		RetailPrice:  product.RetailPrice,
		Price:        product.Price,
		Cv:           product.Cv,
		Priority:     product.Priority,
		IsActive:     product.IsActive,
		Multiplier:   product.Multiplier,
		Limit:        product.Limit,
		Count:        product.Count,
		Sort:         product.Sort,
	}
}

func (mr *Repo) productEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	planIndexes := []mongo.IndexModel{
		{Keys: bson.M{"code": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"category": 1}, Options: options.Index()},
		{Keys: bson.M{"priority": -1}, Options: options.Index()},
		{Keys: bson.M{"sort": 1}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(productTable).Indexes().CreateMany(ctx, planIndexes)
	if err != nil {
		return domain.NewError(productErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
