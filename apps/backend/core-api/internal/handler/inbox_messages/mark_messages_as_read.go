package inboxmessages

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/usecase/inbox"

	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MarkMessageAsReadRequest struct {
	MessageID uuid.UUID `json:"message_id" validate:"required"`
}

func (r *MarkMessageAsReadRequest) Parse(c *fiber.Ctx) error {
	if err := c.BodyParser(r); err != nil {
		return errors.Wrap(customerror.Parse(&customerror.ErrInternalServer, err), "failed to parse request body")
	}
	return nil
}

func (r *MarkMessageAsReadRequest) IsValid() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

type MarkMessageAsReadResponse struct {
	InboxMessage inbox.InboxMessagesViewModel `json:"inbox_message"`
}

// @Name MarkMessageAsRead
// @Tags Inbox Messages
// @Accept json
// @Produce json
// @Success 200 {object} inboxmessages.MarkMessageAsReadResponse
// @Router /inbox-messages/read [put]
// @Security ApiKeyAuth
// @Param message_id path string true "Message ID"
// @Param request body inboxmessages.MarkMessageAsReadRequest true "Request body"
func (h *Handler) MarkMessageAsRead(c *fiber.Ctx) error {
	user, err := h.AuthService.GetUserContext(c)
	if err != nil {
		return errors.Wrap(err, "failed to get user context")
	}
	if user == nil {
		return customerror.Parse(&customerror.ErrUnauthenticated, err)
	}

	var request MarkMessageAsReadRequest
	if err := request.Parse(c); err != nil {
		return errors.Wrap(err, "failed to parse request body")
	}
	if err := request.IsValid(); err != nil {
		return errors.Wrap(err, "invalid request body")
	}

	message, err := h.InboxUc.MarkMessageAsRead(c.Context(), *user, request.MessageID)
	if err != nil {
		return errors.Wrap(err, "failed to mark message as read")
	}

	viewModel, err := h.InboxUc.ToViewModel(c.Context(), *message)
	if err != nil {
		return errors.Wrap(err, "failed to convert inbox message to view model")
	}

	return c.Status(fiber.StatusOK).JSON(MarkMessageAsReadResponse{
		InboxMessage: *viewModel,
	})
}

type MarkAllMessagesAsReadResponse struct {
	InboxMessages []inbox.InboxMessagesViewModel `json:"inbox_messages"`
}

// @Name MarkAllMessagesAsRead
// @Tags Inbox Messages
// @Accept json
// @Produce json
// @Success 200 {object} inboxmessages.MarkAllMessagesAsReadResponse
// @Router /inbox-messages/read-all [put]
// @Security ApiKeyAuth
func (h *Handler) MarkAllMessagesAsRead(c *fiber.Ctx) error {
	user, err := h.AuthService.GetUserContext(c)
	if err != nil {
		return errors.Wrap(err, "failed to get user context")
	}
	if user == nil {
		return customerror.Parse(&customerror.ErrUnauthenticated, err)
	}

	messages, err := h.InboxUc.MarkAllMessagesAsRead(c.Context(), *user)
	if err != nil {
		return errors.Wrap(err, "failed to mark all messages as read")
	}

	viewModels := make([]inbox.InboxMessagesViewModel, len(messages))
	for i, message := range messages {
		viewModel, err := h.InboxUc.ToViewModel(c.Context(), *message)
		if err != nil {
			return errors.Wrap(err, "failed to convert inbox message to view model")
		}
		viewModels[i] = *viewModel
	}

	return c.Status(fiber.StatusOK).JSON(MarkAllMessagesAsReadResponse{
		InboxMessages: viewModels,
	})
}
