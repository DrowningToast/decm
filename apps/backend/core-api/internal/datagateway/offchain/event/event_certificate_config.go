package offchain_datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

type EventCertificateConfigDataGateway interface {
	CreateEventCertificateConfig(ctx context.Context, params generated.CreateEventCertificateConfigParams) (*entity.EventCertificateConfig, error)
	GetEventCertificateConfigByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificateConfig, error)
	GetEventCertificateConfigByEventID(ctx context.Context, eventId uuid.UUID) (*entity.EventCertificateConfig, error)
	UpdateEventCertificateConfig(ctx context.Context, params generated.UpdateEventCertificateConfigParams) (*entity.EventCertificateConfig, error)
	UpdateEventCertificateTextConfig(ctx context.Context, params generated.UpdateEventCertificateTextConfigParams) (*entity.EventCertificateConfig, error)
	ToggleEventCertificateConfigPublished(ctx context.Context, params generated.ToggleEventCertificateConfigPublishedParams) (*entity.EventCertificateConfig, error)
	DeleteEventCertificateConfig(ctx context.Context, eventID uuid.UUID) error
}
