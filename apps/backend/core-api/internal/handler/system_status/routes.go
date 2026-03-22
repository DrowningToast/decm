package system_status

import (
	"apps/backend/core-api/internal/handler/metrics"
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted system status routes")

	systemStatusGroup := r.Group("/system-status")

	systemStatusGroup.Get("/viewmodel", metrics.MeasureHandlerDurationWrapper("systemStatus.GetSystemStatusViewModel", h.GetSystemStatusViewModel))
	systemStatusGroup.Get("/latest", metrics.MeasureHandlerDurationWrapper("systemStatus.GetLatestSchedules", h.GetLatestSchedules))
	systemStatusGroup.Get("/period", metrics.MeasureHandlerDurationWrapper("systemStatus.GetSchedulesBetween", h.GetSchedulesBetween))
	systemStatusGroup.Get("/closest-incoming", metrics.MeasureHandlerDurationWrapper("systemStatus.GetClosestIncomingSchedule", h.GetClosestIncomingSchedule))
	systemStatusGroup.Get("/planned-maintenance", metrics.MeasureHandlerDurationWrapper("systemStatus.GetPlannedMaintenanceSchedules", h.GetPlannedMaintenanceSchedules))
}
