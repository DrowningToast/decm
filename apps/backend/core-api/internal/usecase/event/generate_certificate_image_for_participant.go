package event

import (
	"apps/backend/common/customerror"
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
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

	// 3. Generate the certificate image using the main method
	// The GenerateCertificateImage method handles all the template processing internally
	pngBytes, err := uc.GenerateCertificateImage(ctx, certificateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate certificate image")
	}

	return pngBytes, nil
}
