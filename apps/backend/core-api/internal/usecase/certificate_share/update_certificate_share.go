package certificate_share

import (
	"apps/backend/common/customerror"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

// UpdateCertificateShare updates the password and/or active status on an existing certificate share.
// passwordUpdate encodes whether to keep, remove, or replace the password (see PasswordUpdate).
// active, when non-nil, sets whether the share link is publicly accessible.
func (uc *CertificateShareUsecase) UpdateCertificateShare(
	ctx context.Context,
	currentUser *auth.JwtClaims,
	shareID uuid.UUID,
	passwordUpdate event_datagateway.PasswordUpdate,
	active *bool,
) (*entity.CertificateShare, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	share, err := uc.CertificateShareDg.GetCertificateShareByID(ctx, shareID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if share == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("share not found"))
	}

	certificate, err := uc.EventCertificateDataGateway.GetEventCertificateByID(ctx, share.EventCertificateId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if certificate == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("certificate not found"))
	}

	ownerByCredential := certificate.ReceiverCredentialId != nil && *certificate.ReceiverCredentialId == currentUser.UserId
	ownerByEmail := certificate.ReceiverEmail != nil && currentUser.Email != nil && *certificate.ReceiverEmail == *currentUser.Email
	if !ownerByCredential && !ownerByEmail {
		return nil, customerror.Parse(&customerror.ErrForbidden, errors.New("not authorized to update this share"))
	}

	updated, err := uc.CertificateShareDg.UpdateCertificateShare(ctx, shareID, event_datagateway.UpdateCertificateShareParameters{
		Password: passwordUpdate,
		Active:   active,
	})
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to update certificate share"))
	}

	return updated, nil
}
