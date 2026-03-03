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

type CertificateShareDataGateway interface {
	CreateCertificateShare(ctx context.Context, params CreateCertificateShareParameters) (*entity.CertificateShare, error)
	// GetCertificateShareByHandle returns nil, nil when the handle does not exist.
	GetCertificateShareByHandle(ctx context.Context, handle string) (*entity.CertificateShare, error)
}
