package repository

import (
	"context"
	"server/internal/repository/ton"

	"server/internal/repository/mongodb"
	"server/internal/repository/temporal"
	"server/internal/service/auth"
	"server/internal/service/buy"
	"server/internal/service/seed"
	"server/internal/service/user"
)

var (
	// test services interfaces
	_ auth.DbRepository = (*mongodb.Repo)(nil)
	_ user.DbRepository = (*mongodb.Repo)(nil)
	_ buy.DbRepository  = (*mongodb.Repo)(nil)
	_ seed.DbRepository = (*mongodb.Repo)(nil)

	_ buy.WfRepository = (*temporal.Repo)(nil)
)

func NewDbRepo(ctx context.Context) (*mongodb.Repo, error) {
	// connect and initialize db
	dbr, err := mongodb.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return dbr, nil
}

func NewWfRepo(ctx context.Context) (*temporal.Repo, error) {
	// connect and initialize workflow client
	tmp, err := temporal.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return tmp, nil
}

func NewTonRepo(ctx context.Context) (*ton.Repo, error) {
	// connect and initialize ton client
	tn, err := ton.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return tn, nil
}
