package auth

import (
	"errors"
	"time"

	"apps/backend/common/customerror"

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

func (s *AuthService) GetUserContext(ctx *fiber.Ctx) (*JwtClaims, error) {
	user := ctx.Locals("user")
	if user == nil {
		return nil, nil
	}

	return s.RequireUserContext(ctx)
}

func (s *AuthService) RequireUserContext(ctx *fiber.Ctx) (*JwtClaims, error) {
	user := ctx.Locals("user")
	if user == nil {
		return nil, customerror.AsPresetError(customerror.ErrUnauthenticated, errors.New("user not found"))
	}
	// check parsing
	userClaims, ok := user.(*JwtClaims)
	if !ok {
		return nil, customerror.AsPresetError(customerror.ErrUnauthenticated, errors.New("invalid user claims"))
	}
	return userClaims, nil
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
