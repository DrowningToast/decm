package auth

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthService_GetUserContext_WithRoles tests GetUserContext with different role combinations
func TestAuthService_GetUserContext_WithRoles(t *testing.T) {
	service := NewAuthService("test", "secret", time.Hour, "")
	app := fiber.New()

	tests := []struct {
		name                string
		isVerifiedOrganizer *bool
		isVerifiedIssuer    *bool
		description         string
	}{
		{
			name:                "regular user",
			isVerifiedOrganizer: nil,
			isVerifiedIssuer:    nil,
			description:         "User with no special roles",
		},
		{
			name:                "verified organizer",
			isVerifiedOrganizer: boolPtr(true),
			isVerifiedIssuer:    nil,
			description:         "User with organizer role only",
		},
		{
			name:                "verified issuer",
			isVerifiedOrganizer: nil,
			isVerifiedIssuer:    boolPtr(true),
			description:         "User with issuer role only",
		},
		{
			name:                "both roles",
			isVerifiedOrganizer: boolPtr(true),
			isVerifiedIssuer:    boolPtr(true),
			description:         "User with both organizer and issuer roles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &JwtClaims{
				UserId:              uuid.New(),
				WalletAddress:       "0x123",
				IsVerifiedOrganizer: tt.isVerifiedOrganizer,
				IsVerifiedIssuer:    tt.isVerifiedIssuer,
			}

			app.Get("/test", func(c *fiber.Ctx) error {
				c.Locals("user", user)
				retrieved, err := service.GetUserContext(c)
				require.NoError(t, err)
				assert.Equal(t, user, retrieved)
				assert.Equal(t, tt.isVerifiedOrganizer, retrieved.IsVerifiedOrganizer)
				assert.Equal(t, tt.isVerifiedIssuer, retrieved.IsVerifiedIssuer)
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			_, err := app.Test(req)
			assert.NoError(t, err)
		})
	}
}

// TestAuthService_RequireUserContext_WithRoles tests RequireUserContext with different role combinations
func TestAuthService_RequireUserContext_WithRoles(t *testing.T) {
	service := NewAuthService("test", "secret", time.Hour, "")
	app := fiber.New()

	tests := []struct {
		name                string
		isVerifiedOrganizer *bool
		isVerifiedIssuer    *bool
		description         string
	}{
		{
			name:                "regular user",
			isVerifiedOrganizer: nil,
			isVerifiedIssuer:    nil,
			description:         "User with no special roles",
		},
		{
			name:                "verified organizer",
			isVerifiedOrganizer: boolPtr(true),
			isVerifiedIssuer:    nil,
			description:         "User with organizer role only",
		},
		{
			name:                "verified issuer",
			isVerifiedOrganizer: nil,
			isVerifiedIssuer:    boolPtr(true),
			description:         "User with issuer role only",
		},
		{
			name:                "both roles",
			isVerifiedOrganizer: boolPtr(true),
			isVerifiedIssuer:    boolPtr(true),
			description:         "User with both organizer and issuer roles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &JwtClaims{
				UserId:              uuid.New(),
				WalletAddress:       "0x123",
				IsVerifiedOrganizer: tt.isVerifiedOrganizer,
				IsVerifiedIssuer:    tt.isVerifiedIssuer,
			}

			app.Get("/test", func(c *fiber.Ctx) error {
				c.Locals("user", user)
				retrieved, err := service.RequireUserContext(c)
				require.NoError(t, err)
				assert.Equal(t, user, retrieved)
				assert.Equal(t, tt.isVerifiedOrganizer, retrieved.IsVerifiedOrganizer)
				assert.Equal(t, tt.isVerifiedIssuer, retrieved.IsVerifiedIssuer)
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			_, err := app.Test(req)
			assert.NoError(t, err)
		})
	}
}

// TestAuthService_SetUserContext_WithRoles tests SetUserContext preserves all role information
func TestAuthService_SetUserContext_WithRoles(t *testing.T) {
	service := NewAuthService("test", "secret", time.Hour, "")
	app := fiber.New()

	// Test with organizer role
	t.Run("organizer role preserved", func(t *testing.T) {
		user := &JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x123",
			IsVerifiedOrganizer: boolPtr(true),
			IsVerifiedIssuer:    nil,
		}

		app.Get("/test", func(c *fiber.Ctx) error {
			service.SetUserContext(c, user)
			retrieved := c.Locals("user").(*JwtClaims)
			assert.True(t, *retrieved.IsVerifiedOrganizer)
			assert.Nil(t, retrieved.IsVerifiedIssuer)
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	// Test with issuer role
	t.Run("issuer role preserved", func(t *testing.T) {
		user := &JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x123",
			IsVerifiedOrganizer: nil,
			IsVerifiedIssuer:    boolPtr(true),
		}

		app.Get("/test2", func(c *fiber.Ctx) error {
			service.SetUserContext(c, user)
			retrieved := c.Locals("user").(*JwtClaims)
			assert.Nil(t, retrieved.IsVerifiedOrganizer)
			assert.True(t, *retrieved.IsVerifiedIssuer)
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test2", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	// Test with both roles
	t.Run("both roles preserved", func(t *testing.T) {
		user := &JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x123",
			IsVerifiedOrganizer: boolPtr(true),
			IsVerifiedIssuer:    boolPtr(true),
		}

		app.Get("/test3", func(c *fiber.Ctx) error {
			service.SetUserContext(c, user)
			retrieved := c.Locals("user").(*JwtClaims)
			assert.True(t, *retrieved.IsVerifiedOrganizer)
			assert.True(t, *retrieved.IsVerifiedIssuer)
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test3", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	// Test with no roles (regular user)
	t.Run("regular user - no roles", func(t *testing.T) {
		user := &JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x123",
			IsVerifiedOrganizer: nil,
			IsVerifiedIssuer:    nil,
		}

		app.Get("/test4", func(c *fiber.Ctx) error {
			service.SetUserContext(c, user)
			retrieved := c.Locals("user").(*JwtClaims)
			assert.Nil(t, retrieved.IsVerifiedOrganizer)
			assert.Nil(t, retrieved.IsVerifiedIssuer)
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test4", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})
}




