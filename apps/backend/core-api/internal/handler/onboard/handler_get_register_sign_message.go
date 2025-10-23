package onboard

import "github.com/gofiber/fiber/v2"

// @name getSignMessageResponse
// @description Response for the client to sign to register
type getSignMessageResponse struct {
	Message string `json:"message" validate:"required"`
}

// @Summary Get preset message for the client to sign to register
// @Description Retrieve preset message for the client to sign to register
// @ID get-sign-message
// @Tags Onboard
// @Produce json
// @Success 200 {object} getSignMessageResponse
// @Router /api/v1/onboard/sign-message [get]
func (h Handler) GetSignMessage(ctx *fiber.Ctx) error {
	message := h.OnboardUc.GetRegisterSignMessage()

	response := getSignMessageResponse{
		Message: message,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
