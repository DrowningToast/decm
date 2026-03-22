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
	InboxMessageType entity.InboxMessageType      `json:"inbox_message_type"`
	InboxMessage     inbox.InboxMessagesViewModel `json:"inbox_message"`

	EventRegistrationInvitation *inbox.InboxMessagesEventRegistrationInvitationViewModel `json:"event_registration_invitation,omitempty"`
	EventCertificate            *inbox.InboxMessageCertificateInvitationViewModel        `json:"event_certificate,omitempty"`
}

// @Summary Get inbox message
// @Description Get inbox message
// @Name GetInboxMessage
// @Tags Inbox Messages
// @Accept json
// @Produce json
// @Success 200 {object} inboxmessages.GetInboxMessageResponse
// @Router /api/v1/inbox-messages/{inbox_message_id} [get]
// @Security ApiKeyAuth
func (h *Handler) GetInboxMessage(c *fiber.Ctx) error {
	user, err := h.AuthService.GetUserContext(c)
	if err != nil {
		return errors.Wrap(err, "failed to get user context")
	}
	if user == nil {
		return customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user not authenticated"))
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
			InboxMessageType: entity.InboxMessageTypeGeneral,
			InboxMessage:     *inboxMessageViewModel,
		})
	case entity.InboxMessageTypeEventRegistrationInvitation:
		if result.EventRegistrationInvitation == nil || result.Event == nil {
			return customerror.Parse(&customerror.ErrNotFound, errors.New("event registration invitation or event not found"))
		}
		eventRegistrationInvitationViewModel, err := h.InboxUc.ToWithEventRegistrationInvitationViewModel(c.Context(), result.InboxMessage, *result.EventRegistrationInvitation, *result.Event, result.EventAttendee, result.UserSignature)
		if err != nil {
			return errors.Wrap(err, "failed to convert event registration invitation to view model")
		}
		return c.Status(fiber.StatusOK).JSON(GetInboxMessageResponse{
			InboxMessageType:            entity.InboxMessageTypeEventRegistrationInvitation,
			InboxMessage:                *inboxMessageViewModel,
			EventRegistrationInvitation: eventRegistrationInvitationViewModel,
		})
	case entity.InboxMessageTypeEventCertificateInvitation:
		eventCertificateViewModel, err := h.InboxUc.ToWithCertificateInvitationViewModel(c.Context(), result.InboxMessage, *result.EventCertificate, *user)
		if err != nil {
			return errors.Wrap(err, "failed to convert event certificate to view model")
		}
		return c.Status(fiber.StatusOK).JSON(GetInboxMessageResponse{
			InboxMessageType: entity.InboxMessageTypeEventCertificateInvitation,
			InboxMessage:     *inboxMessageViewModel,
			EventCertificate: eventCertificateViewModel,
		})
	default:
		return customerror.Parse(&customerror.ErrNotFound, errors.New("inbox message not found"))
	}
}
