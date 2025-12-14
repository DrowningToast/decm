package auth

import (
	"testing"
	"time"

	"apps/backend/common"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthService_CreateToken_WithAllRoles tests JWT creation with all possible role combinations
func TestAuthService_CreateToken_WithAllRoles(t *testing.T) {
	service := NewAuthService("test-issuer", "test-secret-key-123456789012345678901234567890", time.Hour, "")

	tests := []struct {
		name                   string
		isVerifiedOrganizer    *bool
		isVerifiedIssuer       *bool
		expectedOrganizerInJWT bool
		expectedIssuerInJWT    bool
	}{
		{
			name:                   "regular user - no roles",
			isVerifiedOrganizer:    nil,
			isVerifiedIssuer:       nil,
			expectedOrganizerInJWT: false,
			expectedIssuerInJWT:    false,
		},
		{
			name:                   "verified organizer only",
			isVerifiedOrganizer:    boolPtr(true),
			isVerifiedIssuer:       nil,
			expectedOrganizerInJWT: true,
			expectedIssuerInJWT:    false,
		},
		{
			name:                   "verified issuer only",
			isVerifiedOrganizer:    nil,
			isVerifiedIssuer:       boolPtr(true),
			expectedOrganizerInJWT: false,
			expectedIssuerInJWT:    true,
		},
		{
			name:                   "both roles - organizer and issuer",
			isVerifiedOrganizer:    boolPtr(true),
			isVerifiedIssuer:       boolPtr(true),
			expectedOrganizerInJWT: true,
			expectedIssuerInJWT:    true,
		},
		{
			name:                   "organizer explicitly false",
			isVerifiedOrganizer:    boolPtr(false),
			isVerifiedIssuer:       nil,
			expectedOrganizerInJWT: false,
			expectedIssuerInJWT:    false,
		},
		{
			name:                   "issuer explicitly false",
			isVerifiedOrganizer:    nil,
			isVerifiedIssuer:       boolPtr(false),
			expectedOrganizerInJWT: false,
			expectedIssuerInJWT:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := JwtPayload{
				UserId:              uuid.New(),
				WalletAddress:       "0x1234567890123456789012345678901234567890",
				Email:               stringPtr("test@example.com"),
				SolutionStatus:      common.SolutionStatusBYOK,
				IsVerifiedOrganizer: tt.isVerifiedOrganizer,
				IsVerifiedIssuer:    tt.isVerifiedIssuer,
			}

			token, err := service.CreateToken(payload)
			require.NoError(t, err)
			assert.NotEmpty(t, token)

			// Verify token and check roles
			claims, err := service.VerifyToken(token)
			require.NoError(t, err)

			// Check organizer role
			if tt.expectedOrganizerInJWT {
				require.NotNil(t, claims.IsVerifiedOrganizer, "IsVerifiedOrganizer should not be nil")
				assert.True(t, *claims.IsVerifiedOrganizer, "IsVerifiedOrganizer should be true")
			} else {
				if claims.IsVerifiedOrganizer != nil {
					assert.False(t, *claims.IsVerifiedOrganizer, "IsVerifiedOrganizer should be false or nil")
				}
			}

			// Check issuer role
			if tt.expectedIssuerInJWT {
				require.NotNil(t, claims.IsVerifiedIssuer, "IsVerifiedIssuer should not be nil")
				assert.True(t, *claims.IsVerifiedIssuer, "IsVerifiedIssuer should be true")
			} else {
				if claims.IsVerifiedIssuer != nil {
					assert.False(t, *claims.IsVerifiedIssuer, "IsVerifiedIssuer should be false or nil")
				}
			}
		})
	}
}

// TestAuthService_VerifyToken_RolePersistence tests that roles persist through token round-trip
func TestAuthService_VerifyToken_RolePersistence(t *testing.T) {
	service := NewAuthService("test-issuer", "test-secret-key-123456789012345678901234567890", time.Hour, "")

	// Test organizer role persistence
	t.Run("organizer role persists", func(t *testing.T) {
		payload := JwtPayload{
			UserId:              uuid.New(),
			WalletAddress:       "0x1234567890123456789012345678901234567890",
			SolutionStatus:      common.SolutionStatusBYOK,
			IsVerifiedOrganizer: boolPtr(true),
		}

		token, err := service.CreateToken(payload)
		require.NoError(t, err)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)
		require.NotNil(t, claims.IsVerifiedOrganizer)
		assert.True(t, *claims.IsVerifiedOrganizer)
	})

	// Test issuer role persistence
	t.Run("issuer role persists", func(t *testing.T) {
		payload := JwtPayload{
			UserId:           uuid.New(),
			WalletAddress:    "0x1234567890123456789012345678901234567890",
			SolutionStatus:   common.SolutionStatusBYOK,
			IsVerifiedIssuer: boolPtr(true),
		}

		token, err := service.CreateToken(payload)
		require.NoError(t, err)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)
		require.NotNil(t, claims.IsVerifiedIssuer)
		assert.True(t, *claims.IsVerifiedIssuer)
	})

	// Test both roles persist
	t.Run("both roles persist", func(t *testing.T) {
		payload := JwtPayload{
			UserId:              uuid.New(),
			WalletAddress:       "0x1234567890123456789012345678901234567890",
			SolutionStatus:      common.SolutionStatusBYOK,
			IsVerifiedOrganizer: boolPtr(true),
			IsVerifiedIssuer:    boolPtr(true),
		}

		token, err := service.CreateToken(payload)
		require.NoError(t, err)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)
		require.NotNil(t, claims.IsVerifiedOrganizer)
		require.NotNil(t, claims.IsVerifiedIssuer)
		assert.True(t, *claims.IsVerifiedOrganizer)
		assert.True(t, *claims.IsVerifiedIssuer)
	})
}

// TestAuthService_CreateToken_RoleIndependence tests that roles are independent
func TestAuthService_CreateToken_RoleIndependence(t *testing.T) {
	service := NewAuthService("test-issuer", "test-secret-key-123456789012345678901234567890", time.Hour, "")

	// Create token with only organizer role
	payload1 := JwtPayload{
		UserId:              uuid.New(),
		WalletAddress:       "0x1234567890123456789012345678901234567890",
		SolutionStatus:      common.SolutionStatusBYOK,
		IsVerifiedOrganizer: boolPtr(true),
		IsVerifiedIssuer:    nil,
	}

	token1, err := service.CreateToken(payload1)
	require.NoError(t, err)

	claims1, err := service.VerifyToken(token1)
	require.NoError(t, err)
	assert.True(t, *claims1.IsVerifiedOrganizer)
	assert.Nil(t, claims1.IsVerifiedIssuer)

	// Create token with only issuer role
	payload2 := JwtPayload{
		UserId:              uuid.New(),
		WalletAddress:       "0x1234567890123456789012345678901234567890",
		SolutionStatus:      common.SolutionStatusBYOK,
		IsVerifiedOrganizer: nil,
		IsVerifiedIssuer:    boolPtr(true),
	}

	token2, err := service.CreateToken(payload2)
	require.NoError(t, err)

	claims2, err := service.VerifyToken(token2)
	require.NoError(t, err)
	assert.Nil(t, claims2.IsVerifiedOrganizer)
	assert.True(t, *claims2.IsVerifiedIssuer)

	// Verify tokens are different
	assert.NotEqual(t, token1, token2)
}







