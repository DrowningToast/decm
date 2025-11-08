package event

import (
	"errors"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ImportCertificateReceiversRequest struct {
	EventID   uuid.UUID                          `json:"event_id" validate:"required"`
	Receivers []ImportCertificateReceiverRequest `json:"receivers" validate:"required,min=1"`
	HostPin   string                             `json:"host_pin" validate:"required"`
}

type ImportCertificateReceiverRequest struct {
	FirstName           string `json:"first_name" validate:"required"`
	LastName            string `json:"last_name" validate:"required"`
	AcademicInstitution string `json:"academic_institution" validate:"required"`
	CertificateTitle    string `json:"certificate_title" validate:"required"`
	CertificateSubtitle string `json:"certificate_subtitle" validate:"required"`
}

type ImportCertificateReceiversResponse struct {
	EventID                 uuid.UUID                  `json:"event_id"`
	EventCertificateAddress string                     `json:"event_certificate_address"`
	Certificates            []*entity.EventCertificate `json:"certificates"`
}

// @Summary Import certificate receivers for an event
// @Description Import certificate receivers for an event, deploys event certificate contract, and creates certificate records
// @ID import-certificate-receivers
// @Tags Event Certificates
// @Accept json
// @Produce json
// @Param request body ImportCertificateReceiversRequest true "Import certificate receivers request"
// @Success 201 {object} ImportCertificateReceiversResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/certificates/import [post]
func (h Handler) ImportCertificateReceivers(ctx *fiber.Ctx) error {
	requestBody := ImportCertificateReceiversRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	// Get current user from JWT
	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	// Convert request to usecase format
	requests := make([]event.ImportCertificateReceiversRequest, 0, len(requestBody.Receivers))
	for _, receiver := range requestBody.Receivers {
		requests = append(requests, event.ImportCertificateReceiversRequest{
			FirstName:           receiver.FirstName,
			LastName:            receiver.LastName,
			AcademicInstitution: receiver.AcademicInstitution,
			CertificateTitle:    receiver.CertificateTitle,
			CertificateSubtitle: receiver.CertificateSubtitle,
			HostPin:             requestBody.HostPin,
		})
	}

	response, err := h.EventUc.ImportCertificateReceivers(ctx.UserContext(), requestBody.EventID, requests, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (r *ImportCertificateReceiversRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func (r *ImportCertificateReceiversRequest) IsValid() error {
	if len(r.Receivers) == 0 {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("at least one receiver is required"))
	}

	for _, receiver := range r.Receivers {
		if receiver.FirstName == "" {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("first_name is required"))
		}
		if receiver.LastName == "" {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("last_name is required"))
		}
	}

	return nil
}
