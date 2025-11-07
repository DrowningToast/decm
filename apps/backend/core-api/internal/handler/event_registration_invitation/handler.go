package event_registration_invitation

import (
	eventRegistrationInvitationUc "apps/backend/core-api/internal/usecase/event_registration_invitation"
	"apps/backend/services/auth"
)

type Handler struct {
	AuthenticationService         auth.AuthService
	EventRegistrationInvitationUc *eventRegistrationInvitationUc.EventRegistrationInvitationUsecase
}

func NewHandler(
	authenticationService auth.AuthService,
	eventRegistrationInvitationUc *eventRegistrationInvitationUc.EventRegistrationInvitationUsecase,
) *Handler {
	return &Handler{
		AuthenticationService:         authenticationService,
		EventRegistrationInvitationUc: eventRegistrationInvitationUc,
	}
}
