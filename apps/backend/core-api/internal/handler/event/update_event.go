package event

import (
	customerror "apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	eventUc "apps/backend/core-api/internal/usecase/event"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UpdateEventRequest struct {
	Name             *string
	ShortDescription *string
	Description      *string
	StartDate        *string
	EndDate          *string
	SeatsCount       *string
	ContactNumber    *string
	ContactAddress   *string
	Location         *string
	GoogleMapQuery   *string
}

// @Summary Update an event
// @Description Update an existing event with optional banner and icon image upload
// @ID update-event
// @Tags Event
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Event name"
// @Param short_description formData string true "Event short description"
// @Param description formData string true "Event description"
// @Param start_date formData string true "Start date (RFC3339 format)"
// @Param end_date formData string true "End date (RFC3339 format)"
// @Param seats_count formData integer true "Number of seats"
// @Param contact_number formData string true "Contact number"
// @Param contact_address formData string true "Contact address"
// @Param location formData string true "Location"
// @Param google_map_query formData string true "Google map query"
// @Param banner formData file false "Event banner image (JPEG, PNG, WebP, max 10MB) - optional"
// @Param icon formData file false "Event icon image (JPEG, PNG, WebP, max 10MB) - optional"
// @Success 200 {object} entity.Event
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/events/:event_id [put]
func (h *Handler) UpdateEvent(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	params := &UpdateEventRequest{}
	if err := params.Parse(ctx); err != nil {
		return err
	}
	if err := params.IsValid(); err != nil {
		return err
	}

	// Check if required fields are not nil
	if params.StartDate == nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("start_date is required"))
	}
	if params.EndDate == nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("end_date is required"))
	}
	if params.SeatsCount == nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("seats_count is required"))
	}

	startDate, err := time.Parse(time.RFC3339, *params.StartDate)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid start_date format: %w", err))
	}

	endDate, err := time.Parse(time.RFC3339, *params.EndDate)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid end_date format: %w", err))
	}

	seatsCount, err := strconv.Atoi(*params.SeatsCount)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid seats_count: %w", err))
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	updateEventParams := eventUc.UpdateEventParameters{
		Name:             params.Name,
		ShortDescription: params.ShortDescription,
		Description:      params.Description,
		StartDate:        &startDate,
		EndDate:          &endDate,
		SeatsCount:       &seatsCount,
		ContactNumber:    params.ContactNumber,
		ContactAddress:   params.ContactAddress,
		Location:         params.Location,
		GoogleMapQuery:   params.GoogleMapQuery,
	}

	bannerFile, _ := ctx.FormFile("banner") // Optional field
	iconFile, _ := ctx.FormFile("icon")     // Optional field

	if bannerFile != nil {
		if err := validatorutils.ValidateImageFile(bannerFile); err != nil {
			return err
		}
		updateEventParams.EventBanner = bannerFile
	}
	if iconFile != nil {
		if err := validatorutils.ValidateImageFile(iconFile); err != nil {
			return err
		}
		updateEventParams.EventIcon = iconFile
	}

	event, err := h.EventUc.UpdateEvent(ctx.UserContext(), eventID, updateEventParams, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(event)
}

// Parse - Parse form data from request
func (r *UpdateEventRequest) Parse(ctx *fiber.Ctx) error {
	// For multipart/form-data requests, we need to parse the form fields manually
	if _, err := ctx.Context().MultipartForm(); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Get form values
	if name := ctx.FormValue("name"); name != "" {
		r.Name = &name
	}
	if shortDesc := ctx.FormValue("short_description"); shortDesc != "" {
		r.ShortDescription = &shortDesc
	}
	if desc := ctx.FormValue("description"); desc != "" {
		r.Description = &desc
	}
	if startDate := ctx.FormValue("start_date"); startDate != "" {
		r.StartDate = &startDate
	}
	if endDate := ctx.FormValue("end_date"); endDate != "" {
		r.EndDate = &endDate
	}
	if seatsCount := ctx.FormValue("seats_count"); seatsCount != "" {
		r.SeatsCount = &seatsCount
	}
	if contactNumber := ctx.FormValue("contact_number"); contactNumber != "" {
		r.ContactNumber = &contactNumber
	}
	if contactAddress := ctx.FormValue("contact_address"); contactAddress != "" {
		r.ContactAddress = &contactAddress
	}
	if location := ctx.FormValue("location"); location != "" {
		r.Location = &location
	}
	if googleMapQuery := ctx.FormValue("google_map_query"); googleMapQuery != "" {
		r.GoogleMapQuery = &googleMapQuery
	}

	return nil
}

// IsValid - Validate request fields
func (r *UpdateEventRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}
