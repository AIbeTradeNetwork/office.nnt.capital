package mongodb

import (
	"context"
	"fmt"
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
	buyTable = "buy"

	// errors prefix
	buyErrorSource = "[repository.mongodb.buy]"
)

type buyDB struct {
	UID          string         `bson:"uid"`
	UserUID      string         `bson:"userUid"`
	RefUID       string         `bson:"refUid"`
	MatchUID     string         `bson:"matchUid"`
	Row          *BigInt        `bson:"row"`
	Col          *BigInt        `bson:"col"`
	Type         domain.BuyType `bson:"type"`
	PlanCode     string         `bson:"planCode"`
	ProductCode  string         `bson:"productCode"`
	CurrencyCode string         `bson:"currencyCode"`
	Amount       int64          `bson:"amount"`
	Cv           int64          `bson:"cv"`
	CreatedAt    time.Time      `bson:"createdAt"`
	PaidAt       time.Time      `bson:"paidAt"`
	ApprovedAt   time.Time      `bson:"approvedAt"`
	ChargedAt    time.Time      `bson:"chargedAt"`
	CancelledAt  time.Time      `bson:"cancelledAt"`
	RefundedAt   time.Time      `bson:"refundedAt"`
	FlowUID      string         `bson:"flowUid"`
	PayUID       string         `bson:"payUid"`
}

func (mr *Repo) BuyGetByUID(ctx context.Context, uid string) (*domain.Buy, error) {
	buyDb := &buyDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(buyTable).
		FindOne(ctx, bson.M{"uid": uid}).Decode(&buyDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.buyConvertFromDb(buyDb), nil
}

func (mr *Repo) BuyGetByPayUID(ctx context.Context, uid string) (*domain.Buy, error) {
	buyDb := &buyDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(buyTable).
		FindOne(ctx, bson.M{"payUid": uid}).Decode(&buyDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.buyConvertFromDb(buyDb), nil
}

func (mr *Repo) BuySumAllByField(ctx context.Context, filter interface{}, field string) (int64, error) {
	cfg := config.Get()

	// TODO: add multi currency support
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": "", "sum": bson.M{"$sum": fmt.Sprintf("$%s", field)}}},
	}
	type Result struct {
		ID  string `bson:"_id"`
		Sum int64  `bson:"sum"`
	}
	opts := options.Aggregate().SetCollation(&options.Collation{Locale: "en_US", NumericOrdering: true})
	cur, err := mr.db.Database(cfg.MongoDB).Collection(buyTable).
		Aggregate(ctx, pipeline, opts)
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)
	var res Result
	for cur.Next(ctx) {
		err = cur.Decode(&res)
		if err != nil {
			return 0, err
		}
	}
	return res.Sum, nil
}

func (mr *Repo) BuyGetAllWith(ctx context.Context, filter interface{}, options *options.FindOptions) ([]*domain.Buy, error) {
	cfg := config.Get()

	cur, err := mr.db.Database(cfg.MongoDB).Collection(buyTable).
		Find(ctx, filter, options)
	if err != nil {
		//if errors.Is(err, mongo.ErrNoDocuments) {
		//	return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		//}
		return nil, domain.NewError(userPlaceErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cur.Close(ctx)
	all := make([]*domain.Buy, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		dbObj := &buyDB{}
		err = cur.Decode(dbObj)

		if err != nil {
			return nil, err
		}

		all = append(all, mr.buyConvertFromDb(dbObj))
	}
	err = cur.Err()
	if err != nil {
		return nil, domain.NewError(buyErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return all, nil
}

func (mr *Repo) BuyCreate(ctx context.Context, buy *domain.Buy) error {
	buyDb := mr.buyConvertToDB(buy)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(buyTable).
		InsertOne(ctx, buyDb)
	if err != nil {
		fmt.Println(err)
		return domain.NewError(buyErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) BuyUpdate(ctx context.Context, buy *domain.Buy) error {
	buyDb := mr.buyConvertToDB(buy)
	update := bson.M{"$set": buyDb}
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(buyTable).
		UpdateOne(ctx, bson.D{{Key: "uid", Value: buyDb.UID}}, update)
	if err != nil {
		return err
	}
	//if result.ModifiedCount == 0 {
	//	return domain.NewError(buyErrorSource).SetCode(domain.ErrUpdate)
	//}
	return nil
}

func (mr *Repo) buyConvertFromDb(buy *buyDB) *domain.Buy {
	return &domain.Buy{
		UID:          buy.UID,
		UserUID:      buy.UserUID,
		RefUID:       buy.RefUID,
		MatchUID:     buy.MatchUID,
		Row:          buy.Row.Int(),
		Col:          buy.Col.Int(),
		Type:         buy.Type,
		PlanCode:     buy.PlanCode,
		ProductCode:  buy.ProductCode,
		CurrencyCode: buy.CurrencyCode,
		Amount:       buy.Amount,
		Cv:           buy.Cv,
		CreatedAt:    buy.CreatedAt,
		PaidAt:       buy.PaidAt,
		ApprovedAt:   buy.ApprovedAt,
		ChargedAt:    buy.ChargedAt,
		CancelledAt:  buy.CancelledAt,
		RefundedAt:   buy.RefundedAt,
		FlowUID:      buy.FlowUID,
		PayUID:       buy.PayUID,
	}
}

func (mr *Repo) buyConvertFromDbList(buys []*buyDB) []*domain.Buy {
	buyList := make([]*domain.Buy, 0, len(buys))
	for _, buy := range buys {
		buyList = append(buyList, mr.buyConvertFromDb(buy))
	}
	return buyList
}

func (mr *Repo) buyConvertToDB(buy *domain.Buy) *buyDB {
	return &buyDB{
		UID:          buy.UID,
		UserUID:      buy.UserUID,
		RefUID:       buy.RefUID,
		MatchUID:     buy.MatchUID,
		Row:          NewBigInt(buy.Row),
		Col:          NewBigInt(buy.Col),
		Type:         buy.Type,
		PlanCode:     buy.PlanCode,
		ProductCode:  buy.ProductCode,
		CurrencyCode: buy.CurrencyCode,
		Amount:       buy.Amount,
		Cv:           buy.Cv,
		CreatedAt:    buy.CreatedAt,
		PaidAt:       buy.PaidAt,
		ApprovedAt:   buy.ApprovedAt,
		ChargedAt:    buy.ChargedAt,
		CancelledAt:  buy.CancelledAt,
		RefundedAt:   buy.RefundedAt,
		FlowUID:      buy.FlowUID,
		PayUID:       buy.PayUID,
	}
}

func (mr *Repo) buyEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userIndexes := []mongo.IndexModel{
		{Keys: bson.M{"uid": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"flowUid": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"payUid": 1}, Options: options.Index()},
		{Keys: bson.D{{Key: "userUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "refUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "matchUid", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "row", Value: 1}, {Key: "col", Value: 1}}, Options: options.Index().SetCollation(&options.Collation{
			Locale:          "en_US",
			NumericOrdering: true,
		})},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(buyTable).Indexes().CreateMany(ctx, userIndexes)
	if err != nil {
		return domain.NewError(buyErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
