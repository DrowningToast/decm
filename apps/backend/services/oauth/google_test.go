package oauth

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestMain(m *testing.M) {
	// Set required environment variables before any tests run
	_ = os.Setenv("ENVIRONMENT", "development")
	_ = os.Setenv("PII_ENCRYPTION_KEY", "test-encryption-key-32-bytes-long!!")
	_ = os.Setenv("JWT_SECRET", "test-jwt-secret-key-32-bytes-long!!")
	_ = os.Setenv("GOOGLE_OAUTH_CLIENT_ID", "test-client-id")
	_ = os.Setenv("GOOGLE_OAUTH_CLIENT_SECRET", "test-client-secret")
	_ = os.Setenv("GOOGLE_OAUTH_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback")

	// Run tests
	code := m.Run()

	// Cleanup if needed
	os.Exit(code)
}

func TestNewGoogleOAuthService(t *testing.T) {
	service := NewGoogleOAuthService()

	assert.NotNil(t, service)
	assert.NotNil(t, service.googleConfig)
	assert.NotNil(t, service.SessionStore)
	assert.NotNil(t, service.httpClient)
	assert.NotNil(t, service.logger)
}

func TestGoogleOAuthService_Login(t *testing.T) {
	service := NewGoogleOAuthService()
	app := fiber.New()

	var url *string
	app.Get("/test", func(c *fiber.Ctx) error {
		sess, err := service.SessionStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		url, err = service.Login(sess)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	assert.NotNil(t, url)
	assert.NotEmpty(t, *url)
	assert.Contains(t, *url, "accounts.google.com")
	assert.Contains(t, *url, "oauth2/auth")
}

func TestGoogleOAuthService_Login_StateSaved(t *testing.T) {
	service := NewGoogleOAuthService()
	app := fiber.New()

	var url *string
	var testErr error
	app.Get("/test", func(c *fiber.Ctx) error {
		// Use the service's session store, not a new one
		sess, err := service.SessionStore.Get(c)
		if err != nil {
			testErr = err
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		url, err = service.Login(sess)
		if err != nil {
			testErr = err
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// The Login method sets the state in the session and saves it
		// We verify this indirectly by checking that Login succeeds and generates a URL
		// The state is also encoded in the URL's state parameter
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.NoError(t, testErr)

	assert.NotNil(t, url)
	assert.NotEmpty(t, *url)
	// Verify the URL contains a state parameter (indirect verification that state was set)
	assert.Contains(t, *url, "state=", "URL should contain state parameter")
}

func TestGoogleOAuthService_Callback_StateMismatch(t *testing.T) {
	service := NewGoogleOAuthService()
	app := fiber.New()

	var callbackErr error
	app.Get("/test", func(c *fiber.Ctx) error {
		sess, err := service.SessionStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Set a state in session
		sess.Set("state", "saved-state-123")
		err = sess.Save()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Try callback with different state (should fail)
		_, callbackErr = service.Callback(context.Background(), sess, "code-123", "different-state")
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	assert.Error(t, callbackErr)
}

func TestGoogleOAuthService_validateToken(t *testing.T) {
	service := NewGoogleOAuthService()

	t.Run("valid token", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "valid-token",
			Expiry:      time.Now().Add(time.Hour),
		}

		err := service.validateToken(token)
		assert.NoError(t, err)
	})

	t.Run("nil token", func(t *testing.T) {
		err := service.validateToken(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nil")
	})

	t.Run("empty access token", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "",
			Expiry:      time.Now().Add(time.Hour),
		}

		err := service.validateToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty")
	})

	t.Run("expired token", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "token",
			Expiry:      time.Now().Add(-time.Hour),
		}

		err := service.validateToken(token)
		assert.Error(t, err)
	})

	t.Run("token expiring soon", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "token",
			Expiry:      time.Now().Add(3 * time.Minute), // Less than 5 minutes
		}

		// Should not error, but may log warning
		err := service.validateToken(token)
		assert.NoError(t, err)
	})

	t.Run("token with zero expiry", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "token",
			Expiry:      time.Time{}, // Zero time
		}

		// Should not error if expiry is zero
		err := service.validateToken(token)
		assert.NoError(t, err)
	})
}

func TestGoogleOAuthService_GetUserInfo_RequiresRealAPI(t *testing.T) {
	// This test requires actual Google OAuth API access
	// In a real scenario, you'd mock the HTTP client
	t.Skip("Skipping test that requires actual Google OAuth API")
}

func TestGoogleOAuthService_GetUserInfo_InvalidToken(t *testing.T) {
	service := NewGoogleOAuthService()

	t.Run("nil token", func(t *testing.T) {
		_, err := service.GetUserInfo(context.Background(), nil)
		assert.Error(t, err)
	})

	t.Run("expired token", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "token",
			Expiry:      time.Now().Add(-time.Hour),
		}

		_, err := service.GetUserInfo(context.Background(), token)
		assert.Error(t, err)
	})
}
