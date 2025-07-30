package mongodb

import (
	"context"
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
	payoutTable = "payout"

	// errors prefix
	payoutErrorSource = "[repository.mongodb.payout]"
)

type payoutDB struct {
	UID           string    `bson:"uid"`
	UserUID       string    `bson:"userUid"`
	Amount        int64     `bson:"amount"`
	Fee           int64     `bson:"fee"`
	CurrencyCode  string    `bson:"currencyCode"`
	MethodCode    string    `bson:"methodCode"`
	AccountNumber string    `bson:"accountNumber"`
	AccountName   string    `bson:"accountName"`
	CreatedAt     time.Time `bson:"createdAt"`
	ApprovedAt    time.Time `bson:"approvedAt"`
	ChargedAt     time.Time `bson:"chargedAt"`
	CancelledAt   time.Time `bson:"cancelledAt"`
	Reason        string    `bson:"reason"`
}

func (mr *Repo) PayoutGetByUID(ctx context.Context, uid string) (*domain.Payout, error) {
	payoutDb := &payoutDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(payoutTable).
		FindOne(ctx, bson.M{"uid": uid}).Decode(&payoutDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(payoutErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(payoutErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.payoutConvertFromDb(payoutDb), nil
}

func (mr *Repo) PayoutGetAllByUserUID(ctx context.Context, userUID string, limit int64, skip int64) ([]*domain.Payout, error) {
	payouts := make([]*domain.Payout, 0)
	cfg := config.Get()

	opts := options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit).SetSkip(skip)
	cursor, err := mr.db.Database(cfg.MongoDB).Collection(payoutTable).
		Find(ctx, bson.M{"userUid": userUID}, opts)
	if err != nil {
		return nil, domain.NewError(payoutErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		payoutDb := &payoutDB{}
		err = cursor.Decode(&payoutDb)
		if err != nil {
			return nil, domain.NewError(payoutErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		payouts = append(payouts, mr.payoutConvertFromDb(payoutDb))
	}
	return payouts, nil
}

func (mr *Repo) PayoutCreate(ctx context.Context, payout *domain.Payout) error {
	payoutDb := mr.payoutConvertToDb(payout)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(payoutTable).
		InsertOne(ctx, payoutDb)
	if err != nil {
		return domain.NewError(payoutErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) payoutConvertToDb(p *domain.Payout) *payoutDB {
	return &payoutDB{
		UID:           p.UID,
		UserUID:       p.UserUID,
		Amount:        p.Amount,
		Fee:           p.Fee,
		CurrencyCode:  p.CurrencyCode,
		MethodCode:    p.MethodCode,
		AccountNumber: p.AccountNumber,
		AccountName:   p.AccountName,
		CreatedAt:     p.CreatedAt,
		ApprovedAt:    p.ApprovedAt,
		ChargedAt:     p.ChargedAt,
		CancelledAt:   p.CancelledAt,
		Reason:        p.Reason,
	}
}

func (mr *Repo) payoutConvertFromDb(p *payoutDB) *domain.Payout {
	return &domain.Payout{
		UID:           p.UID,
		UserUID:       p.UserUID,
		Amount:        p.Amount,
		Fee:           p.Fee,
		CurrencyCode:  p.CurrencyCode,
		MethodCode:    p.MethodCode,
		AccountNumber: p.AccountNumber,
		AccountName:   p.AccountName,
		CreatedAt:     p.CreatedAt,
		ApprovedAt:    p.ApprovedAt,
		ChargedAt:     p.ChargedAt,
		CancelledAt:   p.CancelledAt,
		Reason:        p.Reason,
	}
}

func (mr *Repo) payoutEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	indexes := []mongo.IndexModel{
		{Keys: bson.M{"uid": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"userUid": 1}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(payoutTable).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return domain.NewError(payoutErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
