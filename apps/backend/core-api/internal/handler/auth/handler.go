package auth

import (
	"apps/backend/core-api/internal/usecase/auth"
	oauth_services "apps/backend/services/oauth"
)

type Handler struct {
	AuthUc *auth.AuthUsecase

	GoogleOAuthService *oauth_services.GoogleOAuthService
}

func NewHandler(authUc *auth.AuthUsecase, googleOAuthService *oauth_services.GoogleOAuthService) *Handler {
	return &Handler{AuthUc: authUc, GoogleOAuthService: googleOAuthService}
}
