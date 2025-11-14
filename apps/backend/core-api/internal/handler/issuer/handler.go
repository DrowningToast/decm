package issuer

import (
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	roleguard "apps/backend/core-api/internal/middleware/role_guard"
	issuer_usecase "apps/backend/core-api/internal/usecase/issuer"
	"apps/backend/services/auth"
)

type Handler struct {
	IssuerUc                      *issuer_usecase.IssuerUsecase
	AuthenticationService         *auth.AuthService
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
	RoleGuardMiddleware           *roleguard.RoleGuardMiddleware
}

func NewHandler(issuerUc *issuer_usecase.IssuerUsecase, authenticationService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware, roleGuardMiddleware *roleguard.RoleGuardMiddleware) *Handler {
	return &Handler{
		IssuerUc:                      issuerUc,
		AuthenticationService:         authenticationService,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
		RoleGuardMiddleware:           roleGuardMiddleware,
	}
}
