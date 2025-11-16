package datagateway

import (
	"context"
	"time"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type CreateEventCertificateParameters struct {
	EventID                 uuid.UUID
	ReceiverCredentialID    *uuid.UUID
	ReceiverEmail           *string
	Name                    *string
	AcademicInstitution     *string
	CertificateTitle        *string
	CertificateSubtitle     *string
	EventContractAddress    string
	EventCertificateAddress *string
	CertificateTokenID      *string
	Digest                  *string
}

type UpdateEventCertificateParameters struct {
	ReceiverCredentialID    *uuid.UUID
	ReceiverEmail           *string
	Name                    *string
	AcademicInstitution     *string
	CertificateTitle        *string
	CertificateSubtitle     *string
	EventContractAddress    *string
	EventCertificateAddress *string
	CertificateTokenID      *string
	RevokedAt               *time.Time
	Digest                  *string
}

type EventCertificateDataGateway interface {
	CreateEventCertificate(ctx context.Context, params CreateEventCertificateParameters) (*entity.EventCertificate, error)
	GetEventCertificateByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificate, error)
	GetEventCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error)
	GetAllEventCertificateIDsByEventID(ctx context.Context, eventID uuid.UUID) ([]uuid.UUID, error)
	UpdateEventCertificate(ctx context.Context, id uuid.UUID, params UpdateEventCertificateParameters) (*entity.EventCertificate, error)
	DeleteEventCertificate(ctx context.Context, id uuid.UUID) error
}
