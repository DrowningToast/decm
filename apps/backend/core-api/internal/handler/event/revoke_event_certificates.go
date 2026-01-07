package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RevokeEventCertificatesRequest struct {
	CertificateIDs []uuid.UUID `json:"certificate_ids" validate:"required,min=1"`
}

type RevokeEventCertificatesResponse struct {
	RevokedCertificates []*entity.EventCertificate `json:"revoked_certificates"`
}

// @Summary Revoke event certificates
// @Description Revoke event certificates by certificate IDs
// @ID revoke-event-certificates
// @Tags Event Certificates
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body RevokeEventCertificatesRequest true "Revoke certificates request"
// @Success 200 {object} RevokeEventCertificatesResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/certificates/revoke [post]
func (h Handler) RevokeEventCertificates(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	requestBody := RevokeEventCertificatesRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	// Get current user from JWT
	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	// Convert handler request to usecase request
	usecaseRequest := event.RevokeEventCertificatesRequest{
		CertificateIDs: requestBody.CertificateIDs,
	}

	response, err := h.EventUc.RevokeEventCertificates(ctx.UserContext(), eventID, usecaseRequest, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (r *RevokeEventCertificatesRequest) Parse(ctx *fiber.Ctx) error {
	// Define a temporary struct to parse JSON with string UUIDs
	type tempRequest struct {
		CertificateIDs []string `json:"certificate_ids"`
	}

	var temp tempRequest
	if err := ctx.BodyParser(&temp); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Convert string UUIDs to uuid.UUID type
	r.CertificateIDs = make([]uuid.UUID, len(temp.CertificateIDs))
	for i, certIDStr := range temp.CertificateIDs {
		certID, err := uuid.Parse(certIDStr)
		if err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, err)
		}
		r.CertificateIDs[i] = certID
	}

	return nil
}

func (r *RevokeEventCertificatesRequest) IsValid() error {
	if len(r.CertificateIDs) == 0 {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("at least one certificate ID is required"))
	}
	return nil
}
