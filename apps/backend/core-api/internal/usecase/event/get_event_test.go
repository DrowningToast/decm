package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"errors"
	"testing"
	"time"

	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListEventsByOwnerCredentialID(t *testing.T) {
	ctx := context.Background()
	ownerCredentialId := uuid.New()
	limitCount := int32(10)
	offsetCount := int32(0)

	t.Run("should successfully return list of events", func(t *testing.T) {
		// Arrange
		eventId1 := uuid.New()
		eventId2 := uuid.New()

		expectedEvents := []*entity.Event{
			{
				Id:                eventId1,
				OwnerCredentialId: ownerCredentialId,
				Title:             "Test Event 1",
			},
			{
				Id:                eventId2,
				OwnerCredentialId: ownerCredentialId,
				Title:             "Test Event 2",
			},
		}

		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("ListEventsByOwnerCredentialID", ctx, ownerCredentialId, limitCount, offsetCount).
			Return(expectedEvents, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		// Act
		events, err := uc.ListEventsByOwnerCredentialID(ctx, ownerCredentialId, limitCount, offsetCount)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.Len(t, events, 2)
		assert.Equal(t, eventId1, events[0].Id)
		assert.Equal(t, eventId2, events[1].Id)
		assert.Equal(t, ownerCredentialId, events[0].OwnerCredentialId)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return empty list when no events exist", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("ListEventsByOwnerCredentialID", ctx, ownerCredentialId, limitCount, offsetCount).
			Return([]*entity.Event{}, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		// Act
		events, err := uc.ListEventsByOwnerCredentialID(ctx, ownerCredentialId, limitCount, offsetCount)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.Empty(t, events)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should fail when database query fails", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("ListEventsByOwnerCredentialID", ctx, ownerCredentialId, limitCount, offsetCount).
			Return(nil, errors.New("database error"))

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		// Act
		events, err := uc.ListEventsByOwnerCredentialID(ctx, ownerCredentialId, limitCount, offsetCount)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "database error")
		mockEventDg.AssertExpectations(t)
	})
}

func TestGetEventById(t *testing.T) {
	ctx := context.Background()
	eventId := uuid.New()

	t.Run("should successfully return event", func(t *testing.T) {
		// Arrange
		expectedEvent := &entity.Event{
			Id:    eventId,
			Title: "Test Event",
		}

		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("GetEventById", ctx, eventId).Return(expectedEvent, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		// Act
		event, err := uc.GetEventById(ctx, eventId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, eventId, event.Id)
		assert.Equal(t, "Test Event", event.Title)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should fail when event does not exist", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("GetEventById", ctx, eventId).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event not found")))

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		// Act
		event, err := uc.GetEventById(ctx, eventId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
		mockEventDg.AssertExpectations(t)
	})
}

func TestGetEventCertificatesByEventId(t *testing.T) {
	ctx := context.Background()
	eventId := uuid.New()
	userId := uuid.New()

	t.Run("should successfully return certificates", func(t *testing.T) {
		// Arrange
		certId1 := uuid.New()
		certId2 := uuid.New()

		expectedCertificates := []*entity.EventCertificate{
			{
				Id:      certId1,
				EventId: eventId,
			},
			{
				Id:      certId2,
				EventId: eventId,
			},
		}

		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventId).
			Return(expectedCertificates, nil)

		uc := &EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		certificates, err := uc.GetEventCertificatesByEventId(ctx, eventId, currentUser)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, certificates)
		assert.Len(t, certificates, 2)
		assert.Equal(t, certId1, certificates[0].Id)
		assert.Equal(t, certId2, certificates[1].Id)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return empty list when no certificates exist", func(t *testing.T) {
		// Arrange
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventId).
			Return([]*entity.EventCertificate{}, nil)

		uc := &EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		certificates, err := uc.GetEventCertificatesByEventId(ctx, eventId, currentUser)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, certificates)
		assert.Empty(t, certificates)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should fail when database query fails", func(t *testing.T) {
		// Arrange
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventId).
			Return(nil, errors.New("database error"))

		uc := &EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		certificates, err := uc.GetEventCertificatesByEventId(ctx, eventId, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, certificates)
		assert.Contains(t, err.Error(), "database error")
		mockCertDg.AssertExpectations(t)
	})
}

func TestGetEventViewModelByEventId(t *testing.T) {
	ctx := context.Background()
	eventId := uuid.New()
	userId := uuid.New()
	email := "test@example.com"
	walletAddress := "0x1234567890abcdef"

	currentUser := &auth.JwtClaims{
		UserId:        userId,
		Email:         &email,
		WalletAddress: walletAddress,
	}

	t.Run("should successfully return event view model for invited and joined user", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockInvitationDg := new(MockEventRegistrationInvitationDg)
		mockAttendeeDg := new(MockEventAttendeeDg)
		mockS3Dg := new(MockS3DataGateway)

		// Setup event data
		event := &entity.Event{
			Id:                eventId,
			Title:             "Test Event",
			BannerStorageKey:  "banner-key",
			IconStorageKey:    "icon-key",
			MaxAttendees:      100,
			AttendeesCount:    50,
			OwnerCredentialId: uuid.New(),
		}

		registrationConfig := &entity.EventRegistrationConfig{
			ID:                             uuid.New(),
			EventID:                        eventId,
			IsIdentityVerificationRequired: true,
		}

		eventContract := &entity.EventContract{
			ID:                           uuid.New(),
			EventId:                      eventId,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      uuid.New(),
			EventId: eventId,
		}

		broadcastedAt := time.Now()
		joinedAt := time.Now()
		deadlineBlock := int32(1000)
		estimatedDeadline := time.Now().Add(1 * time.Hour)

		attendeeWithSignature := &event_datagateway.EventAttendeeWithSignature{
			Id:            uuid.New(),
			EventId:       eventId,
			CredentialId:  userId,
			WalletAddress: walletAddress,
			JoinedAt:      joinedAt,
			IsAccepted:    true,
			UserSignature: &entity.UserSignature{
				Id:                         uuid.New(),
				Signature:                  "test-signature",
				SignMessage:                "test-sign-message",
				BroadcastedAt:              &broadcastedAt,
				DeadlineBlock:              &deadlineBlock,
				EstimatedDeadline:          &estimatedDeadline,
				AuthenticationCredentialId: userId,
			},
		}

		mockEventDg.On("GetViewModelById", ctx, eventId).
			Return(event, registrationConfig, eventContract, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventId, userId, &email, &walletAddress).
			Return(invitation, nil, nil)
		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventId, userId).
			Return(attendeeWithSignature, nil)
		mockS3Dg.On("GetPresignedURL", ctx, "banner-key").Return("https://banner-url", nil)
		mockS3Dg.On("GetPresignedURL", ctx, "icon-key").Return("https://icon-url", nil)

		uc := &EventUsecase{
			EventDataGateway:              mockEventDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			EventAttendeeDg:               mockAttendeeDg,
			S3DataGateway:                 mockS3Dg,
		}

		// Act
		result, err := uc.GetEventViewModelByEventId(ctx, eventId, currentUser)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, eventId, result.EventResponse.Id)
		assert.Equal(t, "Test Event", result.EventResponse.Title)
		assert.Equal(t, "https://banner-url", result.EventResponse.BannerPresignedURL)
		assert.Equal(t, "https://icon-url", result.EventResponse.IconPresignedURL)
		assert.True(t, result.IsInvited)
		assert.True(t, result.IsJoined)
		assert.False(t, result.IsFull) // 50 < 100
		assert.Equal(t, "test-signature", *result.JoinedSignature)
		assert.Equal(t, "test-sign-message", *result.JoinedSignMessage)
		assert.Equal(t, joinedAt, *result.JoinedAt)
		assert.True(t, *result.JoinedIsAccepted)
		assert.Equal(t, broadcastedAt, *result.BroadcastedAt)
		assert.Equal(t, deadlineBlock, *result.BroadcastedAtDeadlineBlock)
		assert.Equal(t, estimatedDeadline, *result.BroadcastedAtEstimatedDeadline)

		mockEventDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
		mockAttendeeDg.AssertExpectations(t)
		mockS3Dg.AssertExpectations(t)
	})

	t.Run("should successfully return event view model for invited but not joined user", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockInvitationDg := new(MockEventRegistrationInvitationDg)
		mockAttendeeDg := new(MockEventAttendeeDg)
		mockS3Dg := new(MockS3DataGateway)

		event := &entity.Event{
			Id:                eventId,
			Title:             "Test Event",
			BannerStorageKey:  "banner-key",
			IconStorageKey:    "icon-key",
			MaxAttendees:      100,
			AttendeesCount:    50,
			OwnerCredentialId: uuid.New(),
		}

		registrationConfig := &entity.EventRegistrationConfig{
			ID:      uuid.New(),
			EventID: eventId,
		}

		eventContract := &entity.EventContract{
			ID:                           uuid.New(),
			EventId:                      eventId,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      uuid.New(),
			EventId: eventId,
		}

		mockEventDg.On("GetViewModelById", ctx, eventId).
			Return(event, registrationConfig, eventContract, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventId, userId, &email, &walletAddress).
			Return(invitation, nil, nil)
		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventId, userId).
			Return(nil, nil) // Not joined
		mockS3Dg.On("GetPresignedURL", ctx, "banner-key").Return("https://banner-url", nil)
		mockS3Dg.On("GetPresignedURL", ctx, "icon-key").Return("https://icon-url", nil)

		uc := &EventUsecase{
			EventDataGateway:              mockEventDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			EventAttendeeDg:               mockAttendeeDg,
			S3DataGateway:                 mockS3Dg,
		}

		// Act
		result, err := uc.GetEventViewModelByEventId(ctx, eventId, currentUser)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsInvited)
		assert.False(t, result.IsJoined)
		assert.Nil(t, result.JoinedSignature)
		assert.Nil(t, result.JoinedSignMessage)
		assert.Nil(t, result.JoinedAt)

		mockEventDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
		mockAttendeeDg.AssertExpectations(t)
		mockS3Dg.AssertExpectations(t)
	})

	t.Run("should successfully return event view model for not invited user", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockInvitationDg := new(MockEventRegistrationInvitationDg)
		mockS3Dg := new(MockS3DataGateway)

		event := &entity.Event{
			Id:                eventId,
			Title:             "Test Event",
			BannerStorageKey:  "banner-key",
			IconStorageKey:    "icon-key",
			MaxAttendees:      0, // No max
			AttendeesCount:    50,
			OwnerCredentialId: uuid.New(),
		}

		registrationConfig := &entity.EventRegistrationConfig{
			ID:      uuid.New(),
			EventID: eventId,
		}

		eventContract := &entity.EventContract{
			ID:                           uuid.New(),
			EventId:                      eventId,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
		}

		mockEventDg.On("GetViewModelById", ctx, eventId).
			Return(event, registrationConfig, eventContract, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventId, userId, &email, &walletAddress).
			Return(nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("invitation not found")))
		mockS3Dg.On("GetPresignedURL", ctx, "banner-key").Return("https://banner-url", nil)
		mockS3Dg.On("GetPresignedURL", ctx, "icon-key").Return("https://icon-url", nil)

		uc := &EventUsecase{
			EventDataGateway:              mockEventDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			S3DataGateway:                 mockS3Dg,
		}

		// Act
		result, err := uc.GetEventViewModelByEventId(ctx, eventId, currentUser)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsInvited)
		assert.False(t, result.IsJoined)
		assert.False(t, result.IsFull) // MaxAttendees is 0

		mockEventDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
		mockS3Dg.AssertExpectations(t)
	})

	t.Run("should calculate IsFull correctly when event is at capacity", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockInvitationDg := new(MockEventRegistrationInvitationDg)
		mockS3Dg := new(MockS3DataGateway)

		event := &entity.Event{
			Id:                eventId,
			Title:             "Test Event",
			BannerStorageKey:  "banner-key",
			IconStorageKey:    "icon-key",
			MaxAttendees:      50,
			AttendeesCount:    50, // Full
			OwnerCredentialId: uuid.New(),
		}

		registrationConfig := &entity.EventRegistrationConfig{
			ID:      uuid.New(),
			EventID: eventId,
		}

		eventContract := &entity.EventContract{
			ID:                           uuid.New(),
			EventId:                      eventId,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
		}

		mockEventDg.On("GetViewModelById", ctx, eventId).
			Return(event, registrationConfig, eventContract, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventId, userId, &email, &walletAddress).
			Return(nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("invitation not found")))
		mockS3Dg.On("GetPresignedURL", ctx, "banner-key").Return("https://banner-url", nil)
		mockS3Dg.On("GetPresignedURL", ctx, "icon-key").Return("https://icon-url", nil)

		uc := &EventUsecase{
			EventDataGateway:              mockEventDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			S3DataGateway:                 mockS3Dg,
		}

		// Act
		result, err := uc.GetEventViewModelByEventId(ctx, eventId, currentUser)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsFull) // 50 >= 50

		mockEventDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
		mockS3Dg.AssertExpectations(t)
	})

	t.Run("should fail when GetViewModelById fails", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)

		mockEventDg.On("GetViewModelById", ctx, eventId).
			Return(nil, nil, nil, errors.New("database error"))

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		// Act
		result, err := uc.GetEventViewModelByEventId(ctx, eventId, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")

		mockEventDg.AssertExpectations(t)
	})

	t.Run("should fail when GetEventRegistrationInvitationByEventIDAndCredential fails with non-NotFound error", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockInvitationDg := new(MockEventRegistrationInvitationDg)

		event := &entity.Event{
			Id:                eventId,
			Title:             "Test Event",
			BannerStorageKey:  "banner-key",
			IconStorageKey:    "icon-key",
			OwnerCredentialId: uuid.New(),
		}

		registrationConfig := &entity.EventRegistrationConfig{
			ID:      uuid.New(),
			EventID: eventId,
		}

		eventContract := &entity.EventContract{
			ID:                           uuid.New(),
			EventId:                      eventId,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
		}

		mockEventDg.On("GetViewModelById", ctx, eventId).
			Return(event, registrationConfig, eventContract, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventId, userId, &email, &walletAddress).
			Return(nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("database error")))

		uc := &EventUsecase{
			EventDataGateway:              mockEventDg,
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		result, err := uc.GetEventViewModelByEventId(ctx, eventId, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)

		mockEventDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should fail when IsUserJoined fails with non-NotFound error", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockInvitationDg := new(MockEventRegistrationInvitationDg)
		mockAttendeeDg := new(MockEventAttendeeDg)

		event := &entity.Event{
			Id:                eventId,
			Title:             "Test Event",
			BannerStorageKey:  "banner-key",
			IconStorageKey:    "icon-key",
			OwnerCredentialId: uuid.New(),
		}

		registrationConfig := &entity.EventRegistrationConfig{
			ID:      uuid.New(),
			EventID: eventId,
		}

		eventContract := &entity.EventContract{
			ID:                           uuid.New(),
			EventId:                      eventId,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      uuid.New(),
			EventId: eventId,
		}

		mockEventDg.On("GetViewModelById", ctx, eventId).
			Return(event, registrationConfig, eventContract, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventId, userId, &email, &walletAddress).
			Return(invitation, nil, nil)
		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventId, userId).
			Return(nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("database error")))

		uc := &EventUsecase{
			EventDataGateway:              mockEventDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			EventAttendeeDg:               mockAttendeeDg,
		}

		// Act
		result, err := uc.GetEventViewModelByEventId(ctx, eventId, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)

		mockEventDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
		mockAttendeeDg.AssertExpectations(t)
	})

	t.Run("should fail when ToEventResponse fails (S3 GetPresignedURL for banner fails)", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockInvitationDg := new(MockEventRegistrationInvitationDg)
		mockS3Dg := new(MockS3DataGateway)

		event := &entity.Event{
			Id:                eventId,
			Title:             "Test Event",
			BannerStorageKey:  "banner-key",
			IconStorageKey:    "icon-key",
			OwnerCredentialId: uuid.New(),
		}

		registrationConfig := &entity.EventRegistrationConfig{
			ID:      uuid.New(),
			EventID: eventId,
		}

		eventContract := &entity.EventContract{
			ID:                           uuid.New(),
			EventId:                      eventId,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
		}

		mockEventDg.On("GetViewModelById", ctx, eventId).
			Return(event, registrationConfig, eventContract, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventId, userId, &email, &walletAddress).
			Return(nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("invitation not found")))
		mockS3Dg.On("GetPresignedURL", ctx, "banner-key").Return("", errors.New("S3 error"))

		uc := &EventUsecase{
			EventDataGateway:              mockEventDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			S3DataGateway:                 mockS3Dg,
		}

		// Act
		result, err := uc.GetEventViewModelByEventId(ctx, eventId, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "S3 error")

		mockEventDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
		mockS3Dg.AssertExpectations(t)
	})

	t.Run("should fail when ToEventResponse fails (S3 GetPresignedURL for icon fails)", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockInvitationDg := new(MockEventRegistrationInvitationDg)
		mockS3Dg := new(MockS3DataGateway)

		event := &entity.Event{
			Id:                eventId,
			Title:             "Test Event",
			BannerStorageKey:  "banner-key",
			IconStorageKey:    "icon-key",
			OwnerCredentialId: uuid.New(),
		}

		registrationConfig := &entity.EventRegistrationConfig{
			ID:      uuid.New(),
			EventID: eventId,
		}

		eventContract := &entity.EventContract{
			ID:                           uuid.New(),
			EventId:                      eventId,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
		}

		mockEventDg.On("GetViewModelById", ctx, eventId).
			Return(event, registrationConfig, eventContract, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventId, userId, &email, &walletAddress).
			Return(nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("invitation not found")))
		mockS3Dg.On("GetPresignedURL", ctx, "banner-key").Return("https://banner-url", nil)
		mockS3Dg.On("GetPresignedURL", ctx, "icon-key").Return("", errors.New("S3 error for icon"))

		uc := &EventUsecase{
			EventDataGateway:              mockEventDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			S3DataGateway:                 mockS3Dg,
		}

		// Act
		result, err := uc.GetEventViewModelByEventId(ctx, eventId, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "S3 error for icon")

		mockEventDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
		mockS3Dg.AssertExpectations(t)
	})

	t.Run("should propagate SignatureExpiredAt when user signature is marked as expired", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockEventDataGateway)
		mockInvitationDg := new(MockEventRegistrationInvitationDg)
		mockAttendeeDg := new(MockEventAttendeeDg)
		mockS3Dg := new(MockS3DataGateway)

		event := &entity.Event{
			Id:                eventId,
			Title:             "Test Event",
			BannerStorageKey:  "banner-key",
			IconStorageKey:    "icon-key",
			MaxAttendees:      100,
			AttendeesCount:    50,
			OwnerCredentialId: uuid.New(),
		}

		registrationConfig := &entity.EventRegistrationConfig{
			ID:      uuid.New(),
			EventID: eventId,
		}

		eventContract := &entity.EventContract{
			ID:                           uuid.New(),
			EventId:                      eventId,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      uuid.New(),
			EventId: eventId,
		}

		joinedAt := time.Now()
		expiredAt := time.Now().Add(-1 * time.Hour)
		estimatedDeadline := time.Now().Add(-1 * time.Hour)
		deadlineBlock := int32(99999999)

		attendeeWithSignature := &event_datagateway.EventAttendeeWithSignature{
			Id:            uuid.New(),
			EventId:       eventId,
			CredentialId:  userId,
			WalletAddress: walletAddress,
			JoinedAt:      joinedAt,
			IsAccepted:    true,
			UserSignature: &entity.UserSignature{
				Id:                         uuid.New(),
				Signature:                  "0xexpired-sig",
				SignMessage:                "message",
				BroadcastedAt:              nil,
				DeadlineBlock:              &deadlineBlock,
				EstimatedDeadline:          &estimatedDeadline,
				MarkAsExpiredAt:            &expiredAt,
				AuthenticationCredentialId: userId,
			},
		}

		mockEventDg.On("GetViewModelById", ctx, eventId).
			Return(event, registrationConfig, eventContract, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventId, userId, &email, &walletAddress).
			Return(invitation, nil, nil)
		mockAttendeeDg.On("GetEventAttendeeWithSignature", ctx, eventId, userId).
			Return(attendeeWithSignature, nil)
		mockS3Dg.On("GetPresignedURL", ctx, "banner-key").Return("https://banner-url", nil)
		mockS3Dg.On("GetPresignedURL", ctx, "icon-key").Return("https://icon-url", nil)

		uc := &EventUsecase{
			EventDataGateway:              mockEventDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			EventAttendeeDg:               mockAttendeeDg,
			S3DataGateway:                 mockS3Dg,
		}

		// Act
		result, err := uc.GetEventViewModelByEventId(ctx, eventId, currentUser)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsInvited)
		assert.True(t, result.IsJoined)
		assert.False(t, result.BroadcastedAt != nil) // not broadcasted
		assert.NotNil(t, result.SignatureExpiredAt)
		assert.Equal(t, expiredAt, *result.SignatureExpiredAt)

		mockEventDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
		mockAttendeeDg.AssertExpectations(t)
		mockS3Dg.AssertExpectations(t)
	})
}
