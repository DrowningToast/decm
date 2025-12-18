package issuer

import (
	"strconv"

	issuer_usecase "apps/backend/core-api/internal/usecase/issuer"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get verified issuers
// @Description Get verified issuers with optional search query
// @ID get-verified-issuers
// @Tags Issuer
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param search query string false "Search query (searches first name, last name, email, academic email, wallet address)"
// @Success 200 {object} []entity.Profile
// @Failure 400 {object} customerror.Err
// @Failure 500 {object} customerror.Err
// @Router /api/v1/issuers [get]
func (h *Handler) GetVerifiedIssuers(c *fiber.Ctx) error {
	queryLimit := c.Query("limit", "25")
	queryOffset := c.Query("offset", "0")
	searchQuery := c.Query("search")

	limit, err := strconv.Atoi(queryLimit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid limit")
	}
	offset, err := strconv.Atoi(queryOffset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid offset")
	}

	issuers, err := h.IssuerUc.GetVerifiedIssuers(c.Context(), issuer_usecase.GetVerifiedIssuersRequest{
		SearchQuery: searchQuery,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(issuers)
}
