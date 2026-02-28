package roleguard

import (
	"apps/backend/common/customerror"
	"apps/backend/services/auth"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Role represents a user role that can be checked
type Role string

const (
	RoleVerifiedOrganizer Role = "verified_organizer"
	RoleVerifiedIssuer    Role = "verified_issuer"
)

type RoleGuardMiddleware struct {
	authService *auth.AuthService
}

func New(authService *auth.AuthService) *RoleGuardMiddleware {
	return &RoleGuardMiddleware{
		authService: authService,
	}
}

// RequireRole checks if the authenticated user has at least one of the specified roles
// User must have ONE OR MORE of the roles in the array to pass
func (m *RoleGuardMiddleware) RequireRole(roles ...Role) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims, err := m.authService.RequireUserContext(ctx)
		if err != nil {
			return err
		}

		// Check if user has at least one of the required roles
		hasRole := false
		for _, role := range roles {
			switch role {
			case RoleVerifiedOrganizer:
				if claims.IsVerifiedOrganizer != nil && *claims.IsVerifiedOrganizer {
					hasRole = true
				}
			case RoleVerifiedIssuer:
				if claims.IsVerifiedIssuer != nil && *claims.IsVerifiedIssuer {
					hasRole = true
				}
			}

			if hasRole {
				break
			}
		}

		if !hasRole {
			// Build error message with required roles
			roleNames := make([]string, len(roles))
			for i, role := range roles {
				roleNames[i] = string(role)
			}
			errorMsg := fmt.Sprintf("user does not have required role. Required one of: %s", strings.Join(roleNames, ", "))

			return customerror.Parse(
				&customerror.ErrForbidden,
				errors.New(errorMsg),
			)
		}

		return ctx.Next()
	}
}

// RequireVerifiedOrganizer checks if the authenticated user is a verified organizer
// Convenience method for RequireRole(RoleVerifiedOrganizer)
func (m *RoleGuardMiddleware) RequireVerifiedOrganizer() fiber.Handler {
	return m.RequireRole(RoleVerifiedOrganizer)
}

// RequireVerifiedIssuer checks if the authenticated user is a verified issuer
// Convenience method for RequireRole(RoleVerifiedIssuer)
func (m *RoleGuardMiddleware) RequireVerifiedIssuer() fiber.Handler {
	return m.RequireRole(RoleVerifiedIssuer)
}

// RequireVerifiedOrganizerOrIssuer checks if the authenticated user is either a verified organizer or issuer
// Convenience method for RequireRole(RoleVerifiedOrganizer, RoleVerifiedIssuer)
func (m *RoleGuardMiddleware) RequireVerifiedOrganizerOrIssuer() fiber.Handler {
	return m.RequireRole(RoleVerifiedOrganizer, RoleVerifiedIssuer)
}
