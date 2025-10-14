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

// Logout removes authentication cookies from the user's session
func (s *AuthService) Logout(ctx *fiber.Ctx) {
	// Clear session cookie
	sessionCookie := new(fiber.Cookie)
	sessionCookie.Name = "session"
	sessionCookie.Value = ""
	sessionCookie.Expires = time.Now().Add(-1 * time.Hour) // Set to past time to delete
	sessionCookie.Path = "/"
	sessionCookie.HTTPOnly = true
	sessionCookie.SameSite = "Lax"
	ctx.Cookie(sessionCookie)

	// Clear google_oauth_session cookie
	googleOAuthCookie := new(fiber.Cookie)
	googleOAuthCookie.Name = "google_oauth_session"
	googleOAuthCookie.Value = ""
	googleOAuthCookie.Expires = time.Now().Add(-1 * time.Hour) // Set to past time to delete
	googleOAuthCookie.Path = "/"
	googleOAuthCookie.HTTPOnly = true
	googleOAuthCookie.SameSite = "Lax"
	ctx.Cookie(googleOAuthCookie)
}
