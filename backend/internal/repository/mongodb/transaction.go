package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"server/internal/config"
	"server/internal/domain"
)

const (
	// table name in DB
	transactionTable = "transaction"

	// errors prefix
	transactionErrorSource = "[repository.mongodb.transaction]"
)

type transactionDB struct {
	UID         string                      `bson:"uid"`
	UserUID     string                      `bson:"userUid"`
	FromUID     string                      `bson:"fromUid"`
	Percent     int64                       `bson:"percent"`
	Level       int                         `bson:"level"`
	Type        domain.TransactionType      `bson:"type"`
	RankCode    string                      `bson:"rankCode"`
	TaskCode    string                      `bson:"taskCode"`
	ComboUID    string                      `bson:"comboUid"`
	Amount      int64                       `bson:"amount"`
	PosAmount   int64                       `bson:"posAmount"`
	FullAmount  int64                       `bson:"fullAmount"`
	Coefficient int64                       `bson:"coefficient"`
	BuyUID      string                      `bson:"buyUid"`
	PayoutUID   string                      `bson:"payoutUid"`
	DepositUID  string                      `bson:"depositUid"`
	CreatedAt   time.Time                   `bson:"createdAt"`
	ChargedAt   time.Time                   `bson:"chargedAt"`
	MsgCodes    []domain.TransactionMsgCode `bson:"msgCode"`
}

func (mr *Repo) TransactionGetByUserUIDAndType(ctx context.Context, userUid string, transType domain.TransactionType) (*domain.Transaction, error) {
	transactionDb := &transactionDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
		"type":    transType,
	}

	err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		FindOne(ctx, filter).Decode(&transactionDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.transactionConvertFromDb(transactionDb), nil
}

func (mr *Repo) TransactionGetByUserUIDAndTypeAndRankCode(ctx context.Context, userUid string, transType domain.TransactionType, rankCode string) (*domain.Transaction, error) {
	transactionDb := &transactionDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":  userUid,
		"type":     transType,
		"rankCode": rankCode,
	}

	err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		FindOne(ctx, filter).Decode(&transactionDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.transactionConvertFromDb(transactionDb), nil
}

func (mr *Repo) TransactionGetByUserUIDAndTypeAndTaskCode(ctx context.Context, userUid string, transType domain.TransactionType, taskCode string) (*domain.Transaction, error) {
	transactionDb := &transactionDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":  userUid,
		"type":     transType,
		"taskCode": taskCode,
	}

	err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		FindOne(ctx, filter).Decode(&transactionDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.transactionConvertFromDb(transactionDb), nil
}

func (mr *Repo) TransactionGetByUserUIDAndTypeAndDepositUID(ctx context.Context, userUid string, transType domain.TransactionType, depUid string) (*domain.Transaction, error) {
	transactionDb := &transactionDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":    userUid,
		"type":       transType,
		"depositUid": depUid,
	}

	err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		FindOne(ctx, filter).Decode(&transactionDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.transactionConvertFromDb(transactionDb), nil
}

func (mr *Repo) TransactionGetByUserUIDAndTypeAndComboCode(ctx context.Context, userUid string, transType domain.TransactionType, comboCode string) (*domain.Transaction, error) {
	transactionDb := &transactionDB{}
	cfg := config.Get()

	filter := bson.M{
		"userUid":  userUid,
		"type":     transType,
		"comboUid": comboCode,
	}

	err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		FindOne(ctx, filter).Decode(&transactionDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userRankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.transactionConvertFromDb(transactionDb), nil
}

func (mr *Repo) TransactionSumByUserUIDAndTypeAndDates(ctx context.Context, userUid string, transType domain.TransactionType, dateFrom time.Time, dateTo time.Time) (int64, error) {
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
		"type":    transType,
		"createdAt": bson.M{
			"$gte": dateFrom,
			"$lte": dateTo,
		},
	}

	// TODO: add multi currency support
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": "", "sum": bson.M{"$sum": "$amount"}}},
	}

	type Result struct {
		ID  string `bson:"_id"`
		Sum int64  `bson:"sum"`
	}
	cur, err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, domain.NewError(transactionErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cur.Close(ctx)
	var res Result
	for cur.Next(ctx) {
		err = cur.Decode(&res)
		if err != nil {
			return 0, domain.NewError(transactionErrorSource).SetCode(domain.ErrFind).Add(err)
		}
	}
	return res.Sum, nil
}

func (mr *Repo) TransactionGetAllByUserUID(ctx context.Context, userUid string, curTime time.Time, limit int64, skip int64) ([]*domain.Transaction, error) {
	cfg := config.Get()

	filter := bson.M{
		"userUid": userUid,
		"chargedAt": bson.M{
			"$ne": nil,
			"$lt": curTime,
		},
	}

	opts := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetLimit(limit).SetSkip(skip)

	cur, err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(transactionErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(transactionErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cur.Close(ctx)
	all := make([]*domain.Transaction, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		dbObj := &transactionDB{}
		err = cur.Decode(dbObj)

		if err != nil {
			return nil, err
		}

		all = append(all, mr.transactionConvertFromDb(dbObj))
	}
	err = cur.Err()
	if err != nil {
		return nil, domain.NewError(transactionErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return all, nil
}

func (mr *Repo) TransactionBalanceByUserUIDAndDate(ctx context.Context, userUid string, dateTo time.Time) (int64, error) {
	cfg := config.Get()

	filter := bson.M{
		"payoutUid": "",
		"userUid":   userUid,
		"chargedAt": bson.M{
			"$lte": dateTo,
		},
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": "", "sum": bson.M{"$sum": "$amount"}}},
	}

	type Result struct {
		ID  string `bson:"_id"`
		Sum int64  `bson:"sum"`
	}
	cur, err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, domain.NewError(transactionErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cur.Close(ctx)
	var res Result
	for cur.Next(ctx) {
		err = cur.Decode(&res)
		if err != nil {
			return 0, domain.NewError(transactionErrorSource).SetCode(domain.ErrFind).Add(err)
		}
	}
	return res.Sum, nil
}

func (mr *Repo) TransactionUpdateAllWithUserUIDAndPayoutUIDAndDate(ctx context.Context, userUid string, payoutUid string, dateTo time.Time) error {
	cfg := config.Get()

	_, err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		UpdateMany(ctx, bson.M{
			"payoutUid": "",
			"userUid":   userUid,
			"chargedAt": bson.M{
				"$lte": dateTo,
			},
		}, bson.M{"$set": bson.M{"payoutUid": payoutUid}})
	if err != nil {
		return domain.NewError(transactionErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) TransactionCreate(ctx context.Context, trans *domain.Transaction) error {

	err := mr.UserBalanceChange(ctx, trans.UserUID, "usd", 2, trans.Amount)
	if err != nil {
		return domain.NewError(transactionErrorSource).SetCode(domain.ErrUserBalanceChange).Add(err)
	}

	transactionDb := mr.transactionConvertToDB(trans)
	cfg := config.Get()
	_, err = mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		InsertOne(ctx, transactionDb)
	if err != nil {
		errRollBack := mr.UserBalanceChange(ctx, trans.UserUID, "usd", 2, -trans.Amount)
		if errRollBack != nil {
			return domain.NewError(transactionErrorSource).SetCode(domain.ErrUserBalanceRollback).Add(err).Add(errRollBack)
		}
		return domain.NewError(transactionErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) TransactionDelete(ctx context.Context, trans *domain.Transaction) error {

	err := mr.UserBalanceChange(ctx, trans.UserUID, "usd", 2, -trans.Amount)
	if err != nil {
		return domain.NewError(transactionErrorSource).SetCode(domain.ErrUserBalanceChange).Add(err)
	}

	cfg := config.Get()
	_, err = mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		DeleteOne(ctx, bson.M{"uid": trans.UID})
	if err != nil {
		errRollBack := mr.UserBalanceChange(ctx, trans.UserUID, "usd", 2, trans.Amount)
		if errRollBack != nil {
			return domain.NewError(transactionErrorSource).SetCode(domain.ErrUserBalanceRollback).Add(err).Add(errRollBack)
		}
		return domain.NewError(transactionErrorSource).SetCode(domain.ErrDelete).Add(err)
	}
	return nil
}

func (mr *Repo) TransactionDeleteAllByBuyUIDAndType(ctx context.Context, buyUid string, transType domain.TransactionType) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		DeleteMany(ctx, bson.D{{"buyUid", buyUid}, {"type", transType}})
	if err != nil {
		return domain.NewError(transactionErrorSource).SetCode(domain.ErrDelete).Add(err)
	}
	return nil
}

func (mr *Repo) TransactionDeleteAllByUserUIDAndBuyUIDAndType(ctx context.Context, userUid string, buyUid string, transType domain.TransactionType) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).
		DeleteMany(ctx, bson.D{{"userUid", userUid}, {"buyUid", buyUid}, {"type", transType}})
	if err != nil {
		return domain.NewError(transactionErrorSource).SetCode(domain.ErrDelete).Add(err)
	}
	return nil
}

func (mr *Repo) transactionConvertFromDb(transactionDb *transactionDB) *domain.Transaction {
	return &domain.Transaction{
		UID:         transactionDb.UID,
		UserUID:     transactionDb.UserUID,
		FromUID:     transactionDb.FromUID,
		Percent:     transactionDb.Percent,
		Level:       transactionDb.Level,
		Type:        transactionDb.Type,
		RankCode:    transactionDb.RankCode,
		TaskCode:    transactionDb.TaskCode,
		ComboUID:    transactionDb.ComboUID,
		Amount:      transactionDb.Amount,
		PosAmount:   transactionDb.PosAmount,
		FullAmount:  transactionDb.FullAmount,
		Coefficient: transactionDb.Coefficient,
		BuyUID:      transactionDb.BuyUID,
		PayoutUID:   transactionDb.PayoutUID,
		DepositUID:  transactionDb.DepositUID,
		CreatedAt:   transactionDb.CreatedAt,
		ChargedAt:   transactionDb.ChargedAt,
		MsgCodes:    transactionDb.MsgCodes,
	}
}

func (mr *Repo) transactionConvertFromDbList(transactionsDb []*transactionDB) []*domain.Transaction {
	transactions := make([]*domain.Transaction, len(transactionsDb))
	for i, transactionDb := range transactionsDb {
		transactions[i] = mr.transactionConvertFromDb(transactionDb)
	}
	return transactions
}

func (mr *Repo) transactionConvertToDB(transaction *domain.Transaction) *transactionDB {
	if len(transaction.MsgCodes) == 0 {
		transaction.MsgCodes = nil
	}
	return &transactionDB{
		UID:         transaction.UID,
		UserUID:     transaction.UserUID,
		FromUID:     transaction.FromUID,
		Percent:     transaction.Percent,
		Level:       transaction.Level,
		Type:        transaction.Type,
		RankCode:    transaction.RankCode,
		TaskCode:    transaction.TaskCode,
		ComboUID:    transaction.ComboUID,
		Amount:      transaction.Amount,
		PosAmount:   transaction.PosAmount,
		FullAmount:  transaction.FullAmount,
		Coefficient: transaction.Coefficient,
		BuyUID:      transaction.BuyUID,
		PayoutUID:   transaction.PayoutUID,
		DepositUID:  transaction.DepositUID,
		CreatedAt:   transaction.CreatedAt,
		ChargedAt:   transaction.ChargedAt,
		MsgCodes:    transaction.MsgCodes,
	}
}

func (mr *Repo) transactionEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userIndexes := []mongo.IndexModel{
		{Keys: bson.M{"uid": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"userUid": 1}, Options: options.Index()},
		{Keys: bson.M{"buyUid": 1}, Options: options.Index()},
		{Keys: bson.M{"type": 1}, Options: options.Index()},
		{Keys: bson.M{"rankCode": 1}, Options: options.Index()},
		{Keys: bson.M{"taskCode": 1}, Options: options.Index()},
		{Keys: bson.M{"comboUid": 1}, Options: options.Index()},
		{Keys: bson.M{"createdAt": 1}, Options: options.Index()},
		{Keys: bson.M{"chargedAt": 1}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(transactionTable).Indexes().CreateMany(ctx, userIndexes)
	if err != nil {
		return domain.NewError(transactionErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
