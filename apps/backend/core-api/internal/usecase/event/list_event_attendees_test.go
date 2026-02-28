package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListEventAttendees(t *testing.T) {
	ctx := context.Background()
	eventId := uuid.New()

	t.Run("should fail when event does not exist", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("GetEventById", ctx, eventId).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event not found")))

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		// Act
		attendees, err := uc.ListEventAttendees(ctx, eventId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, attendees)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return empty list when no attendees exist", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id: eventId,
		}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		mockAttendeeDg := new(MockEventAttendeeDg)
		mockAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventId).
			Return([]entity.EventAttendee{}, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
			EventAttendeeDg:  mockAttendeeDg,
		}

		// Act
		attendees, err := uc.ListEventAttendees(ctx, eventId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, attendees)
		assert.Empty(t, attendees)
		mockEventDg.AssertExpectations(t)
		mockAttendeeDg.AssertExpectations(t)
	})

	t.Run("should successfully return list of attendees", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id: eventId,
		}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		attendee1Id := uuid.New()
		attendee2Id := uuid.New()
		credentialId1 := uuid.New()
		credentialId2 := uuid.New()

		expectedAttendees := []entity.EventAttendee{
			{
				Id:                   attendee1Id,
				EventId:              eventId,
				AttendeeCredentialId: credentialId1,
			},
			{
				Id:                   attendee2Id,
				EventId:              eventId,
				AttendeeCredentialId: credentialId2,
			},
		}

		mockAttendeeDg := new(MockEventAttendeeDg)
		mockAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventId).
			Return(expectedAttendees, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
			EventAttendeeDg:  mockAttendeeDg,
		}

		// Act
		attendees, err := uc.ListEventAttendees(ctx, eventId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, attendees)
		assert.Len(t, attendees, 2)
		assert.Equal(t, attendee1Id, attendees[0].Id)
		assert.Equal(t, attendee2Id, attendees[1].Id)
		assert.Equal(t, eventId, attendees[0].EventId)
		assert.Equal(t, eventId, attendees[1].EventId)
		mockEventDg.AssertExpectations(t)
		mockAttendeeDg.AssertExpectations(t)
	})

	t.Run("should fail when attendee retrieval fails", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id: eventId,
		}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		mockAttendeeDg := new(MockEventAttendeeDg)
		mockAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventId).
			Return(nil, errors.New("database error"))

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
			EventAttendeeDg:  mockAttendeeDg,
		}

		// Act
		attendees, err := uc.ListEventAttendees(ctx, eventId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, attendees)
		assert.Contains(t, err.Error(), "database error")
		mockEventDg.AssertExpectations(t)
		mockAttendeeDg.AssertExpectations(t)
	})
}
