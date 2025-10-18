package event

import (
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	"apps/backend/core-api/internal/usecase/event"
	eventconfig "apps/backend/core-api/internal/usecase/eventconfig"
	"apps/backend/services/auth"
	"log/slog"
)

type Handler struct {
	EventUc                       *event.EventUsecase
	EventConfigUc                 *eventconfig.EventConfigUsecase
	AuthenticationService         *auth.AuthService
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
	Logger                        *slog.Logger
}

func NewHandler(eventUc *event.EventUsecase, eventConfigUc *eventconfig.EventConfigUsecase, authenticationService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware, logger *slog.Logger) *Handler {
	return &Handler{
		EventUc:                       eventUc,
		EventConfigUc:                 eventConfigUc,
		AuthenticationService:         authenticationService,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
		Logger:                        logger,
	}
}
