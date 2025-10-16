package verifyjwt

import (
	"apps/backend/common/customerror"
	"apps/backend/common/log"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
)

type VerifyJwtMiddleware struct {
	authService *auth.AuthService
}

func New(authService *auth.AuthService) *VerifyJwtMiddleware {
	return &VerifyJwtMiddleware{
		authService: authService,
	}
}

func (m *VerifyJwtMiddleware) Middleware(ctx *fiber.Ctx) error {
	token := m.authService.GetJwtCookie(ctx)
	logger := log.LoadLogger()
	logger.Info("Verifying JWT", "token", token)
	if len(token) != 0 {
		claims, err := m.authService.VerifyToken(token)
		if err != nil {
			return customerror.Parse(&customerror.ErrUnauthenticated, err)
		}
		m.authService.SetUserContext(ctx, claims)
	}

	return ctx.Next()
}
