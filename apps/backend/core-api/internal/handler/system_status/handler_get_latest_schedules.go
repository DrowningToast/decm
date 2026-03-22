package system_status

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type GetLatestSchedulesQuery struct {
	PageSize int32 `query:"page_size" validate:"required,min=1,max=100"`
}

type GetLatestSchedulesResponse struct {
	Schedules []*entity.SystemStatusSchedule `json:"schedules"`
}

func (q *GetLatestSchedulesQuery) Parse(ctx *fiber.Ctx) error {
	return ctx.QueryParser(q)
}

func (q *GetLatestSchedulesQuery) IsValid() error {
	validate := validator.New()
	err := validate.Struct(q)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

// GetLatestSchedules godoc
// @Summary Get latest system status schedules
// @Description Get latest system status schedules in the past
// @ID get-latest-schedules
// @Tags SystemStatus
// @Param page_size query int true "Number of records to return"
// @Produce json
// @Success 200 {object} GetLatestSchedulesResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/system-status/latest [get]
func (h *Handler) GetLatestSchedules(ctx *fiber.Ctx) error {
	var query GetLatestSchedulesQuery
	if err := query.Parse(ctx); err != nil {
		return err
	}
	if err := query.IsValid(); err != nil {
		return err
	}

	schedules, err := h.systemStatusUc.GetLatestSchedules(ctx.UserContext(), query.PageSize)
	if err != nil {
		return err
	}

	return ctx.JSON(&GetLatestSchedulesResponse{
		Schedules: schedules,
	})
}
