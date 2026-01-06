package system_status

import (
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.NewLogger()
	defer logger.Info("Mounted system status routes")

	systemStatusGroup := r.Group("/system-status")

	systemStatusGroup.Get("/viewmodel", h.GetSystemStatusViewModel)
	systemStatusGroup.Get("/latest", h.GetLatestSchedules)
	systemStatusGroup.Get("/period", h.GetSchedulesBetween)
	systemStatusGroup.Get("/closest-incoming", h.GetClosestIncomingSchedule)
	systemStatusGroup.Get("/planned-maintenance", h.GetPlannedMaintenanceSchedules)
}
