package auth

import (
	"apps/backend/core-api/internal/usecase/oauth"
	"apps/backend/services/auth"
	oauth_services "apps/backend/services/oauth"
)

type Handler struct {
	AuthUc *oauth.OAuthUsecase

	GoogleOAuthService *oauth_services.GoogleOAuthService
	AuthService        *auth.AuthService
}

func NewHandler(authUc *oauth.OAuthUsecase, googleOAuthService *oauth_services.GoogleOAuthService, authService *auth.AuthService) *Handler {
	return &Handler{AuthUc: authUc, GoogleOAuthService: googleOAuthService, AuthService: authService}
}
