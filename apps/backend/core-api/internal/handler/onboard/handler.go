package onboard

import (
	"apps/backend/core-api/config"
	"apps/backend/services/auth"
	"time"

	customerror "apps/backend/common/customerror"

	verifyjwt "apps/backend/core-api/internal/middleware/verify_jwt"
	onboard_usecase "apps/backend/core-api/internal/usecase/onboard"
	profile_usecase "apps/backend/core-api/internal/usecase/profile"

	oauth_services "apps/backend/services/oauth"
)

type Handler struct {
	SessionExpiration  time.Duration
	AuthService        *auth.AuthService
	GoogleOAuthService *oauth_services.GoogleOAuthService

	OnboardUc *onboard_usecase.OnboardUsecase
	ProfileUc *profile_usecase.ProfileUsecase

	VerifyJwtMiddleware *verifyjwt.VerifyJwtMiddleware
}

func NewHandler(onboardUc *onboard_usecase.OnboardUsecase, profileUc *profile_usecase.ProfileUsecase, authService *auth.AuthService, googleOAuthService *oauth_services.GoogleOAuthService, verifyJwtMiddleware *verifyjwt.VerifyJwtMiddleware) (*Handler, error) {
	cfg := config.LoadConfig()

	expiration, err := time.ParseDuration(cfg.Jwt.Expiration)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	return &Handler{
		SessionExpiration:   expiration,
		AuthService:         authService,
		GoogleOAuthService:  googleOAuthService,
		OnboardUc:           onboardUc,
		VerifyJwtMiddleware: verifyjwt.New(authService),
	}, nil
}

type registerResponse struct {
	CredentialId string `json:"credential_id" validate:"required"`
	Jwt          string `json:"jwt" validate:"required"`
}
