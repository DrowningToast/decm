package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ImportCertificateReceiversRequest struct {
	EventID   uuid.UUID                          `json:"event_id" validate:"required"`
	Receivers []ImportCertificateReceiverRequest `json:"receivers" validate:"required,min=1"`
	HostPin   string                             `json:"host_pin" validate:"required"`
}

type ImportCertificateReceiverRequest struct {
	Email               string  `json:"email" validate:"required,email"`
	FirstName           *string `json:"first_name"`
	LastName            *string `json:"last_name"`
	AcademicInstitution *string `json:"academic_institution"`
	CertificateTitle    *string `json:"certificate_title"`
	CertificateSubtitle *string `json:"certificate_subtitle"`
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
			Email:               receiver.Email,               // Required string
			FirstName:           receiver.FirstName,           // Already *string
			LastName:            receiver.LastName,            // Already *string
			AcademicInstitution: receiver.AcademicInstitution, // Already *string
			CertificateTitle:    receiver.CertificateTitle,    // Already *string
			CertificateSubtitle: receiver.CertificateSubtitle, // Already *string
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

	// Validate that each receiver has a required email
	for i, receiver := range r.Receivers {
		if receiver.Email == "" {
			return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("receiver at index %d: email is required", i))
		}
	}

	return nil
}
