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
	planTable = "plan"

	// errors prefix
	planErrorSource = "[repository.mongodb.plan]"
)

type planDB struct {
	Code         string        `bson:"code"`
	Period       time.Duration `bson:"period"`
	Price        int64         `bson:"price"`
	RetailPrice  int64         `bson:"retailPrice"`
	Cv           int64         `bson:"cv"`
	CurrencyCode string        `bson:"currencyCode"`
	PaymentWait  time.Duration `bson:"paymentWait"`
	ChargeWait   time.Duration `bson:"chargeWait"`
	RankCode     string        `bson:"rankCode"`
	RankPeriod   time.Duration `bson:"rankPeriod"`
	Priority     int64         `bson:"priority"`
	BotCount     int64         `bson:"botCount"`
	MaxDeposit   int64         `bson:"maxDeposit"`
	CoinBonus    int64         `bson:"coinBonus"`
}

func (mr *Repo) PlanGetByCode(ctx context.Context, code string) (*domain.Plan, error) {
	planDb := &planDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(planTable).
		FindOne(ctx, bson.M{"code": code}).Decode(&planDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(planErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(planErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.planConvertFromDb(planDb), nil
}

func (mr *Repo) PlanGetAll(ctx context.Context) ([]*domain.Plan, error) {
	plans := make([]*domain.Plan, 0)
	cfg := config.Get()
	cursor, err := mr.db.Database(cfg.MongoDB).Collection(planTable).Find(ctx, bson.M{})
	if err != nil {
		return nil, domain.NewError(planErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		planDb := &planDB{}
		err = cursor.Decode(&planDb)
		if err != nil {
			return nil, domain.NewError(planErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		plans = append(plans, mr.planConvertFromDb(planDb))
	}
	return plans, nil
}

func (mr *Repo) PlanCreate(ctx context.Context, plan *domain.Plan) error {
	planDb := mr.planConvertToDb(plan)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(planTable).
		InsertOne(ctx, planDb)
	if err != nil {
		return domain.NewError(planErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) planConvertToDb(plan *domain.Plan) *planDB {
	return &planDB{
		Code:         plan.Code,
		Period:       plan.Period,
		Price:        plan.Price,
		RetailPrice:  plan.RetailPrice,
		Cv:           plan.Cv,
		CurrencyCode: plan.CurrencyCode,
		PaymentWait:  plan.PaymentWait,
		ChargeWait:   plan.ChargeWait,
		RankCode:     plan.RankCode,
		RankPeriod:   plan.RankPeriod,
		Priority:     plan.Priority,
		BotCount:     plan.BotCount,
		MaxDeposit:   plan.MaxDeposit,
		CoinBonus:    plan.CoinBonus,
	}
}

func (mr *Repo) planConvertFromDb(plan *planDB) *domain.Plan {
	return &domain.Plan{
		Code:         plan.Code,
		Period:       plan.Period,
		Price:        plan.Price,
		RetailPrice:  plan.RetailPrice,
		Cv:           plan.Cv,
		CurrencyCode: plan.CurrencyCode,
		PaymentWait:  plan.PaymentWait,
		ChargeWait:   plan.ChargeWait,
		RankCode:     plan.RankCode,
		RankPeriod:   plan.RankPeriod,
		Priority:     plan.Priority,
		BotCount:     plan.BotCount,
		MaxDeposit:   plan.MaxDeposit,
		CoinBonus:    plan.CoinBonus,
	}
}

func (mr *Repo) planEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	planIndexes := []mongo.IndexModel{
		{Keys: bson.M{"code": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"priority": -1}, Options: options.Index().SetUnique(true)},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(planTable).Indexes().CreateMany(ctx, planIndexes)
	if err != nil {
		return domain.NewError(planErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
