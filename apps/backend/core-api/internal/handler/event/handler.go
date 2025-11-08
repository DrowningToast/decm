package event

import (
	"log/slog"

	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	roleguard "apps/backend/core-api/internal/middleware/role_guard"
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/core-api/internal/usecase/event_registration_invitation"
	eventconfig "apps/backend/core-api/internal/usecase/eventconfig"
	profile "apps/backend/core-api/internal/usecase/profile"
	"apps/backend/services/auth"
)

type Handler struct {
	EventUc                       *event.EventUsecase
	EventConfigUc                 *eventconfig.EventConfigUsecase
	ProfileUc                     *profile.ProfileUsecase
	EventRegistrationInvitationUc *event_registration_invitation.EventRegistrationInvitationUsecase
	AuthenticationService         *auth.AuthService
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
	RoleGuardMiddleware           *roleguard.RoleGuardMiddleware
	Logger                        *slog.Logger
}

func NewHandler(eventUc *event.EventUsecase, eventConfigUc *eventconfig.EventConfigUsecase, profileUc *profile.ProfileUsecase, eventRegistrationInvitationUc *event_registration_invitation.EventRegistrationInvitationUsecase, authenticationService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware, logger *slog.Logger) *Handler {
	return &Handler{
		EventUc:                       eventUc,
		EventConfigUc:                 eventConfigUc,
		ProfileUc:                     profileUc,
		EventRegistrationInvitationUc: eventRegistrationInvitationUc,
		AuthenticationService:         authenticationService,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
		Logger:                        logger,
	}
}
