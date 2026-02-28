package system_status

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type GetSchedulesBetweenQuery struct {
	StartTime int64 `query:"start_time" validate:"required"`
	EndTime   int64 `query:"end_time" validate:"required,gtfield=StartTime"`
}

type GetSchedulesBetweenResponse struct {
	Schedules []*entity.SystemStatusSchedule `json:"schedules"`
}

func (q *GetSchedulesBetweenQuery) Parse(ctx *fiber.Ctx) error {
	return ctx.QueryParser(q)
}

func (q *GetSchedulesBetweenQuery) IsValid() error {
	validate := validator.New()
	err := validate.Struct(q)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

// GetSchedulesBetween godoc
// @Summary Get system status schedules between time period
// @Description Get system status schedules updated within a specific time period using unix timestamps
// @ID get-schedules-between
// @Tags SystemStatus
// @Param start_time query int true "Start time unix timestamp"
// @Param end_time query int true "End time unix timestamp"
// @Produce json
// @Success 200 {object} GetSchedulesBetweenResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/system-status/period [get]
func (h *Handler) GetSchedulesBetween(ctx *fiber.Ctx) error {
	var query GetSchedulesBetweenQuery
	if err := query.Parse(ctx); err != nil {
		return err
	}
	if err := query.IsValid(); err != nil {
		return err
	}

	startDate := time.Unix(query.StartTime, 0)
	endDate := time.Unix(query.EndTime, 0)

	schedules, err := h.systemStatusUc.GetSchedulesBetween(ctx.UserContext(), startDate, endDate)
	if err != nil {
		return err
	}

	return ctx.JSON(&GetSchedulesBetweenResponse{
		Schedules: schedules,
	})
}
