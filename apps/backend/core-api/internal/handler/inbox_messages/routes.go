package inboxmessages

import (
	"apps/backend/common/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted inbox messages routes")

	inboxMessagesGroup := r.Group("/inbox-messages")
	inboxMessagesGroup.Use(h.AuthenticationGuardMiddleware.Middleware).Get("/my", h.GetMyInboxMessages)
	inboxMessagesGroup.Use(h.AuthenticationGuardMiddleware.Middleware).Get("/:inbox_message_id", h.GetInboxMessage)
	inboxMessagesGroup.Use(h.AuthenticationGuardMiddleware.Middleware).Put("/read", h.MarkMessageAsRead)
	inboxMessagesGroup.Use(h.AuthenticationGuardMiddleware.Middleware).Put("/read-all", h.MarkAllMessagesAsRead)
}
