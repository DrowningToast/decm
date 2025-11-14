package eventconfig

import (
	"apps/backend/common/customerror"
	"apps/backend/common/validatorutils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CheckEventPasswordBody struct {
	Password string `json:"password" validate:"required"`
}

func (r *CheckEventPasswordBody) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

type CheckEventPasswordResponse struct {
	IsValid bool `json:"is_valid"`
}

func (r *CheckEventPasswordBody) Parse(ctx *fiber.Ctx) error {
	return ctx.BodyParser(r)
}

// @Summary Check event password
// @ID check-event-password
// @Description Check if the password is correct for an event
// @Tags Event Config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body CheckEventPasswordBody true "Check event password request"
// @Success 200 {object} CheckEventPasswordResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/config/password-check [post]
func (h *Handler) CheckEventPassword(ctx *fiber.Ctx) error {
	eventId, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req CheckEventPasswordBody
	if err := req.Parse(ctx); err != nil {
		return err
	}
	if err := req.IsValid(); err != nil {
		return err
	}

	isValid, err := h.EventConfigUc.CheckEventPassword(ctx.UserContext(), eventId, req.Password)
	if err != nil {
		return err
	}

	return ctx.JSON(CheckEventPasswordResponse{
		IsValid: isValid,
	})
}
