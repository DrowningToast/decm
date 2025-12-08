package auth

import (
	"context"
	"time"

	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
)

type AuthUsecase struct {
	AuthenticationCredentialDg datagateway.AuthenticationCredentialDataGateway
}

func NewAuthUsecase(authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway) *AuthUsecase {
	return &AuthUsecase{
		AuthenticationCredentialDg: authenticationCredentialDg,
	}
}

// CheckRoleResult contains the role check results
type CheckRoleResult struct {
	IsAuthenticated *bool
	IsHost          *bool
	IsIssuer        *bool
}

// CheckRoleParams contains the parameters for role checking
type CheckRoleParams struct {
	JwtClaims   *auth.JwtClaims // Full JWT claims for token validation and user context
	CheckAuth   bool
	CheckHost   bool
	CheckIssuer bool
}

// CheckRole checks user roles based on provided parameters
// Returns only the fields that were requested in params
// Validates JWT token expiration and reads roles directly from JWT claims (no database access)
func (u *AuthUsecase) CheckRole(ctx context.Context, params CheckRoleParams) (*CheckRoleResult, error) {
	result := &CheckRoleResult{}

	// Check if token is valid and not expired
	isValidToken := params.JwtClaims != nil && u.isTokenValid(params.JwtClaims)

	// Check authentication status
	if params.CheckAuth {
		result.IsAuthenticated = &isValidToken

		// If not authenticated, return early with only auth status
		if !isValidToken {
			// If roles were requested but user is not authenticated, return false for them
			if params.CheckHost {
				falseVal := false
				result.IsHost = &falseVal
			}
			if params.CheckIssuer {
				falseVal := false
				result.IsIssuer = &falseVal
			}
			return result, nil
		}
	}

	// Read roles directly from JWT claims (no database query needed!)
	if params.CheckHost && isValidToken {
		// Read IsVerifiedOrganizer from JWT claims
		if params.JwtClaims.IsVerifiedOrganizer != nil {
			result.IsHost = params.JwtClaims.IsVerifiedOrganizer
		} else {
			// If field is not in JWT, default to false
			falseVal := false
			result.IsHost = &falseVal
		}
	}

	if params.CheckIssuer && isValidToken {
		// Read IsVerifiedIssuer from JWT claims
		if params.JwtClaims.IsVerifiedIssuer != nil {
			result.IsIssuer = params.JwtClaims.IsVerifiedIssuer
		} else {
			// If field is not in JWT, default to false
			falseVal := false
			result.IsIssuer = &falseVal
		}
	}

	// Handle case where user is not authenticated but roles were requested
	if !isValidToken {
		if params.CheckHost && result.IsHost == nil {
			falseVal := false
			result.IsHost = &falseVal
		}
		if params.CheckIssuer && result.IsIssuer == nil {
			falseVal := false
			result.IsIssuer = &falseVal
		}
	}

	return result, nil
}

// isTokenValid checks if the JWT token is not expired
func (u *AuthUsecase) isTokenValid(claims *auth.JwtClaims) bool {
	if claims == nil {
		return false
	}

	// Check if token has expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return false
	}

	return true
}

func (u *AuthUsecase) GetAuthenticationCredentialId(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	return u.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, id)
}
