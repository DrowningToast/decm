package event

import (
	"context"

	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/cockroachdb/errors"
)

// MyCertificatesListViewModel represents current user's certificates separated by claimed status
type MyCertificatesListViewModel struct {
	ClaimedCertificates   []*entity.EventCertificate `json:"claimed_certificates"`
	UnclaimedCertificates []*entity.EventCertificate `json:"unclaimed_certificates"`
	TotalClaimed          int                        `json:"total_claimed"`
	TotalUnclaimed        int                        `json:"total_unclaimed"`
}

// GetMyCertificatesListViewModel retrieves current user's certificates separated by claimed/unclaimed status
// Claimed certificates are those with token_id populated
// Unclaimed certificates are those with token_id null and certificate config is published
// Fetches certificates by both credential_id and email (for OAuth users)
func (uc *EventUsecase) GetMyCertificatesListViewModel(ctx context.Context, currentUser *auth.JwtClaims) (*MyCertificatesListViewModel, error) {
	// Get claimed certificates for current user (token_id IS NOT NULL)
	// Query by both credential_id and email to catch certificates assigned before user registration
	claimedCertificates, err := uc.EventCertificateDataGateway.GetClaimedCertificatesByCredentialID(ctx, currentUser.UserId, currentUser.Email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get claimed certificates")
	}

	// Get unclaimed ready certificates for current user (token_id IS NULL and config is published)
	// Query by both credential_id and email to catch certificates assigned before user registration
	unclaimedCertificates, err := uc.EventCertificateDataGateway.GetUnclaimedReadyCertificatesByCredentialID(ctx, currentUser.UserId, currentUser.Email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get unclaimed ready certificates")
	}

	return &MyCertificatesListViewModel{
		ClaimedCertificates:   claimedCertificates,
		UnclaimedCertificates: unclaimedCertificates,
		TotalClaimed:          len(claimedCertificates),
		TotalUnclaimed:        len(unclaimedCertificates),
	}, nil
}
