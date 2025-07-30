package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/internal/config"
	"server/internal/domain"
	"time"
)

const (
	// table name in DB
	notificationTable = "notification"

	// errors prefix
	notificationErrorSource = "[repository.mongodb.notification]"
)

type notificationDB struct {
	UID       string            `bson:"uid"`
	Texts     map[string]string `bson:"texts"`
	CreatedAt time.Time         `bson:"createdAt"`
	SentAt    time.Time         `bson:"sentAt"`
	ToUserUID []string          `bson:"toUserUid"`
	FlowUID   string            `bson:"flowUid"`
}

func (mr *Repo) NotificationCreate(ctx context.Context, notification *domain.Notification) error {
	cfg := config.Get()

	notificationDb := mr.notificationConvertToDB(notification)
	_, err := mr.db.Database(cfg.MongoDB).Collection(notificationTable).InsertOne(ctx, notificationDb)
	if err != nil {
		return domain.NewError(distErrorSource).SetCode(domain.ErrCreate).Add(err)
	}
	return nil
}

func (mr *Repo) NotificationUpdate(ctx context.Context, notification *domain.Notification) error {
	cfg := config.Get()

	notificationDb := mr.notificationConvertToDB(notification)

	_, err := mr.db.Database(cfg.MongoDB).Collection(notificationTable).
		UpdateOne(ctx, bson.M{"uid": notificationDb.UID}, bson.M{"$set": notificationDb})
	if err != nil {
		return domain.NewError(distErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return nil
}

func (mr *Repo) NotificationGetByUID(ctx context.Context, uid string) (*domain.Notification, error) {
	notificationDb := &notificationDB{}
	cfg := config.Get()
	err := mr.db.Database(cfg.MongoDB).Collection(notificationTable).
		FindOne(ctx, bson.M{"uid": uid}).Decode(&notificationDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(distErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(distErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	return mr.notificationConvertFromDB(notificationDb), nil
}

func (mr *Repo) NotificationGetAllByUserUID(ctx context.Context, userUID string) ([]*domain.Notification, error) {
	cfg := config.Get()
	opts := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := mr.db.Database(cfg.MongoDB).Collection(notificationTable).
		Find(ctx, bson.M{"$or": []bson.M{
			{"toUserUid": userUID},
			{"toUserUid": bson.M{"$exists": false}},
			{"toUserUid": bson.M{"$size": 0}},
		}}, opts)
	if err != nil {
		return nil, domain.NewError(distErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)
	var notifications []*domain.Notification
	for cursor.Next(ctx) {
		notificationDb := &notificationDB{}
		err := cursor.Decode(&notificationDb)
		if err != nil {
			return nil, domain.NewError(distErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		notifications = append(notifications, mr.notificationConvertFromDB(notificationDb))
	}
	return notifications, nil
}

func (mr *Repo) notificationConvertToDB(n *domain.Notification) *notificationDB {
	return &notificationDB{
		UID:       n.UID,
		Texts:     n.Texts,
		CreatedAt: n.CreatedAt,
		SentAt:    n.SentAt,
		ToUserUID: n.ToUserUID,
		FlowUID:   n.FlowUID,
	}
}

func (mr *Repo) notificationConvertFromDB(n *notificationDB) *domain.Notification {
	return &domain.Notification{
		UID:       n.UID,
		Texts:     n.Texts,
		CreatedAt: n.CreatedAt,
		SentAt:    n.SentAt,
		ToUserUID: n.ToUserUID,
		FlowUID:   n.FlowUID,
	}
}

func (mr *Repo) notificationEnsureIndexes(ctx context.Context) error {
	cfg := config.Get()
	userIndexes := []mongo.IndexModel{
		{Keys: bson.M{"uid": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"toUserUid": 1}, Options: options.Index()},
		{Keys: bson.M{"sentAt": -1}, Options: options.Index()},
		{Keys: bson.M{"createdAt": -1}, Options: options.Index()},
	}
	_, err := mr.db.Database(cfg.MongoDB).Collection(distTable).Indexes().CreateMany(ctx, userIndexes)
	if err != nil {
		return domain.NewError(distErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
