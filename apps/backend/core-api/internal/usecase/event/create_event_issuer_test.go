package event

import (
	"apps/backend/common/customerror"
	"context"
	"decm-database/go/generated"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateEventIssuer(t *testing.T) {
	ctx := context.Background()
	eventId := uuid.New()
	issuerCredentialId := uuid.New()

	t.Run("should successfully create event issuer", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockEventIssuerDataGateway)

		expectedIssuer := &generated.EventIssuer{
			ID:                 uuid.New(),
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
			IsSigned:           0,
		}

		mockIssuerDg.On("CreateEventIssuer", ctx, generated.CreateEventIssuerParams{
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
			IsSigned:           0,
		}).Return(expectedIssuer, nil)

		uc := &EventUsecase{
			EventIssuerDataGateway: mockIssuerDg,
		}

		params := CreateEventIssuerParams{
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
			IsSigned:           0,
		}

		// Act
		result, err := uc.CreateEventIssuer(ctx, params)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, eventId, result.EventID)
		assert.Equal(t, issuerCredentialId, result.IssuerCredentialID)
		assert.Equal(t, int32(0), result.IsSigned)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should fail when database creation fails", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockEventIssuerDataGateway)

		mockIssuerDg.On("CreateEventIssuer", ctx, generated.CreateEventIssuerParams{
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
			IsSigned:           0,
		}).Return(nil, errors.New("database error"))

		uc := &EventUsecase{
			EventIssuerDataGateway: mockIssuerDg,
		}

		params := CreateEventIssuerParams{
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
			IsSigned:           0,
		}

		// Act
		result, err := uc.CreateEventIssuer(ctx, params)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrInternalServer.Code, *customErr.Code)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should create issuer with signed status 1", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockEventIssuerDataGateway)

		expectedIssuer := &generated.EventIssuer{
			ID:                 uuid.New(),
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
			IsSigned:           1,
		}

		mockIssuerDg.On("CreateEventIssuer", ctx, generated.CreateEventIssuerParams{
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
			IsSigned:           1,
		}).Return(expectedIssuer, nil)

		uc := &EventUsecase{
			EventIssuerDataGateway: mockIssuerDg,
		}

		params := CreateEventIssuerParams{
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
			IsSigned:           1,
		}

		// Act
		result, err := uc.CreateEventIssuer(ctx, params)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.IsSigned)
		mockIssuerDg.AssertExpectations(t)
	})
}
