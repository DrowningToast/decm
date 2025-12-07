package inboxmessages

import (
	"apps/backend/core-api/internal/usecase/inbox"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

type GetInboxMessagesResponse struct {
	InboxMessages []inbox.InboxMessagesViewModel `json:"inbox_messages"`
}

// @Summary Get my inbox messages
// @Description Get my inbox messages
// @Name GetMyInboxMessages
// @Tags Inbox Messages
// @Accept json
// @Produce json
// @Success 200 {object} GetInboxMessagesResponse
// @Router /api/v1/inbox-messages [get]
// @Security ApiKeyAuth
func (h *Handler) GetMyInboxMessages(c *fiber.Ctx) error {
	user, err := h.AuthService.GetUserContext(c)
	if err != nil {
		return errors.Wrap(err, "failed to get user context")
	}
	if user == nil {
		return errors.New("user not found")
	}

	messages, err := h.InboxUc.GetMyInboxMessages(c.Context(), *user)
	if err != nil {
		return errors.Wrap(err, "failed to get my inbox messages")
	}
	inboxMessagesViewModels := make([]inbox.InboxMessagesViewModel, len(messages))
	for i, message := range messages {
		viewModel, err := h.InboxUc.ToViewModel(c.Context(), *message)
		if err != nil {
			return errors.Wrap(err, "failed to convert inbox message to view model")
		}
		inboxMessagesViewModels[i] = *viewModel
	}

	return c.Status(fiber.StatusOK).JSON(GetInboxMessagesResponse{
		InboxMessages: inboxMessagesViewModels,
	})
}
