package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	Issuer     string
	SecretKey  string
	Expiration time.Duration
}

func NewAuthService(issuer string, secretKey string, expiration time.Duration) *AuthService {
	return &AuthService{
		Issuer:     issuer,
		SecretKey:  secretKey,
		Expiration: expiration,
	}
}

func (s *AuthService) SetUserContext(ctx *fiber.Ctx, user *JwtClaims) {
	ctx.Locals("user", user)
}

func (s *AuthService) GetUserContext(ctx *fiber.Ctx) *JwtClaims {
	return ctx.Locals("user").(*JwtClaims)
}
