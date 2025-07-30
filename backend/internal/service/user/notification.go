package user

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"server/internal/domain"
	"server/internal/provider/telegram"
)

func (s *Service) Notifications(ctx context.Context, user *domain.User) ([]*domain.Notification, error) {
	return s.db.NotificationGetAllByUserUID(ctx, user.UID)
}

func (s *Service) Notify(ctx context.Context, notification *domain.Notification) error {
	return s.wf.NotifyCreate(ctx, notification)
}

func (s *Service) NotifyGetAllTgID(ctx context.Context, notification *domain.Notification, page int64) ([]*domain.UserTg, error) {
	var limit int64 = 5000
	skip := limit * (page - 1)
	userAuths, err := s.db.UserAuthGetAllUserTg(ctx, notification.ToUserUID, domain.AuthTypeTelegram, limit, skip)
	if err != nil {
		return nil, err
	}

	return userAuths, nil
}

func (s *Service) NotifyCreate(ctx context.Context, notification *domain.Notification) (*domain.Notification, error) {
	ex, err := s.db.NotificationGetByUID(ctx, notification.UID)
	if ex != nil {
		return ex, nil
	}

	err = s.db.NotificationCreate(ctx, notification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *Service) NotifyUpdate(ctx context.Context, notification *domain.Notification) (*domain.Notification, error) {
	err := s.db.NotificationUpdate(ctx, notification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *Service) NotifyTgSend(ctx context.Context, telegramID int64, text string) (int64, error) {
	c := telegram.NewClient()
	defer c.Close()

	_, err := c.SendMessage(tg.UserID(telegramID), text).
		ParseMode(tg.HTML).
		LinkPreviewOptions(tg.LinkPreviewOptions{IsDisabled: true}).
		Do(ctx)

	if err != nil {
		fmt.Println(err)
	}

	return telegramID, nil
}
