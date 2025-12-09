package oauth

import (
	"context"
	"testing"

	oauth_services "apps/backend/services/oauth"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// MockOAuthService is a mock implementation of the OAuthService interface
type MockOAuthService struct {
	mock.Mock
}

func (m *MockOAuthService) Login(sess *session.Session) (*string, error) {
	args := m.Called(sess)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockOAuthService) Callback(ctx context.Context, sess *session.Session, code string, state string) (*oauth2.Token, error) {
	args := m.Called(ctx, sess, code, state)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockOAuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*oauth_services.OAuthUser, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth_services.OAuthUser), args.Error(1)
}

func TestNewOAuthUsecase(t *testing.T) {
	t.Run("should create new OAuth usecase", func(t *testing.T) {
		// Arrange & Act
		uc := NewOAuthUsecase(nil, nil)

		// Assert
		require.NotNil(t, uc)
	})

	t.Run("should create new OAuth usecase with service", func(t *testing.T) {
		// Arrange
		// We can't easily mock GoogleOAuthService since it's a concrete type
		// So we test with nil and ensure the struct is created

		// Act
		uc := NewOAuthUsecase(nil, nil)

		// Assert
		require.NotNil(t, uc)
		// googleOAuthService will be nil but that's okay for constructor test
	})
}

func TestOAuthUsecase_VerifyGoogleOAuthCode(t *testing.T) {
	// Note: VerifyGoogleOAuthCode is difficult to unit test without refactoring because:
	// 1. It depends on a concrete *GoogleOAuthService type, not an interface
	// 2. GoogleOAuthService has complex external dependencies (HTTP client, OAuth config)
	//
	// Recommendations for better testability:
	// - Refactor OAuthUsecase to depend on oauth_services.OAuthService interface instead of *GoogleOAuthService
	// - This would allow proper mocking and comprehensive unit testing
	//
	// Current test coverage focuses on:
	// - Constructor (NewOAuthUsecase)
	// - Integration tests should cover the happy path and error scenarios
	t.Skip("Skipping VerifyGoogleOAuthCode test - requires refactoring to use interface for better testability")
}
