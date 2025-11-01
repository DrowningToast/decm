package event

import (
	"context"
	"time"

	"apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetEventById godoc
// @Summary Get event by ID
// @Description Get event by ID
// @Tags Events
// @ID get-event-by-id
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} EventResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id} [get]
func (h *Handler) GetEventById(ctx *fiber.Ctx) error {
	eventId, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Create a context with 30-second timeout
	base := ctx.UserContext()
	ctxWithTimeout, cancel := context.WithTimeout(base, 30*time.Second)
	defer cancel()

	event, err := h.EventUc.GetEventById(ctxWithTimeout, eventId)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	if event == nil {
		return customerror.Parse(&customerror.ErrNotFound, nil)
	}

	bannerPresignedURL, err := h.EventUc.S3Service.GetPresignedURL(ctxWithTimeout, event.BannerStorageKey)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	iconPresignedURL, err := h.EventUc.S3Service.GetPresignedURL(ctxWithTimeout, event.IconStorageKey)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.JSON(EventResponse{
		ID:                       event.ID,
		ChainID:                  int32(event.ChainID),
		ContactNumber:            event.ContactNumber,
		OwnerCredentialID:        event.OwnerCredentialID,
		BannerStorageKey:         event.BannerStorageKey,
		IconStorageKey:           event.IconStorageKey,
		BannerPresignedURL:       bannerPresignedURL,
		IconPresignedURL:         iconPresignedURL,
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
	})
}
