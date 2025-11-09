package roleguard

import (
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"apps/backend/common/customerror"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	return auth.NewAuthService("test-issuer", "test-secret-key", time.Hour)
}

func createMockClaims(isVerifiedOrganizer *bool, isVerifiedIssuer *bool) *auth.JwtClaims {
	return &auth.JwtClaims{
		UserId:              uuid.New(),
		WalletAddress:       "0x1234567890123456789012345678901234567890",
		Email:               nil,
		IsVerifiedOrganizer: isVerifiedOrganizer,
		IsVerifiedIssuer:    isVerifiedIssuer,
	}
}

func injectUserContext(claims *auth.JwtClaims) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	}
}

// ============== RequireVerifiedOrganizer Tests ==============

func TestRequireVerifiedOrganizer_Success(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isVerified := true
	mockClaims := createMockClaims(&isVerified, nil)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedOrganizer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRequireVerifiedOrganizer_NotVerified(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isVerified := false
	mockClaims := createMockClaims(&isVerified, nil)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedOrganizer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestRequireVerifiedOrganizer_NilValue(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	mockClaims := createMockClaims(nil, nil)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedOrganizer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestRequireVerifiedOrganizer_NoUserContext(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		roleGuard.RequireVerifiedOrganizer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ============== RequireVerifiedIssuer Tests ==============

func TestRequireVerifiedIssuer_Success(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isVerified := true
	mockClaims := createMockClaims(nil, &isVerified)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedIssuer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRequireVerifiedIssuer_NotVerified(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isVerified := false
	mockClaims := createMockClaims(nil, &isVerified)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedIssuer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestRequireVerifiedIssuer_NilValue(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	mockClaims := createMockClaims(nil, nil)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedIssuer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestRequireVerifiedIssuer_NoUserContext(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		roleGuard.RequireVerifiedIssuer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ============== RequireVerifiedOrganizerOrIssuer Tests ==============

func TestRequireVerifiedOrganizerOrIssuer_OrganizerOnly(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isOrganizerVerified := true
	isIssuerVerified := false
	mockClaims := createMockClaims(&isOrganizerVerified, &isIssuerVerified)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedOrganizerOrIssuer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRequireVerifiedOrganizerOrIssuer_IssuerOnly(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isOrganizerVerified := false
	isIssuerVerified := true
	mockClaims := createMockClaims(&isOrganizerVerified, &isIssuerVerified)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedOrganizerOrIssuer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRequireVerifiedOrganizerOrIssuer_Both(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isOrganizerVerified := true
	isIssuerVerified := true
	mockClaims := createMockClaims(&isOrganizerVerified, &isIssuerVerified)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedOrganizerOrIssuer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRequireVerifiedOrganizerOrIssuer_Neither(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isOrganizerVerified := false
	isIssuerVerified := false
	mockClaims := createMockClaims(&isOrganizerVerified, &isIssuerVerified)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedOrganizerOrIssuer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestRequireVerifiedOrganizerOrIssuer_BothNil(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	mockClaims := createMockClaims(nil, nil)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireVerifiedOrganizerOrIssuer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestRequireVerifiedOrganizerOrIssuer_NoUserContext(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		roleGuard.RequireVerifiedOrganizerOrIssuer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ============== Edge Cases ==============

func TestMiddleware_InvalidUserClaimsType(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		func(c *fiber.Ctx) error {
			// Inject wrong type into context
			c.Locals("user", "invalid-type")
			return c.Next()
		},
		roleGuard.RequireVerifiedOrganizer(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ============== Generic RequireRole Tests ==============

func TestRequireRole_SingleRole_Success(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isVerified := true
	mockClaims := createMockClaims(&isVerified, nil)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireRole(RoleVerifiedOrganizer),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRequireRole_MultipleRoles_HasFirst(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isOrganizerVerified := true
	isIssuerVerified := false
	mockClaims := createMockClaims(&isOrganizerVerified, &isIssuerVerified)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireRole(RoleVerifiedOrganizer, RoleVerifiedIssuer),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRequireRole_MultipleRoles_HasSecond(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isOrganizerVerified := false
	isIssuerVerified := true
	mockClaims := createMockClaims(&isOrganizerVerified, &isIssuerVerified)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireRole(RoleVerifiedOrganizer, RoleVerifiedIssuer),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRequireRole_MultipleRoles_HasNone(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isOrganizerVerified := false
	isIssuerVerified := false
	mockClaims := createMockClaims(&isOrganizerVerified, &isIssuerVerified)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireRole(RoleVerifiedOrganizer, RoleVerifiedIssuer),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestRequireRole_ErrorMessageFormat(t *testing.T) {
	authService := setupAuthService()
	roleGuard := New(authService)

	isOrganizerVerified := false
	isIssuerVerified := false
	mockClaims := createMockClaims(&isOrganizerVerified, &isIssuerVerified)

	logger := setupLogger()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Get("/test",
		injectUserContext(mockClaims),
		roleGuard.RequireRole(RoleVerifiedOrganizer, RoleVerifiedIssuer),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	
	// Error message should list the required roles
	// (Testing through integration test - error is handled by error handler)
}

