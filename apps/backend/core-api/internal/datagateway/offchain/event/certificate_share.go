package offchain_datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"

	"github.com/google/uuid"
)

type CreateCertificateShareParameters struct {
	EventCertificateId uuid.UUID
	Active             bool
	Handle             string
	Password           *string
}

type UpdateCertificateShareParameters struct {
	Password *string
	Active   *bool
}

type CertificateShareDataGateway interface {
	CreateCertificateShare(ctx context.Context, params CreateCertificateShareParameters) (*entity.CertificateShare, error)
	// GetCertificateShareByHandle returns nil, nil when the handle does not exist.
	GetCertificateShareByHandle(ctx context.Context, handle string) (*entity.CertificateShare, error)
	// GetCertificateShareByID returns nil, nil when the share does not exist.
	GetCertificateShareByID(ctx context.Context, id uuid.UUID) (*entity.CertificateShare, error)
	// GetCertificateShareByEventCertificateID returns nil, nil when no share exists for the certificate.
	GetCertificateShareByEventCertificateID(ctx context.Context, eventCertificateID uuid.UUID) (*entity.CertificateShare, error)
	UpdateCertificateShare(ctx context.Context, id uuid.UUID, params UpdateCertificateShareParameters) (*entity.CertificateShare, error)
}
