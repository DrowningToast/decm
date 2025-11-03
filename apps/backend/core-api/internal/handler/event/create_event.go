package event

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/pgmapper"
	"apps/backend/common/validatorutils"
	eventUc "apps/backend/core-api/internal/usecase/event"
	eventconfig "apps/backend/core-api/internal/usecase/eventconfig"

	"github.com/gofiber/fiber/v2"
)

type CreateEventRequest struct {
	Name             string `form:"name"`
	ShortDescription string `form:"short_description"`
	Description      string `form:"description"`
	StartDate        string `form:"start_date"`
	EndDate          string `form:"end_date"`
	SeatsCount       string `form:"seats_count"`
	ContactNumber    string `form:"contact_number"`
	ContactAddress   string `form:"contact_address"`
	Location         string `form:"location"`
	GoogleMapQuery   string `form:"google_map_query"`
	HostPassword     string `form:"host_password"`
}

// @Summary Create a new event
// @Description Create a new event with banner image upload
// @ID create-event
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
// @Param banner formData file true "Event banner image (JPEG, PNG, WebP, max 10MB)"
// @Param icon formData file true "Event icon image (JPEG, PNG, WebP, max 10MB)"
// @Param host_password formData string true "Host password"
// @Success 201 {object} entity.Event
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/events [post]
func (h *Handler) CreateEvent(ctx *fiber.Ctx) error {
	// 1. Parse form fields
	requestBody := CreateEventRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	// 2. Get banner file from form data
	bannerFile, err := ctx.FormFile("banner")
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("banner file is required"))
	}

	iconFile, err := ctx.FormFile("icon")
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("icon file is required"))
	}

	// 3. Validate banner file
	if err := validatorutils.ValidateImageFile(bannerFile); err != nil {
		return err
	}

	// 3. Validate icon file
	if err := validatorutils.ValidateImageFile(iconFile); err != nil {
		return err
	}

	// 4. Parse dates and integers
	startDate, err := time.Parse(time.RFC3339, requestBody.StartDate)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid start_date format: %w", err))
	}

	endDate, err := time.Parse(time.RFC3339, requestBody.EndDate)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid end_date format: %w", err))
	}

	seatsCount, err := strconv.Atoi(requestBody.SeatsCount)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid seats_count: %w", err))
	}

	// 5. Prepare parameters for usecase
	params := eventUc.CreateEventParameters{
		Name:             requestBody.Name,
		ShortDescription: requestBody.ShortDescription,
		Description:      requestBody.Description,
		StartDate:        startDate,
		EndDate:          endDate,
		SeatsCount:       seatsCount,
		ContactNumber:    requestBody.ContactNumber,
		ContactAddress:   requestBody.ContactAddress,
		Location:         requestBody.Location,
		GoogleMapQuery:   requestBody.GoogleMapQuery,
		EventBanner:      bannerFile,
		EventIcon:        iconFile,
		HostPassword:     requestBody.HostPassword,
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	// 6. Call usecase
	event, eventContractAddress, _, err := h.EventUc.CreateEvent(ctx.UserContext(), params, currentUser)
	if err != nil {
		return err
	}

	eventConfigParams := eventconfig.CreateEventRegistrationConfigParams{
		FirstNameRequirementStatus:           pgmapper.Int32ToPgInt4(0),
		LastNameRequirementStatus:            pgmapper.Int32ToPgInt4(0),
		EmailRequirementStatus:               pgmapper.Int32ToPgInt4(0),
		BioRequirementStatus:                 pgmapper.Int32ToPgInt4(0),
		PhoneNumberRequirementStatus:         pgmapper.Int32ToPgInt4(0),
		AddressRequirementStatus:             pgmapper.Int32ToPgInt4(0),
		AcademicInstitutionRequirementStatus: pgmapper.Int32ToPgInt4(0),
		AcademicEmailRequirementStatus:       pgmapper.Int32ToPgInt4(0),
	}

	_, err = h.EventConfigUc.CreateEventRegistrationConfig(ctx.UserContext(), event.ID, eventConfigParams)
	if err != nil {
		return err
	}

	createEventContractParams := eventUc.CreateEventContractParams{
		EventContractAddress:         eventContractAddress.Hex(),
		AccessManagerContractAddress: "",
	}

	_, err = h.EventUc.CreateEventContract(ctx.UserContext(), event.ID, createEventContractParams)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(event)
}

// Parse - Parse form data from request
func (r *CreateEventRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

// IsValid - Validate request fields
func (r *CreateEventRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}
