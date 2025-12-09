package event

import (
	"apps/backend/common/customerror"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PublishEventCertificatesResponse struct {
	EventID              uuid.UUID `json:"event_id"`
	PublishedCount       int       `json:"published_count"`
	IsPublished          bool      `json:"is_published"`
	InboxMessagesCreated int       `json:"inbox_messages_created"`
}

// @Summary Publish event certificates
// @Description Publish certificates for an event. This will create inbox messages for all certificate receivers. All issuers must have signed before publishing.
// @ID publish-event-certificates
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} PublishEventCertificatesResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Security Bearer
// @Router /api/v1/events/{event_id}/certificates/publish [post]
func (h Handler) PublishEventCertificates(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Get current user from JWT
	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	response, err := h.EventUc.PublishEventCertificates(ctx.UserContext(), eventID, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(PublishEventCertificatesResponse{
		EventID:              response.EventID,
		PublishedCount:       response.PublishedCount,
		IsPublished:          response.IsPublished,
		InboxMessagesCreated: response.InboxMessagesCreated,
	})
}













