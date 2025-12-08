package event

import (
	"context"
	"decm-database/go/generated"
	"errors"

	"apps/backend/common/customerror"
	eventdatagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/services/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateEventIssuerParams struct {
	EventID            uuid.UUID
	IssuerCredentialID uuid.UUID
}

func (u *EventUsecase) UpdateEventIssuer(ctx context.Context, eventID uuid.UUID, params []UpdateEventIssuerParams, currentUser *auth.JwtClaims) ([]generated.EventIssuer, error) {
	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	isVerifiedOrganizer := credential.IsVerifiedOrganizer
	if !isVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not a verified organizer"))
	}

	// Get all current issuers for this event
	currentIssuers, err := u.EventIssuerDataGateway.GetEventIssuersByEventID(ctx, eventID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Create a map of desired issuer credential IDs for quick lookup
	desiredIssuerCredentialIDs := make(map[uuid.UUID]bool)
	for _, param := range params {
		desiredIssuerCredentialIDs[param.IssuerCredentialID] = true
	}

	// Create a map of current issuer credential IDs to issuer records
	currentIssuerMap := make(map[uuid.UUID]generated.EventIssuer)
	for _, issuer := range currentIssuers {
		currentIssuerMap[issuer.IssuerCredentialID] = issuer
	}

	// Get all existing certificates for this event (needed for signature management)
	certificates, err := u.EventCertificateDataGateway.GetEventCertificatesByEventID(ctx, eventID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Get sign message and host signature from existing certificate signatures (if any)
	var signMessage *string
	var signMessageDigest *string
	var hostSignature string
	if len(certificates) > 0 {
		existingSigs, err := u.EventCertificateSignatureDataGateway.GetEventCertificateSignaturesByEventCertificateID(
			ctx, certificates[0].Id)
		if err == nil && len(existingSigs) > 0 {
			signMessage = existingSigs[0].SignMessage
			signMessageDigest = existingSigs[0].SignMessageDigest
			hostSignature = existingSigs[0].HostSignature
		}
	}

	// Step 1: Add new issuers and create their certificate signatures
	for _, param := range params {
		issuerCredentialID := param.IssuerCredentialID

		// Check if issuer already exists
		_, exists := currentIssuerMap[issuerCredentialID]
		if !exists {
			// Create new issuer
			_, err := u.EventIssuerDataGateway.CreateEventIssuer(ctx, generated.CreateEventIssuerParams{
				EventID:            eventID,
				IssuerCredentialID: issuerCredentialID,
				IsSigned:           0,
				Signature:          pgtype.Text{},
			})
			if err != nil {
				return nil, customerror.Parse(&customerror.ErrInternalServer, err)
			}

			// Create certificate signatures for existing certificates
			if len(certificates) > 0 && signMessage != nil && signMessageDigest != nil {
				for _, cert := range certificates {
					_, err := u.EventCertificateSignatureDataGateway.CreateEventCertificateSignature(ctx,
						eventdatagateway.CreateEventCertificateSignatureParameters{
							EventCertificateID: cert.Id,
							IssuerCredentialID: issuerCredentialID,
							IssuerSignature:    nil, // Will be set when issuer signs
							HostSignature:      hostSignature,
							SignMessage:        signMessage,
							SignMessageDigest:  signMessageDigest,
						})
					if err != nil {
						return nil, customerror.Parse(&customerror.ErrInternalServer, err)
					}
				}
			}
		}
	}

	// Step 2: Remove issuers that are not in the desired list
	for _, currentIssuer := range currentIssuers {
		if !desiredIssuerCredentialIDs[currentIssuer.IssuerCredentialID] {
			// Delete certificate signatures for this issuer
			for _, cert := range certificates {
				signatures, err := u.EventCertificateSignatureDataGateway.GetEventCertificateSignaturesByEventCertificateID(
					ctx, cert.Id)
				if err != nil {
					// Continue with other certificates even if one fails
					continue
				}

				for _, sig := range signatures {
					if sig.IssuerCredentialId == currentIssuer.IssuerCredentialID {
						err := u.EventCertificateSignatureDataGateway.DeleteEventCertificateSignature(ctx, sig.Id)
						if err != nil {
							// Log error but continue deleting other signatures
							// We want to clean up as much as possible
							continue
						}
					}
				}
			}

			// Delete the issuer
			err := u.EventIssuerDataGateway.DeleteEventIssuer(ctx, currentIssuer.ID)
			if err != nil {
				return nil, customerror.Parse(&customerror.ErrInternalServer, err)
			}
		}
	}

	// Return updated issuer list
	updatedIssuers, err := u.EventIssuerDataGateway.GetEventIssuersByEventID(ctx, eventID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return updatedIssuers, nil
}
