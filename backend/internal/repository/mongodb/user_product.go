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
	userProductTable = "user_product"

	// errors prefix
	userProductErrorSource = "[repository.mongodb.user_product]"
)

type userProductDB struct {
	UID             string    `bson:"uid"`
	UserUID         string    `bson:"userUid"`
	BuyUID          string    `bson:"buyUid"`
	ProductCode     string    `bson:"productCode"`
	ProductCategory string    `bson:"productCategory"`
	StartAt         time.Time `bson:"startAt"`
	EndAt           time.Time `bson:"endAt"`
	Priority        int64     `bson:"priority"`
	Multiplier      int64     `bson:"multiplier"`
	ComboUID        string    `bson:"comboUid"`
	FlowUID         string    `bson:"flowUid"`
}

func (mr *Repo) UserProductGetByDate(ctx context.Context, userUid string, currentDate time.Time) (*domain.UserProduct, error) {
	userProductDb := &userProductDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
		"startAt": bson.M{"$lte": currentDate},
		"endAt":   bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		FindOne(ctx, filter, opts).Decode(&userProductDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userProductConvertFromDb(userProductDb), nil
}

func (mr *Repo) UserProductGetByCodeAndDate(ctx context.Context, userUid string, productCode string, currentDate time.Time) (*domain.UserProduct, error) {
	userProductDb := &userProductDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":     userUid,
		"productCode": productCode,
		"startAt":     bson.M{"$lte": currentDate},
		"endAt":       bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		FindOne(ctx, filter, opts).Decode(&userProductDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userProductConvertFromDb(userProductDb), nil
}

func (mr *Repo) UserProductGetByCode(ctx context.Context, userUid string, productCode string) (*domain.UserProduct, error) {
	userProductDb := &userProductDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":     userUid,
		"productCode": productCode,
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		FindOne(ctx, filter, opts).Decode(&userProductDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userProductConvertFromDb(userProductDb), nil
}

func (mr *Repo) UserProductGetByCodes(ctx context.Context, userUid string, productCodes []string) (*domain.UserProduct, error) {
	userProductDb := &userProductDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
		"productCode": bson.M{
			"$in": productCodes,
		},
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		FindOne(ctx, filter, opts).Decode(&userProductDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userProductConvertFromDb(userProductDb), nil
}

func (mr *Repo) UserProductGetByCategoryAndDate(ctx context.Context, userUid string, productCategory string, currentDate time.Time) (*domain.UserProduct, error) {
	userProductDb := &userProductDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":         userUid,
		"productCategory": productCategory,
		"startAt":         bson.M{"$lte": currentDate},
		"endAt":           bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		FindOne(ctx, filter, opts).Decode(&userProductDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userProductConvertFromDb(userProductDb), nil
}

func (mr *Repo) UserProductGetByComboUID(ctx context.Context, userUid string, comboUid string) (*domain.UserProduct, error) {
	userProductDb := &userProductDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":  userUid,
		"comboUid": comboUid,
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		FindOne(ctx, filter, opts).Decode(&userProductDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userProductConvertFromDb(userProductDb), nil
}

func (mr *Repo) UserProductGetLastByCategoryAndDate(ctx context.Context, userUid string, productCategory string, currentDate time.Time) (*domain.UserProduct, error) {
	userProductDb := &userProductDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":         userUid,
		"productCategory": productCategory,
		"endAt":           bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.D{{"endAt", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		FindOne(ctx, filter, opts).Decode(&userProductDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userProductConvertFromDb(userProductDb), nil
}

func (mr *Repo) UserProductGetLastBoostAndDate(ctx context.Context, userUid string, currentDate time.Time) (*domain.UserProduct, error) {
	userProductDb := &userProductDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":         userUid,
		"productCategory": bson.M{"$in": []string{"boost_x3", "boost_x5", "boost_x10", "boost_x20", "boost_x50"}},
		"endAt":           bson.M{"$gte": currentDate},
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		FindOne(ctx, filter, opts).Decode(&userProductDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userProductConvertFromDb(userProductDb), nil
}

func (mr *Repo) UserProductGetByUID(ctx context.Context, uid string) (*domain.UserProduct, error) {
	userProductDb := &userProductDB{}
	cfg := config.Get()

	filter := bson.M{
		"uid": uid,
	}
	opts := options.FindOne().SetSort(bson.D{{"priority", -1}})

	err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		FindOne(ctx, filter, opts).Decode(&userProductDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userProductConvertFromDb(userProductDb), nil
}

func (mr *Repo) UserProductGetAllWith(ctx context.Context, filter interface{}, options *options.FindOptions) ([]*domain.UserProduct, error) {
	cfg := config.Get()

	cur, err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		Find(ctx, filter, options)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cur.Close(ctx)
	all := make([]*domain.UserProduct, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		dbObj := &userProductDB{}
		err = cur.Decode(dbObj)

		if err != nil {
			return nil, err
		}

		all = append(all, mr.userProductConvertFromDb(dbObj))
	}
	err = cur.Err()
	if err != nil {
		return nil, domain.NewError(userProductErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return all, nil
}

func (mr *Repo) UserProductGetAllByUserUIDAndDate(ctx context.Context, userUid string, curDate time.Time) ([]*domain.UserProduct, error) {
	filter := bson.M{
		"userUid": userUid,
		"endAt":   bson.M{"$gte": curDate},
	}
	opts := options.Find().SetSort(bson.M{"priority": -1})
	return mr.UserProductGetAllWith(ctx, filter, opts)
}

func (mr *Repo) UserProductGetAllByUserUIDAndProductCategoryAndDate(ctx context.Context, userUid string, productCategory string, curDate time.Time) ([]*domain.UserProduct, error) {
	filter := bson.M{
		"userUid":         userUid,
		"productCategory": productCategory,
		"endAt":           bson.M{"$gte": curDate},
	}
	opts := options.Find().SetSort(bson.M{"endAt": 1})
	return mr.UserProductGetAllWith(ctx, filter, opts)
}

func (mr *Repo) UserProductCreate(ctx context.Context, product *domain.UserProduct) error {
	userProductDb := mr.userProductConvertToDb(product)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		InsertOne(ctx, userProductDb)
	if err != nil {
		return domain.NewError(userProductErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) UserProductUpdate(ctx context.Context, product *domain.UserProduct) error {
	userProductDb := mr.userProductConvertToDb(product)
	update := bson.M{"$set": userProductDb}
	cfg := config.Get()
	l := logger.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).
		UpdateOne(ctx, bson.D{{Key: "uid", Value: userProductDb.UID}}, update)
	if err != nil {
		l.Error(err)
		return domain.NewError(userProductErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) userProductConvertToDb(product *domain.UserProduct) *userProductDB {
	return &userProductDB{
		UID:             product.UID,
		UserUID:         product.UserUID,
		BuyUID:          product.BuyUID,
		ProductCode:     product.ProductCode,
		ProductCategory: product.ProductCategory,
		StartAt:         product.StartAt,
		EndAt:           product.EndAt,
		Priority:        product.Priority,
		ComboUID:        product.ComboUID,
		Multiplier:      product.Multiplier,
		FlowUID:         product.FlowUID,
	}
}

func (mr *Repo) userProductConvertFromDb(product *userProductDB) *domain.UserProduct {
	return &domain.UserProduct{
		UID:             product.UID,
		UserUID:         product.UserUID,
		BuyUID:          product.BuyUID,
		ProductCode:     product.ProductCode,
		ProductCategory: product.ProductCategory,
		StartAt:         product.StartAt,
		EndAt:           product.EndAt,
		Priority:        product.Priority,
		ComboUID:        product.ComboUID,
		Multiplier:      product.Multiplier,
		FlowUID:         product.FlowUID,
	}
}

func (mr *Repo) userProductConvertFromDbList(products []*userProductDB) []*domain.UserProduct {
	converted := make([]*domain.UserProduct, len(products))
	for i, product := range products {
		converted[i] = mr.userProductConvertFromDb(product)
	}
	return converted
}

func (mr *Repo) userProductEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userProductIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "uid", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "userUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "buyUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "priority", Value: -1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "productCode", Value: -1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "comboUid", Value: -1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "productCategory", Value: -1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "startAt", Value: 1}, {Key: "endAt", Value: 1}}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(userProductTable).Indexes().CreateMany(ctx, userProductIndexes)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
