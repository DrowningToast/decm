package authenticationguard

import (
	"errors"

	"apps/backend/common/customerror"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
)

type AuthenticationGuardMiddleware struct {
	authService *auth.AuthService
}

func NewMiddleware(authService *auth.AuthService) *AuthenticationGuardMiddleware {
	return &AuthenticationGuardMiddleware{
		authService: authService,
	}
}

func (m *AuthenticationGuardMiddleware) Middleware(ctx *fiber.Ctx) error {
	token := m.authService.GetJwtCookie(ctx)
	if token == "" {
		return customerror.Parse(&customerror.ErrUnauthenticated, errors.New("token is required"))
	}

	claims, err := m.authService.VerifyToken(token)
	if err != nil {
		return customerror.Parse(&customerror.ErrUnauthenticated, err)
	}

	// set local of request
	ctx.Locals("user", claims)

	return ctx.Next()
}
