package event_registration

import (
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	eventcontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/event_contract"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// --- Reusable Mocks for Event Registration Tests ---

// MockEventRegistrationConfigDg implements eventdatagateway.EventRegistrationConfigDataGateway
type MockEventRegistrationConfigDg struct {
	mock.Mock
}

var _ eventdatagateway.EventRegistrationConfigDataGateway = (*MockEventRegistrationConfigDg)(nil)

func (m *MockEventRegistrationConfigDg) CreateEventRegistrationConfig(ctx context.Context, params generated.CreateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationConfig), args.Error(1)
}

func (m *MockEventRegistrationConfigDg) GetEventRegistrationConfigByEventId(ctx context.Context, eventId uuid.UUID) (*entity.EventRegistrationConfig, error) {
	args := m.Called(ctx, eventId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationConfig), args.Error(1)
}

func (m *MockEventRegistrationConfigDg) GetEventRegistrationConfigPasswordByEventId(ctx context.Context, eventId uuid.UUID) (*entity.EventRegistrationConfig, error) {
	args := m.Called(ctx, eventId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationConfig), args.Error(1)
}

func (m *MockEventRegistrationConfigDg) UpdateEventRegistrationConfig(ctx context.Context, params generated.UpdateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationConfig), args.Error(1)
}

func (m *MockEventRegistrationConfigDg) DeleteEventRegistrationConfig(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

// MockEventContractDg implements eventdatagateway.EventContractDataGateway
type MockEventContractDg struct {
	mock.Mock
}

var _ eventdatagateway.EventContractDataGateway = (*MockEventContractDg)(nil)

func (m *MockEventContractDg) CreateEventContract(ctx context.Context, params generated.CreateEventContractParams) (*entity.EventContract, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventContract), args.Error(1)
}

func (m *MockEventContractDg) GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*entity.EventContract, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventContract), args.Error(1)
}

func (m *MockEventContractDg) UpdateEventContract(ctx context.Context, params generated.UpdateEventContractParams) (*entity.EventContract, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventContract), args.Error(1)
}

func (m *MockEventContractDg) DeleteEventContract(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

// MockEventAttendeeDg implements eventdatagateway.EventAttendeeDataGateway
type MockEventAttendeeDg struct {
	mock.Mock
}

var _ eventdatagateway.EventAttendeeDataGateway = (*MockEventAttendeeDg)(nil)

func (m *MockEventAttendeeDg) GetEventAttendeeByEventIdAndCredentialId(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (*entity.EventAttendee, error) {
	args := m.Called(ctx, eventId, credentialId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventAttendee), args.Error(1)
}

func (m *MockEventAttendeeDg) GetEventAttendeeWithSignature(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (*eventdatagateway.EventAttendeeWithSignature, error) {
	args := m.Called(ctx, eventId, credentialId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eventdatagateway.EventAttendeeWithSignature), args.Error(1)
}

func (m *MockEventAttendeeDg) AddParticipant(ctx context.Context, params eventdatagateway.AddParticipantParameters) (*entity.EventAttendee, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventAttendee), args.Error(1)
}

func (m *MockEventAttendeeDg) ListEventAttendeesByEventID(ctx context.Context, eventId uuid.UUID) ([]entity.EventAttendee, error) {
	args := m.Called(ctx, eventId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.EventAttendee), args.Error(1)
}

func (m *MockEventAttendeeDg) DeleteEventAttendeeById(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventAttendeeDg) DeleteEventAttendeeByEventIDAndCredentialID(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) error {
	args := m.Called(ctx, eventId, credentialId)
	return args.Error(0)
}

func (m *MockEventAttendeeDg) HasPendingEventJoinByEventAndCredential(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (bool, error) {
	args := m.Called(ctx, eventId, credentialId)
	return args.Get(0).(bool), args.Error(1)
}

// MockUserSignatureDg implements offchain_datagateway.UserSignatureDataGateway
type MockUserSignatureDg struct {
	mock.Mock
}

var _ offchain_datagateway.UserSignatureDataGateway = (*MockUserSignatureDg)(nil)

func (m *MockUserSignatureDg) CreateUserSignature(ctx context.Context, params offchain_datagateway.CreateUserSignatureParameters) (*entity.UserSignature, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetUserSignatureByID(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetUserSignaturesByCredentialID(ctx context.Context, id uuid.UUID) ([]entity.UserSignature, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) UpdateUserSignatureBroadcastedAt(ctx context.Context, id uuid.UUID, broadcastedAt *time.Time) (*entity.UserSignature, error) {
	args := m.Called(ctx, id, broadcastedAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) UpdateUserSignatureMarkAsExpiredAt(ctx context.Context, id uuid.UUID, markAsExpiredAt *time.Time) (*entity.UserSignature, error) {
	args := m.Called(ctx, id, markAsExpiredAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetPendingUserSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetStaleUserSignatures(ctx context.Context, createdBefore time.Time) ([]entity.UserSignature, error) {
	args := m.Called(ctx, createdBefore)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetBroadcastedUserSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetUserSignaturesByDeadlineBlockRange(ctx context.Context, minBlock int32, maxBlock int32) ([]entity.UserSignature, error) {
	args := m.Called(ctx, minBlock, maxBlock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetUserSignaturesExpiringBefore(ctx context.Context, deadline time.Time) ([]entity.UserSignature, error) {
	args := m.Called(ctx, deadline)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetEventJoinSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetPendingEventJoinSignatures(ctx context.Context) ([]entity.PendingEventJoin, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.PendingEventJoin), args.Error(1)
}

func (m *MockUserSignatureDg) GetPendingCertificateClaimSignatures(ctx context.Context) ([]entity.PendingCertificateClaim, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.PendingCertificateClaim), args.Error(1)
}

func (m *MockUserSignatureDg) GetCertificateClaimSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetOrphanedUserSignatures(ctx context.Context, createdBefore time.Time) ([]entity.UserSignature, error) {
	args := m.Called(ctx, createdBefore)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetUserSignatureWithUsageDetails(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) GetRecentUserSignatureActivity(ctx context.Context, since time.Time) ([]entity.UserSignature, error) {
	args := m.Called(ctx, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDg) DeleteUserSignature(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserSignatureDg) CountUserSignaturesByCredentialID(ctx context.Context, id uuid.UUID) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserSignatureDg) CountPendingUserSignatures(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserSignatureDg) CountBroadcastedUserSignatures(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserSignatureDg) UpdateUserSignatureAbortedAt(ctx context.Context, id uuid.UUID, abortedAt time.Time, reason entity.UserSignatureAbortReason) (*entity.UserSignature, error) {
	args := m.Called(ctx, id, abortedAt, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

// MockEventRegistrationInvitationDataGateway implements eventdatagateway.EventRegistrationInvitationDataGateway
type MockEventRegistrationInvitationDataGateway struct {
	mock.Mock
}

var _ eventdatagateway.EventRegistrationInvitationDataGateway = (*MockEventRegistrationInvitationDataGateway)(nil)

func (m *MockEventRegistrationInvitationDataGateway) CreateEventRegistrationInvitation(ctx context.Context, params eventdatagateway.CreateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) GetEventRegistrationInvitationByID(ctx context.Context, invitationID uuid.UUID) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, invitationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) GetEventRegistrationInvitationsByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) GetEventRegistrationInvitationByInboxMessageID(ctx context.Context, inboxMessageID uuid.UUID) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, inboxMessageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) GetEventRegistrationInvitationByEventIDAndCredential(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID, email *string, walletAddress *string) (*entity.EventRegistrationInvitation, *entity.InboxMessage, error) {
	args := m.Called(ctx, eventID, credentialID, email, walletAddress)
	var invitation *entity.EventRegistrationInvitation
	var inbox *entity.InboxMessage

	if args.Get(0) != nil {
		invitation = args.Get(0).(*entity.EventRegistrationInvitation)
	}
	if args.Get(1) != nil {
		inbox = args.Get(1).(*entity.InboxMessage)
	}

	return invitation, inbox, args.Error(2)
}

func (m *MockEventRegistrationInvitationDataGateway) UpdateEventRegistrationInvitation(ctx context.Context, id uuid.UUID, params eventdatagateway.UpdateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) UpdateEventRegistrationInvitationAcceptedStatus(ctx context.Context, invitationID uuid.UUID, acceptedAt *time.Time) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, invitationID, acceptedAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) DeleteEventRegistrationInvitation(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockBlockchainClientDg implements blockchainclient_datagateway.BlockchainClientDataGateway
type MockBlockchainClientDg struct {
	mock.Mock
}

var _ blockchainclient_datagateway.BlockchainClientDataGateway = (*MockBlockchainClientDg)(nil)

func (m *MockBlockchainClientDg) GetCurrentBlockNumber(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockBlockchainClientDg) GetCalculatedDeadlineBlock(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockBlockchainClientDg) EstimateDeadlineTime(ctx context.Context, deadlineBlock uint64) (*time.Time, error) {
	args := m.Called(ctx, deadlineBlock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*time.Time), args.Error(1)
}

func (m *MockBlockchainClientDg) GetTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bind.TransactOpts), args.Error(1)
}

func (m *MockBlockchainClientDg) GetGasPrice(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*blockchainclient_datagateway.GasPriceInfo), args.Error(1)
}

func (m *MockBlockchainClientDg) WeiToGwei(wei *big.Int) float64 {
	args := m.Called(wei)
	return args.Get(0).(float64)
}

// MockEventContractFactoryDg implements eventcontract_datagateway.EventContractFactoryDataGateway
type MockEventContractFactoryDg struct {
	mock.Mock
}

var _ eventcontract_datagateway.EventContractFactoryDataGateway = (*MockEventContractFactoryDg)(nil)

func (m *MockEventContractFactoryDg) GetContract(address common.Address) (eventcontract_datagateway.EventContractDataGateway, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(eventcontract_datagateway.EventContractDataGateway), args.Error(1)
}

func (m *MockEventContractFactoryDg) CreateContract(ctx context.Context, params eventcontract_datagateway.CreateContractParams) (*eventcontract_datagateway.CreateContractResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eventcontract_datagateway.CreateContractResponse), args.Error(1)
}

// MockOnchainEventContractDg implements eventcontract_datagateway.EventContractDataGateway (onchain operations)
type MockOnchainEventContractDg struct {
	mock.Mock
}

var _ eventcontract_datagateway.EventContractDataGateway = (*MockOnchainEventContractDg)(nil)

func (m *MockOnchainEventContractDg) CheckIfParticipantHasJoinedEvent(ctx context.Context, walletAddress common.Address) (bool, error) {
	args := m.Called(ctx, walletAddress)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockOnchainEventContractDg) GetMaxAttendeeCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockOnchainEventContractDg) GetCurrentAttendeeCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// MockInboxMessageDg implements offchain_datagateway.InboxMessageDataGateway
type MockInboxMessageDg struct {
	mock.Mock
}

var _ offchain_datagateway.InboxMessageDataGateway = (*MockInboxMessageDg)(nil)

func (m *MockInboxMessageDg) CreateInboxMessage(ctx context.Context, params offchain_datagateway.CreateInboxMessageParameters) (*entity.InboxMessage, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetInboxMessageByID(ctx context.Context, id uuid.UUID) (*entity.InboxMessage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetInboxMessagesByCredentialID(ctx context.Context, params offchain_datagateway.GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetUnreadInboxMessageCountByCredentialID(ctx context.Context, params offchain_datagateway.GetInboxMessagesByCredentialIDParameters) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockInboxMessageDg) GetInboxMessagesByReceiverEmail(ctx context.Context, receiverEmail string) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, receiverEmail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetInboxMessagesByReceiverWalletAddress(ctx context.Context, walletAddress string) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, walletAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetInboxMessagesBySenderCredentialID(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, credentialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) UpdateInboxMessageReadStatus(ctx context.Context, id uuid.UUID, isRead int) (*entity.InboxMessage, error) {
	args := m.Called(ctx, id, isRead)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) UpdateInboxMessageReadStatusAll(ctx context.Context, params offchain_datagateway.GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

// MockEventDg implements eventdatagateway.EventDataGateway
type MockEventDg struct {
	mock.Mock
}

var _ eventdatagateway.EventDataGateway = (*MockEventDg)(nil)

func (m *MockEventDg) CreateEvent(ctx context.Context, params eventdatagateway.CreateEventParameters) (*entity.Event, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDg) GetEventById(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDg) GetViewModelById(ctx context.Context, id uuid.UUID) (*entity.Event, *entity.EventRegistrationConfig, *entity.EventContract, error) {
	args := m.Called(ctx, id)
	var event *entity.Event
	var config *entity.EventRegistrationConfig
	var contract *entity.EventContract

	if args.Get(0) != nil {
		event = args.Get(0).(*entity.Event)
	}
	if args.Get(1) != nil {
		config = args.Get(1).(*entity.EventRegistrationConfig)
	}
	if args.Get(2) != nil {
		contract = args.Get(2).(*entity.EventContract)
	}

	return event, config, contract, args.Error(3)
}

func (m *MockEventDg) ListEventsByOwnerCredentialID(ctx context.Context, ownerCredentialID uuid.UUID, limitCount int32, offsetCount int32) ([]*entity.Event, error) {
	args := m.Called(ctx, ownerCredentialID, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func (m *MockEventDg) UpdateEvent(ctx context.Context, id uuid.UUID, params eventdatagateway.UpdateEventParameters) (*entity.Event, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDg) DeleteEvent(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDg) ListEvents(ctx context.Context, limitCount *int32, offsetCount *int32) ([]*entity.Event, error) {
	args := m.Called(ctx, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func (m *MockEventDg) ListEventsByEventAttendeeCredentialID(ctx context.Context, eventAttendeeCredentialID uuid.UUID, limitCount *int32, offsetCount *int32) ([]*entity.Event, error) {
	args := m.Called(ctx, eventAttendeeCredentialID, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func (m *MockOnchainEventContractDg) AddParticipant(ctx context.Context, participantAddress common.Address, signMessage string, signature []byte) error {
	args := m.Called(ctx, participantAddress, signMessage, signature)
	return args.Error(0)
}

func (m *MockOnchainEventContractDg) GetCurrentSeatsCount(ctx context.Context) (*int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int64), args.Error(1)
}

func (m *MockOnchainEventContractDg) GetMaxSeatsCount(ctx context.Context) (*int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int64), args.Error(1)
}

func (m *MockOnchainEventContractDg) GetParticipants(ctx context.Context) ([]common.Address, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]common.Address), args.Error(1)
}

func (m *MockOnchainEventContractDg) UpdateEvent(ctx context.Context, eventName string, eventDescription string, seatsCount int64, signMessage string, signature []byte) error {
	args := m.Called(ctx, eventName, eventDescription, seatsCount, signMessage, signature)
	return args.Error(0)
}

// MockAuthenticationCredentialDg implements offchain_datagateway.AuthenticationCredentialDataGateway
type MockAuthenticationCredentialDg struct {
	mock.Mock
}

var _ offchain_datagateway.AuthenticationCredentialDataGateway = (*MockAuthenticationCredentialDg)(nil)

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, walletAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByGoogleConnectorRef(ctx context.Context, googleConnectorRef string) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, googleConnectorRef)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddress(ctx context.Context, params offchain_datagateway.GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddressParameters) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) UpdateAuthenticationCredential(ctx context.Context, id uuid.UUID, params offchain_datagateway.UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) DeleteAuthenticationCredential(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
