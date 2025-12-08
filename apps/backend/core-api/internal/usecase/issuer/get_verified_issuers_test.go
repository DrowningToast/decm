package issuer

import (
	"context"
	"decm-database/go/generated"
	"errors"
	"testing"
	"time"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockIssuerDataGateway struct {
	mock.Mock
}

func (m *MockIssuerDataGateway) ListVerifiedIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error) {
	args := m.Called(ctx, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Profile), args.Error(1)
}

func (m *MockIssuerDataGateway) ListIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error) {
	args := m.Called(ctx, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Profile), args.Error(1)
}

func (m *MockIssuerDataGateway) SearchIssuerCredentialsByWalletAddress(ctx context.Context, searchQuery string, limitCount int, offsetCount int) ([]entity.AuthenticationCredential, error) {
	args := m.Called(ctx, searchQuery, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.AuthenticationCredential), args.Error(1)
}

func (m *MockIssuerDataGateway) ListAllIssuerCredentials(ctx context.Context, limitCount int) ([]entity.AuthenticationCredential, error) {
	args := m.Called(ctx, limitCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.AuthenticationCredential), args.Error(1)
}

func (m *MockIssuerDataGateway) GetEventsByIssuerCredentialID(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetEventIssuersByCredentialIDRow, error) {
	return nil, errors.New("not implemented")
}

func (m *MockIssuerDataGateway) GetIssuerEventsWithDetails(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetIssuerEventsWithDetailsRow, error) {
	args := m.Called(ctx, issuerCredentialID, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]generated.GetIssuerEventsWithDetailsRow), args.Error(1)
}

func TestGetVerifiedIssuers(t *testing.T) {
	ctx := context.Background()

	// Helper function to create string pointer
	strPtr := func(s string) *string {
		return &s
	}

	// Mock data
	credId1 := uuid.New()
	credId2 := uuid.New()
	credId3 := uuid.New()

	mockCredentials := []entity.AuthenticationCredential{
		{
			Id:                  credId1,
			WalletAddress:       "0x1234567890abcdef",
			GoogleConnectorRef:  strPtr("alice@gmail.com"),
			IsVerifiedIssuer:    true,
			IsVerifiedOrganizer: false,
			IsVerifiedStudent:   false,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			Id:                  credId2,
			WalletAddress:       "0xabcdef1234567890",
			GoogleConnectorRef:  strPtr("bob@gmail.com"),
			IsVerifiedIssuer:    true,
			IsVerifiedOrganizer: false,
			IsVerifiedStudent:   false,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			Id:                  credId3,
			WalletAddress:       "0x9876543210fedcba",
			GoogleConnectorRef:  strPtr("charlie@gmail.com"),
			IsVerifiedIssuer:    true,
			IsVerifiedOrganizer: false,
			IsVerifiedStudent:   false,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	}

	mockProfiles := []entity.Profile{
		{
			Id:                          uuid.New(),
			AuthenticationCredentialId:  credId1,
			FirstName:                   strPtr("Alice"),
			LastName:                    strPtr("Johnson"),
			Email:                       strPtr("alice@example.com"),
			AcademicEmail:               strPtr("alice@university.edu"),
			AcademicInstitution:         strPtr("Stanford University"),
			IsFirstNamePublic:           true,
			IsLastNamePublic:            true,
			IsEmailPublic:               true,
			IsAcademicEmailPublic:       true,
			IsAcademicInstitutionPublic: true,
			CreatedAt:                   time.Now(),
			UpdatedAt:                   time.Now(),
		},
		{
			Id:                          uuid.New(),
			AuthenticationCredentialId:  credId2,
			FirstName:                   strPtr("Bob"),
			LastName:                    strPtr("Smith"),
			Email:                       strPtr("bob@example.com"),
			AcademicEmail:               strPtr("bob@university.edu"),
			AcademicInstitution:         strPtr("MIT"),
			IsFirstNamePublic:           true,
			IsLastNamePublic:            true,
			IsEmailPublic:               true,
			IsAcademicEmailPublic:       true,
			IsAcademicInstitutionPublic: true,
			CreatedAt:                   time.Now(),
			UpdatedAt:                   time.Now(),
		},
		{
			Id:                          uuid.New(),
			AuthenticationCredentialId:  credId3,
			FirstName:                   strPtr("Charlie"),
			LastName:                    strPtr("Brown"),
			Email:                       strPtr("charlie@example.com"),
			AcademicEmail:               strPtr("charlie@university.edu"),
			AcademicInstitution:         strPtr("Harvard University"),
			IsFirstNamePublic:           true,
			IsLastNamePublic:            true,
			IsEmailPublic:               true,
			IsAcademicEmailPublic:       true,
			IsAcademicInstitutionPublic: true,
			CreatedAt:                   time.Now(),
			UpdatedAt:                   time.Now(),
		},
	}

	t.Run("should return paginated verified issuers when no search query", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("ListVerifiedIssuerProfiles", ctx, 10, 0).Return(mockProfiles[:2], nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "Alice", *result[0].FirstName)
		assert.Equal(t, "Bob", *result[1].FirstName)
		assert.Equal(t, "alice@gmail.com", *result[0].GoogleConnectorRef)
		assert.Equal(t, "0x1234567890abcdef", result[0].WalletAddress)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should return empty list when no verified issuers exist", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return([]entity.AuthenticationCredential{}, nil)
		mockIssuerDg.On("ListVerifiedIssuerProfiles", ctx, 10, 0).Return([]entity.Profile{}, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 0, len(result))
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should search by first name (case-insensitive)", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "alice", 1000, 0).
			Return([]entity.AuthenticationCredential{}, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).Return(mockProfiles, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "alice",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "Alice", *result[0].FirstName)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should search by last name (case-insensitive)", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "smith", 1000, 0).
			Return([]entity.AuthenticationCredential{}, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).Return(mockProfiles, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "smith",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "Bob", *result[0].FirstName)
		assert.Equal(t, "Smith", *result[0].LastName)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should search by email (case-insensitive)", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "charlie@example.com", 1000, 0).
			Return([]entity.AuthenticationCredential{}, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).Return(mockProfiles, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "charlie@example.com",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "Charlie", *result[0].FirstName)
		assert.Equal(t, "charlie@example.com", *result[0].Email)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should search by academic email", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "alice@university.edu", 1000, 0).
			Return([]entity.AuthenticationCredential{}, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).Return(mockProfiles, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "alice@university.edu",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "Alice", *result[0].FirstName)
		assert.Equal(t, "alice@university.edu", *result[0].AcademicEmail)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should search by Google connector ref", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "bob@gmail.com", 1000, 0).
			Return([]entity.AuthenticationCredential{}, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).Return(mockProfiles, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "bob@gmail.com",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "Bob", *result[0].FirstName)
		assert.Equal(t, "bob@gmail.com", *result[0].GoogleConnectorRef)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should search by wallet address", func(t *testing.T) {
		// Arrange
		walletSearchResult := []entity.AuthenticationCredential{mockCredentials[0]}

		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "0x1234", 1000, 0).
			Return(walletSearchResult, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).Return(mockProfiles, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "0x1234",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "Alice", *result[0].FirstName)
		assert.Equal(t, "0x1234567890abcdef", result[0].WalletAddress)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should apply pagination to search results", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "university", 1000, 0).
			Return([]entity.AuthenticationCredential{}, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).Return(mockProfiles, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		// Search for "university" which matches all academic emails
		// Then apply pagination: offset 1, limit 1 (should get second result)
		req := GetVerifiedIssuersRequest{
			SearchQuery: "university",
			Limit:       1,
			Offset:      1,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		// Second profile in alphabetical order of academic emails
		assert.Equal(t, "Bob", *result[0].FirstName)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should return empty list when offset exceeds results", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "alice", 1000, 0).
			Return([]entity.AuthenticationCredential{}, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).Return(mockProfiles, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "alice",
			Limit:       10,
			Offset:      100, // Offset beyond results
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 0, len(result))
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should return empty list when no matches found", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "nonexistent", 1000, 0).
			Return([]entity.AuthenticationCredential{}, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).Return(mockProfiles, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "nonexistent",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 0, len(result))
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should handle error from ListAllIssuerCredentials", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).
			Return(nil, errors.New("database error"))

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should handle error from ListVerifiedIssuerProfiles", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("ListVerifiedIssuerProfiles", ctx, 10, 0).
			Return(nil, errors.New("database error"))

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should handle error from SearchIssuerCredentialsByWalletAddress", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "alice", 1000, 0).
			Return(nil, errors.New("database error"))

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "alice",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should handle error from ListIssuerProfiles", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "alice", 1000, 0).
			Return([]entity.AuthenticationCredential{}, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).
			Return(nil, errors.New("database error"))

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "alice",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should prevent duplicate results in search", func(t *testing.T) {
		// Arrange
		// Create a profile that matches both by wallet and by name
		walletSearchResult := []entity.AuthenticationCredential{mockCredentials[0]}

		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("SearchIssuerCredentialsByWalletAddress", ctx, "alice", 1000, 0).
			Return(walletSearchResult, nil)
		mockIssuerDg.On("ListIssuerProfiles", ctx, 1000, 0).Return(mockProfiles, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		// Search for "alice" which matches both wallet and first name
		req := GetVerifiedIssuersRequest{
			SearchQuery: "alice",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result)) // Should only return one result, not duplicates
		assert.Equal(t, "Alice", *result[0].FirstName)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should handle profiles with nil optional fields", func(t *testing.T) {
		// Arrange
		profileWithNilFields := []entity.Profile{
			{
				Id:                         uuid.New(),
				AuthenticationCredentialId: credId1,
				FirstName:                  nil,
				LastName:                   nil,
				Email:                      nil,
				AcademicEmail:              nil,
				CreatedAt:                  time.Now(),
				UpdatedAt:                  time.Now(),
			},
		}

		mockIssuerDg := new(MockIssuerDataGateway)
		mockIssuerDg.On("ListAllIssuerCredentials", ctx, 1000).Return(mockCredentials, nil)
		mockIssuerDg.On("ListVerifiedIssuerProfiles", ctx, 10, 0).Return(profileWithNilFields, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		req := GetVerifiedIssuersRequest{
			SearchQuery: "",
			Limit:       10,
			Offset:      0,
		}

		// Act
		result, err := uc.GetVerifiedIssuers(ctx, req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Nil(t, result[0].FirstName)
		assert.Nil(t, result[0].LastName)
		assert.Nil(t, result[0].Email)
		mockIssuerDg.AssertExpectations(t)
	})
}
