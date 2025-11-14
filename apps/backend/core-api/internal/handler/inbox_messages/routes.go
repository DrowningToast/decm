package inboxmessages

import (
	"apps/backend/common/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted inbox messages routes")

	inboxMessagesGroup := r.Group("/inbox-messages")
	inboxMessagesGroup.Use(h.AuthenticationGuardMiddleware.Middleware)
	inboxMessagesGroup.Get("/my", h.GetMyInboxMessages)
	inboxMessagesGroup.Get("/:inbox_message_id", h.GetInboxMessage)
	inboxMessagesGroup.Put("/read", h.MarkMessageAsRead)
	inboxMessagesGroup.Put("/read-all", h.MarkAllMessagesAsRead)
}
