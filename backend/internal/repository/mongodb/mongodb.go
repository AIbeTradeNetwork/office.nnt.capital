package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/internal/config"
	"server/internal/domain"
)

const (
	mongodbErrorSource = "[repository.mongodb]"
)

type Repo struct {
	db *mongo.Client
}

func Connect(ctx context.Context) (*Repo, error) {
	cfg := config.Get()

	if cfg.MongoURL == "" {
		return nil, domain.NewError(mongodbErrorSource).SetCode(domain.ErrConfig)
	}

	// Debug
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			//log.Print(evt.Command)
		},
	}

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.MongoURL).SetServerAPIOptions(serverAPI).SetMonitor(cmdMonitor)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, domain.NewError(mongodbErrorSource).SetCode(domain.ErrConnect).Add(err)
	}

	repo := &Repo{client}

	// Send a ping to confirm a successful connection
	if err := repo.ping(ctx); err != nil {
		return nil, domain.NewError(mongodbErrorSource).SetCode(domain.ErrConnect).Add(err)
	}

	return repo, nil
}

func (mr *Repo) DropTest(ctx context.Context) error {
	cfg := config.Get()
	if cfg.Env != "test" {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrConfig)
	}
	return mr.db.Database(cfg.MongoDB).Drop(ctx)
}

func (mr *Repo) Disconnect(ctx context.Context) error {
	return mr.db.Disconnect(ctx)
}

func (mr *Repo) ping(ctx context.Context) error {
	// Send a ping to check connection
	var result bson.M
	if err := mr.db.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrConnect).Add(err)
	}
	return nil
}

func (mr *Repo) EnsureIndexes(ctx context.Context) error {
	// user
	err := mr.userEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// user_auth
	err = mr.userAuthEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// user_place
	err = mr.userPlaceEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// user_config
	err = mr.userConfigEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// buy
	err = mr.buyEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// plan
	err = mr.planEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// rank
	err = mr.rankEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// user_plan
	err = mr.userPlanEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// user_rank
	err = mr.userRankEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// user_activity
	err = mr.userActivityEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// transaction
	err = mr.transactionEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// distributor
	err = mr.distEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// claim
	err = mr.claimEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// user_claim
	err = mr.userClaimEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// user_balance
	err = mr.userBalanceEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// currency
	err = mr.currencyEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// notification
	err = mr.notificationEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// deposit
	err = mr.depositEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// product
	err = mr.productEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// task
	err = mr.taskEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// user_product
	err = mr.userProductEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// combo
	err = mr.comboEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// safe
	err = mr.safeEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// user safe
	err = mr.ensureUserSafeIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// TODO: Add indexes

	return nil
}
