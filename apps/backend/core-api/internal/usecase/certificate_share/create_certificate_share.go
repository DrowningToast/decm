package certificate_share

import (
	"apps/backend/common/customerror"
	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

// CreateCertificateShare creates an inactive certificate_share row with a randomly generated handle.
// Only certificates with PENDING or CLAIMED status can be shared.
// hashedPassword should be the Argon2id-hashed password (or nil for no password protection).
func (uc *CertificateShareUsecase) CreateCertificateShare(ctx context.Context, currentUser *auth.JwtClaims, certificateID uuid.UUID, hashedPassword *string) (*entity.CertificateShare, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	certificate, err := uc.EventCertificateDataGateway.GetEventCertificateByID(ctx, certificateID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if certificate == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("certificate not found"))
	}

	// Only the receiver may create a share link.
	// Ownership is proven by credential ID match (primary) or email match (fallback, for
	// certificates imported before the user registered).
	isOwner := false
	if certificate.ReceiverCredentialId != nil {
		isOwner = *certificate.ReceiverCredentialId == currentUser.UserId
	} else if certificate.ReceiverEmail != nil {
		isOwner = currentUser.Email != nil && *currentUser.Email == *certificate.ReceiverEmail
	}
	if !isOwner {
		return nil, customerror.Parse(&customerror.ErrForbidden, errors.New("not authorized to share this certificate"))
	}

	status := entity.GetClaimCertificateStatus(certificate)
	if status != entity.ClaimCertificateStatusPending && status != entity.ClaimCertificateStatusClaimed {
		return nil, customerror.NewWithPreset(
			&customerror.ErrInvalidArgument,
			errors.Errorf("certificate cannot be shared in status %s", status),
		)
	}

	// Return the existing share if one already exists for this certificate.
	existing, err := uc.CertificateShareDg.GetCertificateShareByEventCertificateID(ctx, certificateID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to check existing certificate share"))
	}
	if existing != nil {
		return existing, nil
	}

	handle, err := generateShareHandle()
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to generate share handle"))
	}

	share, err := uc.CertificateShareDg.CreateCertificateShare(ctx, eventdatagateway.CreateCertificateShareParameters{
		EventCertificateId: certificateID,
		Active:             false,
		Handle:             handle,
		Password:           hashedPassword,
	})
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to create certificate share"))
	}

	return share, nil
}

// generateShareHandle returns a 32-character random hex string suitable for use as a URL handle.
func generateShareHandle() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
