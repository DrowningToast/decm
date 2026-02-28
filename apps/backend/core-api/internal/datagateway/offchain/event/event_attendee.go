package offchain_datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"time"

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
	UserSignatureId       *uuid.UUID
}

type EventAttendeeWithSignature struct {
	Id            uuid.UUID
	EventId       uuid.UUID
	CredentialId  uuid.UUID
	WalletAddress string
	JoinedAt      time.Time
	IsAccepted    bool
	UserSignature *entity.UserSignature
	// IsBroadcasted     bool
	// BroadcastedAt     *time.Time
	// Signature         string
	// SignMessage       string
	// DeadlineBlock     *int32
	// EstimatedDeadline *time.Time
}

type EventAttendeeDataGateway interface {
	GetEventAttendeeByEventIdAndCredentialId(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (*entity.EventAttendee, error)
	GetEventAttendeeWithSignature(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (*EventAttendeeWithSignature, error)
	AddParticipant(ctx context.Context, params AddParticipantParameters) (*entity.EventAttendee, error)
	ListEventAttendeesByEventID(ctx context.Context, eventId uuid.UUID) ([]entity.EventAttendee, error)
	DeleteEventAttendeeById(ctx context.Context, id uuid.UUID) error
	DeleteEventAttendeeByEventIDAndCredentialID(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) error
	HasPendingEventJoinByEventAndCredential(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (bool, error)
}
