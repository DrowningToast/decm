package certificate

import (
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"
	"log/slog"

	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
)

type Handler struct {
	EventUc                      *event.EventUsecase
	AuthenticationService        *auth.AuthService
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
	Logger                       *slog.Logger
}

func NewHandler(eventUc *event.EventUsecase, authenticationService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware, logger *slog.Logger) *Handler {
	return &Handler{
		EventUc:                      eventUc,
		AuthenticationService:        authenticationService,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
		Logger:                       logger,
	}
}
