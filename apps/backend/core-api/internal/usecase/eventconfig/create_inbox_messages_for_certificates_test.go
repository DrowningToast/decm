package eventconfig

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ---- Mocks ----

type MockEventDataGateway struct {
	mock.Mock
}

func (m *MockEventDataGateway) GetEventById(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDataGateway) CreateEvent(ctx context.Context, params eventdatagateway.CreateEventParameters) (*entity.Event, error) {
	return nil, nil
}

func (m *MockEventDataGateway) GetViewModelById(ctx context.Context, id uuid.UUID) (*entity.Event, *entity.EventRegistrationConfig, *entity.EventContract, error) {
	return nil, nil, nil, nil
}

func (m *MockEventDataGateway) ListEventsByOwnerCredentialID(ctx context.Context, ownerCredentialID uuid.UUID, limitCount int32, offsetCount int32) ([]*entity.Event, error) {
	return nil, nil
}

func (m *MockEventDataGateway) UpdateEvent(ctx context.Context, id uuid.UUID, params eventdatagateway.UpdateEventParameters) (*entity.Event, error) {
	return nil, nil
}

func (m *MockEventDataGateway) DeleteEvent(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	return nil, nil
}

func (m *MockEventDataGateway) ListEvents(ctx context.Context, limitCount *int32, offsetCount *int32) ([]*entity.Event, error) {
	return nil, nil
}

func (m *MockEventDataGateway) ListEventsByEventAttendeeCredentialID(ctx context.Context, eventAttendeeCredentialID uuid.UUID, limitCount *int32, offsetCount *int32) ([]*entity.Event, error) {
	return nil, nil
}

type MockEventCertificateDataGateway struct {
	mock.Mock
}

func (m *MockEventCertificateDataGateway) CreateEventCertificate(ctx context.Context, params eventdatagateway.CreateEventCertificateParameters) (*entity.EventCertificate, error) {
	return nil, nil
}

func (m *MockEventCertificateDataGateway) GetEventCertificateByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificate, error) {
	return nil, nil
}

func (m *MockEventCertificateDataGateway) GetEventCertificateWithSignature(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*eventdatagateway.EventCertificateWithSignature, error) {
	return nil, nil
}

func (m *MockEventCertificateDataGateway) GetEventCertificateByInboxMessageID(ctx context.Context, inboxMessageID uuid.UUID) (*entity.EventCertificate, error) {
	return nil, nil
}

func (m *MockEventCertificateDataGateway) GetEventCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetAllEventCertificateIDsByEventID(ctx context.Context, eventID uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

func (m *MockEventCertificateDataGateway) GetClaimedCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	return nil, nil
}

func (m *MockEventCertificateDataGateway) GetUnclaimedReadyCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	return nil, nil
}

func (m *MockEventCertificateDataGateway) GetClaimedCertificatesByCredentialID(ctx context.Context, credentialID uuid.UUID, email *string, walletAddress *string) ([]*entity.EventCertificate, error) {
	return nil, nil
}

func (m *MockEventCertificateDataGateway) GetUnclaimedReadyCertificatesByCredentialID(ctx context.Context, credentialID uuid.UUID, email *string, walletAddress *string) ([]*entity.EventCertificate, error) {
	return nil, nil
}

func (m *MockEventCertificateDataGateway) UpdateEventCertificate(ctx context.Context, id uuid.UUID, params eventdatagateway.UpdateEventCertificateParameters) (*entity.EventCertificate, error) {
	return nil, nil
}

func (m *MockEventCertificateDataGateway) UpdateEventCertificateInboxMessageID(ctx context.Context, id uuid.UUID, inboxMessageID uuid.UUID) (*entity.EventCertificate, error) {
	args := m.Called(ctx, id, inboxMessageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) DeleteEventCertificate(ctx context.Context, id uuid.UUID) error {
	return nil
}

type MockInboxMessageDgForCerts struct {
	mock.Mock
}

func (m *MockInboxMessageDgForCerts) CreateInboxMessage(ctx context.Context, params offchain_datagateway.CreateInboxMessageParameters) (*entity.InboxMessage, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDgForCerts) GetInboxMessageByID(ctx context.Context, id uuid.UUID) (*entity.InboxMessage, error) {
	return nil, nil
}

func (m *MockInboxMessageDgForCerts) GetInboxMessagesByCredentialID(ctx context.Context, params offchain_datagateway.GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error) {
	return nil, nil
}

func (m *MockInboxMessageDgForCerts) GetUnreadInboxMessageCountByCredentialID(ctx context.Context, params offchain_datagateway.GetInboxMessagesByCredentialIDParameters) (int, error) {
	return 0, nil
}

func (m *MockInboxMessageDgForCerts) GetInboxMessagesByReceiverEmail(ctx context.Context, receiverEmail string) ([]*entity.InboxMessage, error) {
	return nil, nil
}

func (m *MockInboxMessageDgForCerts) GetInboxMessagesByReceiverWalletAddress(ctx context.Context, walletAddress string) ([]*entity.InboxMessage, error) {
	return nil, nil
}

func (m *MockInboxMessageDgForCerts) GetInboxMessagesBySenderCredentialID(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	return nil, nil
}

func (m *MockInboxMessageDgForCerts) UpdateInboxMessageReadStatus(ctx context.Context, id uuid.UUID, isRead int) (*entity.InboxMessage, error) {
	return nil, nil
}

func (m *MockInboxMessageDgForCerts) UpdateInboxMessageReadStatusAll(ctx context.Context, params offchain_datagateway.GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error) {
	return nil, nil
}

type MockAuthCredentialDgForCerts struct {
	mock.Mock
}

func (m *MockAuthCredentialDgForCerts) GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthCredentialDgForCerts) GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	return nil, nil
}

func (m *MockAuthCredentialDgForCerts) GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, error) {
	return nil, nil
}

func (m *MockAuthCredentialDgForCerts) GetAuthenticationCredentialByGoogleConnectorRef(ctx context.Context, googleConnectorRef string) (*entity.AuthenticationCredential, error) {
	return nil, nil
}

func (m *MockAuthCredentialDgForCerts) GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddress(ctx context.Context, params offchain_datagateway.GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddressParameters) (*entity.AuthenticationCredential, error) {
	return nil, nil
}

func (m *MockAuthCredentialDgForCerts) CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error) {
	return nil, nil
}

func (m *MockAuthCredentialDgForCerts) UpdateAuthenticationCredential(ctx context.Context, id uuid.UUID, params offchain_datagateway.UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error) {
	return nil, nil
}

func (m *MockAuthCredentialDgForCerts) DeleteAuthenticationCredential(ctx context.Context, id uuid.UUID) error {
	return nil
}

// ---- Tests ----

func TestCreateInboxMessagesForCertificates_Web3User(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	eventID := uuid.New()
	senderID := uuid.New()
	credentialID := uuid.New()
	walletAddress := "0xABC123"
	certID := uuid.New()
	inboxMsgID := uuid.New()
	eventTitle := "Test Event"

	t.Run("should include wallet address in inbox message for web3-only user (no email)", func(t *testing.T) {
		mockEventDg := new(MockEventDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockInboxDg := new(MockInboxMessageDgForCerts)
		mockAuthDg := new(MockAuthCredentialDgForCerts)

		// Certificate with credential ID but no email (web3-only user)
		cert := &entity.EventCertificate{
			Id:                   certID,
			EventId:              eventID,
			ReceiverCredentialId: &credentialID,
			ReceiverEmail:        nil,
		}

		mockEventDg.On("GetEventById", ctx, eventID).Return(&entity.Event{
			Id:    eventID,
			Title: eventTitle,
		}, nil)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).Return([]*entity.EventCertificate{cert}, nil)
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, credentialID).Return(&entity.AuthenticationCredential{
			Id:            credentialID,
			WalletAddress: walletAddress,
		}, nil)
		mockInboxDg.On("CreateInboxMessage", ctx, mock.MatchedBy(func(p offchain_datagateway.CreateInboxMessageParameters) bool {
			return p.ReceiverWalletAddress != nil && *p.ReceiverWalletAddress == walletAddress &&
				p.ReceiverCredentialID != nil && *p.ReceiverCredentialID == credentialID
		})).Return(&entity.InboxMessage{Id: inboxMsgID}, nil)
		mockCertDg.On("UpdateEventCertificateInboxMessageID", ctx, certID, inboxMsgID).Return(cert, nil)

		uc := &EventConfigUsecase{
			EventDataGateway:            mockEventDg,
			EventCertificateDataGateway: mockCertDg,
			InboxMessageDg:              mockInboxDg,
			AuthenticationCredentialDg:  mockAuthDg,
			logger:                      logger,
		}

		count, err := uc.createInboxMessagesForCertificates(ctx, eventID, senderID)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		mockInboxDg.AssertExpectations(t)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should include wallet address in inbox message for user with both email and credential", func(t *testing.T) {
		mockEventDg := new(MockEventDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockInboxDg := new(MockInboxMessageDgForCerts)
		mockAuthDg := new(MockAuthCredentialDgForCerts)

		email := "user@example.com"
		cert := &entity.EventCertificate{
			Id:                   certID,
			EventId:              eventID,
			ReceiverCredentialId: &credentialID,
			ReceiverEmail:        &email,
		}

		mockEventDg.On("GetEventById", ctx, eventID).Return(&entity.Event{
			Id:    eventID,
			Title: eventTitle,
		}, nil)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).Return([]*entity.EventCertificate{cert}, nil)
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, credentialID).Return(&entity.AuthenticationCredential{
			Id:            credentialID,
			WalletAddress: walletAddress,
		}, nil)
		mockInboxDg.On("CreateInboxMessage", ctx, mock.MatchedBy(func(p offchain_datagateway.CreateInboxMessageParameters) bool {
			return p.ReceiverWalletAddress != nil && *p.ReceiverWalletAddress == walletAddress &&
				p.ReceiverEmail == email
		})).Return(&entity.InboxMessage{Id: inboxMsgID}, nil)
		mockCertDg.On("UpdateEventCertificateInboxMessageID", ctx, certID, inboxMsgID).Return(cert, nil)

		uc := &EventConfigUsecase{
			EventDataGateway:            mockEventDg,
			EventCertificateDataGateway: mockCertDg,
			InboxMessageDg:              mockInboxDg,
			AuthenticationCredentialDg:  mockAuthDg,
			logger:                      logger,
		}

		count, err := uc.createInboxMessagesForCertificates(ctx, eventID, senderID)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		mockInboxDg.AssertExpectations(t)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should create inbox message for wallet-only receiver (no email, no credential ID)", func(t *testing.T) {
		mockEventDg := new(MockEventDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockInboxDg := new(MockInboxMessageDgForCerts)
		mockAuthDg := new(MockAuthCredentialDgForCerts)

		directWallet := "0xDEADBEEF"
		cert := &entity.EventCertificate{
			Id:                    certID,
			EventId:               eventID,
			ReceiverCredentialId:  nil,
			ReceiverEmail:         nil,
			ReceiverWalletAddress: &directWallet,
		}

		mockEventDg.On("GetEventById", ctx, eventID).Return(&entity.Event{
			Id:    eventID,
			Title: eventTitle,
		}, nil)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).Return([]*entity.EventCertificate{cert}, nil)
		mockInboxDg.On("CreateInboxMessage", ctx, mock.MatchedBy(func(p offchain_datagateway.CreateInboxMessageParameters) bool {
			return p.ReceiverWalletAddress != nil && *p.ReceiverWalletAddress == directWallet &&
				p.ReceiverCredentialID == nil &&
				p.ReceiverEmail == ""
		})).Return(&entity.InboxMessage{Id: inboxMsgID}, nil)
		mockCertDg.On("UpdateEventCertificateInboxMessageID", ctx, certID, inboxMsgID).Return(cert, nil)

		uc := &EventConfigUsecase{
			EventDataGateway:            mockEventDg,
			EventCertificateDataGateway: mockCertDg,
			InboxMessageDg:              mockInboxDg,
			AuthenticationCredentialDg:  mockAuthDg,
			logger:                      logger,
		}

		count, err := uc.createInboxMessagesForCertificates(ctx, eventID, senderID)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should skip certificate with no email, no credential ID, and no wallet address", func(t *testing.T) {
		mockEventDg := new(MockEventDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockInboxDg := new(MockInboxMessageDgForCerts)
		mockAuthDg := new(MockAuthCredentialDgForCerts)

		cert := &entity.EventCertificate{
			Id:                    certID,
			EventId:               eventID,
			ReceiverCredentialId:  nil,
			ReceiverEmail:         nil,
			ReceiverWalletAddress: nil,
		}

		mockEventDg.On("GetEventById", ctx, eventID).Return(&entity.Event{
			Id:    eventID,
			Title: eventTitle,
		}, nil)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).Return([]*entity.EventCertificate{cert}, nil)

		uc := &EventConfigUsecase{
			EventDataGateway:            mockEventDg,
			EventCertificateDataGateway: mockCertDg,
			InboxMessageDg:              mockInboxDg,
			AuthenticationCredentialDg:  mockAuthDg,
			logger:                      logger,
		}

		count, err := uc.createInboxMessagesForCertificates(ctx, eventID, senderID)

		require.NoError(t, err)
		assert.Equal(t, 0, count)
		mockInboxDg.AssertNotCalled(t, "CreateInboxMessage")
	})

	t.Run("should proceed without wallet address if credential lookup fails", func(t *testing.T) {
		mockEventDg := new(MockEventDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockInboxDg := new(MockInboxMessageDgForCerts)
		mockAuthDg := new(MockAuthCredentialDgForCerts)

		email := "user@example.com"
		cert := &entity.EventCertificate{
			Id:                   certID,
			EventId:              eventID,
			ReceiverCredentialId: &credentialID,
			ReceiverEmail:        &email,
		}

		mockEventDg.On("GetEventById", ctx, eventID).Return(&entity.Event{
			Id:    eventID,
			Title: eventTitle,
		}, nil)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).Return([]*entity.EventCertificate{cert}, nil)
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, credentialID).Return(nil, errors.New("db error"))
		// Should still create inbox message, but without wallet address
		mockInboxDg.On("CreateInboxMessage", ctx, mock.MatchedBy(func(p offchain_datagateway.CreateInboxMessageParameters) bool {
			return p.ReceiverWalletAddress == nil && p.ReceiverEmail == email
		})).Return(&entity.InboxMessage{Id: inboxMsgID}, nil)
		mockCertDg.On("UpdateEventCertificateInboxMessageID", ctx, certID, inboxMsgID).Return(cert, nil)

		uc := &EventConfigUsecase{
			EventDataGateway:            mockEventDg,
			EventCertificateDataGateway: mockCertDg,
			InboxMessageDg:              mockInboxDg,
			AuthenticationCredentialDg:  mockAuthDg,
			logger:                      logger,
		}

		count, err := uc.createInboxMessagesForCertificates(ctx, eventID, senderID)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		mockInboxDg.AssertExpectations(t)
	})
}
