package event_registration

import (
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	roleguard "apps/backend/core-api/internal/middleware/role_guard"
	authUc "apps/backend/core-api/internal/usecase/auth"
	eventUsecase "apps/backend/core-api/internal/usecase/event"
	eventRegistrationUc "apps/backend/core-api/internal/usecase/event_registration"
	onboardUc "apps/backend/core-api/internal/usecase/onboard"
	"apps/backend/services/auth"
)

type Handler struct {
	AuthenticationService auth.AuthService
	EventRegistrationUc   *eventRegistrationUc.EventRegistrationUsecase
	EventUc               *eventUsecase.EventUsecase
	OnboardUc             *onboardUc.OnboardUsecase
	AuthUc                *authUc.AuthUsecase

	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
	RoleGuardMiddleware           *roleguard.RoleGuardMiddleware
}

func NewHandler(
	authenticationService auth.AuthService,
	eventRegistrationInvitationUc *eventRegistrationUc.EventRegistrationUsecase,
	eventUsecase *eventUsecase.EventUsecase,
	authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware,
	roleGuardMiddleware *roleguard.RoleGuardMiddleware,
) *Handler {
	return &Handler{
		AuthenticationService:         authenticationService,
		EventRegistrationUc:           eventRegistrationInvitationUc,
		EventUc:                       eventUsecase,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
		RoleGuardMiddleware:           roleGuardMiddleware,
	}
}
