package inboxmessages

import (
	"apps/backend/core-api/internal/entity"
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

	// Enrich messages with certificate-specific data when applicable
	inboxMessagesViewModels := make([]inbox.InboxMessagesViewModel, 0, len(messages))
	for _, message := range messages {
		viewModel, err := h.InboxUc.ToViewModel(c.Context(), *message)
		if err != nil {
			return errors.Wrap(err, "failed to convert inbox message to view model")
		}

		// For certificate invitations, enrich with certificate data
		if message.MessageType == entity.InboxMessageTypeEventCertificateInvitation {
			enrichedViewModel, err := h.InboxUc.ToEnrichedCertificateViewModel(c.Context(), *message, *user)
			if err != nil {
				// If enrichment fails, fall back to basic view model
				inboxMessagesViewModels = append(inboxMessagesViewModels, *viewModel)
				continue
			}
			inboxMessagesViewModels = append(inboxMessagesViewModels, *enrichedViewModel)
		} else {
			inboxMessagesViewModels = append(inboxMessagesViewModels, *viewModel)
		}
	}

	return c.Status(fiber.StatusOK).JSON(GetInboxMessagesResponse{
		InboxMessages: inboxMessagesViewModels,
	})
}
