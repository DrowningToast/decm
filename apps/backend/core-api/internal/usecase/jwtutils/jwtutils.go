package jwtutils

import (
	"fmt"
	"time"

	customerror "apps/backend/common/customerror"
	"apps/backend/core-api/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	cfg    = config.LoadConfig()
	issuer = "decm-service"
)

type JwtPayload struct {
	UserId        uuid.UUID `json:"id"`
	WalletAddress string    `json:"wallet_address"`
}

type JwtClaims struct {
	UserId        uuid.UUID `json:"id"`
	WalletAddress string    `json:"wallet_address"`
	jwt.RegisteredClaims
}

func CreateToken(payload JwtPayload) (string, error) {
	claims := JwtClaims{
		UserId:        payload.UserId,
		WalletAddress: payload.WalletAddress,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Jwt.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.Jwt.SecretKey))
	if err != nil {
		return "", customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Jwt.SecretKey), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
