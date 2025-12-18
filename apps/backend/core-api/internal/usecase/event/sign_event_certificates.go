package event

import (
	"context"
	"fmt"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
)

type SignEventCertificatesRequest struct {
	IssuerPin string `json:"issuer_pin"`
}

type CertificateSignature struct {
	Certificate *entity.EventCertificate `json:"certificate,omitempty"`
	Signature   string                   `json:"signature"`
}

type SignEventCertificatesResponse struct {
	Certificates []CertificateSignature `json:"certificates"`
}

func (uc *EventUsecase) SignEventCertificates(ctx context.Context, eventID uuid.UUID, request SignEventCertificatesRequest, currentUser *auth.JwtClaims) (*SignEventCertificatesResponse, error) {
	// 1. Check if current user is authorized
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	if !credential.IsVerifiedIssuer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user is not a verified issuer"))
	}

	// 2. Check if event exists
	event, err := uc.EventDataGateway.GetEventById(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event not found"))
	}

	// 3. Get all certificates for the event
	certificates, err := uc.EventCertificateDataGateway.GetEventCertificatesByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if len(certificates) == 0 {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("no certificates found for this event"))
	}

	// 4. Get eventContract from eventContracts table using eventID
	eventContract, err := uc.EventContractDataGateway.GetEventContractByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if eventContract == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event contract not found"))
	}

	// 5. Decrypt private key
	privateKey, _, err := cyptoutils.DecryptPrivateKey(*credential.EncryptedPrivateKey, request.IssuerPin)
	if err != nil {
		return nil, err
	}

	// 6. Get certificate config (same for all certificates in the event)
	certificateConfig, err := uc.EventCertificateConfigDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// 7. Process each certificate
	var signedCertificates []CertificateSignature

	for _, certificate := range certificates {
		// Get certificate signatures for this certificate config
		signatures, err := uc.EventCertificateSignatureDataGateway.GetEventCertificateSignaturesByEventCertificateConfigID(ctx, certificateConfig.ID)
		if err != nil {
			return nil, err
		}

		// Find the signature for this issuer
		var targetSignature *entity.EventCertificateSignature
		for _, sig := range signatures {
			if sig.IssuerCredentialId == currentUser.UserId {
				targetSignature = sig
				break
			}
		}

		// Skip if no signature found for this issuer
		if targetSignature == nil {
			continue
		}

		// Sign the message
		if targetSignature.SignMessage == nil {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("sign message not found for certificate %s", certificate.Id))
		}

		// Use HashEthereumMessage to match contract's recoverSigner which applies Ethereum prefix
		signMessageDigest := cyptoutils.HashEthereumMessage(*targetSignature.SignMessage)
		signature, err := cyptoutils.Sign(signMessageDigest.Bytes(), privateKey)
		if err != nil {
			return nil, err
		}

		// Convert signature to string for database storage
		signatureStr := hexutil.Encode(signature)

		// Update certificate signature in database
		_, err = uc.EventCertificateSignatureDataGateway.UpdateEventCertificateIssuerSignature(ctx, targetSignature.Id, &signatureStr)
		if err != nil {
			return nil, err
		}

		// Add to response
		signedCertificates = append(signedCertificates, CertificateSignature{
			Certificate: certificate,
			Signature:   fmt.Sprintf("0x%x", signature),
		})
	}

	if len(signedCertificates) == 0 {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("no certificates found with signatures for this issuer"))
	}

	// Update issuer signing status to signed (1)
	err = uc.EventIssuerDataGateway.UpdateEventIssuerSigningStatus(ctx, eventID, currentUser.UserId, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to update issuer signing status: %w", err)
	}

	return &SignEventCertificatesResponse{
		Certificates: signedCertificates,
	}, nil
}
