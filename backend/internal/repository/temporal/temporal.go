package temporal

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"

	"server/internal/config"
	"server/internal/domain"
	"server/internal/workflow"
)

const (
	temporalErrorSource = "[repository.temporal]"
)

var (
	WorkflowTimeRetention = 90 * 24 * time.Hour
)

type Repo struct {
	c *client.Client
}

func Connect(ctx context.Context) (*Repo, error) {
	cfg := config.Get()

	c, err := client.Dial(client.Options{
		HostPort:  cfg.TemporalURL,
		Namespace: workflow.BuyQueue,
	})

	if err != nil {
		fmt.Println(err)
		return nil, domain.NewError(temporalErrorSource).SetCode(domain.ErrConfig).Add(err)
	}

	r := &Repo{c: &c}
	err = r.ensureNamespaces(ctx)
	if err != nil {
		return nil, domain.NewError(temporalErrorSource).SetCode(domain.ErrConnect).Add(err)
	}

	return r, nil
}

func (mr *Repo) ensureNamespaces(ctx context.Context) error {
	cfg := config.Get()

	cn, err := client.NewNamespaceClient(client.Options{
		HostPort: cfg.TemporalURL,
	})
	if err != nil {
		fmt.Println(err)
		return domain.NewError(temporalErrorSource).SetCode(domain.ErrConnect).Add(err)
	}
	defer cn.Close()

	_ = mr.buyEnsureNamespaces(ctx, cn)
	return nil
}
