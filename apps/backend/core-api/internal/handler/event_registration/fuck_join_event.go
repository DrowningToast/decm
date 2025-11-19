package event_registration

import (
	customerror "apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/usecase/event_registration"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FuckJoinEventPayloadHandler struct {
	FirstName           *string `json:"first_name,omitempty"`
	LastName            *string `json:"last_name,omitempty"`
	Email               *string `json:"email,omitempty"`
	PhoneNumber         *string `json:"phone_number,omitempty"`
	AcademicInstitution *string `json:"academic_institution,omitempty"`
	AcademicEmail       *string `json:"academic_email,omitempty"`
	Address             *string `json:"address,omitempty"`
	Bio                 *string `json:"bio,omitempty"`
	PinCode             *string `json:"pin_code,omitempty"`
}

func (r *FuckJoinEventPayloadHandler) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

func (r *FuckJoinEventPayloadHandler) Parse(ctx *fiber.Ctx) error {
	return ctx.BodyParser(r)
}

// @Summary Fuck join event
// @Description Fuck join event
// @ID fuck-join-event
// @Tags Event Registration
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param fuckJoinEventPayload body FuckJoinEventPayload true "Fuck join event payload"
// @Success 200 {object} entity.EventAttendee
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/event-registration/join/{event_id}/fuck [post]
func (h *Handler) FuckJoinEvent(ctx *fiber.Ctx) error {

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	payload := FuckJoinEventPayloadHandler{}
	if err := payload.Parse(ctx); err != nil {
		return err
	}

	if err := payload.IsValid(); err != nil {
		return err
	}

	eventIdStr := ctx.Params("event_id")
	if eventIdStr == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("event_id is required"))
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	fuckJoinEventPayload := event_registration.FuckJoinEventPayload{
		PinCode:             payload.PinCode,
		FirstName:           payload.FirstName,
		LastName:            payload.LastName,
		Email:               payload.Email,
		PhoneNumber:         payload.PhoneNumber,
		AcademicInstitution: payload.AcademicInstitution,
		AcademicEmail:       payload.AcademicEmail,
		Address:             payload.Address,
		Bio:                 payload.Bio,
	}

	eventAttendee, err := h.EventRegistrationUc.FuckJoinEvent(ctx.Context(), currentUser, eventId, fuckJoinEventPayload)
	if err != nil {
		return err
	}

	return ctx.JSON(eventAttendee)
}
