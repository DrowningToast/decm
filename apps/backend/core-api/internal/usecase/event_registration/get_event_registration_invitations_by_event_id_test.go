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
	"github.com/stretchr/testify/require"
)

func TestGetEventRegistrationInvitationsByEventId(t *testing.T) {
	t.Run("should get invitations successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		invitations := []*entity.EventRegistrationInvitation{
			{
				Id:      uuid.New(),
				EventId: eventID,
			},
			{
				Id:      uuid.New(),
				EventId: eventID,
			},
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationsByEventID", ctx, eventID).Return(invitations, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		result, err := uc.GetEventRegistrationInvitationsByEventId(ctx, GetEventRegistrationInvitationsByEventIDParameters{
			EventId: eventID,
		}, currentUser)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, eventID, result[0].EventId)
		assert.Equal(t, eventID, result[1].EventId)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return empty list when no invitations found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		invitations := []*entity.EventRegistrationInvitation{}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationsByEventID", ctx, eventID).Return(invitations, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		result, err := uc.GetEventRegistrationInvitationsByEventId(ctx, GetEventRegistrationInvitationsByEventIDParameters{
			EventId: eventID,
		}, currentUser)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when database fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		userID := uuid.New()

		currentUser := &auth.JwtClaims{
			UserId: userID,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationsByEventID", ctx, eventID).Return(nil, errors.New("database error"))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		result, err := uc.GetEventRegistrationInvitationsByEventId(ctx, GetEventRegistrationInvitationsByEventIDParameters{
			EventId: eventID,
		}, currentUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInternalServer.Code, *customError.Code)
		mockInvitationDg.AssertExpectations(t)
	})
}
