package certificate

import (
	"apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
)

// GetMyCertificatesListViewModelResponse represents the response structure for current user's certificates list
type GetMyCertificatesListViewModelResponse struct {
	ClaimedCertificates   interface{} `json:"claimed_certificates"`
	UnclaimedCertificates interface{} `json:"unclaimed_certificates"`
	TotalClaimed          int         `json:"total_claimed"`
	TotalUnclaimed        int         `json:"total_unclaimed"`
}

// @Summary Get my certificates list viewmodel
// @Description Get current user's certificates separated by claimed and unclaimed status. Claimed certificates have token_id populated, unclaimed certificates have token_id null and certificate config is published. Returns all certificates for the authenticated user across all events.
// @ID get-my-certificates-list-viewmodel
// @Tags Certificates
// @Accept json
// @Produce json
// @Success 200 {object} GetMyCertificatesListViewModelResponse
// @Failure 401 {object} customerror.ErrResponse "Unauthorized - missing or invalid authentication"
// @Failure 500 {object} customerror.ErrResponse "Internal server error"
// @Router /api/v1/certificates/my-list-viewmodel [get]
func (h Handler) GetMyCertificatesListViewModel(ctx *fiber.Ctx) error {
	// Get current user from JWT
	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	// Get user's certificates list viewmodel
	viewmodel, err := h.EventUc.GetMyCertificatesListViewModel(ctx.UserContext(), currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(&event.MyCertificatesListViewModel{
		ClaimedCertificates:   viewmodel.ClaimedCertificates,
		UnclaimedCertificates: viewmodel.UnclaimedCertificates,
		TotalClaimed:          viewmodel.TotalClaimed,
		TotalUnclaimed:        viewmodel.TotalUnclaimed,
	})
}
