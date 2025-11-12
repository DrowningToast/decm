package auth

import (
	"apps/backend/common/customerror"
	auth_usecase "apps/backend/core-api/internal/usecase/auth"

	"github.com/gofiber/fiber/v2"
)

// CheckRoleResponse represents the response for role checking
// @description Role verification response
type CheckRoleResponse struct {
	IsAuthenticated *bool `json:"is_authenticated,omitempty"`
	IsHost          *bool `json:"is_host,omitempty"`
	IsIssuer        *bool `json:"is_issuer,omitempty"`
} // @name CheckRoleResponse

// @Summary Check user roles
// @Tags Auth
// @Description Check if the current user has specific roles (authenticated, host, issuer). Only returns requested fields.
// @ID check-role
// @Accept json
// @Produce json
// @Param is_authenticated query boolean false "Check if user is authenticated"
// @Param is_host query boolean false "Check if user is a verified host/organizer"
// @Param is_issuer query boolean false "Check if user is a verified issuer"
// @Success 200 {object} CheckRoleResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/auth/check-role [get]
func (h Handler) CheckRole(ctx *fiber.Ctx) error {
	// Parse query parameters
	checkAuthenticated := ctx.Query("is_authenticated") == "true"
	checkHost := ctx.Query("is_host") == "true"
	checkIssuer := ctx.Query("is_issuer") == "true"

	// Get user from context (may be nil if not authenticated)
	user, _ := h.AuthService.GetUserContext(ctx)

	// Prepare usecase parameters with full JWT claims
	// This allows the usecase to validate token expiration and access all claims
	params := auth_usecase.CheckRoleParams{
		JwtClaims:   user, // Pass entire JWT claims, not just user ID
		CheckAuth:   checkAuthenticated,
		CheckHost:   checkHost,
		CheckIssuer: checkIssuer,
	}

	// Call usecase to perform business logic
	result, err := h.AuthUc.CheckRole(ctx.UserContext(), params)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Map usecase result to response
	response := CheckRoleResponse{
		IsAuthenticated: result.IsAuthenticated,
		IsHost:          result.IsHost,
		IsIssuer:        result.IsIssuer,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
