package auth

import (
	"fmt"
	"time"

	"apps/backend/common"
	customerror "apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtPayload struct {
	UserId         uuid.UUID             `json:"id"`
	WalletAddress  string                `json:"wallet_address"`
	Email          *string               `json:"email,omitempty"`
	SolutionStatus common.SolutionStatus `json:"solution_status"`

	IsVerifiedOrganizer *bool `json:"is_verified_organizer,omitempty"`
	IsVerifiedIssuer    *bool `json:"is_verified_issuer,omitempty"`
}

type JwtClaims struct {
	UserId        uuid.UUID `json:"id"`
	WalletAddress string    `json:"wallet_address"`
	Email         *string   `json:"email,omitempty"`

	IsVerifiedOrganizer *bool `json:"is_verified_organizer,omitempty"`
	IsVerifiedIssuer    *bool `json:"is_verified_issuer,omitempty"`

	jwt.RegisteredClaims
}

func (s *AuthService) CreateToken(payload JwtPayload) (string, error) {
	claims := JwtClaims{
		UserId:              payload.UserId,
		WalletAddress:       payload.WalletAddress,
		Email:               payload.Email,
		IsVerifiedOrganizer: payload.IsVerifiedOrganizer,
		IsVerifiedIssuer:    payload.IsVerifiedIssuer,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return tokenString, nil
}

func (s *AuthService) VerifyToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	// casting
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *AuthService) SetJwtCookie(ctx *fiber.Ctx, token string) {
	cookie := new(fiber.Cookie)
	cookie.Name = "themis-session"
	cookie.Value = token
	cookie.Expires = time.Now().Add(s.Expiration)
	cookie.Path = "/"
	cookie.HTTPOnly = true
	cookie.SameSite = "Lax" // Allows cookies to work across different ports on localhost
	// Note: In production with HTTPS, use SameSite="None" and Secure=true
	if s.CookieDomain != "" {
		cookie.Domain = s.CookieDomain
	}
	ctx.Cookie(cookie)
}

func (s *AuthService) GetJwtCookie(ctx *fiber.Ctx) string {
	return ctx.Cookies("themis-session")
}
