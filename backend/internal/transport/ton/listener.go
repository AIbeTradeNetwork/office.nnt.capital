package ton

import (
	"context"
	"server/internal/domain"
	"server/internal/repository"
	"server/internal/service/user"
)

const (
	listenerErrorSource = "[transport.ton.listener]"
)

func StartListener() error {
	ctx := context.Background()

	tonRepo, err := repository.NewTonRepo(ctx)
	if err != nil {
		return domain.NewError(listenerErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	wfRepo, err := repository.NewWfRepo(ctx)
	if err != nil {
		return domain.NewError(listenerErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	dbRepo, err := repository.NewDbRepo(ctx)
	if err != nil {
		return domain.NewError(listenerErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	userService := user.NewUserService(dbRepo, wfRepo, tonRepo, nil, nil, nil)

	// blocking call
	err = userService.ListenTonDeposits(ctx)
	if err != nil {
		return err
	}

	return nil
}
