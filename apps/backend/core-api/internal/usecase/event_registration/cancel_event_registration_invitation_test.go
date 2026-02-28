package event_registration

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCancelEventRegistrationInvitation(t *testing.T) {
	t.Run("should cancel invitation successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		inboxMessageID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		now := time.Now()

		currentUser := &auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:             invitationID,
			EventId:        eventID,
			InboxMessageId: inboxMessageID,
			Email:          &email,
			CancelledAt:    nil, // Not cancelled yet
		}

		originalInboxMessage := &entity.InboxMessage{
			Id:                   inboxMessageID,
			ReceiverCredentialId: &userID,
			ReceiverEmail:        &email,
		}

		event := &entity.Event{
			Id:    eventID,
			Title: "Test Event",
		}

		updatedInvitation := &entity.EventRegistrationInvitation{
			Id:             invitationID,
			EventId:        eventID,
			InboxMessageId: inboxMessageID,
			Email:          &email,
			CancelledAt:    &now,
		}

		createdNotification := &entity.InboxMessage{
			Id:                   uuid.New(),
			ReceiverCredentialId: &userID,
			ReceiverEmail:        &email,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		mockInvitationDg.On("UpdateEventRegistrationInvitation", ctx, invitationID, mock.AnythingOfType("offchain_datagateway.UpdateEventRegistrationInvitationParameters")).Return(updatedInvitation, nil)

		mockInboxMessageDg := new(MockInboxMessageDg)
		mockInboxMessageDg.On("GetInboxMessageByID", ctx, inboxMessageID).Return(originalInboxMessage, nil)
		mockInboxMessageDg.On("CreateInboxMessage", ctx, mock.AnythingOfType("offchain_datagateway.CreateInboxMessageParameters")).Return(createdNotification, nil)

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
			InboxMessageDg:                mockInboxMessageDg,
			EventDg:                       mockEventDg,
		}

		// Act
		result, err := uc.CancelEventRegistrationInvitation(ctx, CancelEventRegistrationInvitationParameters{
			EventRegistrationInvitationID: invitationID,
		}, currentUser)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, invitationID, result.Id)
		assert.NotNil(t, result.CancelledAt)
		mockInvitationDg.AssertExpectations(t)
		mockInboxMessageDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return error when invitation not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"

		currentUser := &auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		result, err := uc.CancelEventRegistrationInvitation(ctx, CancelEventRegistrationInvitationParameters{
			EventRegistrationInvitationID: invitationID,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrNotFound.Code, *customError.Code)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when invitation already cancelled", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		inboxMessageID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		cancelledAt := time.Now().Add(-24 * time.Hour)

		currentUser := &auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:             invitationID,
			EventId:        eventID,
			InboxMessageId: inboxMessageID,
			Email:          &email,
			CancelledAt:    &cancelledAt, // Already cancelled
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		result, err := uc.CancelEventRegistrationInvitation(ctx, CancelEventRegistrationInvitationParameters{
			EventRegistrationInvitationID: invitationID,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customError.Code)
		assert.Contains(t, err.Error(), "invitation is already cancelled")
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when original inbox message not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		inboxMessageID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"

		currentUser := &auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:             invitationID,
			EventId:        eventID,
			InboxMessageId: inboxMessageID,
			Email:          &email,
			CancelledAt:    nil,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)

		mockInboxMessageDg := new(MockInboxMessageDg)
		mockInboxMessageDg.On("GetInboxMessageByID", ctx, inboxMessageID).Return(nil, errors.New("inbox message not found"))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
			InboxMessageDg:                mockInboxMessageDg,
		}

		// Act
		result, err := uc.CancelEventRegistrationInvitation(ctx, CancelEventRegistrationInvitationParameters{
			EventRegistrationInvitationID: invitationID,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInternalServer.Code, *customError.Code)
		assert.Contains(t, err.Error(), "failed to get original inbox message")
		mockInvitationDg.AssertExpectations(t)
		mockInboxMessageDg.AssertExpectations(t)
	})

	t.Run("should return error when event not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		inboxMessageID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"

		currentUser := &auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:             invitationID,
			EventId:        eventID,
			InboxMessageId: inboxMessageID,
			Email:          &email,
			CancelledAt:    nil,
		}

		originalInboxMessage := &entity.InboxMessage{
			Id:                   inboxMessageID,
			ReceiverCredentialId: &userID,
			ReceiverEmail:        &email,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)

		mockInboxMessageDg := new(MockInboxMessageDg)
		mockInboxMessageDg.On("GetInboxMessageByID", ctx, inboxMessageID).Return(originalInboxMessage, nil)

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(nil, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
			InboxMessageDg:                mockInboxMessageDg,
			EventDg:                       mockEventDg,
		}

		// Act
		result, err := uc.CancelEventRegistrationInvitation(ctx, CancelEventRegistrationInvitationParameters{
			EventRegistrationInvitationID: invitationID,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrNotFound.Code, *customError.Code)
		assert.Contains(t, err.Error(), "event not found")
		mockInvitationDg.AssertExpectations(t)
		mockInboxMessageDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return error when update invitation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		inboxMessageID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"

		currentUser := &auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:             invitationID,
			EventId:        eventID,
			InboxMessageId: inboxMessageID,
			Email:          &email,
			CancelledAt:    nil,
		}

		originalInboxMessage := &entity.InboxMessage{
			Id:                   inboxMessageID,
			ReceiverCredentialId: &userID,
			ReceiverEmail:        &email,
		}

		event := &entity.Event{
			Id:    eventID,
			Title: "Test Event",
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		mockInvitationDg.On("UpdateEventRegistrationInvitation", ctx, invitationID, mock.AnythingOfType("offchain_datagateway.UpdateEventRegistrationInvitationParameters")).Return(nil, errors.New("database error"))

		mockInboxMessageDg := new(MockInboxMessageDg)
		mockInboxMessageDg.On("GetInboxMessageByID", ctx, inboxMessageID).Return(originalInboxMessage, nil)

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
			InboxMessageDg:                mockInboxMessageDg,
			EventDg:                       mockEventDg,
		}

		// Act
		result, err := uc.CancelEventRegistrationInvitation(ctx, CancelEventRegistrationInvitationParameters{
			EventRegistrationInvitationID: invitationID,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInternalServer.Code, *customError.Code)
		mockInvitationDg.AssertExpectations(t)
		mockInboxMessageDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should succeed even if notification creation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		inboxMessageID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		now := time.Now()

		currentUser := &auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:             invitationID,
			EventId:        eventID,
			InboxMessageId: inboxMessageID,
			Email:          &email,
			CancelledAt:    nil,
		}

		originalInboxMessage := &entity.InboxMessage{
			Id:                   inboxMessageID,
			ReceiverCredentialId: &userID,
			ReceiverEmail:        &email,
		}

		event := &entity.Event{
			Id:    eventID,
			Title: "Test Event",
		}

		updatedInvitation := &entity.EventRegistrationInvitation{
			Id:             invitationID,
			EventId:        eventID,
			InboxMessageId: inboxMessageID,
			Email:          &email,
			CancelledAt:    &now,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		mockInvitationDg.On("UpdateEventRegistrationInvitation", ctx, invitationID, mock.AnythingOfType("offchain_datagateway.UpdateEventRegistrationInvitationParameters")).Return(updatedInvitation, nil)

		mockInboxMessageDg := new(MockInboxMessageDg)
		mockInboxMessageDg.On("GetInboxMessageByID", ctx, inboxMessageID).Return(originalInboxMessage, nil)
		mockInboxMessageDg.On("CreateInboxMessage", ctx, mock.AnythingOfType("offchain_datagateway.CreateInboxMessageParameters")).Return(nil, errors.New("failed to create notification"))

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
			InboxMessageDg:                mockInboxMessageDg,
			EventDg:                       mockEventDg,
		}

		// Act
		result, err := uc.CancelEventRegistrationInvitation(ctx, CancelEventRegistrationInvitationParameters{
			EventRegistrationInvitationID: invitationID,
		}, currentUser)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, invitationID, result.Id)
		assert.NotNil(t, result.CancelledAt)
		mockInvitationDg.AssertExpectations(t)
		mockInboxMessageDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should succeed when receiver email is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		inboxMessageID := uuid.New()
		userID := uuid.New()
		now := time.Now()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:             invitationID,
			EventId:        eventID,
			InboxMessageId: inboxMessageID,
			Email:          nil, // No email
			CancelledAt:    nil,
		}

		originalInboxMessage := &entity.InboxMessage{
			Id:                   inboxMessageID,
			ReceiverCredentialId: &userID,
			ReceiverEmail:        nil, // No email
		}

		event := &entity.Event{
			Id:    eventID,
			Title: "Test Event",
		}

		updatedInvitation := &entity.EventRegistrationInvitation{
			Id:             invitationID,
			EventId:        eventID,
			InboxMessageId: inboxMessageID,
			Email:          nil,
			CancelledAt:    &now,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		mockInvitationDg.On("UpdateEventRegistrationInvitation", ctx, invitationID, mock.AnythingOfType("offchain_datagateway.UpdateEventRegistrationInvitationParameters")).Return(updatedInvitation, nil)

		mockInboxMessageDg := new(MockInboxMessageDg)
		mockInboxMessageDg.On("GetInboxMessageByID", ctx, inboxMessageID).Return(originalInboxMessage, nil)

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
			InboxMessageDg:                mockInboxMessageDg,
			EventDg:                       mockEventDg,
		}

		// Act
		result, err := uc.CancelEventRegistrationInvitation(ctx, CancelEventRegistrationInvitationParameters{
			EventRegistrationInvitationID: invitationID,
		}, currentUser)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, invitationID, result.Id)
		assert.NotNil(t, result.CancelledAt)
		mockInvitationDg.AssertExpectations(t)
		mockInboxMessageDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		// CreateInboxMessage should NOT be called
		mockInboxMessageDg.AssertNotCalled(t, "CreateInboxMessage")
	})
}
