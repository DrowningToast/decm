package event

import (
	"context"
	"fmt"
	"time"

	"apps/backend/common/customerror"
	eventdatagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
)

type RevokeAllEventCertificatesResponse struct {
	RevokedCertificates []*entity.EventCertificate `json:"revoked_certificates"`
}

func (uc *EventUsecase) RevokeAllEventCertificates(ctx context.Context, eventID uuid.UUID, currentUser *auth.JwtClaims) (*RevokeAllEventCertificatesResponse, error) {
	// 1. Check if current user is authorized
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	if !credential.IsVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user is not a verified organizer"))
	}

	// 2. Check if event exists
	event, err := uc.EventDataGateway.GetEventById(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if credential.Id != event.OwnerCredentialId {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user is not owner of the event"))
	}

	if event == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event not found"))
	}

	// 3. Get all certificate IDs for the event
	certificateIDs, err := uc.EventCertificateDataGateway.GetAllEventCertificateIDsByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// If no certificates found, return empty response
	if len(certificateIDs) == 0 {
		return &RevokeAllEventCertificatesResponse{
			RevokedCertificates: []*entity.EventCertificate{},
		}, nil
	}

	// 4. Revoke all certificates in the database
	revokedCertificates := make([]*entity.EventCertificate, 0, len(certificateIDs))
	now := time.Now()

	for _, certificateID := range certificateIDs {
		// Get certificate details
		certificate, err := uc.EventCertificateDataGateway.GetEventCertificateByID(ctx, certificateID)
		if err != nil {
			return nil, err
		}

		if certificate == nil {
			return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("certificate not found"))
		}

		// Update certificate in database to mark as revoked
		updatedCertificate, err := uc.EventCertificateDataGateway.UpdateEventCertificate(ctx, certificateID, eventdatagateway.UpdateEventCertificateParameters{
			RevokedAt:               &now,
			ReceiverCredentialID:    certificate.ReceiverCredentialId,
			ReceiverEmail:           certificate.ReceiverEmail,
			Name:                    certificate.Name,
			AcademicInstitution:     certificate.AcademicInstitution,
			CertificateTitle:        certificate.CertificateTitle,
			CertificateSubtitle:     certificate.CertificateSubtitle,
			EventContractAddress:    &certificate.EventContractAddress,
			EventCertificateAddress: certificate.EventCertificateAddress,
			CertificateTokenID:      certificate.CertificateTokenId,
		})
		if err != nil {
			return nil, err
		}

		revokedCertificates = append(revokedCertificates, updatedCertificate)
	}

	return &RevokeAllEventCertificatesResponse{
		RevokedCertificates: revokedCertificates,
	}, nil
}
