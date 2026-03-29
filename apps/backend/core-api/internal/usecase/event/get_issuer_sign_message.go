package event

import (
	"apps/backend/common/customerror"
	"apps/backend/services/auth"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type GetIssuerSignMessageResponse struct {
	SignMessage string `json:"sign_message"`
}

// GetIssuerSignMessage returns the stored sign_message for the current BYOK issuer.
// This is used by BYOK wallet users who need to know what message to sign before calling SignEventCertificates.
func (uc *EventUsecase) GetIssuerSignMessage(ctx context.Context, eventID uuid.UUID, currentUser *auth.JwtClaims) (*GetIssuerSignMessageResponse, error) {
	// 1. Verify user is an issuer
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	if !credential.IsVerifiedIssuer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user is not a verified issuer"))
	}

	// Only BYOK users need this endpoint
	if credential.EncryptedPrivateKey != nil {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("only BYOK users need to fetch the sign message"))
	}

	// 2. Get certificate config for the event
	certificateConfig, err := uc.EventCertificateConfigDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// 3. Get signature records for this config
	signatures, err := uc.EventCertificateSignatureDataGateway.GetEventCertificateSignaturesByEventCertificateConfigID(ctx, certificateConfig.ID)
	if err != nil {
		return nil, err
	}

	// Find the signature record for this issuer
	for _, sig := range signatures {
		if sig.IssuerCredentialId == currentUser.UserId {
			if sig.SignMessage == nil {
				return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("sign message not found for this issuer"))
			}
			return &GetIssuerSignMessageResponse{
				SignMessage: *sig.SignMessage,
			}, nil
		}
	}

	return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("no signature record found for this issuer"))
}
