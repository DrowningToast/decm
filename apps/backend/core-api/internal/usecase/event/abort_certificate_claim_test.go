package event

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAbortCertificateClaim_CallsUpdateAbortedAt(t *testing.T) {
	ctx := context.Background()
	signatureId := uuid.New()

	userSigDg := new(MockUserSignatureDataGateway)
	userSigDg.On(
		"UpdateUserSignatureAbortedAt",
		ctx,
		signatureId,
		mock.AnythingOfType("time.Time"),
		entity.AbortReasonParticipantNotJoined,
	).Return(&entity.UserSignature{Id: signatureId, AbortedAt: func() *time.Time { t := time.Now(); return &t }()}, nil)

	uc := &EventUsecase{UserSignatureDg: userSigDg}
	err := uc.AbortCertificateClaim(ctx, signatureId, entity.AbortReasonParticipantNotJoined)

	require.NoError(t, err)
	userSigDg.AssertExpectations(t)
}

func TestAbortCertificateClaim_PropagatesError_OnDataGatewayFailure(t *testing.T) {
	ctx := context.Background()
	signatureId := uuid.New()

	userSigDg := new(MockUserSignatureDataGateway)
	userSigDg.On(
		"UpdateUserSignatureAbortedAt",
		ctx,
		signatureId,
		mock.AnythingOfType("time.Time"),
		entity.AbortReasonParticipantNotJoined,
	).Return(nil, errors.New("db error"))

	uc := &EventUsecase{UserSignatureDg: userSigDg}
	err := uc.AbortCertificateClaim(ctx, signatureId, entity.AbortReasonParticipantNotJoined)

	assert.Error(t, err)
	userSigDg.AssertExpectations(t)
}
