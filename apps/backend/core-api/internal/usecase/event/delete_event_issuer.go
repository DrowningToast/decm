package event

import (
	"apps/backend/common/customerror"
	"apps/backend/services/auth"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (u *EventUsecase) DeleteEventIssuer(ctx context.Context, id uuid.UUID, currentUser *auth.JwtClaims) error {
	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, currentUser.UserId)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	isVerifiedOrganizer := credential.IsVerifiedOrganizer
	if !isVerifiedOrganizer {
		return customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not a verified organizer"))
	}

	// Get the issuer to retrieve event ID and issuer credential ID
	issuer, err := u.EventIssuerDataGateway.GetEventIssuerByID(ctx, id)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Get certificate config for the event
	certificateConfig, err := u.EventCertificateConfigDg.GetEventCertificateConfigByEventID(ctx, issuer.EventID)
	if err != nil {
		// If config doesn't exist, that's okay - just proceed with issuer deletion
		// The signature might not exist either
	} else {
		// Find and delete the certificate signature for this issuer (even if not signed yet)
		signatures, err := u.EventCertificateSignatureDataGateway.GetEventCertificateSignaturesByEventCertificateConfigID(
			ctx, certificateConfig.ID)
		if err == nil {
			for _, sig := range signatures {
				if sig.IssuerCredentialId == issuer.IssuerCredentialID {
					err := u.EventCertificateSignatureDataGateway.DeleteEventCertificateSignature(ctx, sig.Id)
					if err != nil {
						return customerror.Parse(&customerror.ErrInternalServer,
							fmt.Errorf("failed to delete certificate signature for issuer: %w", err))
					}
					break // Only one signature per issuer per config
				}
			}
		}
	}

	// Delete the issuer
	err = u.EventIssuerDataGateway.DeleteEventIssuer(ctx, id)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return nil
}
