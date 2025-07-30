package auth

import (
	"context"
)

func (s *Service) CheckRefTasks(ctx context.Context, userUid string) error {
	ctx = context.Background()

	if userUid == "" {
		return nil
	}

	dTask, _ := s.db.TaskGetByRefUID(ctx, userUid)
	if dTask == nil {
		return nil
	}

	err := s.db.TaskIncRefCount(ctx, dTask.Code)
	if err != nil {
		return err
	}

	return nil
}
