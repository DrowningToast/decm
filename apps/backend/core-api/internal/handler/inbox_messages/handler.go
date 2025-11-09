package inboxmessages

import (
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	auth_usecase "apps/backend/core-api/internal/usecase/auth"
	"apps/backend/core-api/internal/usecase/inbox"
	auth_service "apps/backend/services/auth"
)

type Handler struct {
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
	AuthService                   *auth_service.AuthService

	AuthUc  *auth_usecase.AuthUsecase
	InboxUc *inbox.InboxUsecase
}

func NewHandler(inboxUsecase *inbox.InboxUsecase) *Handler {
	return &Handler{InboxUc: inboxUsecase}
}
