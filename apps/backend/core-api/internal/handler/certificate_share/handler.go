package certificate_share_handler

import (
	certificate_share_usecase "apps/backend/core-api/internal/usecase/certificate_share"
	"apps/backend/services/auth"
	"log/slog"
	"time"

	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"

	"github.com/google/uuid"
)

// CertificateShareViewModel is the handler-layer response shape for a CertificateShare.
// Defined as a concrete struct (not an alias) so swag can resolve it without
// cross-package alias lookups — both handler and usecase declare package certificate_share.
type CertificateShareViewModel struct {
	Id                 uuid.UUID `json:"id"`
	EventCertificateId uuid.UUID `json:"event_certificate_id"`
	Active             bool      `json:"active"`
	Handle             string    `json:"handle"`
	HasPassword        bool      `json:"has_password"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func newCertificateShareViewModel(vm *certificate_share_usecase.CertificateShareViewModel) *CertificateShareViewModel {
	if vm == nil {
		return nil
	}
	return &CertificateShareViewModel{
		Id:                 vm.Id,
		EventCertificateId: vm.EventCertificateId,
		Active:             vm.Active,
		Handle:             vm.Handle,
		HasPassword:        vm.HasPassword,
		CreatedAt:          vm.CreatedAt,
		UpdatedAt:          vm.UpdatedAt,
	}
}

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
