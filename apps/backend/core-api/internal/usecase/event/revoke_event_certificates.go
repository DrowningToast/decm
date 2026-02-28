package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"fmt"
	"time"

	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"

	"github.com/google/uuid"
)

type RevokeEventCertificatesRequest struct {
	CertificateIDs []uuid.UUID `json:"certificate_ids"`
}

type RevokeEventCertificatesResponse struct {
	RevokedCertificates []*entity.EventCertificate `json:"revoked_certificates"`
}

func (uc *EventUsecase) RevokeEventCertificates(ctx context.Context, eventID uuid.UUID, request RevokeEventCertificatesRequest, currentUser *auth.JwtClaims) (*RevokeEventCertificatesResponse, error) {
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

	if event == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event not found"))
	}

	// 3. Check if certificate configuration is published
	certificateConfig, err := uc.EventCertificateConfigDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if certificateConfig.IsPublished {
		return nil, customerror.Parse(&customerror.ErrForbidden, fmt.Errorf("cannot revoke certificates after certificate configuration has been published"))
	}

	// 3. Revoke each certificate in the database
	revokedCertificates := make([]*entity.EventCertificate, 0, len(request.CertificateIDs))
	now := time.Now()

	for _, certificateID := range request.CertificateIDs {
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

	// 4. Reset all issuers' signing status since certificate configuration needs re-approval
	err = uc.EventIssuerDataGateway.ResetAllEventIssuersSigningStatus(ctx, eventID)
	if err != nil {
		uc.logger.Error("failed to reset issuers signing status", "error", err, "event_id", eventID)
		// Log error but don't fail the revocation operation
	}

	return &RevokeEventCertificatesResponse{
		RevokedCertificates: revokedCertificates,
	}, nil
}
