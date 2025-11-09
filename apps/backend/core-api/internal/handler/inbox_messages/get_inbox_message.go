package inboxmessages

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/inbox"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetInboxMessageResponse struct {
	InboxMessage inbox.InboxMessagesViewModel `json:"inbox_message"`
}

type GetInboxMessageEventRegistrationInvitationResponse struct {
	InboxMessage inbox.InboxMessagesEventRegistrationInvitationViewModel `json:"inbox_message"`
}

// @Summary Get inbox message
// @Description Get inbox message
// @Name GetInboxMessage
// @Tags Inbox Messages
// @Accept json
// @Produce json
// @Success 200 {object} inboxmessages.GetInboxMessageResponse
// @Success 200 {object} inboxmessages.GetInboxMessageEventRegistrationInvitationResponse
// @Router /inbox-messages/{inbox_message_id} [get]
// @Security ApiKeyAuth
func (h *Handler) GetInboxMessage(c *fiber.Ctx) error {
	user, err := h.AuthService.GetUserContext(c)
	if err != nil {
		return errors.Wrap(err, "failed to get user context")
	}
	if user == nil {
		return errors.New("user not found")
	}

	inboxMessageId, err := uuid.Parse(c.Params("inbox_message_id"))
	if err != nil {
		return errors.Wrap(err, "failed to parse inbox message id")
	}

	result, err := h.InboxUc.GetInboxMessage(c.Context(), *user, inboxMessageId)
	if err != nil {
		return errors.Wrap(err, "failed to get inbox message")
	}
	if result == nil {
		return customerror.Parse(&customerror.ErrNotFound, errors.New("inbox message not found"))
	}

	inboxMessageViewModel, err := h.InboxUc.ToViewModel(c.Context(), result.InboxMessage)
	if err != nil {
		return errors.Wrap(err, "failed to convert inbox message to view model")
	}

	switch inboxMessageViewModel.MessageType {
	case entity.InboxMessageTypeGeneral:
		return c.Status(fiber.StatusOK).JSON(GetInboxMessageResponse{
			InboxMessage: *inboxMessageViewModel,
		})
	case entity.InboxMessageTypeEventRegistrationInvitation:
		if result.EventRegistrationInvitation == nil || result.Event == nil {
			return customerror.Parse(&customerror.ErrNotFound, errors.New("event registration invitation or event not found"))
		}
		eventRegistrationInvitationViewModel, err := h.InboxUc.ToWithEventRegistrationInvitationViewModel(c.Context(), result.InboxMessage, *result.EventRegistrationInvitation, *result.Event)
		if err != nil {
			return errors.Wrap(err, "failed to convert event registration invitation to view model")
		}
		return c.Status(fiber.StatusOK).JSON(GetInboxMessageEventRegistrationInvitationResponse{
			InboxMessage: *eventRegistrationInvitationViewModel,
		})
	default:
		return customerror.Parse(&customerror.ErrNotFound, errors.New("inbox message not found"))
	}
}
