package team

import (
	"context"
	"fmt"
	"log"
	"time"

	"server/internal/domain"
)

// ProcessExpiredApplications обрабатывает просроченные заявки (старше 4 часов)
func (s *Service) ProcessExpiredApplications(ctx context.Context) error {
	log.Println("Starting processing of expired partner applications...")

	// Получаем все pending заявки старше 4 часов
	expiredApplications, err := s.db.PartnerApplicationGetExpired(ctx, time.Hour*4)
	if err != nil {
		log.Printf("Error getting expired applications: %v", err)
		return domain.NewError(teamErrorSource).SetCode(domain.ErrFind).Add(err)
	}

	log.Printf("Found %d expired applications to process", len(expiredApplications))

	for _, application := range expiredApplications {
		err := s.processExpiredApplication(ctx, application)
		if err != nil {
			log.Printf("Error processing expired application %s: %v", application.UID, err)
			// Продолжаем обработку остальных заявок даже если одна не обработалась
			continue
		}
	}

	log.Println("Finished processing expired partner applications")
	return nil
}

// processExpiredApplication обрабатывает одну просроченную заявку
func (s *Service) processExpiredApplication(ctx context.Context, application *domain.PartnerApplication) error {
	log.Printf("Processing expired application %s: %s -> %s", 
		application.UID, application.ApplicantUID, application.PartnerUID)

	// Пытаемся автоматически принять заявку у текущего партнера
	canAccept, err := s.canPartnerAcceptApplication(ctx, application.PartnerUID)
	if err != nil {
		return fmt.Errorf("error checking partner capacity: %w", err)
	}

	if canAccept {
		// У партнера есть свободный слот - принимаем заявку
		log.Printf("Partner %s has available slot, auto-accepting application %s", 
			application.PartnerUID, application.UID)
		
		return s.autoAcceptApplication(ctx, application, application.PartnerUID)
	}

	// У партнера нет свободного слота - эскалируем по спонсорской линии
	log.Printf("Partner %s has no available slots, escalating application %s up the sponsor line", 
		application.PartnerUID, application.UID)

	return s.escalateApplicationUpSponsorLine(ctx, application)
}

// canPartnerAcceptApplication проверяет, может ли партнер принять заявку
func (s *Service) canPartnerAcceptApplication(ctx context.Context, partnerUID string) (bool, error) {
	// Получаем количество текущих партнеров
	partnersCount, err := s.db.UserPlaceCountByMatchUID(ctx, partnerUID)
	if err != nil {
		return false, fmt.Errorf("error counting partners: %w", err)
	}

	// Получаем лимит партнеров из абонемента
	subscriptions, err := s.db.UserProductGetAllByUserUIDAndProductCategoryAndDate(ctx, partnerUID, "subscription", time.Now().UTC())
	if err != nil {
		return false, fmt.Errorf("error getting subscriptions: %w", err)
	}

	var partnersLimit int64 = 1 // По умолчанию Researcher (1 партнер)
	if len(subscriptions) > 0 {
		lastSubscription := subscriptions[len(subscriptions)-1]
		subscriptionInfo := domain.GetSubscriptionInfo(lastSubscription.ProductCode)
		if subscriptionInfo != nil {
			partnersLimit = subscriptionInfo.PartnersLimit
		}
	} else {
		subscriptionInfo := domain.GetSubscriptionInfo(domain.SubscriptionResearcher)
		partnersLimit = subscriptionInfo.PartnersLimit
	}

	return partnersCount < partnersLimit, nil
}

// escalateApplicationUpSponsorLine эскалирует заявку вверх по спонсорской линии
func (s *Service) escalateApplicationUpSponsorLine(ctx context.Context, application *domain.PartnerApplication) error {
	currentPartnerUID := application.PartnerUID
	maxLevels := 10 // Максимальное количество уровней для поиска

	for level := 0; level < maxLevels; level++ {
		// Получаем спонсора текущего партнера
		partner, err := s.db.UserGetByUID(ctx, currentPartnerUID)
		if err != nil {
			return fmt.Errorf("error getting partner %s: %w", currentPartnerUID, err)
		}

		// Если у партнера нет спонсора, останавливаемся
		if partner.RefUID == "" {
			log.Printf("Reached top of sponsor line for application %s, no one can accept", application.UID)
			return s.rejectExpiredApplication(ctx, application, "No available slots in sponsor line")
		}

		sponsorUID := partner.RefUID
		log.Printf("Checking sponsor %s (level %d) for application %s", sponsorUID, level+1, application.UID)

		// Проверяем, может ли спонсор принять заявку
		canAccept, err := s.canPartnerAcceptApplication(ctx, sponsorUID)
		if err != nil {
			log.Printf("Error checking sponsor %s capacity: %v", sponsorUID, err)
			currentPartnerUID = sponsorUID
			continue
		}

		if canAccept {
			// Спонсор может принять заявку
			log.Printf("Sponsor %s can accept application %s, auto-accepting", sponsorUID, application.UID)
			
			// Обновляем заявку - меняем партнера на спонсора
			application.PartnerUID = sponsorUID
			return s.autoAcceptApplication(ctx, application, sponsorUID)
		}

		// Переходим к следующему уровню
		currentPartnerUID = sponsorUID
	}

	// Если дошли до максимального количества уровней
	log.Printf("Reached maximum escalation levels for application %s, rejecting", application.UID)
	return s.rejectExpiredApplication(ctx, application, "Maximum escalation levels reached")
}

// autoAcceptApplication автоматически принимает заявку
func (s *Service) autoAcceptApplication(ctx context.Context, application *domain.PartnerApplication, acceptingPartnerUID string) error {
	// Создаем запрос на обработку заявки
	req := &domain.PartnerApplicationResponseReq{
		ApplicationUID: application.UID,
		Status:         domain.PartnerApplicationStatusApproved,
		Response:       "Automatically approved after 4 hours",
	}

	// Используем существующую логику обработки заявки
	_, err := s.ProcessPartnerApplication(ctx, acceptingPartnerUID, req)
	if err != nil {
		return fmt.Errorf("error auto-accepting application: %w", err)
	}

	log.Printf("Successfully auto-accepted application %s by partner %s", 
		application.UID, acceptingPartnerUID)
	
	return nil
}

// rejectExpiredApplication отклоняет просроченную заявку
func (s *Service) rejectExpiredApplication(ctx context.Context, application *domain.PartnerApplication, reason string) error {
	// Обновляем статус заявки
	application.Status = domain.PartnerApplicationStatusRejected
	application.Response = fmt.Sprintf("Automatically rejected: %s", reason)
	now := time.Now().UTC()
	application.ProcessedAt = &now
	application.ProcessedBy = "system"

	// Сохраняем обновленную заявку
	err := s.db.PartnerApplicationUpdate(ctx, application)
	if err != nil {
		return fmt.Errorf("error rejecting expired application: %w", err)
	}

	log.Printf("Rejected expired application %s: %s", application.UID, reason)
	return nil
}

// ProcessExpiredApplicationByUID processes exactly one pending application by its UID.
// It mirrors the logic used by the cron-like processor, so we can reuse it as a Temporal activity.
func (s *Service) ProcessExpiredApplicationByUID(ctx context.Context, applicationUID string) error {
    // Load application
    application, err := s.db.PartnerApplicationGetByUID(ctx, applicationUID)
    if err != nil {
        return err
    }

    // If already processed, nothing to do (idempotent)
    if application.Status != domain.PartnerApplicationStatusPending {
        return nil
    }

    // Try current partner first
    canAccept, err := s.canPartnerAcceptApplication(ctx, application.PartnerUID)
    if err != nil {
        return err
    }
    if canAccept {
        return s.autoAcceptApplication(ctx, application, application.PartnerUID)
    }

    // Escalate along sponsor line
    return s.escalateApplicationUpSponsorLine(ctx, application)
}
