package event

import (
	"encoding/hex"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	event_uc "apps/backend/core-api/internal/usecase/event"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetClaimCertificateSignMessageResponse struct {
	SignMessage string `json:"sign_message"`
}

// @Summary Get claim certificate sign message
// @Description Get sign message for claiming a certificate
// @ID get-claim-certificate-sign-message
// @Tags Certificates
// @Accept json
// @Produce json
// @Param certificate_id path string true "Certificate ID"
// @Success 200 {object} GetClaimCertificateSignMessageResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificates/claim/{certificate_id}/sign-message [get]
func (h *Handler) GetClaimCertificateSignMessage(ctx *fiber.Ctx) error {
	certificateIdStr := ctx.Params("certificate_id")
	if certificateIdStr == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate_id is required"))
	}
	certificateId, err := uuid.Parse(certificateIdStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Get certificate to extract contract address
	certificate, err := h.EventUc.EventCertificateDataGateway.GetEventCertificateByID(ctx.UserContext(), certificateId)
	if err != nil {
		return errors.Wrap(err, "failed to get certificate by id")
	}
	if certificate == nil {
		return customerror.Parse(&customerror.ErrNotFound, errors.New("certificate not found"))
	}
	if certificate.EventCertificateAddress == nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate is not published yet"))
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	client, err := cyptoutils.GetEthereumClient()
	if err != nil {
		return errors.Wrap(err, "failed to get ethereum client")
	}

	signMessage, _, err := h.EventUc.GetClaimCertificateSignMessage(ctx.UserContext(), client, common.HexToAddress(currentUser.WalletAddress), *currentUser, common.HexToAddress(*certificate.EventCertificateAddress), nil)
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

// @Summary Claim certificate
// @Description Claim a certificate by signing with wallet
// @ID claim-certificate
// @Tags Certificates
// @Accept json
// @Produce json
// @Param certificate_id path string true "Certificate ID"
// @Param claimCertificateBody body ClaimCertificateBody true "Claim certificate body"
// @Success 200 {object} entity.EventCertificate
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificates/claim/{certificate_id} [post]
func (h *Handler) ClaimCertificate(ctx *fiber.Ctx) error {
	certificateIdStr := ctx.Params("certificate_id")
	if certificateIdStr == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate_id is required"))
	}
	certificateId, err := uuid.Parse(certificateIdStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req ClaimCertificateBody
	if err := req.Parse(ctx); err != nil {
		return err
	}
	if err := req.IsValid(); err != nil {
		return err
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	client, err := cyptoutils.GetEthereumClient()
	if err != nil {
		return errors.Wrap(err, "failed to get ethereum client")
	}

	proof := event_uc.CheckClaimEligibilityParams{
		CertificatePassword: req.CertificatePassword,
	}

	if req.Signature != nil {
		signature, err := hex.DecodeString(*req.Signature)
		if err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, err)
		}
		if req.SignMessage == nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("original sign message is required"))
		}
		certificate, err := h.EventUc.ClaimCertificateWithSignature(ctx.UserContext(), client, currentUser, certificateId, proof, signature, *req.SignMessage)
		if err != nil {
			return errors.Wrap(err, "failed to claim certificate with signature")
		}
		return ctx.Status(fiber.StatusOK).JSON(certificate)
	} else {
		if req.AccountPassword == nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("account password is required when signature is not provided"))
		}
		certificate, err := h.EventUc.ClaimCertificateWithPin(ctx.UserContext(), client, currentUser, certificateId, proof, *req.AccountPassword)
		if err != nil {
			return errors.Wrap(err, "failed to claim certificate with pin")
		}
		return ctx.Status(fiber.StatusOK).JSON(certificate)
	}
}
