package mongodb

import (
	"context"
	"errors"
	"fmt"
	"server/internal/config"
	"server/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// table name in DB
	userBalanceTable = "user_balance"

	// errors prefix
	userBalanceErrorSource = "[repository.mongodb.user_balance]"
)

type userBalanceDB struct {
	UserUID      string `bson:"userUid"`
	CurrencyCode string `bson:"currencyCode"`
	Precision    uint8  `bson:"precision"`
	Amount       int64  `bson:"amount"`
}

func (mr *Repo) UserBalanceCreate(ctx context.Context, userBalance *domain.UserBalance) error {
	cfg := config.Get()

	userBalanceDb := mr.userBalanceConvertToDB(userBalance)

	_, err := mr.db.Database(cfg.MongoDB).Collection(userBalanceTable).InsertOne(ctx, userBalanceDb)
	if err != nil {
		return domain.NewError(userBalanceErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return nil
}

func (mr *Repo) UserBalanceUpdate(ctx context.Context, userBalance *domain.UserBalance) error {
	cfg := config.Get()

	userBalanceDb := mr.userBalanceConvertToDB(userBalance)

	_, err := mr.db.Database(cfg.MongoDB).Collection(userBalanceTable).
		UpdateOne(ctx, bson.M{"userUid": userBalanceDb.UserUID, "currencyCode": userBalanceDb.CurrencyCode}, bson.M{"$set": userBalanceDb})
	if err != nil {
		return domain.NewError(userBalanceErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return nil
}

func (mr *Repo) UserBalanceChange(ctx context.Context, userUid string, currencyCode string, precision uint8, value int64) error {
	cfg := config.Get()

	filter := bson.M{"userUid": userUid, "currencyCode": currencyCode}
	opts := options.Update().SetUpsert(true)
	if value < 0 {
		filter["amount"] = bson.M{"$gte": value * -1}
		opts = options.Update()
	}

	r, err := mr.db.Database(cfg.MongoDB).Collection(userBalanceTable).
		UpdateOne(ctx, filter, bson.M{
			"$inc": bson.M{"amount": value},
			"$setOnInsert": bson.M{
				"userUid":      userUid,
				"currencyCode": currencyCode,
				"precision":    precision,
			},
		}, opts)

	if err != nil {
		return domain.NewError(userBalanceErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	if r.ModifiedCount != 1 && r.UpsertedCount != 1 && value != 0 {
		return domain.NewError(userBalanceErrorSource).SetCode(domain.ErrUpdate).Add(errors.New("balance not modified"))
	}

	return nil
}

func (mr *Repo) UserBalanceGetByUserUIDAndCurrencyCode(ctx context.Context, userUid string, currencyCode string) (*domain.UserBalance, error) {
	cfg := config.Get()

	fmt.Println("DEBUG: Поиск баланса:", userUid, currencyCode)

	var userBalanceDb *userBalanceDB
	err := mr.db.Database(cfg.MongoDB).Collection(userBalanceTable).
		FindOne(ctx, bson.M{"userUid": userUid, "currencyCode": currencyCode}).
		Decode(&userBalanceDb)
	if err != nil {
		fmt.Println("DEBUG: Баланс не найден или ошибка:", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userBalanceErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}

		return nil, domain.NewError(userBalanceErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	fmt.Println("DEBUG: Найденный баланс:", userBalanceDb)

	return mr.userBalanceConvertFromDB(userBalanceDb), nil
}

func (mr *Repo) userBalanceConvertToDB(userBalance *domain.UserBalance) *userBalanceDB {
	return &userBalanceDB{
		UserUID:      userBalance.UserUID,
		CurrencyCode: userBalance.CurrencyCode,
		Precision:    userBalance.Precision,
		Amount:       userBalance.Amount,
	}
}

func (mr *Repo) userBalanceConvertFromDB(userBalanceDB *userBalanceDB) *domain.UserBalance {
	return &domain.UserBalance{
		UserUID:      userBalanceDB.UserUID,
		CurrencyCode: userBalanceDB.CurrencyCode,
		Precision:    userBalanceDB.Precision,
		Amount:       userBalanceDB.Amount,
	}
}

func (mr *Repo) userBalanceEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()

	_, err := mr.db.Database(cfg.MongoDB).Collection(userBalanceTable).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"userUid", 1}, {"currencyCode", 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return domain.NewError(userBalanceErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	return nil
}
