package team

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"server/internal/domain"
)

// CreatePartnerApplication создаёт заявку на партнёрство
func (s *Service) CreatePartnerApplication(ctx context.Context, applicantUID string, req *domain.PartnerApplicationReq) (*domain.PartnerApplication, error) {
	// Проверяем, что заявитель существует
	applicant, err := s.db.UserGetByUID(ctx, applicantUID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// Проверяем, что партнёр существует
	_, err = s.db.UserGetByUID(ctx, req.PartnerUID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// Проверяем, что заявитель не пытается подать заявку самому себе
	if applicantUID == req.PartnerUID {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("cannot apply to yourself"))
	}

	// Проверяем, что заявитель является клиентом партнёра
	if applicant.RefUID != req.PartnerUID {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("applicant must be a client of the partner"))
	}

	// Проверяем, что у заявителя нет активной заявки к этому партнёру
	existingApplications, err := s.db.PartnerApplicationGetAllByApplicantUID(ctx, applicantUID, 100, 0)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	for _, app := range existingApplications {
		if app.PartnerUID == req.PartnerUID && app.Status == domain.PartnerApplicationStatusPending {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("application already exists"))
		}
	}

	// Проверяем, что заявитель ещё не является партнёром
	existingPartnerPlace, err := s.db.UserPlaceGetByUserUID(ctx, applicantUID)
	if err == nil && existingPartnerPlace != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("applicant already has a place in binary structure"))
	}

	// Создаём заявку
	application := &domain.PartnerApplication{
		UID:          domain.GenUID(12),
		ApplicantUID: applicantUID,
		PartnerUID:   req.PartnerUID,
		Status:       domain.PartnerApplicationStatusPending,
		Message:      req.Message,
		CreatedAt:    time.Now().UTC(),
	}

	err = s.db.PartnerApplicationCreate(ctx, application)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
	}

	return application, nil
}

// ProcessPartnerApplication обрабатывает заявку на партнёрство (одобряет или отклоняет)
func (s *Service) ProcessPartnerApplication(ctx context.Context, partnerUID string, req *domain.PartnerApplicationResponseReq) (*domain.PartnerApplication, error) {
	// Получаем заявку
	application, err := s.db.PartnerApplicationGetByUID(ctx, req.ApplicationUID)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	// Проверяем, что заявка адресована этому партнёру
	if application.PartnerUID != partnerUID {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("application is not addressed to this partner"))
	}

	// Проверяем, что заявка ещё не обработана
	if application.Status != domain.PartnerApplicationStatusPending {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("application already processed"))
	}

	// Обновляем статус заявки
	application.Status = req.Status
	application.Response = req.Response
	now := time.Now().UTC()
	application.ProcessedAt = &now
	application.ProcessedBy = partnerUID

	// Если заявка одобрена, добавляем заявителя как партнёра
	if req.Status == domain.PartnerApplicationStatusApproved {
		// Проверяем лимит партнёров
		partnersCount, err := s.db.UserPlaceCountByMatchUID(ctx, partnerUID)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCount).Add(err)
		}

		// Получаем лимит партнёров из абонемента
		subscriptions, err := s.db.UserProductGetAllByUserUIDAndProductCategoryAndDate(ctx, partnerUID, "subscription", time.Now().UTC())
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
		}

		var partnersLimit int64 = 1 // По умолчанию Researcher (1 партнёр)
		if len(subscriptions) > 0 {
			lastSubscription := subscriptions[len(subscriptions)-1]
			subscriptionInfo := domain.GetSubscriptionInfo(lastSubscription.ProductCode)
			if subscriptionInfo != nil {
				partnersLimit = subscriptionInfo.PartnersLimit
			}
		} else {
			// Если нет активного абонемента, используем базовый Researcher
			subscriptionInfo := domain.GetSubscriptionInfo(domain.SubscriptionResearcher)
			partnersLimit = subscriptionInfo.PartnersLimit
		}

		// Проверяем лимит партнёров
		if partnersCount >= partnersLimit {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrValidation).Add(errors.New("maximum partners limit reached for current subscription"))
		}

		// Получаем место партнёра в бинарной структуре
		partnerPlace, err := s.db.UserPlaceGetByUserUID(ctx, partnerUID)
		if err != nil && domain.ErrorIs(err, domain.ErrNoDocuments) {
			// Если у партнёра нет места, создаём его автоматически (корневой партнёр)
			partnerPlace = &domain.UserPlace{
				UserUID:   partnerUID,
				MatchUID:  "", // Корневой партнёр
				Row:       big.NewInt(0),
				Col:       big.NewInt(0),
				CreatedAt: time.Now().UTC(),
			}
			err = s.db.UserPlaceCreate(ctx, partnerPlace)
			if err != nil {
				return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
			}
		} else if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
		}

		// Создаём место для заявителя
		newPlace, err := s.PlaceGetNew(ctx, partnerPlace, domain.UserTeamTypeUndefined)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
		}

		newPlace.UserUID = application.ApplicantUID
		newPlace.MatchUID = partnerUID

		err = s.db.UserPlaceCreate(ctx, newPlace)
		if err != nil {
			return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrCreate).Add(err)
		}
	}

	// Сохраняем обновлённую заявку
	err = s.db.PartnerApplicationUpdate(ctx, application)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrUpdate).Add(err)
	}

	return application, nil
}

// GetPartnerApplications получает заявки на партнёрство для партнёра
func (s *Service) GetPartnerApplications(ctx context.Context, partnerUID string, limit int64, skip int64) ([]*domain.PartnerApplication, error) {
	fmt.Printf("DEBUG: GetPartnerApplications called for partnerUID: %s\n", partnerUID)
	applications, err := s.db.PartnerApplicationGetAllByPartnerUID(ctx, partnerUID, limit, skip)
	if err != nil {
		fmt.Printf("DEBUG: Error getting applications: %v\n", err)
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	fmt.Printf("DEBUG: Found %d applications for partner %s\n", len(applications), partnerUID)
	for i, app := range applications {
		fmt.Printf("DEBUG: Application %d: %s -> %s (status: %s)\n", i, app.ApplicantUID, app.PartnerUID, app.Status)
	}

	return applications, nil
}

// GetMyApplications получает заявки пользователя на партнёрство
func (s *Service) GetMyApplications(ctx context.Context, applicantUID string, limit int64, skip int64) ([]*domain.PartnerApplication, error) {
	applications, err := s.db.PartnerApplicationGetAllByApplicantUID(ctx, applicantUID, limit, skip)
	if err != nil {
		return nil, domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	return applications, nil
}

// GetPartnerApplicationsCount получает количество заявок на партнёрство для партнёра
func (s *Service) GetPartnerApplicationsCount(ctx context.Context, partnerUID string) (int64, error) {
	count, err := s.db.PartnerApplicationCountByPartnerUID(ctx, partnerUID)
	if err != nil {
		return 0, domain.NewError(teamErrorSource).SetCode(domain.ErrCount).Add(err)
	}

	return count, nil
}

// GetMyApplicationsCount получает количество заявок пользователя на партнёрство
func (s *Service) GetMyApplicationsCount(ctx context.Context, applicantUID string) (int64, error) {
	count, err := s.db.PartnerApplicationCountByApplicantUID(ctx, applicantUID)
	if err != nil {
		return 0, domain.NewError(teamErrorSource).SetCode(domain.ErrCount).Add(err)
	}

	return count, nil
}
