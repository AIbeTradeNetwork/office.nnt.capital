package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"server/internal/config"
	"server/internal/domain"
)

const (
	// table name in DB
	configTable = "config"

	// errors prefix
	configErrorSource = "[repository.mongodb.config]"
)

type fastStartBonusDB struct {
	Amount int64 `bson:"amount"`
	MinCv  int64 `bson:"minCv"`
}

type configDB struct {
	ID                        string             `bson:"_id"`
	PayoutAmountMin           int64              `bson:"payoutAmountMin"`
	PayoutFeeMin              int64              `bson:"payoutFeeMin"`
	PayoutFeePercent          int64              `bson:"payoutFeePercent"`
	FastStartBonuses          []fastStartBonusDB `bson:"fastStartBonuses"`
	FastStartDuration         time.Duration      `bson:"fastStartDuration"`
	ActivityCvAmount          int64              `bson:"activityCvAmount"`
	ActivityMaxDuration       time.Duration      `bson:"activityMaxDuration"`
	ClientRefBonus            map[string]int64   `bson:"clientRefBonus"`
	DistributorRefBonus       map[string]int64   `bson:"distributorRefBonus"`
	DistributorPrice          int64              `bson:"distributorPrice"`
	DefaultCurrencyCode       string             `bson:"defaultCurrencyCode"`
	AllowTeamTypeSwitchAmount int64              `bson:"allowTeamTypeSwitchAmount"`
	DefaultRefUid             string             `bson:"defaultRefUid"`
	CoinCode                  string             `bson:"coinCode"`
	CoinRefBonus              int64              `bson:"coinRefBonus"`
	CoinToRefBonus            int64              `bson:"coinToRefBonus"`
	CoinRefLinePercent        []int64            `bson:"coinRefLinePercent"`
	InitializedAt             time.Time          `bson:"initializedAt"`
	UnlimInvite               bool               `bson:"unlimInvite"`
	TonLastLT                 uint64             `bson:"tonLastLT"`
	TonWallet                 string             `bson:"tonWallet"`
	RefSafeUID                string             `bson:"refSafeUid"`
	Tier1SafeUID              string             `bson:"tier1SafeUid"`
	Tier2SafeUID              string             `bson:"tier2SafeUid"`
	Tier1CoinSafeUID          string             `bson:"tier1CoinSafeUid"`
	Tier2CoinSafeUID          string             `bson:"tier2CoinSafeUid"`
}

func (mr *Repo) ConfigGet(ctx context.Context) (*domain.Config, error) {
	configDb := &configDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(configTable).
		FindOne(ctx, bson.M{"_id": "CONFIG"}).Decode(&configDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(configErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(configErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	bonuses := make([]domain.FastStartBonus, 0, len(configDb.FastStartBonuses))
	for _, fsb := range configDb.FastStartBonuses {
		bonuses = append(bonuses, domain.FastStartBonus{
			Amount: fsb.Amount,
			MinCv:  fsb.MinCv,
		})
	}
	return &domain.Config{
		PayoutAmountMin:           configDb.PayoutAmountMin,
		PayoutFeeMin:              configDb.PayoutFeeMin,
		PayoutFeePercent:          configDb.PayoutFeePercent,
		FastStartBonuses:          bonuses,
		FastStartDuration:         configDb.FastStartDuration,
		ActivityCvAmount:          configDb.ActivityCvAmount,
		ActivityMaxDuration:       configDb.ActivityMaxDuration,
		ClientRefBonus:            configDb.ClientRefBonus,
		DistributorRefBonus:       configDb.DistributorRefBonus,
		DistributorPrice:          configDb.DistributorPrice,
		DefaultCurrencyCode:       configDb.DefaultCurrencyCode,
		AllowTeamTypeSwitchAmount: configDb.AllowTeamTypeSwitchAmount,
		DefaultRefUid:             configDb.DefaultRefUid,
		CoinCode:                  configDb.CoinCode,
		CoinRefBonus:              configDb.CoinRefBonus,
		CoinToRefBonus:            configDb.CoinToRefBonus,
		CoinRefLinePercent:        configDb.CoinRefLinePercent,
		InitializedAt:             configDb.InitializedAt,
		UnlimInvite:               configDb.UnlimInvite,
		TonLastLT:                 configDb.TonLastLT,
		TonWallet:                 configDb.TonWallet,
		RefSafeUID:                configDb.RefSafeUID,   // Default safe for referrals
		Tier1SafeUID:              configDb.Tier1SafeUID, // Safe for buyers of 6th digit
		Tier2SafeUID:              configDb.Tier2SafeUID, // Safe for buyers of 7th digit
		Tier1CoinSafeUID:          configDb.Tier1CoinSafeUID,
		Tier2CoinSafeUID:          configDb.Tier2CoinSafeUID,
	}, nil
}

func (mr *Repo) ConfigUpdateTonLastLT(ctx context.Context, tonLastLT uint64) error {
	cfg := config.Get()
	filter := bson.M{"_id": "CONFIG"}
	update := bson.M{"$set": bson.M{"tonLastLT": tonLastLT}}
	_, err := mr.db.Database(cfg.MongoDB).Collection(configTable).UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.NewError(configErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) ConfigCreate(ctx context.Context, dConfig *domain.Config) error {
	cfg := config.Get()
	bonuses := make([]fastStartBonusDB, 0, len(dConfig.FastStartBonuses))
	for _, fsb := range dConfig.FastStartBonuses {
		bonuses = append(bonuses, fastStartBonusDB{
			Amount: fsb.Amount,
			MinCv:  fsb.MinCv,
		})
	}
	configDb := &configDB{
		ID:                        "CONFIG",
		PayoutAmountMin:           dConfig.PayoutAmountMin,
		PayoutFeeMin:              dConfig.PayoutFeeMin,
		PayoutFeePercent:          dConfig.PayoutFeePercent,
		FastStartBonuses:          bonuses,
		FastStartDuration:         dConfig.FastStartDuration,
		ActivityCvAmount:          dConfig.ActivityCvAmount,
		ActivityMaxDuration:       dConfig.ActivityMaxDuration,
		ClientRefBonus:            dConfig.ClientRefBonus,
		DistributorRefBonus:       dConfig.DistributorRefBonus,
		DistributorPrice:          dConfig.DistributorPrice,
		DefaultCurrencyCode:       dConfig.DefaultCurrencyCode,
		AllowTeamTypeSwitchAmount: dConfig.AllowTeamTypeSwitchAmount,
		DefaultRefUid:             dConfig.DefaultRefUid,
		CoinCode:                  dConfig.CoinCode,
		CoinRefBonus:              dConfig.CoinRefBonus,
		CoinToRefBonus:            dConfig.CoinToRefBonus,
		CoinRefLinePercent:        dConfig.CoinRefLinePercent,
		InitializedAt:             dConfig.InitializedAt,
		UnlimInvite:               dConfig.UnlimInvite,
		TonLastLT:                 dConfig.TonLastLT,
		TonWallet:                 dConfig.TonWallet,
		RefSafeUID:                dConfig.RefSafeUID,
		Tier1SafeUID:              dConfig.Tier1SafeUID,
		Tier2SafeUID:              dConfig.Tier2SafeUID,
		Tier1CoinSafeUID:          dConfig.Tier1CoinSafeUID,
		Tier2CoinSafeUID:          dConfig.Tier2CoinSafeUID,
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(configTable).
		InsertOne(ctx, configDb)
	if err != nil {
		return domain.NewError(configErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}
