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
	rankTable = "rank"

	// errors prefix
	rankErrorSource = "[repository.mongodb.rank]"
)

type rankTeamConditionDB struct {
	RankCode string              `bson:"rankCode"`
	TeamType domain.UserTeamType `bson:"teamType"`
	IsRef    bool                `bson:"isRef"`
	Count    uint64              `bson:"count"`
}

type rankGetBonusDB struct {
	Amount int64 `bson:"amount"`
	Months int   `bson:"months"`
}

type rankDB struct {
	Code          string                `bson:"code"`
	MinCv         int64                 `bson:"minCv"`
	TeamCondition []rankTeamConditionDB `bson:"teamCondition"`
	StartAt       time.Time             `bson:"startAt"`
	EndAt         time.Time             `bson:"endAt"`
	BinBonus      int64                 `bson:"binBonus"`
	RefBonus      map[string]int64      `bson:"refBonus"`
	BinBonusLimit int64                 `bson:"binBonusLimit"`
	MatchBonus    []int64               `bson:"matchBonus"`
	Priority      int64                 `bson:"priority"`
	FirstBonus    rankGetBonusDB        `bson:"firstBonus"`
	ApproveBonus  rankGetBonusDB        `bson:"approveBonus"`
}

func (mr *Repo) RankGetAll(ctx context.Context) ([]*domain.Rank, error) {
	ranks := make([]*domain.Rank, 0)
	cfg := config.Get()
	cursor, err := mr.db.Database(cfg.MongoDB).Collection(rankTable).Find(ctx, bson.M{})
	if err != nil {
		return nil, domain.NewError(rankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		rankDb := &rankDB{}
		err = cursor.Decode(&rankDb)
		if err != nil {
			return nil, domain.NewError(rankErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		ranks = append(ranks, mr.rankConvertFromDb(rankDb))
	}
	return ranks, nil
}

func (mr *Repo) RankGetByCode(ctx context.Context, code string) (*domain.Rank, error) {
	rankDb := &rankDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(rankTable).
		FindOne(ctx, bson.M{"code": code}).Decode(&rankDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(rankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(rankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.rankConvertFromDb(rankDb), nil
}

func (mr *Repo) RankGetNext(ctx context.Context, rank *domain.Rank) (*domain.Rank, error) {
	rankDb := &rankDB{}
	cfg := config.Get()
	opts := options.FindOne().SetSort(bson.D{{"priority", 1}})
	err := mr.db.Database(cfg.MongoDB).Collection(rankTable).
		FindOne(ctx, bson.M{"priority": bson.M{"$gt": rank.Priority}}, opts).Decode(&rankDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(rankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(rankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.rankConvertFromDb(rankDb), nil
}

func (mr *Repo) RankGetByMinCv(ctx context.Context, minCv int64) (*domain.Rank, error) {
	rankDb := &rankDB{}
	cfg := config.Get()

	filter := bson.D{
		{"minCv", bson.M{"$lte": minCv, "$gt": 0}},
	}
	opts := options.FindOne().SetSort(bson.M{"priority": -1})

	err := mr.db.Database(cfg.MongoDB).Collection(rankTable).
		FindOne(ctx, filter, opts).Decode(&rankDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(rankErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(rankErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.rankConvertFromDb(rankDb), nil
}

func (mr *Repo) RankCreate(ctx context.Context, rank *domain.Rank) error {
	rankDb := mr.rankConvertToDb(rank)
	cfg := config.Get()
	_, err := mr.db.Database(cfg.MongoDB).Collection(rankTable).
		InsertOne(ctx, rankDb)
	if err != nil {
		return domain.NewError(rankErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) rankConvertFromDb(db *rankDB) *domain.Rank {
	teamCondition := make([]domain.RankTeamCondition, 0, len(db.TeamCondition))
	for _, tc := range db.TeamCondition {
		teamCondition = append(teamCondition, domain.RankTeamCondition{
			RankCode: tc.RankCode,
			TeamType: tc.TeamType,
			IsRef:    tc.IsRef,
			Count:    tc.Count,
		})
	}
	return &domain.Rank{
		Code:          db.Code,
		MinCv:         db.MinCv,
		TeamCondition: teamCondition,
		BinBonus:      db.BinBonus,
		RefBonus:      db.RefBonus,
		MatchBonus:    db.MatchBonus,
		Priority:      db.Priority,
		BinBonusLimit: db.BinBonusLimit,
		FirstBonus: domain.RankGetBonus{
			Amount: db.FirstBonus.Amount,
			Months: db.FirstBonus.Months,
		},
		ApproveBonus: domain.RankGetBonus{
			Amount: db.ApproveBonus.Amount,
			Months: db.ApproveBonus.Months,
		},
	}
}

func (mr *Repo) rankConvertToDb(rank *domain.Rank) *rankDB {
	teamCondition := make([]rankTeamConditionDB, 0, len(rank.TeamCondition))
	for _, tc := range rank.TeamCondition {
		teamCondition = append(teamCondition, rankTeamConditionDB{
			RankCode: tc.RankCode,
			TeamType: tc.TeamType,
			IsRef:    tc.IsRef,
			Count:    tc.Count,
		})
	}
	return &rankDB{
		Code:          rank.Code,
		MinCv:         rank.MinCv,
		TeamCondition: teamCondition,
		BinBonus:      rank.BinBonus,
		RefBonus:      rank.RefBonus,
		MatchBonus:    rank.MatchBonus,
		Priority:      rank.Priority,
		BinBonusLimit: rank.BinBonusLimit,
		FirstBonus: rankGetBonusDB{
			Amount: rank.FirstBonus.Amount,
			Months: rank.FirstBonus.Months,
		},
		ApproveBonus: rankGetBonusDB{
			Amount: rank.ApproveBonus.Amount,
			Months: rank.ApproveBonus.Months,
		},
	}
}

func (mr *Repo) rankEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	planIndexes := []mongo.IndexModel{
		{Keys: bson.M{"code": 1}, Options: options.Index().SetUnique(true)},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(rankTable).Indexes().CreateMany(ctx, planIndexes)
	if err != nil {
		return domain.NewError(rankErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
