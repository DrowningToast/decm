package inboxmessages

import (
	"apps/backend/core-api/internal/handler/metrics"
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted inbox messages routes")

	inboxMessagesGroup := r.Group("/inbox-messages")
	inboxMessagesGroup.Use(h.AuthenticationGuardMiddleware.Middleware)
	inboxMessagesGroup.Get("/", metrics.MeasureHandlerDurationWrapper("inboxMessages.GetMyInboxMessages", h.GetMyInboxMessages))
	inboxMessagesGroup.Get("/:inbox_message_id", metrics.MeasureHandlerDurationWrapper("inboxMessages.GetInboxMessage", h.GetInboxMessage))
	inboxMessagesGroup.Put("/read", metrics.MeasureHandlerDurationWrapper("inboxMessages.MarkMessageAsRead", h.MarkMessageAsRead))
	inboxMessagesGroup.Put("/read-all", metrics.MeasureHandlerDurationWrapper("inboxMessages.MarkAllMessagesAsRead", h.MarkAllMessagesAsRead))
}
