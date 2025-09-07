package onboard

import (
	"time"

	"apps/backend/core-api/config"
	usecase "apps/backend/core-api/internal/usecase/onboard"
)

type Handler struct {
	SessionExpiration time.Duration

	OnboardUc *usecase.OnboardUsecase
}

func NewHandler(onboardUc *usecase.OnboardUsecase) Handler {
	cfg := config.LoadConfig()

	return Handler{
		SessionExpiration: cfg.Jwt.Expiration,
		OnboardUc:         onboardUc,
	}
}
