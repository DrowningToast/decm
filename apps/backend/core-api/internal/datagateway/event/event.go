package datagateway

import (
	"context"
	"mime/multipart"
	"time"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type CreateEventParameters struct {
	Name             string
	ShortDescription string
	Description      string
	StartDate        time.Time
	EndDate          time.Time
	SeatsCount       int
	ContactNumber    string
	ContactAddress   string
	Location         string
	GoogleMapQuery   string
	EventBanner      *multipart.FileHeader
	EventIcon        *multipart.FileHeader
}

type EventDataGateway interface {
	CreateEvent(ctx context.Context, params CreateEventParameters) (*entity.Event, error)
	GetEventById(ctx context.Context, id uuid.UUID) (*entity.Event, error)
}
