package event_registration

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasPendingEventJoin_ReturnsTrue_WhenPendingJoinExists(t *testing.T) {
	ctx := context.Background()
	eventId := uuid.New()
	credentialId := uuid.New()

	attendeeDg := new(MockEventAttendeeDg)
	attendeeDg.On("HasPendingEventJoinByEventAndCredential", ctx, eventId, credentialId).Return(true, nil)

	uc := &EventRegistrationUsecase{EventAttendeeDg: attendeeDg}
	result, err := uc.HasPendingEventJoin(ctx, eventId, credentialId)

	require.NoError(t, err)
	assert.True(t, result)
	attendeeDg.AssertExpectations(t)
}

func TestHasPendingEventJoin_ReturnsFalse_WhenNoPendingJoin(t *testing.T) {
	ctx := context.Background()
	eventId := uuid.New()
	credentialId := uuid.New()

	attendeeDg := new(MockEventAttendeeDg)
	attendeeDg.On("HasPendingEventJoinByEventAndCredential", ctx, eventId, credentialId).Return(false, nil)

	uc := &EventRegistrationUsecase{EventAttendeeDg: attendeeDg}
	result, err := uc.HasPendingEventJoin(ctx, eventId, credentialId)

	require.NoError(t, err)
	assert.False(t, result)
	attendeeDg.AssertExpectations(t)
}

func TestHasPendingEventJoin_PropagatesError_OnDataGatewayFailure(t *testing.T) {
	ctx := context.Background()
	eventId := uuid.New()
	credentialId := uuid.New()

	attendeeDg := new(MockEventAttendeeDg)
	attendeeDg.On("HasPendingEventJoinByEventAndCredential", ctx, eventId, credentialId).Return(false, errors.New("db error"))

	uc := &EventRegistrationUsecase{EventAttendeeDg: attendeeDg}
	result, err := uc.HasPendingEventJoin(ctx, eventId, credentialId)

	require.Error(t, err)
	assert.False(t, result)
	attendeeDg.AssertExpectations(t)
}
