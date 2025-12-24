package datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"

	"github.com/google/uuid"
)

type AddParticipantParameters struct {
	EventId               uuid.UUID
	CredentialId          uuid.UUID
	ContractAddress       string
	Signature             []byte
	IsParticipantAccepted bool
	FirstName             *string
	LastName              *string
	Email                 *string
	PhoneNumber           *string
	AcademicInstitution   *string
	AcademicEmail         *string
	Address               *string
	Bio                   *string
}

type EventAttendeeDataGateway interface {
	GetEventAttendeeByEventIdAndCredentialId(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (*entity.EventAttendee, error)
	AddParticipant(ctx context.Context, params AddParticipantParameters) (*entity.EventAttendee, error)
	ListEventAttendeesByEventID(ctx context.Context, eventId uuid.UUID) ([]entity.EventAttendee, error)
}
