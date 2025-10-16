package eventconfig

import (
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	usecase "apps/backend/core-api/internal/usecase/eventconfig"
	"apps/backend/services/auth"
)

type Handler struct {
	EventConfigUc *usecase.EventConfigUsecase

	AuthenticationService         *auth.AuthService
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
}

func NewHandler(eventConfigUc *usecase.EventConfigUsecase, authenticationService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware) *Handler {
	return &Handler{
		EventConfigUc:                 eventConfigUc,
		AuthenticationService:         authenticationService,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
	}
}
