package certificate_share

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"context"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

// CertificateShareViewStatus describes the outcome of a share link lookup.
type CertificateShareViewStatus string

const (
	// CertificateShareViewStatusReady means the share is valid and the certificate has been claimed.
	// Certificate data is returned alongside this status.
	CertificateShareViewStatusReady CertificateShareViewStatus = "READY"

	// CertificateShareViewStatusValidButPending means the share is valid but the certificate
	// has not been claimed on-chain yet.
	CertificateShareViewStatusValidButPending CertificateShareViewStatus = "VALID_BUT_PENDING"

	// CertificateShareViewStatusPasswordLocked means the share exists but requires a password.
	CertificateShareViewStatusPasswordLocked CertificateShareViewStatus = "PASSWORD_LOCKED"
)

// GetCertificateShareByHandle resolves a public share link by its handle (no password).
//
// Return matrix:
//   - ErrNotFound error                — handle not found (404)
//   - (nil, &PasswordLocked, nil)      — share requires a password
//   - (nil, &ValidButPending, nil)     — share valid, certificate not yet claimed
//   - (cert, &Ready, nil)             — share valid and certificate claimed
func (uc *CertificateShareUsecase) GetCertificateShareByHandle(ctx context.Context, handle string) (*entity.EventCertificate, *CertificateShareViewStatus, error) {
	share, err := uc.CertificateShareDg.GetCertificateShareByHandle(ctx, handle)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if share == nil {
		return nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("certificate share not found"))
	}

	if share.Password != nil {
		status := CertificateShareViewStatusPasswordLocked
		return nil, &status, nil
	}

	return uc.resolveCertificateForShare(ctx, share.EventCertificateId)
}

// GetCertificateShareByHandleWithPassword resolves a password-protected share link.
// If the share has no password the password argument is ignored.
//
// Return matrix:
//   - (nil, nil, nil)                  — handle not found
//   - (nil, &PasswordLocked, nil)      — share requires a password and the supplied password is wrong
//   - (nil, &ValidButPending, nil)     — share valid, certificate not yet claimed
//   - (cert, &Ready, nil)             — share valid and certificate claimed
func (uc *CertificateShareUsecase) GetCertificateShareByHandleWithPassword(ctx context.Context, handle string, password string) (*entity.EventCertificate, *CertificateShareViewStatus, error) {
	share, err := uc.CertificateShareDg.GetCertificateShareByHandle(ctx, handle)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if share == nil {
		return nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("certificate share not found"))
	}

	if share.Password != nil && *share.Password != password {
		status := CertificateShareViewStatusPasswordLocked
		return nil, &status, nil
	}

	return uc.resolveCertificateForShare(ctx, share.EventCertificateId)
}

// resolveCertificateForShare fetches the certificate and maps its claim status to a view status.
func (uc *CertificateShareUsecase) resolveCertificateForShare(ctx context.Context, certID uuid.UUID) (*entity.EventCertificate, *CertificateShareViewStatus, error) {
	cert, err := uc.EventCertificateDataGateway.GetEventCertificateByID(ctx, certID)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get certificate"))
	}
	if cert == nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("certificate not found for share"))
	}

	if entity.GetClaimCertificateStatus(cert) == entity.ClaimCertificateStatusPending {
		status := CertificateShareViewStatusValidButPending
		return nil, &status, nil
	}

	status := CertificateShareViewStatusReady
	return cert, &status, nil
}
