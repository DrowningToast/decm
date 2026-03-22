package system_status

import (
	"apps/backend/core-api/internal/entity"

	"github.com/gofiber/fiber/v2"
)

type GetPlannedMaintenanceSchedulesResponse struct {
	Schedules []*entity.SystemStatusSchedule `json:"schedules"`
}

// GetPlannedMaintenanceSchedules godoc
// @Summary Get planned maintenance schedules
// @Description Get all upcoming planned maintenance schedules
// @ID get-planned-maintenance-schedules
// @Tags SystemStatus
// @Produce json
// @Success 200 {object} GetPlannedMaintenanceSchedulesResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/system-status/planned-maintenance [get]
func (h *Handler) GetPlannedMaintenanceSchedules(ctx *fiber.Ctx) error {
	schedules, err := h.systemStatusUc.GetPlannedMaintenanceSchedules(ctx.UserContext())
	if err != nil {
		return err
	}

	if schedules == nil {
		schedules = []*entity.SystemStatusSchedule{}
	}

	return ctx.JSON(&GetPlannedMaintenanceSchedulesResponse{
		Schedules: schedules,
	})
}
