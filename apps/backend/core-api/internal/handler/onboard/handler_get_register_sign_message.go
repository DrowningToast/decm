package onboard

import "github.com/gofiber/fiber/v2"

// @name getRegisterSignMessageResponse
// @description Response for the client to sign to register
type getRegisterSignMessageResponse struct {
	Message string `json:"message"`
}

// @Summary Get preset message for the client to sign to register
// @Description Retrieve preset message for the client to sign to register
// @ID get-register-sign-message
// @Tag Onboard
// @Produce json
// @Success 200 {object} getRegisterSignMessageResponse
// @Router /api/v1/onboard/sign-message [get]
func (h Handler) GetRegisterSignMessage(ctx *fiber.Ctx) error {
	message := h.OnboardUc.GetRegisterSignMessage()

	response := getRegisterSignMessageResponse{
		Message: message,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
