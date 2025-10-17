package datagateway

import (
	"context"
	"time"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type CreateEventParameters struct {
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
}

type EventDataGateway interface {
	CreateEvent(ctx context.Context, params CreateEventParameters) (*entity.Event, error)
	GetEventById(ctx context.Context, id uuid.UUID) (*entity.Event, error)
}
