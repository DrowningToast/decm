package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SignEventCertificatesRequest struct {
	// For PIN-based (non-BYOK) users
	IssuerPin *string `json:"issuer_pin,omitempty"`
	// For wallet-based (BYOK) users
	IssuerSignMessage *string `json:"issuer_sign_message,omitempty"`
	IssuerSignature   *string `json:"issuer_signature,omitempty"`
}

type CertificateSignature struct {
	Certificate *entity.EventCertificate `json:"certificate,omitempty"`
	Signature   string                   `json:"signature"`
}

type SignEventCertificatesResponse struct {
	Certificates []CertificateSignature `json:"certificates"`
	TotalSigned  int                    `json:"total_signed"`
}

// @Summary Sign event certificates
// @Description Sign all event certificates for an event by issuer. PIN-based users provide issuer_pin; BYOK wallet users provide issuer_sign_message and issuer_signature.
// @ID sign-event-certificates
// @Tags Event Certificates
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body SignEventCertificatesRequest true "Sign certificates request"
// @Success 200 {object} SignEventCertificatesResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/certificates/sign [post]
func (h Handler) SignEventCertificates(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	requestBody := SignEventCertificatesRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	// Get current user from JWT
	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	response, err := h.EventUc.SignEventCertificates(ctx.UserContext(), eventID, event.SignEventCertificatesRequest{
		IssuerPin:         requestBody.IssuerPin,
		IssuerSignMessage: requestBody.IssuerSignMessage,
		IssuerSignature:   requestBody.IssuerSignature,
	}, currentUser)
	if err != nil {
		return err
	}

	// Convert usecase response to handler response
	var certificates []CertificateSignature
	for _, cert := range response.Certificates {
		certificates = append(certificates, CertificateSignature{
			Certificate: cert.Certificate,
			Signature:   cert.Signature,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(SignEventCertificatesResponse{
		Certificates: certificates,
		TotalSigned:  len(certificates),
	})
}

func (r *SignEventCertificatesRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func (r *SignEventCertificatesRequest) IsValid() error {
	hasPinAuth := r.IssuerPin != nil && *r.IssuerPin != ""
	hasWalletAuth := r.IssuerSignMessage != nil && r.IssuerSignature != nil
	if !hasPinAuth && !hasWalletAuth {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("either issuer_pin or (issuer_sign_message and issuer_signature) is required"))
	}
	return nil
}
