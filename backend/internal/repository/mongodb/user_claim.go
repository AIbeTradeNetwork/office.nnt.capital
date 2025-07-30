package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/internal/config"
	"server/internal/domain"
)

const (
	// table name in DB
	userClaimTable = "user_claim"

	// errors prefix
	userClaimErrorSource = "[repository.mongodb.user_claim]"
)

type userClaimDB struct {
	UID          string               `bson:"uid"`
	ClaimCode    string               `bson:"claimCode"`
	UserUID      string               `bson:"userUid"`
	RefUID       string               `bson:"refUid"`
	Level        int                  `bson:"level"`
	CreatedAt    time.Time            `bson:"createdAt"`
	ClaimedAt    time.Time            `bson:"claimedAt"`
	Amount       int64                `bson:"amount"`
	CurrencyCode string               `bson:"currencyCode"`
	Precision    uint8                `bson:"precision"`
	Type         domain.UserClaimType `bson:"type"`
	TaskCode     string               `bson:"taskCode"`
	ComboUID     string               `bson:"comboUid"`
	PartnerCode  string               `bson:"partnerCode"`
}

func (mr *Repo) UserClaimCreate(ctx context.Context, userClaim *domain.UserClaim) error {
	cfg := config.Get()

	userClaimDb := mr.userClaimConvertToDB(userClaim)

	_, err := mr.db.Database(cfg.MongoDB).Collection(userClaimTable).InsertOne(ctx, userClaimDb)
	if err != nil {
		return domain.NewError(userClaimErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return nil
}

func (mr *Repo) UserClaimUpsert(ctx context.Context, userClaim *domain.UserClaim) error {
	cfg := config.Get()

	userClaimDb := mr.userClaimConvertToDB(userClaim)

	_, err := mr.db.Database(cfg.MongoDB).Collection(userClaimTable).
		UpdateOne(ctx, bson.M{"uid": userClaim.UID}, bson.M{"$set": userClaimDb}, options.Update().SetUpsert(true))
	if err != nil {
		return domain.NewError(userClaimErrorSource).SetCode(domain.ErrUpsert).Add(err)
	}
	return nil
}

func (mr *Repo) UserClaimUpdate(ctx context.Context, userClaim *domain.UserClaim) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userClaimTable).
		UpdateOne(ctx, bson.M{"uid": userClaim.UID}, bson.M{"$set": mr.userClaimConvertToDB(userClaim)})
	if err != nil {
		return domain.NewError(userClaimErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) UserClaimDelete(ctx context.Context, userClaim *domain.UserClaim) error {
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userClaimTable).
		DeleteOne(ctx, bson.M{"uid": userClaim.UID})
	if err != nil {
		return domain.NewError(userClaimErrorSource).SetCode(domain.ErrDelete).Add(err)
	}
	return nil
}

func (mr *Repo) UserClaimGetLastByClaimCodeAndTypeAndUserUID(ctx context.Context, claimCode string, claimType domain.UserClaimType, userUid string) (*domain.UserClaim, error) {
	cfg := config.Get()

	var userClaimDb *userClaimDB

	opts := options.FindOne().SetSort(bson.M{"claimedAt": -1})
	err := mr.db.Database(cfg.MongoDB).Collection(userClaimTable).
		FindOne(ctx, bson.M{"claimCode": claimCode, "userUid": userUid, "type": claimType}, opts).Decode(&userClaimDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userClaimErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userClaimErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return mr.userClaimConvertFromDB(userClaimDb), nil
}

func (mr *Repo) UserClaimGetByUserUIDAndClaimCodeAndTypeAndTaskCode(ctx context.Context, userUid string, claimCode string, claimType domain.UserClaimType, taskCode string) (*domain.UserClaim, error) {
	cfg := config.Get()

	var userClaimDb *userClaimDB

	opts := options.FindOne().SetSort(bson.M{"claimedAt": -1})
	err := mr.db.Database(cfg.MongoDB).Collection(userClaimTable).
		FindOne(ctx, bson.M{
			"userUid":   userUid,
			"claimCode": claimCode,
			"type":      claimType,
			"taskCode":  taskCode,
		}, opts).Decode(&userClaimDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userClaimErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userClaimErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return mr.userClaimConvertFromDB(userClaimDb), nil
}

func (mr *Repo) UserClaimGetByUserUIDAndClaimCodeAndTypeAndComboCode(ctx context.Context, userUid string, claimCode string, claimType domain.UserClaimType, comboCode string) (*domain.UserClaim, error) {
	cfg := config.Get()

	var userClaimDb *userClaimDB

	opts := options.FindOne().SetSort(bson.M{"claimedAt": -1})
	err := mr.db.Database(cfg.MongoDB).Collection(userClaimTable).
		FindOne(ctx, bson.M{
			"userUid":   userUid,
			"claimCode": claimCode,
			"type":      claimType,
			"comboUid":  comboCode,
		}, opts).Decode(&userClaimDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userClaimErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userClaimErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return mr.userClaimConvertFromDB(userClaimDb), nil
}

func (mr *Repo) UserClaimGetByUserUIDAndClaimCodeAndTypeAndPartnerCode(ctx context.Context, userUid string, claimCode string, claimType domain.UserClaimType, partnerCode string) (*domain.UserClaim, error) {
	cfg := config.Get()

	var userClaimDb *userClaimDB

	opts := options.FindOne().SetSort(bson.M{"claimedAt": -1})
	err := mr.db.Database(cfg.MongoDB).Collection(userClaimTable).
		FindOne(ctx, bson.M{
			"userUid":     userUid,
			"claimCode":   claimCode,
			"type":        claimType,
			"partnerCode": partnerCode,
		}, opts).Decode(&userClaimDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userClaimErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userClaimErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return mr.userClaimConvertFromDB(userClaimDb), nil
}

func (mr *Repo) userClaimConvertToDB(userClaim *domain.UserClaim) *userClaimDB {
	return &userClaimDB{
		UID:          userClaim.UID,
		ClaimCode:    userClaim.ClaimCode,
		UserUID:      userClaim.UserUID,
		RefUID:       userClaim.RefUID,
		Level:        userClaim.Level,
		CreatedAt:    userClaim.CreatedAt,
		ClaimedAt:    userClaim.ClaimedAt,
		Amount:       userClaim.Amount,
		CurrencyCode: userClaim.CurrencyCode,
		Precision:    userClaim.Precision,
		Type:         userClaim.Type,
		TaskCode:     userClaim.TaskCode,
		ComboUID:     userClaim.ComboUID,
		PartnerCode:  userClaim.PartnerCode,
	}
}

func (mr *Repo) userClaimConvertFromDB(userClaimDB *userClaimDB) *domain.UserClaim {
	return &domain.UserClaim{
		UID:          userClaimDB.UID,
		ClaimCode:    userClaimDB.ClaimCode,
		UserUID:      userClaimDB.UserUID,
		RefUID:       userClaimDB.RefUID,
		Level:        userClaimDB.Level,
		CreatedAt:    userClaimDB.CreatedAt,
		ClaimedAt:    userClaimDB.ClaimedAt,
		Amount:       userClaimDB.Amount,
		CurrencyCode: userClaimDB.CurrencyCode,
		Precision:    userClaimDB.Precision,
		Type:         userClaimDB.Type,
		TaskCode:     userClaimDB.TaskCode,
		ComboUID:     userClaimDB.ComboUID,
		PartnerCode:  userClaimDB.PartnerCode,
	}
}

func (mr *Repo) userClaimEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()

	_, err := mr.db.Database(cfg.MongoDB).Collection(userClaimTable).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"uid", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"claimCode", 1}, {"userUid", 1}}, Options: options.Index()},
		{Keys: bson.D{{"type", 1}}, Options: options.Index()},
		{Keys: bson.D{{"claimedAt", -1}}, Options: options.Index()},
		{Keys: bson.D{{"taskCode", -1}}, Options: options.Index()},
		{Keys: bson.D{{"comboUid", -1}}, Options: options.Index()},
		{Keys: bson.D{{"partnerCode", -1}}, Options: options.Index()},
	})

	if err != nil {
		return domain.NewError(userClaimErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	return nil
}
