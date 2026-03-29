package event

import (
	"apps/backend/common/customerror"
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

// GenerateCertificateImageForParticipant generates a certificate image for a specific participant
// It includes proper authorization checks to ensure the user owns the certificate
// Authorization is checked by credential ID, email (OAuth users), or wallet address (BYOK users)
func (uc *EventUsecase) GenerateCertificateImageForParticipant(
	ctx context.Context,
	certificateID uuid.UUID,
	currentUserID uuid.UUID,
	currentUserEmail *string,
	currentUserWalletAddress *string,
) ([]byte, error) {
	// 1. Get the certificate to verify ownership and get data
	certificate, err := uc.EventCertificateDataGateway.GetEventCertificateByID(ctx, certificateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get certificate")
	}

	// 2. Authorization: Check if the current user owns this certificate
	// Check by credential ID, email (OAuth), or wallet address (BYOK)
	var isAuthorized bool

	// Check by credential ID
	if certificate.ReceiverCredentialId != nil && *certificate.ReceiverCredentialId == currentUserID {
		isAuthorized = true
	}

	// Check by email
	if !isAuthorized && certificate.ReceiverEmail != nil && currentUserEmail != nil {
		if *certificate.ReceiverEmail == *currentUserEmail {
			isAuthorized = true
		}
	}

	// Check by wallet address (BYOK users)
	if !isAuthorized && certificate.ReceiverWalletAddress != nil && currentUserWalletAddress != nil {
		if *certificate.ReceiverWalletAddress == strings.ToLower(*currentUserWalletAddress) {
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
