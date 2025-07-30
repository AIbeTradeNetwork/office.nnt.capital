package mongodb

import (
	"context"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"server/internal/config"
	"server/internal/domain"
)

const (
	// table name in DB
	depositTable = "deposit"

	// errors prefix
	depositErrorSource = "[repository.mongodb.deposit]"
)

type depositDB struct {
	UID           string    `bson:"uid"`
	UserUID       string    `bson:"userUid"`
	Amount        int64     `bson:"amount"`
	Fee           int64     `bson:"fee"`
	CurrencyCode  string    `bson:"currencyCode"`
	Precision     uint8     `bson:"precision"`
	Rate          float64   `bson:"rate"`
	MethodCode    string    `bson:"methodCode"`
	AccountNumber string    `bson:"accountNumber"`
	AccountName   string    `bson:"accountName"`
	CreatedAt     time.Time `bson:"createdAt"`
	ApprovedAt    time.Time `bson:"approvedAt"`
	ChargedAt     time.Time `bson:"chargedAt"`
	CancelledAt   time.Time `bson:"cancelledAt"`
	Comment       string    `bson:"comment"`
	TxUID         string    `bson:"txUid"`
	TxLT          uint64    `bson:"txLt"`
}

func (mr *Repo) DepositGetByUID(ctx context.Context, uid string) (*domain.Deposit, error) {
	depositDb := &depositDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(depositTable).
		FindOne(ctx, bson.M{"uid": uid}).Decode(&depositDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(depositErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(depositErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.depositConvertFromDb(depositDb), nil
}

func (mr *Repo) DepositGetAllByUserUID(ctx context.Context, userUID string, limit int64, skip int64) ([]*domain.Deposit, error) {
	deposits := make([]*domain.Deposit, 0)
	cfg := config.Get()

	opts := options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit).SetSkip(skip)
	cursor, err := mr.db.Database(cfg.MongoDB).Collection(depositTable).
		Find(ctx, bson.M{"userUid": userUID}, opts)
	if err != nil {
		return nil, domain.NewError(depositErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		depositDb := &depositDB{}
		err = cursor.Decode(&depositDb)
		if err != nil {
			return nil, domain.NewError(depositErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		deposits = append(deposits, mr.depositConvertFromDb(depositDb))
	}
	return deposits, nil
}

func (mr *Repo) DepositCreate(ctx context.Context, deposit *domain.Deposit) error {
	depositDb := mr.depositConvertToDb(deposit)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(depositTable).
		InsertOne(ctx, depositDb)
	if err != nil {
		return domain.NewError(depositErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) DepositUpdate(ctx context.Context, deposit *domain.Deposit) error {
	depositDb := mr.depositConvertToDb(deposit)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(depositTable).
		UpdateOne(ctx, bson.M{"uid": deposit.UID}, bson.M{"$set": depositDb})
	if err != nil {
		return domain.NewError(depositErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) DepositSumByUserUID(ctx context.Context, userUID string) (int64, error) {
    cfg := config.Get()
    pipeline := []bson.M{
        {"$match": bson.M{"userUid": userUID}},
        {"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}},
    }
    cur, err := mr.db.Database(cfg.MongoDB).Collection(depositTable).Aggregate(ctx, pipeline)
    if err != nil {
        return 0, err
    }
    defer cur.Close(ctx)
    var result struct{ Total int64 `bson:"total"` }
    if cur.Next(ctx) {
        if err := cur.Decode(&result); err != nil {
            return 0, err
        }
        return result.Total, nil
    }
    return 0, nil
}

func (mr *Repo) depositConvertToDb(p *domain.Deposit) *depositDB {
	return &depositDB{
		UID:           p.UID,
		UserUID:       p.UserUID,
		Amount:        p.Amount,
		Fee:           p.Fee,
		CurrencyCode:  p.CurrencyCode,
		Precision:     p.Precision,
		Rate:          p.Rate.InexactFloat64(),
		MethodCode:    p.MethodCode,
		AccountNumber: p.AccountNumber,
		AccountName:   p.AccountName,
		CreatedAt:     p.CreatedAt,
		ApprovedAt:    p.ApprovedAt,
		ChargedAt:     p.ChargedAt,
		CancelledAt:   p.CancelledAt,
		Comment:       p.Comment,
		TxUID:         p.TxUID,
		TxLT:          p.TxLT,
	}
}

func (mr *Repo) depositConvertFromDb(p *depositDB) *domain.Deposit {
	return &domain.Deposit{
		UID:           p.UID,
		UserUID:       p.UserUID,
		Amount:        p.Amount,
		Fee:           p.Fee,
		CurrencyCode:  p.CurrencyCode,
		Precision:     p.Precision,
		Rate:          decimal.NewFromFloat(p.Rate),
		MethodCode:    p.MethodCode,
		AccountNumber: p.AccountNumber,
		AccountName:   p.AccountName,
		CreatedAt:     p.CreatedAt,
		ApprovedAt:    p.ApprovedAt,
		ChargedAt:     p.ChargedAt,
		CancelledAt:   p.CancelledAt,
		Comment:       p.Comment,
		TxUID:         p.TxUID,
		TxLT:          p.TxLT,
	}
}

func (mr *Repo) depositEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	indexes := []mongo.IndexModel{
		{Keys: bson.M{"uid": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"userUid": 1}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(depositTable).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return domain.NewError(depositErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
