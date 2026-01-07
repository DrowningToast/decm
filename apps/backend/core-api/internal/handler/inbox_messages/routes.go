package inboxmessages

import (
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted inbox messages routes")

	inboxMessagesGroup := r.Group("/inbox-messages")
	inboxMessagesGroup.Use(h.AuthenticationGuardMiddleware.Middleware)
	inboxMessagesGroup.Get("/", h.GetMyInboxMessages)
	inboxMessagesGroup.Get("/:inbox_message_id", h.GetInboxMessage)
	inboxMessagesGroup.Put("/read", h.MarkMessageAsRead)
	inboxMessagesGroup.Put("/read-all", h.MarkAllMessagesAsRead)
}
