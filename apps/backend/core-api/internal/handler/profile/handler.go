package profile

import (
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	usecase "apps/backend/core-api/internal/usecase/profile"
	"apps/backend/services/auth"
)

type Handler struct {
	ProfileUc *usecase.ProfileUsecase

	AuthenticationService         *auth.AuthService
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
}

func NewHandler(profileUc *usecase.ProfileUsecase, authenticationService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware) *Handler {
	return &Handler{
		ProfileUc:                     profileUc,
		AuthenticationService:         authenticationService,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
	}
}
