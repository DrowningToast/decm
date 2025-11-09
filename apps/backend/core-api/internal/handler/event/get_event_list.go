package event

import (
	customerror "apps/backend/common/customerror"
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type GetEventListQuery struct {
	IncludeActiveEvents   bool `query:"include_active_events,omitempty" validate:"omitempty,boolean" default:"true"`
	IncludeInactiveEvents bool `query:"include_inactive_events,omitempty" validate:"omitempty,boolean" default:"true"`
	IncludeClosedEvents   bool `query:"include_closed_events,omitempty" validate:"omitempty,boolean" default:"false"`

	OnlyUserJoinedEvents bool `query:"only_user_joined_events,omitempty" validate:"omitempty,boolean" default:"false"`
}

// GetEventsList godoc
// @Summary Get events list
// @Description Get events list
// @ID get-events-list
// @Tags Event
// @Param include_active_events query bool false "Include active events" default(true)
// @Param include_inactive_events query bool false "Include inactive events" default(true)
// @Param include_closed_events query bool false "Include closed events" default(false)
// @Param only_user_joined_events query bool false "Only user joined events" default(false)
// @Produce json
// @Success 200 {array} []EventResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events [get]
func (h *Handler) GetEventsList(ctx *fiber.Ctx) error {
	var getEventListQuery GetEventListQuery
	if err := getEventListQuery.Parse(ctx); err != nil {
		return err
	}
	if err := getEventListQuery.IsValid(); err != nil {
		return err
	}

	var currentUser *auth.JwtClaims
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	var targetUser *auth.JwtClaims = nil
	if getEventListQuery.OnlyUserJoinedEvents {
		targetUser = currentUser
	}
	events, err := h.EventUc.ListEvents(ctx.UserContext(), event.ListEventsParameters{
		OnlyUserJoinedEvents:  targetUser,
		IncludeActiveEvents:   getEventListQuery.IncludeActiveEvents,
		IncludeInactiveEvents: getEventListQuery.IncludeInactiveEvents,
		IncludeClosedEvents:   getEventListQuery.IncludeClosedEvents,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(events)
}

func (q *GetEventListQuery) Parse(ctx *fiber.Ctx) error {
	return ctx.QueryParser(q)
}

func (q *GetEventListQuery) IsValid() error {
	validate := validator.New()
	err := validate.Struct(q)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}
