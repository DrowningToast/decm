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

func NewHandler(onboardUc *usecase.OnboardUsecase) Handler {
	cfg := config.LoadConfig()

	return Handler{
		SessionExpiration: cfg.Jwt.Expiration,
		OnboardUc:         onboardUc,
	}
}
