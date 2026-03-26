package certificate

import (
	"apps/backend/common/validatorutils"
	"encoding/hex"
	"strings"

	customerror "apps/backend/common/customerror"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetClaimCertificateSignMessageResponse struct {
	SignMessage string `json:"sign_message"`
}

// @Summary Get claim certificate sign message
// @Description Get sign message for claiming a certificate (requires authentication)
// @ID get-claim-certificate-sign-message
// @Tags Certificates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param certificate_id path string true "Certificate ID"
// @Success 200 {object} GetClaimCertificateSignMessageResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificates/claim/{certificate_id}/sign-message [get]
func (h *Handler) GetClaimCertificateSignMessage(ctx *fiber.Ctx) error {
	// 1. Validate and parse certificate ID
	certificateIdStr := ctx.Params("certificate_id")
	if certificateIdStr == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate_id is required"))
	}
	certificateId, err := uuid.Parse(certificateIdStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "invalid certificate_id format"))
	}

	// 2. Get authenticated user context
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	// 3. Get certificate and validate basic requirements
	certificate, err := h.EventUc.EventCertificateDataGateway.GetEventCertificateByID(ctx.UserContext(), certificateId)
	if err != nil {
		return errors.Wrap(err, "failed to get certificate by id")
	}
	if certificate == nil {
		return customerror.Parse(&customerror.ErrNotFound, errors.New("certificate not found"))
	}

	// 4. Validate certificate is published
	if certificate.EventCertificateAddress == nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate is not published yet"))
	}

	// 5. Validate certificate is not already claimed
	if certificate.CertificateTokenId != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate has already been claimed"))
	}

	// 6. Validate certificate is not revoked
	if certificate.RevokedAt != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate has been revoked"))
	}

	// 7. Basic eligibility check - user must match credential ID, email, or wallet_address (if specified)
	if certificate.ReceiverCredentialId != nil {
		if *certificate.ReceiverCredentialId != currentUser.UserId {
			return customerror.Parse(&customerror.ErrForbidden, errors.New("you are not eligible to claim this certificate"))
		}
	} else if certificate.ReceiverEmail != nil {
		if currentUser.Email == nil || *currentUser.Email != *certificate.ReceiverEmail {
			return customerror.Parse(&customerror.ErrForbidden, errors.New("you are not eligible to claim this certificate"))
		}
	} else if certificate.ReceiverWalletAddress != nil {
		if !strings.EqualFold(currentUser.WalletAddress, *certificate.ReceiverWalletAddress) {
			return customerror.Parse(&customerror.ErrForbidden, errors.New("you are not eligible to claim this certificate"))
		}
	}
	// If all are nil, it's an open certificate - anyone can get sign message

	signMessage, _, err := h.EventUc.GetClaimCertificateSignMessage(ctx.UserContext(), common.HexToAddress(currentUser.WalletAddress), *currentUser, common.HexToAddress(*certificate.EventCertificateAddress), nil)
	if err != nil {
		return errors.Wrap(err, "failed to get claim certificate sign message")
	}

	return ctx.Status(fiber.StatusOK).JSON(GetClaimCertificateSignMessageResponse{
		SignMessage: *signMessage,
	})
}

type ClaimCertificateBody struct {
	AccountPassword     *string `json:"account_password,omitempty"`
	CertificatePassword *string `json:"certificate_password,omitempty"`
	Signature           *string `json:"signature,omitempty"`
	SignMessage         *string `json:"sign_message,omitempty"`
}

func (r *ClaimCertificateBody) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

func (r *ClaimCertificateBody) Parse(ctx *fiber.Ctx) error {
	return ctx.BodyParser(r)
}

type ClaimCertificateResponse struct {
	Certificate   interface{} `json:"certificate"`
	UserSignature interface{} `json:"user_signature"`
	Status        string      `json:"status"`
	Message       string      `json:"message"`
}

func NewClaimCertificateResponse(certificate interface{}, userSignature interface{}) ClaimCertificateResponse {
	return ClaimCertificateResponse{
		Certificate:   certificate,
		UserSignature: userSignature,
		Status:        "queued",
		Message:       "Certificate claim has been queued for processing. The certificate will be minted shortly.",
	}
}

// @Summary Claim certificate
// @Description Queue a certificate claim using account password or wallet signature (requires authentication). The certificate will be minted asynchronously.
// @ID claim-certificate
// @Tags Certificates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param certificate_id path string true "Certificate ID"
// @Param claimCertificateBody body ClaimCertificateBody true "Claim certificate body (provide either account_password or signature+sign_message)"
// @Success 202 {object} ClaimCertificateResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificates/claim/{certificate_id} [post]
func (h *Handler) ClaimCertificate(ctx *fiber.Ctx) error {
	// 1. Validate and parse certificate ID
	certificateIdStr := ctx.Params("certificate_id")
	if certificateIdStr == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate_id is required"))
	}
	certificateId, err := uuid.Parse(certificateIdStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "invalid certificate_id format"))
	}

	// 2. Parse and validate request body
	var req ClaimCertificateBody
	if err := req.Parse(ctx); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "failed to parse request body"))
	}
	if err := req.IsValid(); err != nil {
		return err
	}

	// 3. Validate that either password OR signature is provided (not both)
	hasPassword := req.AccountPassword != nil
	hasSignature := req.Signature != nil && req.SignMessage != nil

	if !hasPassword && !hasSignature {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("either account_password or (signature + sign_message) is required"))
	}

	if hasPassword && hasSignature {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("provide either account_password or signature, not both"))
	}

	// 4. Get authenticated user context
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	// 5. Route to appropriate claiming flow
	if hasSignature {
		// Wallet Extension Flow (MetaMask, WalletConnect, etc.)

		// Decode signature from hex string
		signatureHex := strings.TrimPrefix(*req.Signature, "0x")
		signature, err := hex.DecodeString(signatureHex)
		if err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "invalid signature format - must be valid hex"))
		}

		// Validate signature length (65 bytes for ECDSA signature)
		if len(signature) != 65 {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid signature length - expected 65 bytes"))
		}

		// Queue certificate claim with signature
		certificate, userSignature, err := h.EventUc.ClaimCertificateWithSignature(
			ctx.UserContext(),
			currentUser,
			certificateId,
			signature,
			*req.SignMessage,
		)
		if err != nil {
			return err
		}

		response := NewClaimCertificateResponse(certificate, userSignature)
		return ctx.Status(fiber.StatusAccepted).JSON(response)
	} else {
		// PIN/Password Flow (Mobile/Web App)

		// Queue certificate claim with password
		certificate, userSignature, err := h.EventUc.ClaimCertificateWithPin(
			ctx.UserContext(),
			currentUser,
			certificateId,
			*req.AccountPassword,
		)
		if err != nil {
			return err
		}

		response := NewClaimCertificateResponse(certificate, userSignature)
		return ctx.Status(fiber.StatusAccepted).JSON(response)
	}
}
