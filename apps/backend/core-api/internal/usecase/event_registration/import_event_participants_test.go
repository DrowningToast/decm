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

func TestImportEventParticipants(t *testing.T) {
	t.Run("should import participants successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userID,
			EventStatus:       entity.EventStatusActive,
		}

		participants := []Participant{
			{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
			{
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane@example.com",
			},
		}

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		mockInboxMessageDg := new(MockInboxMessageDg)
		mockInboxMessageDg.On("CreateInboxMessage", ctx, mock.AnythingOfType("offchain_datagateway.CreateInboxMessageParameters")).Return(&entity.InboxMessage{Id: uuid.New()}, nil)

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("CreateEventRegistrationInvitation", ctx, mock.AnythingOfType("offchain_datagateway.CreateEventRegistrationInvitationParameters")).Return(&entity.EventRegistrationInvitation{Id: uuid.New()}, nil)

		uc := &EventRegistrationUsecase{
			EventDg:                       mockEventDg,
			InboxMessageDg:                mockInboxMessageDg,
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		result, err := uc.ImportEventParticipants(ctx, ImportEventParticipantsParameters{
			EventID:      eventID,
			Participants: participants,
		}, currentUser)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		mockEventDg.AssertExpectations(t)
		mockInboxMessageDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
		// Verify CreateInboxMessage was called twice (once per participant)
		mockInboxMessageDg.AssertNumberOfCalls(t, "CreateInboxMessage", 2)
		// Verify CreateEventRegistrationInvitation was called twice
		mockInvitationDg.AssertNumberOfCalls(t, "CreateEventRegistrationInvitation", 2)
	})

	t.Run("should return error when event not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		participants := []Participant{
			{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
		}

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(nil, errors.New("event not found"))

		uc := &EventRegistrationUsecase{
			EventDg: mockEventDg,
		}

		// Act
		result, err := uc.ImportEventParticipants(ctx, ImportEventParticipantsParameters{
			EventID:      eventID,
			Participants: participants,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInternalServer.Code, *customError.Code)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return error when user is not event owner", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()
		ownerID := uuid.New() // Different from userID

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: ownerID, // Different owner
			EventStatus:       entity.EventStatusActive,
		}

		participants := []Participant{
			{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
		}

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		uc := &EventRegistrationUsecase{
			EventDg: mockEventDg,
		}

		// Act
		result, err := uc.ImportEventParticipants(ctx, ImportEventParticipantsParameters{
			EventID:      eventID,
			Participants: participants,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customError.Code)
		assert.Contains(t, err.Error(), "user is not the event owner")
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return error when event is not active", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userID,
			EventStatus:       entity.EventStatusInactive, // Not active
		}

		participants := []Participant{
			{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
		}

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		uc := &EventRegistrationUsecase{
			EventDg: mockEventDg,
		}

		// Act
		result, err := uc.ImportEventParticipants(ctx, ImportEventParticipantsParameters{
			EventID:      eventID,
			Participants: participants,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customError.Code)
		assert.Contains(t, err.Error(), "event is not active")
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return error when participants array is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userID,
			EventStatus:       entity.EventStatusActive,
		}

		participants := []Participant{} // Empty array

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		uc := &EventRegistrationUsecase{
			EventDg: mockEventDg,
		}

		// Act
		result, err := uc.ImportEventParticipants(ctx, ImportEventParticipantsParameters{
			EventID:      eventID,
			Participants: participants,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customError.Code)
		assert.Contains(t, err.Error(), "participants array is empty")
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return error when participant email is missing", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userID,
			EventStatus:       entity.EventStatusActive,
		}

		participants := []Participant{
			{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "", // Missing email
			},
		}

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		uc := &EventRegistrationUsecase{
			EventDg: mockEventDg,
		}

		// Act
		result, err := uc.ImportEventParticipants(ctx, ImportEventParticipantsParameters{
			EventID:      eventID,
			Participants: participants,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customError.Code)
		assert.Contains(t, err.Error(), "email is required")
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return error when creating inbox message fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userID,
			EventStatus:       entity.EventStatusActive,
		}

		participants := []Participant{
			{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
		}

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		mockInboxMessageDg := new(MockInboxMessageDg)
		mockInboxMessageDg.On("CreateInboxMessage", ctx, mock.AnythingOfType("offchain_datagateway.CreateInboxMessageParameters")).Return(nil, errors.New("database error"))

		uc := &EventRegistrationUsecase{
			EventDg:        mockEventDg,
			InboxMessageDg: mockInboxMessageDg,
		}

		// Act
		result, err := uc.ImportEventParticipants(ctx, ImportEventParticipantsParameters{
			EventID:      eventID,
			Participants: participants,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInternalServer.Code, *customError.Code)
		mockEventDg.AssertExpectations(t)
		mockInboxMessageDg.AssertExpectations(t)
	})

	t.Run("should return error when creating invitation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userID,
			EventStatus:       entity.EventStatusActive,
		}

		participants := []Participant{
			{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
		}

		mockEventDg := new(MockEventDg)
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		mockInboxMessageDg := new(MockInboxMessageDg)
		mockInboxMessageDg.On("CreateInboxMessage", ctx, mock.AnythingOfType("offchain_datagateway.CreateInboxMessageParameters")).Return(&entity.InboxMessage{Id: uuid.New()}, nil)

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("CreateEventRegistrationInvitation", ctx, mock.AnythingOfType("offchain_datagateway.CreateEventRegistrationInvitationParameters")).Return(nil, errors.New("database error"))

		uc := &EventRegistrationUsecase{
			EventDg:                       mockEventDg,
			InboxMessageDg:                mockInboxMessageDg,
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		result, err := uc.ImportEventParticipants(ctx, ImportEventParticipantsParameters{
			EventID:      eventID,
			Participants: participants,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInternalServer.Code, *customError.Code)
		mockEventDg.AssertExpectations(t)
		mockInboxMessageDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
	})
}
