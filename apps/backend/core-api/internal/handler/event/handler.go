package event

import (
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"
)

type Handler struct {
	EventUc                       *event.EventUsecase
	AuthenticationService         *auth.AuthService
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
}

func NewHandler(eventUc *event.EventUsecase, authenticationService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware) *Handler {
	return &Handler{
		EventUc:                       eventUc,
		AuthenticationService:         authenticationService,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
	}
}
