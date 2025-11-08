package auth

import (
	"context"
	"testing"
	"time"

	"apps/backend/services/auth"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthUsecase_CheckRole_Authentication(t *testing.T) {
	t.Run("should return true when user has valid token", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		validClaims := &auth.JwtClaims{
			UserId:        uuid.New(),
			WalletAddress: "0x1234567890",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		params := CheckRoleParams{
			JwtClaims: validClaims,
			CheckAuth: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsAuthenticated)
		assert.True(t, *result.IsAuthenticated)
		assert.Nil(t, result.IsHost)
		assert.Nil(t, result.IsIssuer)
	})

	t.Run("should return false when token is expired", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		expiredClaims := &auth.JwtClaims{
			UserId:        uuid.New(),
			WalletAddress: "0x1234567890",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			},
		}

		params := CheckRoleParams{
			JwtClaims: expiredClaims,
			CheckAuth: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsAuthenticated)
		assert.False(t, *result.IsAuthenticated)
	})

	t.Run("should return false when claims are nil", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()

		params := CheckRoleParams{
			JwtClaims: nil,
			CheckAuth: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsAuthenticated)
		assert.False(t, *result.IsAuthenticated)
	})
}

func TestAuthUsecase_CheckRole_HostRole(t *testing.T) {
	t.Run("should return true when user is verified organizer", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		isVerifiedOrganizer := true
		validClaims := &auth.JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x1234567890",
			IsVerifiedOrganizer: &isVerifiedOrganizer,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		params := CheckRoleParams{
			JwtClaims: validClaims,
			CheckHost: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsHost)
		assert.True(t, *result.IsHost)
		assert.Nil(t, result.IsAuthenticated)
		assert.Nil(t, result.IsIssuer)
	})

	t.Run("should return false when user is not verified organizer", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		isVerifiedOrganizer := false
		validClaims := &auth.JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x1234567890",
			IsVerifiedOrganizer: &isVerifiedOrganizer,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		params := CheckRoleParams{
			JwtClaims: validClaims,
			CheckHost: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsHost)
		assert.False(t, *result.IsHost)
	})

	t.Run("should return false when IsVerifiedOrganizer is nil in claims", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		validClaims := &auth.JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x1234567890",
			IsVerifiedOrganizer: nil, // Not set in JWT
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		params := CheckRoleParams{
			JwtClaims: validClaims,
			CheckHost: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsHost)
		assert.False(t, *result.IsHost) // Should default to false
	})

	t.Run("should return false when token is expired", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		isVerifiedOrganizer := true
		expiredClaims := &auth.JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x1234567890",
			IsVerifiedOrganizer: &isVerifiedOrganizer,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired
			},
		}

		params := CheckRoleParams{
			JwtClaims: expiredClaims,
			CheckHost: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsHost)
		assert.False(t, *result.IsHost) // Should be false due to expired token
	})

	t.Run("should return false when claims are nil", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()

		params := CheckRoleParams{
			JwtClaims: nil,
			CheckHost: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsHost)
		assert.False(t, *result.IsHost)
	})
}

func TestAuthUsecase_CheckRole_IssuerRole(t *testing.T) {
	t.Run("should return true when user is verified issuer", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		isVerifiedIssuer := true
		validClaims := &auth.JwtClaims{
			UserId:           uuid.New(),
			WalletAddress:    "0x1234567890",
			IsVerifiedIssuer: &isVerifiedIssuer,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		params := CheckRoleParams{
			JwtClaims:   validClaims,
			CheckIssuer: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsIssuer)
		assert.True(t, *result.IsIssuer)
		assert.Nil(t, result.IsAuthenticated)
		assert.Nil(t, result.IsHost)
	})

	t.Run("should return false when user is not verified issuer", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		isVerifiedIssuer := false
		validClaims := &auth.JwtClaims{
			UserId:           uuid.New(),
			WalletAddress:    "0x1234567890",
			IsVerifiedIssuer: &isVerifiedIssuer,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		params := CheckRoleParams{
			JwtClaims:   validClaims,
			CheckIssuer: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsIssuer)
		assert.False(t, *result.IsIssuer)
	})

	t.Run("should return false when IsVerifiedIssuer is nil in claims", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		validClaims := &auth.JwtClaims{
			UserId:           uuid.New(),
			WalletAddress:    "0x1234567890",
			IsVerifiedIssuer: nil, // Not set in JWT
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		params := CheckRoleParams{
			JwtClaims:   validClaims,
			CheckIssuer: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsIssuer)
		assert.False(t, *result.IsIssuer) // Should default to false
	})
}

func TestAuthUsecase_CheckRole_MultipleRoles(t *testing.T) {
	t.Run("should check all requested roles", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		isVerifiedOrganizer := true
		isVerifiedIssuer := false
		validClaims := &auth.JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x1234567890",
			IsVerifiedOrganizer: &isVerifiedOrganizer,
			IsVerifiedIssuer:    &isVerifiedIssuer,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		params := CheckRoleParams{
			JwtClaims:   validClaims,
			CheckAuth:   true,
			CheckHost:   true,
			CheckIssuer: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		
		// All three checks should be present
		require.NotNil(t, result.IsAuthenticated)
		require.NotNil(t, result.IsHost)
		require.NotNil(t, result.IsIssuer)
		
		// Verify values
		assert.True(t, *result.IsAuthenticated)
		assert.True(t, *result.IsHost)
		assert.False(t, *result.IsIssuer)
	})

	t.Run("should only return requested fields", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		isVerifiedOrganizer := true
		isVerifiedIssuer := true
		validClaims := &auth.JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x1234567890",
			IsVerifiedOrganizer: &isVerifiedOrganizer,
			IsVerifiedIssuer:    &isVerifiedIssuer,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		params := CheckRoleParams{
			JwtClaims:   validClaims,
			CheckAuth:   false, // Not checking auth
			CheckHost:   true,  // Only checking host
			CheckIssuer: false, // Not checking issuer
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		
		// Only IsHost should be set
		assert.Nil(t, result.IsAuthenticated)
		require.NotNil(t, result.IsHost)
		assert.Nil(t, result.IsIssuer)
		
		assert.True(t, *result.IsHost)
	})

	t.Run("should return all false when token expired and all roles requested", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		isVerifiedOrganizer := true
		isVerifiedIssuer := true
		expiredClaims := &auth.JwtClaims{
			UserId:              uuid.New(),
			WalletAddress:       "0x1234567890",
			IsVerifiedOrganizer: &isVerifiedOrganizer,
			IsVerifiedIssuer:    &isVerifiedIssuer,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired
			},
		}

		params := CheckRoleParams{
			JwtClaims:   expiredClaims,
			CheckAuth:   true,
			CheckHost:   true,
			CheckIssuer: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		
		// All should be false due to expired token
		require.NotNil(t, result.IsAuthenticated)
		require.NotNil(t, result.IsHost)
		require.NotNil(t, result.IsIssuer)
		
		assert.False(t, *result.IsAuthenticated)
		assert.False(t, *result.IsHost)
		assert.False(t, *result.IsIssuer)
	})
}

func TestAuthUsecase_CheckRole_EdgeCases(t *testing.T) {
	t.Run("should handle no checks requested", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		validClaims := &auth.JwtClaims{
			UserId:        uuid.New(),
			WalletAddress: "0x1234567890",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		params := CheckRoleParams{
			JwtClaims:   validClaims,
			CheckAuth:   false,
			CheckHost:   false,
			CheckIssuer: false,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		
		// Nothing should be set
		assert.Nil(t, result.IsAuthenticated)
		assert.Nil(t, result.IsHost)
		assert.Nil(t, result.IsIssuer)
	})

	t.Run("should handle claims without expiration", func(t *testing.T) {
		// Arrange
		uc := NewAuthUsecase()
		ctx := context.Background()
		
		claimsNoExpiry := &auth.JwtClaims{
			UserId:           uuid.New(),
			WalletAddress:    "0x1234567890",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: nil, // No expiration set
			},
		}

		params := CheckRoleParams{
			JwtClaims: claimsNoExpiry,
			CheckAuth: true,
		}

		// Act
		result, err := uc.CheckRole(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.IsAuthenticated)
		
		// Should be valid since no expiration means it doesn't expire
		assert.True(t, *result.IsAuthenticated)
	})
}

func TestAuthUsecase_isTokenValid(t *testing.T) {
	uc := NewAuthUsecase()

	t.Run("should return false for nil claims", func(t *testing.T) {
		assert.False(t, uc.isTokenValid(nil))
	})

	t.Run("should return true for valid token", func(t *testing.T) {
		claims := &auth.JwtClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}
		assert.True(t, uc.isTokenValid(claims))
	})

	t.Run("should return false for expired token", func(t *testing.T) {
		claims := &auth.JwtClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			},
		}
		assert.False(t, uc.isTokenValid(claims))
	})

	t.Run("should return true when no expiration is set", func(t *testing.T) {
		claims := &auth.JwtClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: nil,
			},
		}
		assert.True(t, uc.isTokenValid(claims))
	})
}

