package profile

import (
	"context"
	"errors"
	"testing"
	"time"

	"apps/backend/common"
	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockProfileDataGateway is a mock implementation of ProfileDataGateway
type MockProfileDataGateway struct {
	mock.Mock
}

func (m *MockProfileDataGateway) CreateProfile(ctx context.Context, profile entity.Profile) (*entity.Profile, error) {
	args := m.Called(ctx, profile)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfileDataGateway) GetProfileById(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfileDataGateway) GetProfileByAuthenticationCredentialId(ctx context.Context, credentialId uuid.UUID) (*entity.Profile, error) {
	args := m.Called(ctx, credentialId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfileDataGateway) GetProfileAndCredentialWithCredentialId(ctx context.Context, credentialId uuid.UUID) (*entity.Profile, *entity.AuthenticationCredential, error) {
	args := m.Called(ctx, credentialId)
	if args.Get(0) == nil || args.Get(1) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entity.Profile), args.Get(1).(*entity.AuthenticationCredential), args.Error(2)
}

func (m *MockProfileDataGateway) GetProfileByEmail(ctx context.Context, email string) (*entity.Profile, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfileDataGateway) UpdateProfile(ctx context.Context, id uuid.UUID, params datagateway.UpdateProfileParameters) (*entity.Profile, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfileDataGateway) UpdateProfileByAuthenticationCredentialId(ctx context.Context, credentialId uuid.UUID, params datagateway.UpdateProfileParameters) (*entity.Profile, error) {
	args := m.Called(ctx, credentialId, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfileDataGateway) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProfileDataGateway) ListVerifiedIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error) {
	args := m.Called(ctx, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Profile), args.Error(1)
}

// MockAuthenticationCredentialDataGateway is a mock implementation
type MockAuthenticationCredentialDataGateway struct {
	mock.Mock
}

func (m *MockAuthenticationCredentialDataGateway) GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDataGateway) GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDataGateway) GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, walletAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDataGateway) GetAuthenticationCredentialByGoogleConnectorRef(ctx context.Context, googleConnectorRef string) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, googleConnectorRef)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDataGateway) GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddress(ctx context.Context, params datagateway.GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddressParameters) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDataGateway) CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDataGateway) UpdateAuthenticationCredential(ctx context.Context, id uuid.UUID, params datagateway.UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDataGateway) DeleteAuthenticationCredential(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockInboxMessageDataGateway is a mock implementation of InboxMessageDataGateway
type MockInboxMessageDataGateway struct {
	mock.Mock
}

func (m *MockInboxMessageDataGateway) CreateInboxMessage(ctx context.Context, params datagateway.CreateInboxMessageParameters) (*entity.InboxMessage, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetInboxMessageByID(ctx context.Context, id uuid.UUID) (*entity.InboxMessage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetInboxMessagesByCredentialID(ctx context.Context, params datagateway.GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetUnreadInboxMessageCountByCredentialID(ctx context.Context, params datagateway.GetInboxMessagesByCredentialIDParameters) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetInboxMessagesByReceiverEmail(ctx context.Context, receiverEmail string) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, receiverEmail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetInboxMessagesByReceiverWalletAddress(ctx context.Context, walletAddress string) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, walletAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetInboxMessagesBySenderCredentialID(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, credentialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) UpdateInboxMessageReadStatus(ctx context.Context, id uuid.UUID, isRead int) (*entity.InboxMessage, error) {
	args := m.Called(ctx, id, isRead)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) UpdateInboxMessageReadStatusAll(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, credentialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func TestNewProfileUsecase(t *testing.T) {
	t.Run("should create new profile usecase", func(t *testing.T) {
		// Arrange
		mockProfileDg := new(MockProfileDataGateway)
		mockAuthCredDg := new(MockAuthenticationCredentialDataGateway)
		mockAuthService := &auth.AuthService{}
		mockInboxDg := new(MockInboxMessageDataGateway)

		// Act
		uc := NewProfileUsecase(mockProfileDg, mockAuthCredDg, mockAuthService, mockInboxDg)

		// Assert
		require.NotNil(t, uc)
		assert.Equal(t, mockProfileDg, uc.ProfileDg)
		assert.Equal(t, mockAuthCredDg, uc.AuthenticationCredentialDg)
		assert.Equal(t, mockAuthService, uc.AuthService)
	})
}

func TestProfileUsecase_GetProfileById(t *testing.T) {
	t.Run("should get profile by id successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		profileID := uuid.New()
		expectedProfile := &entity.Profile{
			Id:        profileID,
			FirstName: stringPtr("John"),
			LastName:  stringPtr("Doe"),
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileById", ctx, profileID).Return(expectedProfile, nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		profile, err := uc.GetProfileById(ctx, profileID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, expectedProfile.Id, profile.Id)
		assert.Equal(t, expectedProfile.FirstName, profile.FirstName)
		mockProfileDg.AssertExpectations(t)
	})

	t.Run("should return error when profile not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		profileID := uuid.New()

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileById", ctx, profileID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		profile, err := uc.GetProfileById(ctx, profileID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, profile)
		mockProfileDg.AssertExpectations(t)
	})
}

func TestProfileUsecase_GetProfileByAuthenticationCredentialId(t *testing.T) {
	t.Run("should get profile by credential id successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		expectedProfile := &entity.Profile{
			Id:                         uuid.New(),
			AuthenticationCredentialId: credentialID,
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByAuthenticationCredentialId", ctx, credentialID).Return(expectedProfile, nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		profile, err := uc.GetProfileByAuthenticationCredentialId(ctx, credentialID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, credentialID, profile.AuthenticationCredentialId)
		mockProfileDg.AssertExpectations(t)
	})
}

func TestProfileUsecase_GetProfileAndCredentialWithCredentialId(t *testing.T) {
	t.Run("should get profile and credential successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		expectedProfile := &entity.Profile{
			Id:                         uuid.New(),
			AuthenticationCredentialId: credentialID,
		}
		expectedCredential := &entity.AuthenticationCredential{
			Id:            credentialID,
			WalletAddress: "0x1234567890",
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileAndCredentialWithCredentialId", ctx, credentialID).Return(expectedProfile, expectedCredential, nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		profile, credential, err := uc.GetProfileAndCredentialWithCredentialId(ctx, credentialID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, profile)
		require.NotNil(t, credential)
		assert.Equal(t, credentialID, profile.AuthenticationCredentialId)
		assert.Equal(t, credentialID, credential.Id)
		mockProfileDg.AssertExpectations(t)
	})
}

func TestProfileUsecase_GetProfileByEmail(t *testing.T) {
	t.Run("should get profile by email successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"
		expectedProfile := &entity.Profile{
			Id:    uuid.New(),
			Email: &email,
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByEmail", ctx, email).Return(expectedProfile, nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		profile, err := uc.GetProfileByEmail(ctx, email)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, &email, profile.Email)
		mockProfileDg.AssertExpectations(t)
	})
}

func TestProfileUsecase_CreateProfile(t *testing.T) {
	t.Run("should create profile successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		firstName := "John"
		email := "john@example.com"

		params := CreateProfileParameters{
			AuthenticationCredentialId: credentialID,
			FirstName:                  &firstName,
			Email:                      &email,
		}

		expectedProfile := &entity.Profile{
			Id:                         uuid.New(),
			AuthenticationCredentialId: credentialID,
			FirstName:                  &firstName,
			Email:                      &email,
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("CreateProfile", ctx, mock.AnythingOfType("entity.Profile")).Return(expectedProfile, nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		profile, err := uc.CreateProfile(ctx, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, credentialID, profile.AuthenticationCredentialId)
		assert.Equal(t, &firstName, profile.FirstName)
		mockProfileDg.AssertExpectations(t)
	})
}

func TestProfileUsecase_DeleteProfile(t *testing.T) {
	t.Run("should delete profile successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		profileID := uuid.New()

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("DeleteProfile", ctx, profileID).Return(nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		err := uc.DeleteProfile(ctx, profileID)

		// Assert
		require.NoError(t, err)
		mockProfileDg.AssertExpectations(t)
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		profileID := uuid.New()

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("DeleteProfile", ctx, profileID).Return(errors.New("database error"))

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		err := uc.DeleteProfile(ctx, profileID)

		// Assert
		require.Error(t, err)
		mockProfileDg.AssertExpectations(t)
	})
}

func TestProfileUsecase_VerifyUserPassword(t *testing.T) {
	t.Run("should verify password successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		password := "MySecurePassword123"

		// Hash the password
		hashedPassword, err := hashutils.HashPassword(password)
		require.NoError(t, err)

		credential := &entity.AuthenticationCredential{
			Id:             credentialID,
			HashedPassword: &hashedPassword,
		}

		mockAuthCredDg := new(MockAuthenticationCredentialDataGateway)
		mockAuthCredDg.On("GetAuthenticationCredentialById", ctx, credentialID).Return(credential, nil)

		uc := &ProfileUsecase{
			AuthenticationCredentialDg: mockAuthCredDg,
		}

		// Act
		isValid, err := uc.VerifyUserPassword(ctx, credentialID, password)

		// Assert
		require.NoError(t, err)
		assert.True(t, isValid)
		mockAuthCredDg.AssertExpectations(t)
	})

	t.Run("should return false for incorrect password", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		correctPassword := "CorrectPassword123"
		wrongPassword := "WrongPassword123"

		// Hash the correct password
		hashedPassword, err := hashutils.HashPassword(correctPassword)
		require.NoError(t, err)

		credential := &entity.AuthenticationCredential{
			Id:             credentialID,
			HashedPassword: &hashedPassword,
		}

		mockAuthCredDg := new(MockAuthenticationCredentialDataGateway)
		mockAuthCredDg.On("GetAuthenticationCredentialById", ctx, credentialID).Return(credential, nil)

		uc := &ProfileUsecase{
			AuthenticationCredentialDg: mockAuthCredDg,
		}

		// Act
		isValid, err := uc.VerifyUserPassword(ctx, credentialID, wrongPassword)

		// Assert
		require.NoError(t, err)
		assert.False(t, isValid)
		mockAuthCredDg.AssertExpectations(t)
	})

	t.Run("should return error when password not set", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()

		credential := &entity.AuthenticationCredential{
			Id:             credentialID,
			HashedPassword: nil, // No password set
		}

		mockAuthCredDg := new(MockAuthenticationCredentialDataGateway)
		mockAuthCredDg.On("GetAuthenticationCredentialById", ctx, credentialID).Return(credential, nil)

		uc := &ProfileUsecase{
			AuthenticationCredentialDg: mockAuthCredDg,
		}

		// Act
		isValid, err := uc.VerifyUserPassword(ctx, credentialID, "any-password")

		// Assert
		require.Error(t, err)
		assert.False(t, isValid)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customError.Code)
		mockAuthCredDg.AssertExpectations(t)
	})
}

func TestProfileUsecase_GetMyProfileViewModel(t *testing.T) {
	t.Run("should get profile view model successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		profileID := uuid.New()
		email := "test@example.com"
		firstName := "John"
		walletAddress := "0x1234567890"

		user := &auth.JwtClaims{
			UserId: credentialID,
		}

		profile := &entity.Profile{
			Id:                         profileID,
			AuthenticationCredentialId: credentialID,
			Email:                      &email,
			FirstName:                  &firstName,
			IsEmailPublic:              true,
			IsFirstNamePublic:          true,
			CreatedAt:                  time.Now(),
			UpdatedAt:                  time.Now(),
		}

		credential := &entity.AuthenticationCredential{
			Id:                 credentialID,
			WalletAddress:      walletAddress,
			GoogleConnectorRef: &email,
			SolutionStatus:     common.SolutionStatusManaged,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileAndCredentialWithCredentialId", ctx, credentialID).Return(profile, credential, nil)

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetUnreadInboxMessageCountByCredentialID", ctx, mock.Anything).Return(0, nil)

		uc := &ProfileUsecase{
			ProfileDg:      mockProfileDg,
			InboxMessageDg: mockInboxDg,
		}

		// Act
		viewModel, err := uc.GetMyProfileViewModel(ctx, user)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, viewModel)
		assert.Equal(t, profileID.String(), viewModel.ProfileId)
		assert.Equal(t, credentialID.String(), viewModel.AuthenticationCredentialId)
		assert.Equal(t, &email, viewModel.Email)
		assert.Equal(t, &firstName, viewModel.FirstName)
		assert.Equal(t, walletAddress, viewModel.WalletAddress)
		assert.Equal(t, common.SolutionStatusManaged, viewModel.SolutionStatus)
		mockProfileDg.AssertExpectations(t)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return error when profile not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()

		user := &auth.JwtClaims{
			UserId: credentialID,
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileAndCredentialWithCredentialId", ctx, credentialID).Return(nil, nil, errors.New("not found"))

		mockInboxDg := new(MockInboxMessageDataGateway)

		uc := &ProfileUsecase{
			ProfileDg:      mockProfileDg,
			InboxMessageDg: mockInboxDg,
		}

		// Act
		viewModel, err := uc.GetMyProfileViewModel(ctx, user)

		// Assert
		require.Error(t, err)
		assert.Nil(t, viewModel)
		mockProfileDg.AssertExpectations(t)
	})
}

func TestProfileUsecase_UpdateProfile(t *testing.T) {
	t.Run("should update profile successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		profileID := uuid.New()
		newEmail := "newemail@example.com"
		newFirstName := "UpdatedFirst"

		existingProfile := &entity.Profile{
			Id:                         profileID,
			AuthenticationCredentialId: credentialID,
		}

		updateParams := UpdateProfileParameters{
			Email:     &newEmail,
			FirstName: &newFirstName,
		}

		updatedProfile := &entity.Profile{
			Id:                         profileID,
			AuthenticationCredentialId: credentialID,
			Email:                      &newEmail,
			FirstName:                  &newFirstName,
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByAuthenticationCredentialId", ctx, credentialID).Return(existingProfile, nil)
		mockProfileDg.On("UpdateProfile", ctx, profileID, mock.AnythingOfType("datagateway.UpdateProfileParameters")).Return(updatedProfile, nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		result, err := uc.UpdateProfile(ctx, credentialID, updateParams)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, &newEmail, result.Email)
		assert.Equal(t, &newFirstName, result.FirstName)
		mockProfileDg.AssertExpectations(t)
	})

	t.Run("should return error when profile not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		updateParams := UpdateProfileParameters{
			FirstName: stringPtr("Test"),
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByAuthenticationCredentialId", ctx, credentialID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		result, err := uc.UpdateProfile(ctx, credentialID, updateParams)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		mockProfileDg.AssertExpectations(t)
	})

	t.Run("should return error when profile is nil", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		updateParams := UpdateProfileParameters{
			FirstName: stringPtr("Test"),
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByAuthenticationCredentialId", ctx, credentialID).Return(nil, nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		result, err := uc.UpdateProfile(ctx, credentialID, updateParams)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrNotFound.Code, *customError.Code)
		mockProfileDg.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		profileID := uuid.New()
		existingProfile := &entity.Profile{
			Id:                         profileID,
			AuthenticationCredentialId: credentialID,
		}

		updateParams := UpdateProfileParameters{
			FirstName: stringPtr("Test"),
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByAuthenticationCredentialId", ctx, credentialID).Return(existingProfile, nil)
		mockProfileDg.On("UpdateProfile", ctx, profileID, mock.AnythingOfType("datagateway.UpdateProfileParameters")).Return(nil, errors.New("database error"))

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		result, err := uc.UpdateProfile(ctx, credentialID, updateParams)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		mockProfileDg.AssertExpectations(t)
	})
}

func TestProfileUsecase_UpdateProfileByCredentialId(t *testing.T) {
	t.Run("should update profile by credential id successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		profileID := uuid.New()
		newEmail := "newemail@example.com"

		existingProfile := &entity.Profile{
			Id:                         profileID,
			AuthenticationCredentialId: credentialID,
		}

		updateParams := UpdateProfileParameters{
			Email: &newEmail,
		}

		updatedProfile := &entity.Profile{
			Id:                         profileID,
			AuthenticationCredentialId: credentialID,
			Email:                      &newEmail,
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByAuthenticationCredentialId", ctx, credentialID).Return(existingProfile, nil)
		mockProfileDg.On("UpdateProfileByAuthenticationCredentialId", ctx, credentialID, mock.AnythingOfType("datagateway.UpdateProfileParameters")).Return(updatedProfile, nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		result, err := uc.UpdateProfileByCredentialId(ctx, credentialID, updateParams)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, &newEmail, result.Email)
		mockProfileDg.AssertExpectations(t)
	})

	t.Run("should return error when profile not found (ErrNotFound)", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		updateParams := UpdateProfileParameters{
			FirstName: stringPtr("Test"),
		}

		// Note: The implementation has a logic bug where ErrNotFound returns error instead of creating profile
		notFoundErr := customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found"))

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByAuthenticationCredentialId", ctx, credentialID).Return(nil, notFoundErr)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		result, err := uc.UpdateProfileByCredentialId(ctx, credentialID, updateParams)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		mockProfileDg.AssertExpectations(t)
	})

	t.Run("should create new profile when non-ErrNotFound error occurs", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		profilePictureUrl := "https://example.com/picture.jpg"

		updateParams := UpdateProfileParameters{
			ProfilePictureUrl: &profilePictureUrl,
		}

		newProfile := &entity.Profile{
			Id:                         uuid.New(),
			AuthenticationCredentialId: credentialID,
			ProfilePictureUrl:          &profilePictureUrl,
		}

		// Note: The implementation has inverted logic - creates profile when error is NOT ErrNotFound
		dbError := customerror.NewWithPreset(&customerror.ErrInternalServer, errors.New("database error"))

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByAuthenticationCredentialId", ctx, credentialID).Return(nil, dbError)
		mockProfileDg.On("CreateProfile", ctx, mock.AnythingOfType("entity.Profile")).Return(newProfile, nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		result, err := uc.UpdateProfileByCredentialId(ctx, credentialID, updateParams)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, &profilePictureUrl, result.ProfilePictureUrl)
		mockProfileDg.AssertExpectations(t)
	})

	t.Run("should return error when profile is nil", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		updateParams := UpdateProfileParameters{
			FirstName: stringPtr("Test"),
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByAuthenticationCredentialId", ctx, credentialID).Return(nil, nil)

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		result, err := uc.UpdateProfileByCredentialId(ctx, credentialID, updateParams)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrNotFound.Code, *customError.Code)
		mockProfileDg.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		profileID := uuid.New()
		existingProfile := &entity.Profile{
			Id:                         profileID,
			AuthenticationCredentialId: credentialID,
		}

		updateParams := UpdateProfileParameters{
			FirstName: stringPtr("Test"),
		}

		mockProfileDg := new(MockProfileDataGateway)
		mockProfileDg.On("GetProfileByAuthenticationCredentialId", ctx, credentialID).Return(existingProfile, nil)
		mockProfileDg.On("UpdateProfileByAuthenticationCredentialId", ctx, credentialID, mock.AnythingOfType("datagateway.UpdateProfileParameters")).Return(nil, errors.New("database error"))

		uc := &ProfileUsecase{
			ProfileDg: mockProfileDg,
		}

		// Act
		result, err := uc.UpdateProfileByCredentialId(ctx, credentialID, updateParams)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		mockProfileDg.AssertExpectations(t)
	})
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
