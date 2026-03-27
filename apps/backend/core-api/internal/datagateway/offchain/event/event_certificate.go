package offchain_datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateEventCertificateParameters struct {
	EventID                 uuid.UUID
	ReceiverCredentialID    *uuid.UUID
	ReceiverEmail           *string
	ReceiverWalletAddress   *string
	Name                    *string
	AcademicInstitution     *string
	CertificateTitle        *string
	CertificateSubtitle     *string
	EventContractAddress    string
	EventCertificateAddress *string
	CertificateTokenID      *string
	Digest                  *string
	UserClaimSignatureId    *uuid.UUID
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
	UserClaimSignatureId    *uuid.UUID
}

type EventCertificateWithSignature struct {
	Id                 uuid.UUID
	EventId            uuid.UUID
	CredentialId       *uuid.UUID
	TokenId            *string
	JoinedAt           time.Time
	IsBroadcasted      bool
	BroadcastedAt      *time.Time
	SignatureCreatedAt *time.Time
	EstimatedDeadline  *time.Time
	Signature          string
	SignMessage        string
}

type EventCertificateDataGateway interface {
	CreateEventCertificate(ctx context.Context, params CreateEventCertificateParameters) (*entity.EventCertificate, error)
	GetEventCertificateByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificate, error)
	GetEventCertificateWithSignature(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*EventCertificateWithSignature, error)
	GetEventCertificateByInboxMessageID(ctx context.Context, inboxMessageID uuid.UUID) (*entity.EventCertificate, error)
	GetEventCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error)
	GetAllEventCertificateIDsByEventID(ctx context.Context, eventID uuid.UUID) ([]uuid.UUID, error)
	GetClaimedCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error)
	GetUnclaimedReadyCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error)
	GetClaimedCertificatesByCredentialID(ctx context.Context, credentialID uuid.UUID, email *string, walletAddress *string) ([]*entity.EventCertificate, error)
	GetUnclaimedReadyCertificatesByCredentialID(ctx context.Context, credentialID uuid.UUID, email *string, walletAddress *string) ([]*entity.EventCertificate, error)
	UpdateEventCertificate(ctx context.Context, id uuid.UUID, params UpdateEventCertificateParameters) (*entity.EventCertificate, error)
	UpdateEventCertificateInboxMessageID(ctx context.Context, id uuid.UUID, inboxMessageID uuid.UUID) (*entity.EventCertificate, error)
	DeleteEventCertificate(ctx context.Context, id uuid.UUID) error
}
