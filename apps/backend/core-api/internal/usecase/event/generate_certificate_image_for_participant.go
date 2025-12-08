package event

import (
	"context"
	"fmt"

	"apps/backend/common/customerror"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// GenerateCertificateImageForParticipant generates a certificate image for a specific participant
// It includes proper authorization checks to ensure the user owns the certificate
// Authorization is checked by either credential ID or email (for Google OAuth users)
func (uc *EventUsecase) GenerateCertificateImageForParticipant(
	ctx context.Context,
	certificateID uuid.UUID,
	currentUserID uuid.UUID,
	currentUserEmail *string,
) ([]byte, error) {
	// 1. Get the certificate to verify ownership and get data
	certificate, err := uc.EventCertificateDataGateway.GetEventCertificateByID(ctx, certificateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get certificate")
	}

	// 2. Authorization: Check if the current user owns this certificate
	// Check either by credential ID or by email (for Google OAuth users)
	var isAuthorized bool

	// Check by credential ID
	if certificate.ReceiverCredentialId != nil && *certificate.ReceiverCredentialId == currentUserID {
		isAuthorized = true
	}

	// If not authorized by credential ID, check by email
	if !isAuthorized && certificate.ReceiverEmail != nil && currentUserEmail != nil {
		if *certificate.ReceiverEmail == *currentUserEmail {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		return nil, customerror.Parse(&customerror.ErrForbidden,
			fmt.Errorf("user does not own this certificate"))
	}

	// 3. Get the event details to include event name
	event, err := uc.EventDataGateway.GetEventById(ctx, certificate.EventId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get event")
	}

	// 4. Get the certificate config to know which template to use
	certConfig, err := uc.EventCertificateConfigDg.GetEventCertificateConfigByEventID(ctx, certificate.EventId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get certificate config")
	}

	// 5. Prepare template variables from certificate data
	templateVars := CertificateTemplateVariables{
		Name:                *certificate.Name,
		EventName:           event.Title,
		AcademicInstitution: certificate.AcademicInstitution,
		CertificateTitle:    certificate.CertificateTitle,
		CertificateSubtitle: certificate.CertificateSubtitle,
	}

	// 6. Generate the certificate image with positions from config
	// Helper to convert pgtype.Float8 to *float64
	pgFloat8ToPtr := func(f pgtype.Float8) *float64 {
		if f.Valid {
			val := f.Float64
			return &val
		}
		return nil
	}

	// Convert float64 to *float64
	namePosX := certConfig.NamePosX
	namePosY := certConfig.NamePosY
	eventNamePosX := certConfig.EventNamePosX
	eventNamePosY := certConfig.EventNamePosY

	params := GenerateCertificateImageParams{
		TemplateVariables:       templateVars,
		NamePosX:                &namePosX,
		NamePosY:                &namePosY,
		EventNamePosX:           &eventNamePosX,
		EventNamePosY:           &eventNamePosY,
		CertificateTitlePosX:    pgFloat8ToPtr(certConfig.CertificateTitlePosX),
		CertificateTitlePosY:    pgFloat8ToPtr(certConfig.CertificateTitlePosY),
		CertificateSubtitlePosX: pgFloat8ToPtr(certConfig.CertificateSubtitlePosX),
		CertificateSubtitlePosY: pgFloat8ToPtr(certConfig.CertificateSubtitlePosY),
		AcademicInstitutionPosX: pgFloat8ToPtr(certConfig.AcademicInstitutionPosX),
		AcademicInstitutionPosY: pgFloat8ToPtr(certConfig.AcademicInstitutionPosY),
	}

	pngBytes, err := uc.GenerateCertificateImage(ctx, certConfig.ID, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate certificate image")
	}

	return pngBytes, nil
}
