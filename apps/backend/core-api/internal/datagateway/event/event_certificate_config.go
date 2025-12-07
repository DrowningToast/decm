package datagateway

import (
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

type EventCertificateConfigDataGateway interface {
	CreateEventCertificateConfig(ctx context.Context, params generated.CreateEventCertificateConfigParams) (*generated.EventCertificateConfig, error)
	GetEventCertificateConfigByEventID(ctx context.Context, eventId uuid.UUID) (*generated.EventCertificateConfig, error)
	UpdateEventCertificateConfig(ctx context.Context, params generated.UpdateEventCertificateConfigParams) (*generated.EventCertificateConfig, error)
	ToggleEventCertificateConfigPublished(ctx context.Context, params generated.ToggleEventCertificateConfigPublishedParams) (*generated.EventCertificateConfig, error)
	DeleteEventCertificateConfig(ctx context.Context, eventID uuid.UUID) error
}
