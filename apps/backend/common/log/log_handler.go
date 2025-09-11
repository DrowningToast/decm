package log

import (
	"context"
	"log/slog"
)

type RequestIDHandler struct {
	slog.Handler
}

func (h *RequestIDHandler) Handle(ctx context.Context, r slog.Record) error {
	// Retrieve the request ID from the context (e.g., via a custom context key)
	if requestID, ok := ctx.Value("request_id").(string); ok {
		r.AddAttrs(slog.String("request_id", requestID))
	}
	return h.Handler.Handle(ctx, r)
}

func NewRequestIDHandler(h slog.Handler) *RequestIDHandler {
	return &RequestIDHandler{h}
}

