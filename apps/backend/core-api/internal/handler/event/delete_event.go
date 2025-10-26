package event

import (
	customerror "apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DeleteEvent godoc
// @Summary Delete event by ID
// @Description Delete event by ID
// @Tags Event
// @ID delete-event-by-id
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} EventResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id} [delete]
func (h *Handler) DeleteEvent(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	event, err := h.EventUc.DeleteEvent(ctx.UserContext(), eventID, currentUser)
	if err != nil {
		return err
	}

	// Map domain event to EventResponse DTO
	response := EventResponse{
		ID:                       event.ID,
		ChainID:                  int32(event.ChainID),
		ContactNumber:            event.ContactNumber,
		OwnerCredentialID:        event.OwnerCredentialID,
		BannerStorageKey:         event.BannerStorageKey,
		IconStorageKey:           event.IconStorageKey,
		BannerPresignedURL:       "", // Empty after deletion
		IconPresignedURL:         "", // Empty after deletion
		Title:                    event.Title,
		ShortDescription:         event.ShortDescription,
		LongDescription:          event.LongDescription,
		StartDate:                event.StartDate,
		EndDate:                  event.EndDate,
		Location:                 event.Location,
		GoogleMapQuery:           event.GoogleMapQuery,
		MaxAttendees:             int32(event.MaxAttendees),
		IsPublic:                 event.IsPublic,
		IsBookingRequestRequired: event.IsBookingRequestRequired,
		IsVerified:               event.IsVerified,
		IsTicketTransferable:     event.IsTicketTransferable,
		CreatedAt:                event.CreatedAt,
		UpdatedAt:                event.UpdatedAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
