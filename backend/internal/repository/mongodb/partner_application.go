package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/internal/config"
	"server/internal/domain"
)

const (
	partnerApplicationTable       = "partner_applications"
	partnerApplicationErrorSource = "[repository.mongodb.partner_application]"
)

type partnerApplicationDB struct {
	UID          string                          `bson:"uid"`
	ApplicantUID string                          `bson:"applicant_uid"`
	PartnerUID   string                          `bson:"partner_uid"`
	Status       domain.PartnerApplicationStatus `bson:"status"`
	Message      string                          `bson:"message"`
	Response     string                          `bson:"response"`
	CreatedAt    time.Time                       `bson:"created_at"`
	ProcessedAt  *time.Time                      `bson:"processed_at"`
	ProcessedBy  string                          `bson:"processed_by"`
}

func (mr *Repo) partnerApplicationConvertToDb(app *domain.PartnerApplication) *partnerApplicationDB {
	return &partnerApplicationDB{
		UID:          app.UID,
		ApplicantUID: app.ApplicantUID,
		PartnerUID:   app.PartnerUID,
		Status:       app.Status,
		Message:      app.Message,
		Response:     app.Response,
		CreatedAt:    app.CreatedAt,
		ProcessedAt:  app.ProcessedAt,
		ProcessedBy:  app.ProcessedBy,
	}
}

func (mr *Repo) partnerApplicationConvertFromDb(appDb *partnerApplicationDB) *domain.PartnerApplication {
	return &domain.PartnerApplication{
		UID:          appDb.UID,
		ApplicantUID: appDb.ApplicantUID,
		PartnerUID:   appDb.PartnerUID,
		Status:       appDb.Status,
		Message:      appDb.Message,
		Response:     appDb.Response,
		CreatedAt:    appDb.CreatedAt,
		ProcessedAt:  appDb.ProcessedAt,
		ProcessedBy:  appDb.ProcessedBy,
	}
}

func (mr *Repo) PartnerApplicationCreate(ctx context.Context, app *domain.PartnerApplication) error {
	cfg := config.Get()
	appDb := mr.partnerApplicationConvertToDb(app)

	_, err := mr.db.Database(cfg.MongoDB).Collection(partnerApplicationTable).InsertOne(ctx, appDb)
	if err != nil {
		return domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return nil
}

func (mr *Repo) PartnerApplicationGetByUID(ctx context.Context, uid string) (*domain.PartnerApplication, error) {
	cfg := config.Get()
	appDb := &partnerApplicationDB{}

	err := mr.db.Database(cfg.MongoDB).Collection(partnerApplicationTable).
		FindOne(ctx, bson.M{"uid": uid}).Decode(&appDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrNoDocuments).Add(err)
		}
		return nil, domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return mr.partnerApplicationConvertFromDb(appDb), nil
}

func (mr *Repo) PartnerApplicationGetAllByPartnerUID(ctx context.Context, partnerUID string, limit int64, skip int64) ([]*domain.PartnerApplication, error) {
	cfg := config.Get()

	filter := bson.M{"partner_uid": partnerUID}
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.D{{"created_at", -1}})

	cursor, err := mr.db.Database(cfg.MongoDB).Collection(partnerApplicationTable).Find(ctx, filter, opts)
	if err != nil {
		return nil, domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)

	var applications []*domain.PartnerApplication
	for cursor.Next(ctx) {
		var appDb partnerApplicationDB
		if err := cursor.Decode(&appDb); err != nil {
			return nil, domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		applications = append(applications, mr.partnerApplicationConvertFromDb(&appDb))
	}

	return applications, nil
}

func (mr *Repo) PartnerApplicationGetAllByApplicantUID(ctx context.Context, applicantUID string, limit int64, skip int64) ([]*domain.PartnerApplication, error) {
	cfg := config.Get()

	filter := bson.M{"applicant_uid": applicantUID}
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.D{{"createdAt", -1}})

	cursor, err := mr.db.Database(cfg.MongoDB).Collection(partnerApplicationTable).Find(ctx, filter, opts)
	if err != nil {
		return nil, domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)

	var applications []*domain.PartnerApplication
	for cursor.Next(ctx) {
		var appDb partnerApplicationDB
		if err := cursor.Decode(&appDb); err != nil {
			return nil, domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		applications = append(applications, mr.partnerApplicationConvertFromDb(&appDb))
	}

	return applications, nil
}

func (mr *Repo) PartnerApplicationUpdate(ctx context.Context, app *domain.PartnerApplication) error {
	cfg := config.Get()
	appDb := mr.partnerApplicationConvertToDb(app)

	_, err := mr.db.Database(cfg.MongoDB).Collection(partnerApplicationTable).
		UpdateOne(ctx, bson.M{"uid": app.UID}, bson.M{"$set": appDb})
	if err != nil {
		return domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return nil
}

func (mr *Repo) PartnerApplicationCountByPartnerUID(ctx context.Context, partnerUID string) (int64, error) {
	cfg := config.Get()

	filter := bson.M{"partner_uid": partnerUID}
	count, err := mr.db.Database(cfg.MongoDB).Collection(partnerApplicationTable).CountDocuments(ctx, filter)
	if err != nil {
		return 0, domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrCount).Add(err)
	}

	return count, nil
}

func (mr *Repo) PartnerApplicationCountByApplicantUID(ctx context.Context, applicantUID string) (int64, error) {
	cfg := config.Get()

	filter := bson.M{"applicant_uid": applicantUID}
	count, err := mr.db.Database(cfg.MongoDB).Collection(partnerApplicationTable).CountDocuments(ctx, filter)
	if err != nil {
		return 0, domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrCount).Add(err)
	}

	return count, nil
}

// PartnerApplicationGetExpired получает просроченные заявки (старше указанного времени)
func (mr *Repo) PartnerApplicationGetExpired(ctx context.Context, expiredAfter time.Duration) ([]*domain.PartnerApplication, error) {
	cfg := config.Get()

	// Вычисляем время, после которого заявка считается просроченной
	expiredBefore := time.Now().UTC().Add(-expiredAfter)

	// Фильтр: статус pending и создана раньше чем expiredBefore
	filter := bson.M{
		"status":     domain.PartnerApplicationStatusPending,
		"created_at": bson.M{"$lt": expiredBefore},
	}

	cursor, err := mr.db.Database(cfg.MongoDB).Collection(partnerApplicationTable).Find(ctx, filter)
	if err != nil {
		return nil, domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrFind).Add(err)
	}
	defer cursor.Close(ctx)

	var applications []*domain.PartnerApplication
	for cursor.Next(ctx) {
		var appDb partnerApplicationDB
		if err := cursor.Decode(&appDb); err != nil {
			return nil, domain.NewError(partnerApplicationErrorSource).SetCode(domain.ErrFind).Add(err)
		}
		applications = append(applications, mr.partnerApplicationConvertFromDb(&appDb))
	}

	return applications, nil
}
