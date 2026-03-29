package event

import (
	"apps/backend/common/customerror"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetIssuerSignMessageResponse struct {
	SignMessage string `json:"sign_message"`
}

// @Summary Get sign message for BYOK issuer certificate signing
// @Description Returns the stored sign_message that a BYOK wallet issuer must sign before calling the sign certificates endpoint.
// @ID get-issuer-sign-message
// @Tags Event Certificates
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} GetIssuerSignMessageResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/certificates/sign/message [get]
func (h Handler) GetIssuerSignMessage(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	response, err := h.EventUc.GetIssuerSignMessage(ctx.UserContext(), eventID, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(GetIssuerSignMessageResponse{
		SignMessage: response.SignMessage,
	})
}
