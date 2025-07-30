package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/provider/logger"
)

const (
	// table name in DB
	userConfigTable = "user_config"

	// errors prefix
	userConfigErrorSource = "[repository.mongodb.user_config]"
)

type userConfigDB struct {
	UserUID      string              `bson:"userUid"`
	TeamType     domain.UserTeamType `bson:"teamType"`
	LastTeamType domain.UserTeamType `bson:"lastTeamType"`
	AllowSwitch  bool                `bson:"allowSwitch"`
}

func (mr *Repo) UserConfigGetByUserUID(ctx context.Context, uid string) (*domain.UserConfig, error) {
	userConfigDb := &userConfigDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(userConfigTable).
		FindOne(ctx, bson.M{"userUid": uid}).Decode(&userConfigDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(userConfigErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(userConfigErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.userConfigConvertFromDb(userConfigDb), nil
}

func (mr *Repo) UserConfigCreate(ctx context.Context, conf *domain.UserConfig) error {
	userConfigDb := mr.userConfigConvertToDb(conf)
	cfg := config.Get()
	l := logger.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userConfigTable).
		InsertOne(ctx, userConfigDb)
	if err != nil {
		l.Error(err)
		return domain.NewError(userConfigErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) UserConfigUpdate(ctx context.Context, conf *domain.UserConfig) error {
	userConfigDb := mr.userConfigConvertToDb(conf)
	update := bson.M{"$set": userConfigDb}
	cfg := config.Get()
	l := logger.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(userConfigTable).
		UpdateOne(ctx, bson.M{"userUid": userConfigDb.UserUID}, update)
	if err != nil {
		l.Error(err)
		return domain.NewError(userConfigErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}
	return nil
}

func (mr *Repo) userConfigConvertFromDb(db *userConfigDB) *domain.UserConfig {
	return &domain.UserConfig{
		UserUID:      db.UserUID,
		TeamType:     db.TeamType,
		LastTeamType: db.LastTeamType,
		AllowSwitch:  db.AllowSwitch,
	}
}

func (mr *Repo) userConfigConvertToDb(conf *domain.UserConfig) *userConfigDB {
	return &userConfigDB{
		UserUID:      conf.UserUID,
		TeamType:     conf.TeamType,
		LastTeamType: conf.LastTeamType,
		AllowSwitch:  conf.AllowSwitch,
	}
}

func (mr *Repo) userConfigEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userConfigIndexes := []mongo.IndexModel{
		{Keys: bson.M{"userUid": 1}, Options: options.Index().SetUnique(true)},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(userConfigTable).Indexes().CreateMany(ctx, userConfigIndexes)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
