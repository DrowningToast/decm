package system_status

import (
	"apps/backend/core-api/internal/entity"

	"github.com/gofiber/fiber/v2"
)

type GetClosestIncomingScheduleResponse struct {
	Schedule *entity.SystemStatusSchedule `json:"schedule,omitempty" extensions:"x-nullable"`
}

// GetClosestIncomingSchedule godoc
// @Summary Get closest incoming status update
// @Description Get the closest incoming status update scheduled for the future
// @ID get-closest-incoming-schedule
// @Tags SystemStatus
// @Produce json
// @Success 200 {object} GetClosestIncomingScheduleResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/system-status/closest-incoming [get]
func (h *Handler) GetClosestIncomingSchedule(ctx *fiber.Ctx) error {
	schedule, err := h.systemStatusUc.GetClosestIncomingSchedule(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(&GetClosestIncomingScheduleResponse{
		Schedule: schedule,
	})
}
