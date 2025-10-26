package issuer

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get verified issuers
// @Description Get verified issuers
// @ID get-verified-issuers
// @Tags Issuer
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} []entity.Profile
// @Failure 400 {object} customerror.Err
// @Failure 500 {object} customerror.Err
// @Router /api/v1/issuers [get]
func (h *Handler) GetVerifiedIssuers(c *fiber.Ctx) error {
	queryLimit := c.Query("limit")
	queryOffset := c.Query("offset")

	limit, err := strconv.Atoi(queryLimit)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit"})
	}
	offset, err := strconv.Atoi(queryOffset)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid offset"})
	}

	issuers, err := h.IssuerUc.GetVerifiedIssuers(c.Context(), limit, offset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(issuers)
}
