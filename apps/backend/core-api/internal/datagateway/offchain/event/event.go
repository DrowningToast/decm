package offchain_datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateEventParameters struct {
	ChainId                  int
	Name                     string
	ShortDescription         string
	Description              string
	StartDate                time.Time
	EndDate                  time.Time
	SeatsCount               int
	ContactNumber            string
	ContactAddress           string
	Location                 string
	GoogleMapQuery           string
	BannerStorageKey         string
	IconStorageKey           string
	OwnerCredentialID        uuid.UUID
	IsPublic                 bool
	IsBookingRequestRequired bool
	IsVerified               bool
	IsTicketTransferable     bool
	EventStatus              entity.EventStatus
}

type UpdateEventParameters struct {
	Name                     *string
	ShortDescription         *string
	Description              *string
	StartDate                *time.Time
	EndDate                  *time.Time
	SeatsCount               *int
	ContactNumber            *string
	ContactAddress           *string
	Location                 *string
	GoogleMapQuery           *string
	BannerStorageKey         *string
	IconStorageKey           *string
	OwnerCredentialID        *uuid.UUID
	EventType                *entity.EventType
	IsBookingRequestRequired *bool
	IsTicketTransferable     *bool
}

type EventDataGateway interface {
	CreateEvent(ctx context.Context, params CreateEventParameters) (*entity.Event, error)
	GetEventById(ctx context.Context, id uuid.UUID) (*entity.Event, error)
	GetViewModelById(ctx context.Context, id uuid.UUID) (*entity.Event, *entity.EventRegistrationConfig, *entity.EventContract, error)
	ListEventsByOwnerCredentialID(ctx context.Context, ownerCredentialID uuid.UUID, limitCount int32, offsetCount int32) ([]*entity.Event, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, params UpdateEventParameters) (*entity.Event, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) (*entity.Event, error)

	ListEvents(ctx context.Context, limitCount *int32, offsetCount *int32) ([]*entity.Event, error)
	ListEventsByEventAttendeeCredentialID(ctx context.Context, eventAttendeeCredentialID uuid.UUID, limitCount *int32, offsetCount *int32) ([]*entity.Event, error)
}
