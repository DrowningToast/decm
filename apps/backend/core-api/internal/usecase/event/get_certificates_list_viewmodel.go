package event

import (
	"context"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

// CertificatesListViewModel represents certificates separated by claimed status
type CertificatesListViewModel struct {
	ClaimedCertificates   []*entity.EventCertificate `json:"claimed_certificates"`
	UnclaimedCertificates []*entity.EventCertificate `json:"unclaimed_certificates"`
	TotalClaimed          int                        `json:"total_claimed"`
	TotalUnclaimed        int                        `json:"total_unclaimed"`
}

// GetCertificatesListViewModel retrieves certificates separated by claimed/unclaimed status
// Claimed certificates are those with token_id populated
// Unclaimed certificates are those with token_id null and certificate config is published
func (uc *EventUsecase) GetCertificatesListViewModel(ctx context.Context, eventID uuid.UUID, currentUser *auth.JwtClaims) (*CertificatesListViewModel, error) {
	// Verify user has access to this event (must be host or issuer)
	event, err := uc.EventDataGateway.GetEventById(ctx, eventID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, err)
	}

	// Check if user is the event host
	if event.OwnerCredentialId != currentUser.UserId {
		// Check if user is an issuer for this event
		issuers, err := uc.EventIssuerDataGateway.GetEventIssuersByEventID(ctx, eventID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get event issuers")
		}

		isIssuer := false
		for _, issuer := range issuers {
			if issuer.IssuerCredentialID == currentUser.UserId {
				isIssuer = true
				break
			}
		}

		if !isIssuer {
			return nil, customerror.Parse(&customerror.ErrForbidden, errors.New("user does not have permission to view certificates for this event"))
		}
	}

	// Get claimed certificates (token_id IS NOT NULL)
	claimedCertificates, err := uc.EventCertificateDataGateway.GetClaimedCertificatesByEventID(ctx, eventID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get claimed certificates")
	}

	// Get unclaimed ready certificates (token_id IS NULL and config is published)
	unclaimedCertificates, err := uc.EventCertificateDataGateway.GetUnclaimedReadyCertificatesByEventID(ctx, eventID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get unclaimed ready certificates")
	}

	return &CertificatesListViewModel{
		ClaimedCertificates:   claimedCertificates,
		UnclaimedCertificates: unclaimedCertificates,
		TotalClaimed:          len(claimedCertificates),
		TotalUnclaimed:        len(unclaimedCertificates),
	}, nil
}
