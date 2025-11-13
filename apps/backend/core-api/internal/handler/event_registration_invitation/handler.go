package event_registration_invitation

import (
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	roleguard "apps/backend/core-api/internal/middleware/role_guard"
	eventRegistrationInvitationUc "apps/backend/core-api/internal/usecase/event_registration_invitation"
	"apps/backend/services/auth"
)

type Handler struct {
	AuthenticationService         auth.AuthService
	EventRegistrationInvitationUc *eventRegistrationInvitationUc.EventRegistrationInvitationUsecase

	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
	RoleGuardMiddleware           *roleguard.RoleGuardMiddleware
}

func NewHandler(
	authenticationService auth.AuthService,
	eventRegistrationInvitationUc *eventRegistrationInvitationUc.EventRegistrationInvitationUsecase,
	authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware,
	roleGuardMiddleware *roleguard.RoleGuardMiddleware,
) *Handler {
	return &Handler{
		AuthenticationService:         authenticationService,
		EventRegistrationInvitationUc: eventRegistrationInvitationUc,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
		RoleGuardMiddleware:           roleGuardMiddleware,
	}
}
