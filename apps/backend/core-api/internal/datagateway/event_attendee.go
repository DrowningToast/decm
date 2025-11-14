package datagateway

import (
	"context"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type EventAttendeeDataGateway interface {
	GetEventAttendeeByEventIdAndCredentialId(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*entity.EventAttendee, error)
}
