package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ErrorHandler is a custom error handler for Fiber
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Default error code
	code := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log error for debugging
	fmt.Printf("Error: %v\n", err)

	// Return error response
	return ctx.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
		"code":    code,
	})
}

// RequestID middleware adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate unique request ID
		requestID := uuid.New().String()

		// Set request ID in context
		c.Set("X-Request-ID", requestID)
		c.Locals("requestID", requestID)

		return c.Next()
	}
}

// AcademicAuth middleware for academic identity verification
func AcademicAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for academic credentials/tokens
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Academic authentication required",
			})
		}

		// TODO: Implement academic identity verification logic
		// This would integrate with LDAP or institutional systems

		return c.Next()
	}
}

// BlockchainAuth middleware for blockchain/Web3 authentication
func BlockchainAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for blockchain wallet signature
		signature := c.Get("X-Wallet-Signature")
		address := c.Get("X-Wallet-Address")

		if signature == "" || address == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Blockchain authentication required",
			})
		}

		// TODO: Implement wallet signature verification
		// This would verify the wallet signature for Web3 auth

		return c.Next()
	}
}

// RoleAuth middleware for role-based access control
func RoleAuth(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by auth middleware)
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Authentication required",
			})
		}

		role := userRole.(string)

		// Check if user has required role
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Insufficient permissions",
		})
	}
}
