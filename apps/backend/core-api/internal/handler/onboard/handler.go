package onboard

import (
	"time"

	"apps/backend/core-api/config"
	usecase "apps/backend/core-api/internal/usecase/onboard"
	"apps/backend/services/auth"
	oauth_services "apps/backend/services/oauth"
)

type Handler struct {
	SessionExpiration  time.Duration
	AuthService        *auth.AuthService
	GoogleOAuthService *oauth_services.GoogleOAuthService

	OnboardUc *usecase.OnboardUsecase
}

func NewHandler(onboardUc *usecase.OnboardUsecase, authService *auth.AuthService) Handler {
	cfg := config.LoadConfig()

	return Handler{
		SessionExpiration: cfg.Jwt.Expiration,
		AuthService:       authService,
		OnboardUc:         onboardUc,
	}
}

type registerResponse struct {
	CredentialId string `json:"credential_id" validate:"required"`
	Jwt          string `json:"jwt" validate:"required"`
}
