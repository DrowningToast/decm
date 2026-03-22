package certificate_share

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"crypto/ecdsa"
	"log/slog"
	"time"

	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	certificatecontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/certificate_contract"

	"github.com/google/uuid"
)

// CertificateImageGenerator generates PNG images for certificates.
// Implemented by EventUsecase.
type CertificateImageGenerator interface {
	GenerateCertificateImage(ctx context.Context, certificateID uuid.UUID) ([]byte, error)
}

// CertificateShareViewModel is a safe representation of a CertificateShare for use in viewmodels.
// It omits the raw hashed password and instead exposes a HasPassword boolean.
type CertificateShareViewModel struct {
	Id                 uuid.UUID `json:"id"`
	EventCertificateId uuid.UUID `json:"event_certificate_id"`
	Active             bool      `json:"active"`
	Handle             string    `json:"handle"`
	HasPassword        bool      `json:"has_password"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// NewCertificateShareViewModel maps a CertificateShare entity to a CertificateShareViewModel.
// Returns nil when share is nil.
func NewCertificateShareViewModel(share *entity.CertificateShare) *CertificateShareViewModel {
	if share == nil {
		return nil
	}
	return &CertificateShareViewModel{
		Id:                 share.Id,
		EventCertificateId: share.EventCertificateId,
		Active:             share.Active,
		Handle:             share.Handle,
		HasPassword:        share.Password != nil,
		CreatedAt:          share.CreatedAt,
		UpdatedAt:          share.UpdatedAt,
	}
}

type CertificateShareUsecase struct {
	EventCertificateDataGateway  eventdatagateway.EventCertificateDataGateway
	CertificateShareDg           eventdatagateway.CertificateShareDataGateway
	CertificateContractFactoryDg certificatecontract_datagateway.CertificateContractFactoryDataGateway
	CertificateImageGenerator    CertificateImageGenerator
	// BackendPrivateKey is used to decrypt BackendEncryptedUserData on the certificate payload.
	// When nil, decryption is skipped and DecryptedUserData in the result will be nil.
	BackendPrivateKey *ecdsa.PrivateKey
	Logger            *slog.Logger
}

func NewCertificateShareUsecase(
	eventCertificateDataGateway eventdatagateway.EventCertificateDataGateway,
	certificateShareDg eventdatagateway.CertificateShareDataGateway,
	certificateContractFactoryDg certificatecontract_datagateway.CertificateContractFactoryDataGateway,
	certificateImageGenerator CertificateImageGenerator,
	backendPrivateKey *ecdsa.PrivateKey,
	logger *slog.Logger,
) *CertificateShareUsecase {
	return &CertificateShareUsecase{
		EventCertificateDataGateway:  eventCertificateDataGateway,
		CertificateShareDg:           certificateShareDg,
		CertificateContractFactoryDg: certificateContractFactoryDg,
		CertificateImageGenerator:    certificateImageGenerator,
		BackendPrivateKey:            backendPrivateKey,
		Logger:                       logger,
	}
}
