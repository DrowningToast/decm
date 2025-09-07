package customerror

import (
	"log/slog"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

type ErrApiResponse struct {
	message string
}

func GetErrFiberHandler(logger *slog.Logger) func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		// Is custom error
		var customErr *Err
		if errors.As(err, &customErr) {
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
					message: customErr.Message,
				},
			)
		}
		// Is unknown general error
		logger.Error(errors.Wrap(err, "an unknown error has occurred").Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrApiResponse{
				message: "An unknown error has occurred. Please try again later.",
			},
		)
	}
}
