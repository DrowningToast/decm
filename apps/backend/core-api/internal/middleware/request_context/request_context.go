package requestcontext

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// New creates a middleware that extracts the request ID from Fiber's locals
// (set by the requestid middleware) and adds it to the request's context.Context
// with the key "request_id". This allows the logger to retrieve it.
func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get request ID from Fiber locals (set by requestid middleware)
		reqID, ok := c.Locals("requestid").(string)
		if !ok || reqID == "" {
			// Fallback to response header if not in locals
			reqID = c.GetRespHeader("X-Request-ID")
		}

		if reqID != "" {
			// Create new context with request_id
			// "request_id" matches the key expected by apps/backend/common/log/log_handler.go
			ctx := context.WithValue(c.UserContext(), "request_id", reqID)
			c.SetUserContext(ctx)
		}

		return c.Next()
	}
}
