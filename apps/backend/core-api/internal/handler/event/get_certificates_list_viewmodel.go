package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetCertificatesListViewModelResponse represents the response structure for certificates list
type GetCertificatesListViewModelResponse struct {
	ClaimedCertificates   interface{} `json:"claimed_certificates"`
	UnclaimedCertificates interface{} `json:"unclaimed_certificates"`
	TotalClaimed          int         `json:"total_claimed"`
	TotalUnclaimed        int         `json:"total_unclaimed"`
}

// @Summary Get certificates list viewmodel
// @Description Get event certificates separated by claimed and unclaimed status. Claimed certificates have token_id populated, unclaimed certificates have token_id null and certificate config is published. Only event hosts and issuers can access this endpoint.
// @ID get-certificates-list-viewmodel
// @Tags Event Certificates
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} GetCertificatesListViewModelResponse
// @Failure 400 {object} customerror.ErrResponse "Invalid event ID"
// @Failure 401 {object} customerror.ErrResponse "Unauthorized - missing or invalid authentication"
// @Failure 403 {object} customerror.ErrResponse "Forbidden - user is not host or issuer of this event"
// @Failure 404 {object} customerror.ErrResponse "Event not found"
// @Failure 500 {object} customerror.ErrResponse "Internal server error"
// @Router /api/v1/events/{event_id}/certificates/list-viewmodel [get]
func (h Handler) GetCertificatesListViewModel(ctx *fiber.Ctx) error {
	// Get event ID from path parameter
	eventIDStr := ctx.Params("event_id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Get current user from JWT
	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	// Get certificates list viewmodel
	viewmodel, err := h.EventUc.GetCertificatesListViewModel(ctx.UserContext(), eventID, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(&event.CertificatesListViewModel{
		ClaimedCertificates:   viewmodel.ClaimedCertificates,
		UnclaimedCertificates: viewmodel.UnclaimedCertificates,
		TotalClaimed:          viewmodel.TotalClaimed,
		TotalUnclaimed:        viewmodel.TotalUnclaimed,
	})
}
