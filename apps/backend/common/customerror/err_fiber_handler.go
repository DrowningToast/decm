package customerror

import (
	"log/slog"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

type ErrApiResponse struct {
	Message string `json:"message"`
}

func GetErrFiberHandler(logger *slog.Logger) func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		// Is custom error
		var customErr *Err
		if errors.As(err, &customErr) && customErr != nil {
			// Log the error
			switch customErr.LoggerLevel {
			case slog.LevelError:
				logger.Error(errors.Wrap(err, customErr.Message).Error())
			case slog.LevelWarn:
				logger.Warn(errors.Wrap(err, customErr.Message).Error())
			case slog.LevelInfo:
				logger.Info(errors.Wrap(err, customErr.Message).Error())
			case slog.LevelDebug:
				logger.Debug(errors.Wrap(err, customErr.Message).Error())
			}
			// Return the error
			return ctx.Status(*customErr.HttpStatus).JSON(
				ErrApiResponse{
					Message: customErr.Message,
				},
			)
		}
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			// Log fiber error
			switch fiberErr.Code {
			case fiber.StatusBadRequest:
				logger.WarnContext(ctx.UserContext(), "A Fiber error has occurred", slog.Int("status_code", fiberErr.Code), slog.String("error", fiberErr.Error()))
			case fiber.StatusUnauthorized:
				logger.WarnContext(ctx.UserContext(), "A Fiber error has occurred", slog.Int("status_code", fiberErr.Code), slog.String("error", fiberErr.Error()))
			case fiber.StatusNotFound:
				logger.WarnContext(ctx.UserContext(), "A Fiber error has occurred", slog.Int("status_code", fiberErr.Code), slog.String("error", fiberErr.Error()))
			case fiber.StatusInternalServerError:
				logger.ErrorContext(ctx.UserContext(), "A Fiber error has occurred", slog.Int("status_code", fiberErr.Code), slog.String("error", fiberErr.Error()))
			default:
				logger.ErrorContext(ctx.UserContext(), "A Fiber error has occurred", slog.Int("status_code", fiberErr.Code), slog.String("error", fiberErr.Error()))
			}
			// Return the error
			return ctx.Status(fiberErr.Code).JSON(
				ErrApiResponse{
					Message: fiberErr.Message,
				},
			)
		}
		// Is unknown general error
		logger.Error(errors.Wrap(err, "an unknown error has occurred").Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrApiResponse{
				Message: "An unknown error has occurred. Please try again later.",
			},
		)
	}
}
