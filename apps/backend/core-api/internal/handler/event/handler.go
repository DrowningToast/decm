package event

import (
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/core-api/internal/usecase/event_registration"
	"apps/backend/services/auth"
	"log/slog"

	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	roleguard "apps/backend/core-api/internal/middleware/role_guard"

	eventconfig "apps/backend/core-api/internal/usecase/eventconfig"
	profile "apps/backend/core-api/internal/usecase/profile"
)

type Handler struct {
	EventUc                       *event.EventUsecase
	EventConfigUc                 *eventconfig.EventConfigUsecase
	ProfileUc                     *profile.ProfileUsecase
	EventRegistrationInvitationUc *event_registration.EventRegistrationUsecase
	AuthenticationService         *auth.AuthService
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
	RoleGuardMiddleware           *roleguard.RoleGuardMiddleware
	Logger                        *slog.Logger
}

func NewHandler(eventUc *event.EventUsecase, eventConfigUc *eventconfig.EventConfigUsecase, profileUc *profile.ProfileUsecase, eventRegistrationInvitationUc *event_registration.EventRegistrationUsecase, authenticationService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware, logger *slog.Logger) *Handler {
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
