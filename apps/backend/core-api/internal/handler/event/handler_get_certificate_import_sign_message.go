package event

import (
	"apps/backend/common/customerror"
	eventUsecase "apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetCertificateImportSignMessageRequest struct {
	EventID   uuid.UUID                          `json:"event_id"`
	Receivers []ImportCertificateReceiverRequest `json:"receivers" validate:"required,min=1"`
}

type GetCertificateImportSignMessageResponse struct {
	SignMessage             string `json:"sign_message"`
	EventCertificateAddress string `json:"event_certificate_address"`
}

// @Summary Get sign message for certificate import (BYOK wallet users)
// @Description Deploys the certificate contract if not yet deployed, computes receiver hashes, and returns the sign message JSON that the host must sign with their wallet before calling the import endpoint.
// @ID get-certificate-import-sign-message
// @Tags Event Certificates
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body GetCertificateImportSignMessageRequest true "Receivers to import"
// @Success 200 {object} GetCertificateImportSignMessageResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/certificates/import/sign-message [post]
func (h Handler) GetCertificateImportSignMessage(ctx *fiber.Ctx) error {
	eventIDStr := ctx.Params("event_id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	requestBody := GetCertificateImportSignMessageRequest{}
	if err := ctx.BodyParser(&requestBody); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	if len(requestBody.Receivers) == 0 {
		return customerror.Parse(&customerror.ErrInvalidArgument, fiber.NewError(fiber.StatusBadRequest, "at least one receiver is required"))
	}

	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	requests := make([]eventUsecase.ImportCertificateReceiversRequest, 0, len(requestBody.Receivers))
	for _, receiver := range requestBody.Receivers {
		requests = append(requests, eventUsecase.ImportCertificateReceiversRequest{
			Email:               receiver.Email,
			WalletAddress:       receiver.WalletAddress,
			FirstName:           receiver.FirstName,
			LastName:            receiver.LastName,
			AcademicInstitution: receiver.AcademicInstitution,
			CertificateTitle:    receiver.CertificateTitle,
			CertificateSubtitle: receiver.CertificateSubtitle,
		})
	}

	response, err := h.EventUc.GetCertificateImportSignMessage(ctx.UserContext(), eventID, requests, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
