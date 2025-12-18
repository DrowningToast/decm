package customerror

import (
	"fmt"
	"log/slog"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

type ErrApiResponse struct {
	Message string `json:"message"`
}

func GetErrFiberHandler(logger *slog.Logger) func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		path := ctx.OriginalURL()
		method := ctx.Method()

		// Is custom error
		var customErr *Err
		if errors.As(err, &customErr) && customErr != nil {
			// Log the error
			stackTrace := fmt.Sprintf("%+v", err)
			switch customErr.LoggerLevel {
			case slog.LevelError:
				logger.ErrorContext(ctx.UserContext(), errors.Wrap(err, customErr.Message).Error(), slog.String("method", method), slog.String("path", path), slog.String("stack_trace", stackTrace))
			case slog.LevelWarn:
				logger.WarnContext(ctx.UserContext(), errors.Wrap(err, customErr.Message).Error(), slog.String("method", method), slog.String("path", path), slog.String("stack_trace", stackTrace))
			case slog.LevelInfo:
				logger.InfoContext(ctx.UserContext(), errors.Wrap(err, customErr.Message).Error(), slog.String("method", method), slog.String("path", path), slog.String("stack_trace", stackTrace))
			case slog.LevelDebug:
				logger.DebugContext(ctx.UserContext(), errors.Wrap(err, customErr.Message).Error(), slog.String("method", method), slog.String("path", path), slog.String("stack_trace", stackTrace))
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
			stackTrace := fmt.Sprintf("%+v", err)
			switch fiberErr.Code {
			case fiber.StatusBadRequest:
				logger.WarnContext(ctx.UserContext(), "A Fiber error has occurred", slog.Int("status_code", fiberErr.Code), slog.String("error", fiberErr.Error()), slog.String("method", method), slog.String("path", path), slog.String("stack_trace", stackTrace))
			case fiber.StatusUnauthorized:
				logger.WarnContext(ctx.UserContext(), "A Fiber error has occurred", slog.Int("status_code", fiberErr.Code), slog.String("error", fiberErr.Error()), slog.String("method", method), slog.String("path", path), slog.String("stack_trace", stackTrace))
			case fiber.StatusNotFound:
				logger.WarnContext(ctx.UserContext(), "A Fiber error has occurred", slog.Int("status_code", fiberErr.Code), slog.String("error", fiberErr.Error()), slog.String("method", method), slog.String("path", path), slog.String("stack_trace", stackTrace))
			case fiber.StatusInternalServerError:
				logger.ErrorContext(ctx.UserContext(), "A Fiber error has occurred", slog.Int("status_code", fiberErr.Code), slog.String("error", fiberErr.Error()), slog.String("method", method), slog.String("path", path), slog.String("stack_trace", stackTrace))
			default:
				logger.ErrorContext(ctx.UserContext(), "A Fiber error has occurred", slog.Int("status_code", fiberErr.Code), slog.String("error", fiberErr.Error()), slog.String("method", method), slog.String("path", path), slog.String("stack_trace", stackTrace))
			}
			// Return the error
			return ctx.Status(fiberErr.Code).JSON(
				ErrApiResponse{
					Message: fiberErr.Message,
				},
			)
		}
		// Is unknown general error
		logger.ErrorContext(ctx.UserContext(), errors.Wrap(err, "an unknown error has occurred").Error(), slog.String("method", method), slog.String("path", path), slog.String("stack_trace", fmt.Sprintf("%+v", err)))
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrApiResponse{
				Message: "An unknown error has occurred. Please try again later.",
			},
		)
	}
}
