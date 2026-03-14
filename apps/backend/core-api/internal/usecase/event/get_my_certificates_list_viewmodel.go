package event

import (
	"apps/backend/core-api/internal/entity"
	certificate_share "apps/backend/core-api/internal/usecase/certificate_share"
	"apps/backend/services/auth"
	"context"

	"github.com/cockroachdb/errors"
)

type ClaimedCertificateViewModel struct {
	entity.EventCertificate
	Status           ClaimCertificateStatus                              `json:"status"`
	CertificateShare *certificate_share.CertificateShareViewModel `json:"certificate_share,omitempty"`
}

type UnclaimedCertificateViewModel struct {
	entity.EventCertificate
	CertificateShare *certificate_share.CertificateShareViewModel `json:"certificate_share,omitempty"`
}

// MyCertificatesListViewModel represents current user's certificates separated by claimed status
type MyCertificatesListViewModel struct {
	ClaimedCertificates   []*ClaimedCertificateViewModel   `json:"claimed_certificates"`
	UnclaimedCertificates []*UnclaimedCertificateViewModel `json:"unclaimed_certificates"`
	TotalClaimed          int                              `json:"total_claimed"`
	TotalUnclaimed        int                              `json:"total_unclaimed"`
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
	claimedCertificateViewModels := make([]*ClaimedCertificateViewModel, len(claimedCertificates))
	for i, certificate := range claimedCertificates {
		share, err := uc.CertificateShareDataGateway.GetCertificateShareByEventCertificateID(ctx, certificate.Id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get certificate share")
		}
		claimedCertificateViewModels[i] = &ClaimedCertificateViewModel{
			EventCertificate: *certificate,
			Status:           GetClaimCertificateStatus(certificate),
			CertificateShare: certificate_share.NewCertificateShareViewModel(share),
		}
	}

	// Get unclaimed ready certificates for current user (token_id IS NULL and config is published)
	// Query by both credential_id and email to catch certificates assigned before user registration
	unclaimedCertificates, err := uc.EventCertificateDataGateway.GetUnclaimedReadyCertificatesByCredentialID(ctx, currentUser.UserId, currentUser.Email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get unclaimed ready certificates")
	}
	unclaimedCertificateViewModels := make([]*UnclaimedCertificateViewModel, len(unclaimedCertificates))
	for i, certificate := range unclaimedCertificates {
		share, err := uc.CertificateShareDataGateway.GetCertificateShareByEventCertificateID(ctx, certificate.Id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get certificate share")
		}
		unclaimedCertificateViewModels[i] = &UnclaimedCertificateViewModel{
			EventCertificate: *certificate,
			CertificateShare: certificate_share.NewCertificateShareViewModel(share),
		}
	}

	return &MyCertificatesListViewModel{
		ClaimedCertificates:   claimedCertificateViewModels,
		UnclaimedCertificates: unclaimedCertificateViewModels,
		TotalClaimed:          len(claimedCertificateViewModels),
		TotalUnclaimed:        len(unclaimedCertificateViewModels),
	}, nil
}
