package event_test

import (
	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	"apps/backend/core-api/internal/entity"
	event_usecase "apps/backend/core-api/internal/usecase/event"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Data Gateways
type mockEventAttendeeDg struct {
	mock.Mock
}

func (m *mockEventAttendeeDg) GetEventAttendeeByEventIdAndCredentialId(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*entity.EventAttendee, error) {
	return nil, nil
}

func (m *mockEventAttendeeDg) GetEventAttendeeWithSignature(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*eventdatagateway.EventAttendeeWithSignature, error) {
	args := m.Called(ctx, eventID, credentialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eventdatagateway.EventAttendeeWithSignature), args.Error(1)
}

func (m *mockEventAttendeeDg) AddParticipant(ctx context.Context, params eventdatagateway.AddParticipantParameters) (*entity.EventAttendee, error) {
	return nil, nil
}

func (m *mockEventAttendeeDg) ListEventAttendeesByEventID(ctx context.Context, eventId uuid.UUID) ([]entity.EventAttendee, error) {
	return nil, nil
}

func (m *mockEventAttendeeDg) DeleteEventAttendeeById(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockEventAttendeeDg) DeleteEventAttendeeByEventIDAndCredentialID(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) error {
	args := m.Called(ctx, eventId, credentialId)
	return args.Error(0)
}

func (m *mockEventAttendeeDg) HasPendingEventJoinByEventAndCredential(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (bool, error) {
	args := m.Called(ctx, eventId, credentialId)
	return args.Bool(0), args.Error(1)
}

type mockUserSignatureDg struct {
	mock.Mock
}

func (m *mockUserSignatureDg) GetUserSignatureByID(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

// Implement other methods to satisfy interface
func (m *mockUserSignatureDg) CreateUserSignature(ctx context.Context, params offchain_datagateway.CreateUserSignatureParameters) (*entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetUserSignaturesByCredentialID(ctx context.Context, id uuid.UUID) ([]entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) UpdateUserSignatureBroadcastedAt(ctx context.Context, id uuid.UUID, b *time.Time) (*entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetPendingUserSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetStaleUserSignatures(ctx context.Context, t time.Time) ([]entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetBroadcastedUserSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetUserSignaturesByDeadlineBlockRange(ctx context.Context, min, max int32) ([]entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetUserSignaturesExpiringBefore(ctx context.Context, t time.Time) ([]entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetEventJoinSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetCertificateClaimSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetOrphanedUserSignatures(ctx context.Context, t time.Time) ([]entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetUserSignatureWithUsageDetails(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) GetRecentUserSignatureActivity(ctx context.Context, t time.Time) ([]entity.UserSignature, error) {
	return nil, nil
}

func (m *mockUserSignatureDg) DeleteUserSignature(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockUserSignatureDg) CountUserSignaturesByCredentialID(ctx context.Context, id uuid.UUID) (int64, error) {
	return 0, nil
}

func (m *mockUserSignatureDg) CountPendingUserSignatures(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockUserSignatureDg) CountBroadcastedUserSignatures(ctx context.Context) (int64, error) {
	return 0, nil
}

type mockEventCertificateDg struct {
	mock.Mock
}

func (m *mockEventCertificateDg) CreateEventCertificate(ctx context.Context, params eventdatagateway.CreateEventCertificateParameters) (*entity.EventCertificate, error) {
	return nil, nil
}

func (m *mockEventCertificateDg) GetEventCertificateByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificate), args.Error(1)
}

func (m *mockEventCertificateDg) GetEventCertificateByInboxMessageID(ctx context.Context, id uuid.UUID) (*entity.EventCertificate, error) {
	return nil, nil
}

func (m *mockEventCertificateDg) GetEventCertificatesByEventID(ctx context.Context, id uuid.UUID) ([]*entity.EventCertificate, error) {
	return nil, nil
}

func (m *mockEventCertificateDg) GetAllEventCertificateIDsByEventID(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

func (m *mockEventCertificateDg) GetClaimedCertificatesByEventID(ctx context.Context, id uuid.UUID) ([]*entity.EventCertificate, error) {
	return nil, nil
}

func (m *mockEventCertificateDg) GetUnclaimedReadyCertificatesByEventID(ctx context.Context, id uuid.UUID) ([]*entity.EventCertificate, error) {
	return nil, nil
}

func (m *mockEventCertificateDg) GetClaimedCertificatesByCredentialID(ctx context.Context, id uuid.UUID, email *string) ([]*entity.EventCertificate, error) {
	args := m.Called(ctx, id, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificate), args.Error(1)
}

func (m *mockEventCertificateDg) GetUnclaimedReadyCertificatesByCredentialID(ctx context.Context, id uuid.UUID, email *string) ([]*entity.EventCertificate, error) {
	args := m.Called(ctx, id, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificate), args.Error(1)
}

func (m *mockEventCertificateDg) UpdateEventCertificate(ctx context.Context, id uuid.UUID, params eventdatagateway.UpdateEventCertificateParameters) (*entity.EventCertificate, error) {
	return nil, nil
}

func (m *mockEventCertificateDg) UpdateEventCertificateInboxMessageID(ctx context.Context, id uuid.UUID, inboxID uuid.UUID) (*entity.EventCertificate, error) {
	return nil, nil
}

func (m *mockEventCertificateDg) DeleteEventCertificate(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockEventCertificateDg) GetEventCertificateWithSignature(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*eventdatagateway.EventCertificateWithSignature, error) {
	args := m.Called(ctx, eventID, credentialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eventdatagateway.EventCertificateWithSignature), args.Error(1)
}

type mockCertificateShareDg struct {
	mock.Mock
}

func (m *mockCertificateShareDg) CreateCertificateShare(ctx context.Context, params eventdatagateway.CreateCertificateShareParameters) (*entity.CertificateShare, error) {
	return nil, nil
}

func (m *mockCertificateShareDg) GetCertificateShareByHandle(ctx context.Context, handle string) (*entity.CertificateShare, error) {
	return nil, nil
}

func (m *mockCertificateShareDg) GetCertificateShareByID(ctx context.Context, id uuid.UUID) (*entity.CertificateShare, error) {
	return nil, nil
}

func (m *mockCertificateShareDg) GetCertificateShareByEventCertificateID(ctx context.Context, eventCertificateID uuid.UUID) (*entity.CertificateShare, error) {
	args := m.Called(ctx, eventCertificateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CertificateShare), args.Error(1)
}

func (m *mockCertificateShareDg) UpdateCertificateShare(ctx context.Context, id uuid.UUID, params eventdatagateway.UpdateCertificateShareParameters) (*entity.CertificateShare, error) {
	return nil, nil
}

func TestIsUserJoined(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	credentialID := uuid.New()
	now := time.Now()

	t.Run("User not joined (no attendee row)", func(t *testing.T) {
		mockAttendeeDg := new(mockEventAttendeeDg)
		uc := &event_usecase.EventUsecase{
			EventAttendeeDg: mockAttendeeDg,
		}

		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventID, credentialID).Return(nil, nil)

		status, err := uc.IsUserJoined(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.Nil(t, status)
		mockAttendeeDg.AssertExpectations(t)
	})

	t.Run("User joined with no signature row — broadcast_at must be nil", func(t *testing.T) {
		// Regression: previously this case incorrectly set IsBroadcasted=true
		// and BroadcastedAt=&joinedAt when UserSignature is nil.
		mockAttendeeDg := new(mockEventAttendeeDg)
		uc := &event_usecase.EventUsecase{
			EventAttendeeDg: mockAttendeeDg,
		}

		attendee := &eventdatagateway.EventAttendeeWithSignature{
			Id:            uuid.New(),
			JoinedAt:      now,
			IsAccepted:    true,
			UserSignature: nil,
		}

		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventID, credentialID).Return(attendee, nil)

		status, err := uc.IsUserJoined(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.False(t, status.IsBroadcasted)
		assert.Nil(t, status.BroadcastedAt)
		assert.Equal(t, now, status.JoinedAt)
		assert.True(t, status.IsAccepted)
		assert.Empty(t, status.Signature)
		assert.Empty(t, status.SignMessage)
		assert.Nil(t, status.DeadlineBlock)
		assert.Nil(t, status.EstimatedDeadline)

		mockAttendeeDg.AssertExpectations(t)
	})

	t.Run("User joined with no signature row, not accepted", func(t *testing.T) {
		mockAttendeeDg := new(mockEventAttendeeDg)
		uc := &event_usecase.EventUsecase{
			EventAttendeeDg: mockAttendeeDg,
		}

		attendee := &eventdatagateway.EventAttendeeWithSignature{
			Id:            uuid.New(),
			JoinedAt:      now,
			IsAccepted:    false,
			UserSignature: nil,
		}

		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventID, credentialID).Return(attendee, nil)

		status, err := uc.IsUserJoined(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.False(t, status.IsBroadcasted)
		assert.Nil(t, status.BroadcastedAt)
		assert.False(t, status.IsAccepted)

		mockAttendeeDg.AssertExpectations(t)
	})

	t.Run("User joined, signature present but not yet broadcasted", func(t *testing.T) {
		mockAttendeeDg := new(mockEventAttendeeDg)
		uc := &event_usecase.EventUsecase{
			EventAttendeeDg: mockAttendeeDg,
		}

		attendee := &eventdatagateway.EventAttendeeWithSignature{
			Id:         uuid.New(),
			JoinedAt:   now,
			IsAccepted: true,
			UserSignature: &entity.UserSignature{
				Signature:     "0xsig",
				SignMessage:   "message",
				BroadcastedAt: nil,
			},
		}

		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventID, credentialID).Return(attendee, nil)

		status, err := uc.IsUserJoined(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.False(t, status.IsBroadcasted)
		assert.Nil(t, status.BroadcastedAt)
		assert.Equal(t, "0xsig", status.Signature)
		assert.Equal(t, "message", status.SignMessage)
		assert.Equal(t, now, status.JoinedAt)
		assert.True(t, status.IsAccepted)

		mockAttendeeDg.AssertExpectations(t)
	})

	t.Run("User joined, signature broadcasted", func(t *testing.T) {
		mockAttendeeDg := new(mockEventAttendeeDg)
		uc := &event_usecase.EventUsecase{
			EventAttendeeDg: mockAttendeeDg,
		}

		broadcastTime := now.Add(time.Minute)
		attendee := &eventdatagateway.EventAttendeeWithSignature{
			Id:         uuid.New(),
			JoinedAt:   now,
			IsAccepted: true,
			UserSignature: &entity.UserSignature{
				Signature:     "0xsig",
				SignMessage:   "message",
				BroadcastedAt: &broadcastTime,
			},
		}

		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventID, credentialID).Return(attendee, nil)

		status, err := uc.IsUserJoined(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.True(t, status.IsBroadcasted)
		assert.Equal(t, &broadcastTime, status.BroadcastedAt)
		assert.Equal(t, "0xsig", status.Signature)
		assert.Equal(t, "message", status.SignMessage)
		assert.Equal(t, now, status.JoinedAt)
		assert.True(t, status.IsAccepted)
		assert.Nil(t, status.DeadlineBlock)
		assert.Nil(t, status.EstimatedDeadline)

		mockAttendeeDg.AssertExpectations(t)
	})

	t.Run("User joined, broadcasted with deadline block and estimated deadline", func(t *testing.T) {
		mockAttendeeDg := new(mockEventAttendeeDg)
		uc := &event_usecase.EventUsecase{
			EventAttendeeDg: mockAttendeeDg,
		}

		broadcastTime := now.Add(time.Minute)
		estimatedDeadline := now.Add(24 * time.Hour)
		deadlineBlock := int32(12345678)
		attendee := &eventdatagateway.EventAttendeeWithSignature{
			Id:         uuid.New(),
			JoinedAt:   now,
			IsAccepted: true,
			UserSignature: &entity.UserSignature{
				Signature:         "0xsig",
				SignMessage:       "message",
				BroadcastedAt:     &broadcastTime,
				DeadlineBlock:     &deadlineBlock,
				EstimatedDeadline: &estimatedDeadline,
			},
		}

		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventID, credentialID).Return(attendee, nil)

		status, err := uc.IsUserJoined(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.True(t, status.IsBroadcasted)
		assert.Equal(t, &broadcastTime, status.BroadcastedAt)
		assert.Equal(t, &deadlineBlock, status.DeadlineBlock)
		assert.Equal(t, &estimatedDeadline, status.EstimatedDeadline)

		mockAttendeeDg.AssertExpectations(t)
	})

	t.Run("User joined, pending broadcast with deadline block but no broadcast time", func(t *testing.T) {
		mockAttendeeDg := new(mockEventAttendeeDg)
		uc := &event_usecase.EventUsecase{
			EventAttendeeDg: mockAttendeeDg,
		}

		estimatedDeadline := now.Add(24 * time.Hour)
		deadlineBlock := int32(12345678)
		attendee := &eventdatagateway.EventAttendeeWithSignature{
			Id:         uuid.New(),
			JoinedAt:   now,
			IsAccepted: true,
			UserSignature: &entity.UserSignature{
				Signature:         "0xsig",
				SignMessage:       "message",
				BroadcastedAt:     nil,
				DeadlineBlock:     &deadlineBlock,
				EstimatedDeadline: &estimatedDeadline,
			},
		}

		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventID, credentialID).Return(attendee, nil)

		status, err := uc.IsUserJoined(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.False(t, status.IsBroadcasted)
		assert.Nil(t, status.BroadcastedAt)
		assert.Equal(t, &deadlineBlock, status.DeadlineBlock)
		assert.Equal(t, &estimatedDeadline, status.EstimatedDeadline)

		mockAttendeeDg.AssertExpectations(t)
	})

	t.Run("User joined, signature marked as expired — SignatureExpiredAt is populated and IsBroadcasted is false", func(t *testing.T) {
		mockAttendeeDg := new(mockEventAttendeeDg)
		uc := &event_usecase.EventUsecase{
			EventAttendeeDg: mockAttendeeDg,
		}

		expiredAt := now.Add(-1 * time.Hour) // deadline passed 1 hour ago
		estimatedDeadline := now.Add(-1 * time.Hour)
		deadlineBlock := int32(99999999)
		attendee := &eventdatagateway.EventAttendeeWithSignature{
			Id:         uuid.New(),
			JoinedAt:   now,
			IsAccepted: true,
			UserSignature: &entity.UserSignature{
				Signature:         "0xsig",
				SignMessage:       "message",
				BroadcastedAt:     nil,
				DeadlineBlock:     &deadlineBlock,
				EstimatedDeadline: &estimatedDeadline,
				MarkAsExpiredAt:   &expiredAt,
			},
		}

		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventID, credentialID).Return(attendee, nil)

		status, err := uc.IsUserJoined(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.False(t, status.IsBroadcasted)
		assert.Nil(t, status.BroadcastedAt)
		assert.Equal(t, &expiredAt, status.SignatureExpiredAt)
		assert.Equal(t, "0xsig", status.Signature)
		assert.Equal(t, "message", status.SignMessage)

		mockAttendeeDg.AssertExpectations(t)
	})
}
