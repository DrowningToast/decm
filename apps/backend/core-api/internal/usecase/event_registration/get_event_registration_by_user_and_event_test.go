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

func TestGetEventRegistrationByUserAndEvent(t *testing.T) {
	t.Run("should get invitation successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890abcdef"
		invitationID := uuid.New()
		inboxMessageID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      invitationID,
			EventId: eventID,
		}

		inbox := &entity.InboxMessage{
			Id: inboxMessageID,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, mock.Anything, mock.Anything).Return(invitation, inbox, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		resultInvitation, resultInbox, err := uc.GetEventRegistrationByUserAndEvent(ctx, eventID, currentUser)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, resultInvitation)
		require.NotNil(t, resultInbox)
		assert.Equal(t, invitationID, resultInvitation.Id)
		assert.Equal(t, inboxMessageID, resultInbox.Id)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return nil when invitation not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890abcdef"

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, mock.Anything, mock.Anything).Return(nil, nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		resultInvitation, resultInbox, err := uc.GetEventRegistrationByUserAndEvent(ctx, eventID, currentUser)

		// Assert
		require.NoError(t, err) // Function handles not found gracefully
		assert.Nil(t, resultInvitation)
		assert.Nil(t, resultInbox)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when database fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890abcdef"

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, mock.Anything, mock.Anything).Return(nil, nil, customerror.NewWithPreset(&customerror.ErrInternalServer, errors.New("database error")))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		resultInvitation, resultInbox, err := uc.GetEventRegistrationByUserAndEvent(ctx, eventID, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, resultInvitation)
		assert.Nil(t, resultInbox)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInternalServer.Code, *customError.Code)
		mockInvitationDg.AssertExpectations(t)
	})
}
