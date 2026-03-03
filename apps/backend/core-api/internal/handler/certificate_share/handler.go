package certificate_share

import (
	certificate_share_usecase "apps/backend/core-api/internal/usecase/certificate_share"
	"apps/backend/services/auth"
	"log/slog"

	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
)

type Handler struct {
	CertificateShareUc            *certificate_share_usecase.CertificateShareUsecase
	AuthenticationService         *auth.AuthService
	AuthenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware
	Logger                        *slog.Logger
}

func NewHandler(
	certificateShareUc *certificate_share_usecase.CertificateShareUsecase,
	authenticationService *auth.AuthService,
	authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		CertificateShareUc:            certificateShareUc,
		AuthenticationService:         authenticationService,
		AuthenticationGuardMiddleware: authenticationGuardMiddleware,
		Logger:                        logger,
	}
}
