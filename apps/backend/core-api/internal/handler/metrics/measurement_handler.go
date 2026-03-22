package metrics

import (
	"apps/backend/common/customerror"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var measurement = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "decm_handler_duration_milliseconds",
	Buckets: []float64{5, 10, 50, 100, 200, 400, 800, 2000, 5000, 10000},
}, []string{"handler", "status", "http_status"})

func MeasureHandlerDuration(operation string) func(httpStatus *int, err error) time.Duration {
	startTime := time.Now()

	return func(httpStatus *int, err error) time.Duration {
		elapsed := time.Since(startTime)

		_httpStatus := ""
		if err != nil {
			customErr := customerror.TryParseAsCustomErr(err)
			if customErr != nil {
				_httpStatus = strconv.Itoa(*customErr.HttpStatus)
			} else {
				_httpStatus = strconv.Itoa(fiber.ErrInternalServerError.Code)
			}
		}
		if httpStatus != nil {
			_httpStatus = strconv.Itoa(*httpStatus)
		}

		_status := "success"
		if err != nil {
			_status = "failure"
		}

		measurement.WithLabelValues(operation, _status, _httpStatus).Observe(float64(elapsed.Milliseconds()))

		return elapsed
	}
}

func MeasureHandlerDurationWrapper(operation string, handler func(ctx *fiber.Ctx) error) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		cb := MeasureHandlerDuration(operation)

		err := handler(ctx)
		statusCode := ctx.Response().StatusCode()
		cb(&statusCode, err)

		return err
	}
}
