package datagateway

import (
	"context"
	"time"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type CreateEventRegistrationInvitationParameters struct {
	EventID        uuid.UUID
	InboxMessageID uuid.UUID
	ValidUntil     *time.Time
	Code           *string
}

type EventRegistrationInvitationDataGateway interface {
	CreateEventRegistrationInvitation(ctx context.Context, params CreateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error)
	GetEventRegistrationInvitationByID(ctx context.Context, id uuid.UUID) (*entity.EventRegistrationInvitation, error)
	GetEventRegistrationInvitationsByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventRegistrationInvitation, error)
	UpdateEventRegistrationInvitation(ctx context.Context, id uuid.UUID, params UpdateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error)
	DeleteEventRegistrationInvitation(ctx context.Context, id uuid.UUID) error
}

type UpdateEventRegistrationInvitationParameters struct {
	ValidUntil  *time.Time
	Code        *string
	CancelledAt *time.Time
}
