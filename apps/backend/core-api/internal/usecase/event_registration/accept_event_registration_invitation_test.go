package event_registration

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// All mocks are now in mocks_test.go for reusability across test files

func TestEventRegistrationUsecase_AcceptEventRegistrationInvitation(t *testing.T) {
	t.Run("should accept invitation successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      invitationID,
			EventId: eventID,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, nil, nil)
		mockInvitationDg.On("UpdateEventRegistrationInvitationAcceptedStatus", ctx, invitationID, mock.AnythingOfType("*time.Time")).Return(invitation, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, currentUser)

		// Assert
		require.NoError(t, err)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when user is not authenticated", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()

		uc := &EventRegistrationUsecase{}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, nil)

		// Assert
		require.Error(t, err)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrUnauthenticated.Code, *customError.Code)
	})

	t.Run("should return error when invitation not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		email := "test@example.com"
		currentUser := &auth.JwtClaims{
			UserId:        uuid.New(),
			Email:         &email,
			WalletAddress: "0x1234567890",
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, currentUser)

		// Assert
		require.Error(t, err)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when invitation belongs to different user", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      invitationID,
			EventId: eventID,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		// Return nil to simulate invitation not belonging to current user
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(nil, nil, customerror.NewWithPreset(&customerror.ErrUnauthorized, errors.New("unauthorized")))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, currentUser)

		// Assert
		require.Error(t, err)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customError.Code)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when invitation ID mismatch", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		differentInvitationID := uuid.New()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      invitationID,
			EventId: eventID,
		}

		differentInvitation := &entity.EventRegistrationInvitation{
			Id:      differentInvitationID,
			EventId: eventID,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		// Return different invitation ID
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(differentInvitation, nil, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, currentUser)

		// Assert
		require.Error(t, err)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customError.Code)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      invitationID,
			EventId: eventID,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, nil, nil)
		mockInvitationDg.On("UpdateEventRegistrationInvitationAcceptedStatus", ctx, invitationID, mock.AnythingOfType("*time.Time")).Return(nil, errors.New("database error"))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, currentUser)

		// Assert
		require.Error(t, err)
		mockInvitationDg.AssertExpectations(t)
	})
}
