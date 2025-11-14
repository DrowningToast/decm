package authenticationguard

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"apps/backend/common/customerror"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors in tests
	}))
}

func setupTestApp(middleware fiber.Handler) *fiber.App {
	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		middleware,
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "success",
			})
		},
	)

	return app
}

func setupAuthService() *auth.AuthService {
	return auth.NewAuthService("test-issuer", "test-secret-key", time.Hour, "test-cookie-domain")
}

func createValidJWT(authService *auth.AuthService) (string, *auth.JwtClaims) {
	userId := uuid.New()
	payload := auth.JwtPayload{
		UserId:              userId,
		WalletAddress:       "0x1234567890123456789012345678901234567890",
		Email:               nil,
		IsVerifiedOrganizer: nil,
		IsVerifiedIssuer:    nil,
	}

	token, _ := authService.CreateToken(payload)

	// Create expected claims for comparison
	claims := &auth.JwtClaims{
		UserId:              userId,
		WalletAddress:       payload.WalletAddress,
		Email:               payload.Email,
		IsVerifiedOrganizer: payload.IsVerifiedOrganizer,
		IsVerifiedIssuer:    payload.IsVerifiedIssuer,
	}

	return token, claims
}

// ============== Authentication Guard Middleware Tests ==============

func TestAuthenticationGuard_ValidToken(t *testing.T) {
	authService := setupAuthService()
	token, _ := createValidJWT(authService)

	middleware := New(authService)
	app := setupTestApp(middleware.Middleware)

	// Create request with token in cookie
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAuthenticationGuard_ValidToken_WithUserContext(t *testing.T) {
	authService := setupAuthService()
	token, _ := createValidJWT(authService)

	middleware := New(authService)
	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		middleware.Middleware,
		func(c *fiber.Ctx) error {
			// Verify user context is set
			user := c.Locals("user")
			require.NotNil(t, user)

			claims, ok := user.(*auth.JwtClaims)
			require.True(t, ok)
			require.NotNil(t, claims)
			require.NotZero(t, claims.UserId)

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"userId": claims.UserId.String(),
			})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAuthenticationGuard_MissingToken(t *testing.T) {
	authService := setupAuthService()
	middleware := New(authService)
	app := setupTestApp(middleware.Middleware)

	req := httptest.NewRequest("GET", "/test", nil)
	// No cookie added

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticationGuard_EmptyToken(t *testing.T) {
	authService := setupAuthService()
	middleware := New(authService)
	app := setupTestApp(middleware.Middleware)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: "",
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticationGuard_InvalidToken(t *testing.T) {
	authService := setupAuthService()
	middleware := New(authService)
	app := setupTestApp(middleware.Middleware)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: "invalid-token-format",
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticationGuard_MalformedToken(t *testing.T) {
	authService := setupAuthService()
	middleware := New(authService)
	app := setupTestApp(middleware.Middleware)

	// Create a token signed with different key
	differentAuthService := auth.NewAuthService("different-issuer", "different-secret-key", time.Hour, "test-cookie-domain")
	token, _ := createValidJWT(differentAuthService)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticationGuard_ExpiredToken(t *testing.T) {
	// Create auth service with very short expiration
	shortLivedAuthService := auth.NewAuthService("test-issuer", "test-secret-key", -time.Hour, "test-cookie-domain")
	token, _ := createValidJWT(shortLivedAuthService)

	// Use the original auth service to verify (which has different timing)
	authService := setupAuthService()
	middleware := New(authService)
	app := setupTestApp(middleware.Middleware)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticationGuard_TokenFromDifferentCookie(t *testing.T) {
	authService := setupAuthService()
	token, _ := createValidJWT(authService)

	middleware := New(authService)
	app := setupTestApp(middleware.Middleware)

	req := httptest.NewRequest("GET", "/test", nil)
	// Add token to wrong cookie name
	req.AddCookie(&http.Cookie{
		Name:  "wrong_cookie_name",
		Value: token,
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticationGuard_MultipleRequests(t *testing.T) {
	authService := setupAuthService()
	token1, _ := createValidJWT(authService)
	token2, _ := createValidJWT(authService)

	middleware := New(authService)
	app := setupTestApp(middleware.Middleware)

	// First request
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token1,
	})
	resp1, err := app.Test(req1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

	// Second request
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token2,
	})
	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
}

func TestAuthenticationGuard_ContextPropagation(t *testing.T) {
	authService := setupAuthService()
	token, expectedClaims := createValidJWT(authService)

	middleware := New(authService)
	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		middleware.Middleware,
		func(c *fiber.Ctx) error {
			user := c.Locals("user")
			if user == nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "user context not set",
				})
			}

			claims, ok := user.(*auth.JwtClaims)
			if !ok {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "user context has wrong type",
				})
			}

			// Verify the claims match what we expect
			if claims.UserId != expectedClaims.UserId {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "user ID mismatch",
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "success",
			})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAuthenticationGuard_TokenWithSpecialCharacters(t *testing.T) {
	authService := setupAuthService()
	middleware := New(authService)
	app := setupTestApp(middleware.Middleware)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9!@#$%^&*()",
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticationGuard_ContinuesOnSuccess(t *testing.T) {
	authService := setupAuthService()
	token, _ := createValidJWT(authService)

	middleware := New(authService)
	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	handlerExecuted := false

	app.Get("/test",
		middleware.Middleware,
		func(c *fiber.Ctx) error {
			handlerExecuted = true
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "success",
			})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	// Verify handler was executed
	assert.True(t, handlerExecuted)
}
