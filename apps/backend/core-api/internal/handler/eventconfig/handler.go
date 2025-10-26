package eventconfig

import (
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	event_usecase "apps/backend/core-api/internal/usecase/event"
	usecase "apps/backend/core-api/internal/usecase/eventconfig"
	"apps/backend/services/auth"
)

type Handler struct {
	EventConfigUc *usecase.EventConfigUsecase
	EventUc       *event_usecase.EventUsecase

	AuthenticationService         *auth.AuthService
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
}

func NewHandler(eventConfigUc *usecase.EventConfigUsecase, eventUc *event_usecase.EventUsecase, authenticationService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware) *Handler {
	return &Handler{
		EventConfigUc:                 eventConfigUc,
		EventUc:                       eventUc,
		AuthenticationService:         authenticationService,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
	}
}
