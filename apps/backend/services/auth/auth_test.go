package auth

import (
	"apps/backend/common/customerror"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthService(t *testing.T) {
	issuer := "test-issuer"
	secretKey := "test-secret-key"
	expiration := time.Hour
	cookieDomain := "example.com"

	service := NewAuthService(issuer, secretKey, expiration, cookieDomain)

	assert.NotNil(t, service)
	assert.Equal(t, issuer, service.Issuer)
	assert.Equal(t, secretKey, service.SecretKey)
	assert.Equal(t, expiration, service.Expiration)
	assert.Equal(t, cookieDomain, service.CookieDomain)
}

func TestAuthService_SetUserContext(t *testing.T) {
	service := NewAuthService("test", "secret", time.Hour, "")
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		user := &JwtClaims{
			UserId:        uuid.New(),
			WalletAddress: "0x123",
		}
		service.SetUserContext(c, user)

		retrieved := c.Locals("user")
		assert.NotNil(t, retrieved)
		assert.Equal(t, user, retrieved)
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	_, err := app.Test(req)
	assert.NoError(t, err)
}

func TestAuthService_GetUserContext(t *testing.T) {
	service := NewAuthService("test", "secret", time.Hour, "")
	app := fiber.New()

	t.Run("user exists", func(t *testing.T) {
		user := &JwtClaims{
			UserId:        uuid.New(),
			WalletAddress: "0x123",
		}

		app.Get("/test1", func(c *fiber.Ctx) error {
			c.Locals("user", user)
			retrieved, err := service.GetUserContext(c)
			require.NoError(t, err)
			assert.Equal(t, user, retrieved)
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test1", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	t.Run("user does not exist", func(t *testing.T) {
		app.Get("/test2", func(c *fiber.Ctx) error {
			retrieved, err := service.GetUserContext(c)
			assert.NoError(t, err)
			assert.Nil(t, retrieved)
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test2", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})
}

func TestAuthService_RequireUserContext(t *testing.T) {
	service := NewAuthService("test", "secret", time.Hour, "")
	app := fiber.New()

	t.Run("user exists", func(t *testing.T) {
		user := &JwtClaims{
			UserId:        uuid.New(),
			WalletAddress: "0x123",
		}

		app.Get("/test1", func(c *fiber.Ctx) error {
			c.Locals("user", user)
			retrieved, err := service.RequireUserContext(c)
			require.NoError(t, err)
			assert.Equal(t, user, retrieved)
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test1", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	t.Run("user does not exist", func(t *testing.T) {
		app.Get("/test2", func(c *fiber.Ctx) error {
			retrieved, err := service.RequireUserContext(c)
			assert.Error(t, err)
			assert.Nil(t, retrieved)
			assert.IsType(t, &customerror.Err{}, err)
			return c.SendStatus(fiber.StatusUnauthorized)
		})

		req := httptest.NewRequest("GET", "/test2", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	t.Run("invalid user type", func(t *testing.T) {
		app.Get("/test3", func(c *fiber.Ctx) error {
			c.Locals("user", "invalid-type")
			retrieved, err := service.RequireUserContext(c)
			assert.Error(t, err)
			assert.Nil(t, retrieved)
			assert.IsType(t, &customerror.Err{}, err)
			return c.SendStatus(fiber.StatusUnauthorized)
		})

		req := httptest.NewRequest("GET", "/test3", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})
}

func TestAuthService_Logout(t *testing.T) {
	service := NewAuthService("test", "secret", time.Hour, "example.com")
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		service.Logout(c)
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Verify cookies are set with past expiration
	cookies := resp.Header.Values("Set-Cookie")
	assert.NotEmpty(t, cookies)

	// Check that both cookies are cleared
	cookieStr := cookies[0]
	assert.Contains(t, cookieStr, "themis-session")
}

func TestAuthService_Logout_NoDomain(t *testing.T) {
	service := NewAuthService("test", "secret", time.Hour, "")
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		service.Logout(c)
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	cookies := resp.Header.Values("Set-Cookie")
	assert.NotEmpty(t, cookies)
}
