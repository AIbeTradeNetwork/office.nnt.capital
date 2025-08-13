package domain

import "time"

type PartnerApplicationStatus string

const (
	PartnerApplicationStatusPending   PartnerApplicationStatus = "pending"
	PartnerApplicationStatusApproved  PartnerApplicationStatus = "approved"
	PartnerApplicationStatusRejected  PartnerApplicationStatus = "rejected"
)

type PartnerApplication struct {
	UID           string
	ApplicantUID  string // UID пользователя, подающего заявку
	PartnerUID    string // UID партнёра, к которому подаётся заявка
	Status        PartnerApplicationStatus
	Message       string // Сообщение от заявителя
	Response      string // Ответ партнёра (при отклонении)
	CreatedAt     time.Time
	ProcessedAt   *time.Time
	ProcessedBy   string // UID партнёра, который обработал заявку
}

type PartnerApplicationReq struct {
	PartnerUID string `json:"partner_uid"`
	Message    string `json:"message"`
}

type PartnerApplicationResponseReq struct {
	ApplicationUID string `json:"application_uid"`
	Status         PartnerApplicationStatus `json:"status"`
	Response       string `json:"response"`
} 