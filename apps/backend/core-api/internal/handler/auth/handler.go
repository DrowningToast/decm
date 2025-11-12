package auth

import (
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	auth_usecase "apps/backend/core-api/internal/usecase/auth"
	"apps/backend/core-api/internal/usecase/oauth"
	"apps/backend/services/auth"
	oauth_services "apps/backend/services/oauth"
)

type Handler struct {
	OAuthUc *oauth.OAuthUsecase
	AuthUc  *auth_usecase.AuthUsecase

	GoogleOAuthService *oauth_services.GoogleOAuthService
	AuthService        *auth.AuthService

	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
}

func NewHandler(oauthUc *oauth.OAuthUsecase, authUc *auth_usecase.AuthUsecase, googleOAuthService *oauth_services.GoogleOAuthService, authService *auth.AuthService, authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware) *Handler {
	return &Handler{OAuthUc: oauthUc, AuthUc: authUc, GoogleOAuthService: googleOAuthService, AuthService: authService, AuthenticationGuardMiddleware: authenticationGuardMiddleware}
}
