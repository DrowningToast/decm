package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"apps/backend/common"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_CreateToken(t *testing.T) {
	service := NewAuthService("test-issuer", "test-secret-key-123456789012345678901234567890", time.Hour, "")

	payload := JwtPayload{
		UserId:              uuid.New(),
		WalletAddress:       "0x1234567890123456789012345678901234567890",
		Email:               stringPtr("test@example.com"),
		SolutionStatus:      common.SolutionStatusBYOK,
		IsVerifiedOrganizer: boolPtr(true),
		IsVerifiedIssuer:    boolPtr(false),
	}

	token, err := service.CreateToken(payload)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_CreateToken_VerifyToken_RoundTrip(t *testing.T) {
	service := NewAuthService("test-issuer", "test-secret-key-123456789012345678901234567890", time.Hour, "")

	userId := uuid.New()
	payload := JwtPayload{
		UserId:              userId,
		WalletAddress:       "0x1234567890123456789012345678901234567890",
		Email:               stringPtr("test@example.com"),
		SolutionStatus:      common.SolutionStatusBYOK,
		IsVerifiedOrganizer: boolPtr(true),
		IsVerifiedIssuer:    boolPtr(false),
	}

	token, err := service.CreateToken(payload)
	require.NoError(t, err)

	claims, err := service.VerifyToken(token)
	require.NoError(t, err)
	assert.Equal(t, userId, claims.UserId)
	assert.Equal(t, payload.WalletAddress, claims.WalletAddress)
	assert.Equal(t, payload.Email, claims.Email)
	assert.Equal(t, payload.IsVerifiedOrganizer, claims.IsVerifiedOrganizer)
	assert.Equal(t, payload.IsVerifiedIssuer, claims.IsVerifiedIssuer)
}

func TestAuthService_VerifyToken_InvalidToken(t *testing.T) {
	service := NewAuthService("test-issuer", "test-secret-key-123456789012345678901234567890", time.Hour, "")

	_, err := service.VerifyToken("invalid-token")
	assert.Error(t, err)
}

func TestAuthService_VerifyToken_WrongSecret(t *testing.T) {
	service1 := NewAuthService("test-issuer", "secret-key-1-123456789012345678901234567890", time.Hour, "")
	service2 := NewAuthService("test-issuer", "secret-key-2-123456789012345678901234567890", time.Hour, "")

	payload := JwtPayload{
		UserId:         uuid.New(),
		WalletAddress:  "0x1234567890123456789012345678901234567890",
		SolutionStatus: common.SolutionStatusBYOK,
	}

	token, err := service1.CreateToken(payload)
	require.NoError(t, err)

	_, err = service2.VerifyToken(token)
	assert.Error(t, err)
}

func TestAuthService_VerifyToken_ExpiredToken(t *testing.T) {
	service := NewAuthService("test-issuer", "test-secret-key-123456789012345678901234567890", -time.Hour, "")

	payload := JwtPayload{
		UserId:         uuid.New(),
		WalletAddress:  "0x1234567890123456789012345678901234567890",
		SolutionStatus: common.SolutionStatusBYOK,
	}

	token, err := service.CreateToken(payload)
	require.NoError(t, err)

	_, err = service.VerifyToken(token)
	assert.Error(t, err)
}

func TestAuthService_SetJwtCookie(t *testing.T) {
	service := NewAuthService("test-issuer", "secret", time.Hour, "example.com")
	app := fiber.New()

	token := "test-token"
	app.Get("/test", func(c *fiber.Ctx) error {
		service.SetJwtCookie(c, token)
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Verify cookie is set in response
	cookies := resp.Header.Values("Set-Cookie")
	assert.NotEmpty(t, cookies)
	assert.Contains(t, cookies[0], "themis-session="+token)
}

func TestAuthService_SetJwtCookie_NoDomain(t *testing.T) {
	service := NewAuthService("test-issuer", "secret", time.Hour, "")
	app := fiber.New()

	token := "test-token"
	app.Get("/test", func(c *fiber.Ctx) error {
		service.SetJwtCookie(c, token)
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Verify cookie is set in response
	cookies := resp.Header.Values("Set-Cookie")
	assert.NotEmpty(t, cookies)
	assert.Contains(t, cookies[0], "themis-session="+token)
}

func TestAuthService_GetJwtCookie(t *testing.T) {
	service := NewAuthService("test-issuer", "secret", time.Hour, "")
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		cookieValue := service.GetJwtCookie(c)
		assert.Equal(t, "test-token-value", cookieValue)
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "themis-session",
		Value: "test-token-value",
	})
	_, err := app.Test(req)
	assert.NoError(t, err)
}

func TestAuthService_GetJwtCookie_NoCookie(t *testing.T) {
	service := NewAuthService("test-issuer", "secret", time.Hour, "")
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		cookieValue := service.GetJwtCookie(c)
		assert.Empty(t, cookieValue)
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	_, err := app.Test(req)
	assert.NoError(t, err)
}

func TestAuthService_CreateToken_OptionalFields(t *testing.T) {
	service := NewAuthService("test-issuer", "test-secret-key-123456789012345678901234567890", time.Hour, "")

	t.Run("with all optional fields", func(t *testing.T) {
		payload := JwtPayload{
			UserId:              uuid.New(),
			WalletAddress:       "0x1234567890123456789012345678901234567890",
			Email:               stringPtr("test@example.com"),
			SolutionStatus:      common.SolutionStatusBYOK,
			IsVerifiedOrganizer: boolPtr(true),
			IsVerifiedIssuer:    boolPtr(true),
		}

		token, err := service.CreateToken(payload)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("without optional fields", func(t *testing.T) {
		payload := JwtPayload{
			UserId:         uuid.New(),
			WalletAddress:  "0x1234567890123456789012345678901234567890",
			SolutionStatus: common.SolutionStatusBYOK,
		}

		token, err := service.CreateToken(payload)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)
		assert.Nil(t, claims.Email)
		assert.Nil(t, claims.IsVerifiedOrganizer)
		assert.Nil(t, claims.IsVerifiedIssuer)
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}






