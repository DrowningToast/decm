package eventconfig

import (
	"apps/backend/services/auth"
	"log/slog"

	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	roleguard "apps/backend/core-api/internal/middleware/role_guard"
	event_usecase "apps/backend/core-api/internal/usecase/event"
	usecase "apps/backend/core-api/internal/usecase/eventconfig"
)

type Handler struct {
	EventConfigUc *usecase.EventConfigUsecase
	EventUc       *event_usecase.EventUsecase

	AuthenticationService *auth.AuthService

	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
	RoleGuardMiddleware           *roleguard.RoleGuardMiddleware

	Logger *slog.Logger
}

func NewHandler(eventConfigUc *usecase.EventConfigUsecase, eventUc *event_usecase.EventUsecase, authenticationService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware, roleGuardMiddleware *roleguard.RoleGuardMiddleware, logger *slog.Logger) *Handler {
	return &Handler{
		EventConfigUc:                 eventConfigUc,
		EventUc:                       eventUc,
		AuthenticationService:         authenticationService,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
		RoleGuardMiddleware:           roleGuardMiddleware,
		Logger:                        logger,
	}
}
