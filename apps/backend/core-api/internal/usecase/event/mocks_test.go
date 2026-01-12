package event

import (
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	eventcontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/event_contract"
	storage_datagateway "apps/backend/core-api/internal/datagateway/storage"
	"apps/backend/core-api/config"
	"apps/backend/core-api/config/blockchain"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"
	"math/big"
	"mime/multipart"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// --- Reusable Mocks for Event Package Tests ---

// MockEventDataGateway implements event_datagateway.EventDataGateway
type MockEventDataGateway struct {
	mock.Mock
}

var _ event_datagateway.EventDataGateway = (*MockEventDataGateway)(nil)

func (m *MockEventDataGateway) CreateEvent(ctx context.Context, params event_datagateway.CreateEventParameters) (*entity.Event, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDataGateway) GetEventById(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDataGateway) GetViewModelById(ctx context.Context, id uuid.UUID) (*entity.Event, *entity.EventRegistrationConfig, *entity.EventContract, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, nil, nil, args.Error(3)
	}
	return args.Get(0).(*entity.Event), args.Get(1).(*entity.EventRegistrationConfig), args.Get(2).(*entity.EventContract), args.Error(3)
}

func (m *MockEventDataGateway) ListEventsByOwnerCredentialID(ctx context.Context, ownerCredentialID uuid.UUID, limitCount int32, offsetCount int32) ([]*entity.Event, error) {
	args := m.Called(ctx, ownerCredentialID, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func (m *MockEventDataGateway) UpdateEvent(ctx context.Context, id uuid.UUID, params event_datagateway.UpdateEventParameters) (*entity.Event, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDataGateway) DeleteEvent(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDataGateway) ListEvents(ctx context.Context, limitCount *int32, offsetCount *int32) ([]*entity.Event, error) {
	args := m.Called(ctx, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func (m *MockEventDataGateway) ListEventsByEventAttendeeCredentialID(ctx context.Context, eventAttendeeCredentialID uuid.UUID, limitCount *int32, offsetCount *int32) ([]*entity.Event, error) {
	args := m.Called(ctx, eventAttendeeCredentialID, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

// MockAuthenticationCredentialDg implements offchain_datagateway.AuthenticationCredentialDataGateway
type MockAuthenticationCredentialDg struct {
	mock.Mock
}

var _ offchain_datagateway.AuthenticationCredentialDataGateway = (*MockAuthenticationCredentialDg)(nil)

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx context.Context, credentialId uuid.UUID) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, credentialId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
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

// MockS3DataGateway implements storage_datagateway.StorageDataGateway
type MockS3DataGateway struct {
	mock.Mock
}

var _ storage_datagateway.S3DataGateway = (*MockS3DataGateway)(nil)

func (m *MockS3DataGateway) GetS3UploadRequestObject(entityType storage_datagateway.StorageKeyType, entityID uuid.UUID, fileHeader *multipart.FileHeader) (*storage_datagateway.S3UploadRequestObject, error) {
	args := m.Called(entityType, entityID, fileHeader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage_datagateway.S3UploadRequestObject), args.Error(1)
}

func (m *MockS3DataGateway) PutFile(ctx context.Context, requestObject *storage_datagateway.S3UploadRequestObject) (string, error) {
	args := m.Called(ctx, requestObject)
	return args.String(0), args.Error(1)
}

func (m *MockS3DataGateway) DeleteFile(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockS3DataGateway) GetStorageKey(entityType storage_datagateway.StorageKeyType, entityID uuid.UUID, fileName string, fileExtension string) string {
	args := m.Called(entityType, entityID, fileName, fileExtension)
	return args.String(0)
}

func (m *MockS3DataGateway) GetPresignedURL(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockS3DataGateway) GetFile(ctx context.Context, key string) (interface{}, error) {
	args := m.Called(ctx, key)
	return args.Get(0), args.Error(1)
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

// MockBlockchainClientDataGateway implements blockchainclient_datagateway.BlockchainClientDataGateway
type MockBlockchainClientDataGateway struct {
	mock.Mock
}

var _ blockchainclient_datagateway.BlockchainClientDataGateway = (*MockBlockchainClientDataGateway)(nil)

func (m *MockBlockchainClientDataGateway) GetCurrentBlockNumber(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockBlockchainClientDataGateway) GetCalculatedDeadlineBlock(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockBlockchainClientDataGateway) EstimateDeadlineTime(ctx context.Context, deadlineBlock uint64) (*time.Time, error) {
	args := m.Called(ctx, deadlineBlock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*time.Time), args.Error(1)
}

func (m *MockBlockchainClientDataGateway) GetTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bind.TransactOpts), args.Error(1)
}

func (m *MockBlockchainClientDataGateway) GetGasPrice(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*blockchainclient_datagateway.GasPriceInfo), args.Error(1)
}

func (m *MockBlockchainClientDataGateway) WeiToGwei(wei *big.Int) float64 {
	args := m.Called(wei)
	return args.Get(0).(float64)
}

// MockEventRegistrationInvitationDg implements event_datagateway.EventRegistrationInvitationDataGateway
type MockEventRegistrationInvitationDg struct {
	mock.Mock
}

var _ event_datagateway.EventRegistrationInvitationDataGateway = (*MockEventRegistrationInvitationDg)(nil)

func (m *MockEventRegistrationInvitationDg) CreateEventRegistrationInvitation(ctx context.Context, params event_datagateway.CreateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDg) GetEventRegistrationInvitationByID(ctx context.Context, id uuid.UUID) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDg) GetEventRegistrationInvitationsByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDg) GetEventRegistrationInvitationByInboxMessageID(ctx context.Context, inboxMessageID uuid.UUID) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, inboxMessageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDg) GetEventRegistrationInvitationByEventIDAndCredential(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID, email *string, walletAddress *string) (*entity.EventRegistrationInvitation, *entity.InboxMessage, error) {
	args := m.Called(ctx, eventId, credentialId, email, walletAddress)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), nil, args.Error(2)
}

func (m *MockEventRegistrationInvitationDg) UpdateEventRegistrationInvitation(ctx context.Context, id uuid.UUID, params event_datagateway.UpdateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDg) UpdateEventRegistrationInvitationAcceptedStatus(ctx context.Context, id uuid.UUID, acceptedAt *time.Time) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, id, acceptedAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDg) DeleteEventRegistrationInvitation(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEventAttendeeDg implements event_datagateway.EventAttendeeDataGateway
type MockEventAttendeeDg struct {
	mock.Mock
}

var _ event_datagateway.EventAttendeeDataGateway = (*MockEventAttendeeDg)(nil)

func (m *MockEventAttendeeDg) GetEventAttendeeWithSignature(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (*event_datagateway.EventAttendeeWithSignature, error) {
	args := m.Called(ctx, eventId, credentialId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*event_datagateway.EventAttendeeWithSignature), args.Error(1)
}

func (m *MockEventAttendeeDg) GetEventAttendeeByEventIdAndCredentialId(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (*entity.EventAttendee, error) {
	args := m.Called(ctx, eventId, credentialId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventAttendee), args.Error(1)
}

func (m *MockEventAttendeeDg) AddParticipant(ctx context.Context, params event_datagateway.AddParticipantParameters) (*entity.EventAttendee, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventAttendee), args.Error(1)
}

func (m *MockEventAttendeeDg) ListEventAttendeesByEventID(ctx context.Context, eventID uuid.UUID) ([]entity.EventAttendee, error) {
	args := m.Called(ctx, eventID)
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
	return args.Bool(0), args.Error(1)
}

// MockEventCertificateDataGateway implements event_datagateway.EventCertificateDataGateway
type MockEventCertificateDataGateway struct {
	mock.Mock
}

var _ event_datagateway.EventCertificateDataGateway = (*MockEventCertificateDataGateway)(nil)

func (m *MockEventCertificateDataGateway) CreateEventCertificate(ctx context.Context, params event_datagateway.CreateEventCertificateParameters) (*entity.EventCertificate, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetEventCertificateWithSignature(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*event_datagateway.EventCertificateWithSignature, error) {
	args := m.Called(ctx, eventID, credentialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*event_datagateway.EventCertificateWithSignature), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetEventCertificateByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetEventCertificateByInboxMessageID(ctx context.Context, inboxMessageID uuid.UUID) (*entity.EventCertificate, error) {
	args := m.Called(ctx, inboxMessageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetEventCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetAllEventCertificateIDsByEventID(ctx context.Context, eventID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

func (m *MockEventCertificateDataGateway) UpdateEventCertificate(ctx context.Context, id uuid.UUID, params event_datagateway.UpdateEventCertificateParameters) (*entity.EventCertificate, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) UpdateEventCertificateInboxMessageID(ctx context.Context, id uuid.UUID, inboxMessageID uuid.UUID) (*entity.EventCertificate, error) {
	args := m.Called(ctx, id, inboxMessageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) DeleteEventCertificate(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventCertificateDataGateway) GetClaimedCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetUnclaimedReadyCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetClaimedCertificatesByCredentialID(ctx context.Context, credentialID uuid.UUID, email *string) ([]*entity.EventCertificate, error) {
	args := m.Called(ctx, credentialID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetUnclaimedReadyCertificatesByCredentialID(ctx context.Context, credentialID uuid.UUID, email *string) ([]*entity.EventCertificate, error) {
	args := m.Called(ctx, credentialID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificate), args.Error(1)
}

// MockEventCertificateSignatureDataGateway implements event_datagateway.EventCertificateSignatureDataGateway
type MockEventCertificateSignatureDataGateway struct {
	mock.Mock
}

var _ event_datagateway.EventCertificateSignatureDataGateway = (*MockEventCertificateSignatureDataGateway)(nil)

func (m *MockEventCertificateSignatureDataGateway) CreateEventCertificateSignature(ctx context.Context, params event_datagateway.CreateEventCertificateSignatureParameters) (*entity.EventCertificateSignature, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateSignature), args.Error(1)
}

func (m *MockEventCertificateSignatureDataGateway) GetEventCertificateSignatureByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificateSignature, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateSignature), args.Error(1)
}

func (m *MockEventCertificateSignatureDataGateway) GetEventCertificateSignaturesByEventCertificateConfigID(ctx context.Context, eventCertificateConfigID uuid.UUID) ([]*entity.EventCertificateSignature, error) {
	args := m.Called(ctx, eventCertificateConfigID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificateSignature), args.Error(1)
}

func (m *MockEventCertificateSignatureDataGateway) UpdateEventCertificateSignature(ctx context.Context, id uuid.UUID, params event_datagateway.UpdateEventCertificateSignatureParameters) (*entity.EventCertificateSignature, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateSignature), args.Error(1)
}

func (m *MockEventCertificateSignatureDataGateway) UpdateEventCertificateIssuerSignature(ctx context.Context, id uuid.UUID, issuerSignature *string) (*entity.EventCertificateSignature, error) {
	args := m.Called(ctx, id, issuerSignature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateSignature), args.Error(1)
}

func (m *MockEventCertificateSignatureDataGateway) DeleteEventCertificateSignature(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEventIssuerDataGateway implements event_datagateway.EventIssuerDataGateway
type MockEventIssuerDataGateway struct {
	mock.Mock
}

var _ event_datagateway.EventIssuerDataGateway = (*MockEventIssuerDataGateway)(nil)

func (m *MockEventIssuerDataGateway) CreateEventIssuer(ctx context.Context, params generated.CreateEventIssuerParams) (*generated.EventIssuer, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*generated.EventIssuer), args.Error(1)
}

func (m *MockEventIssuerDataGateway) GetEventIssuerByID(ctx context.Context, id uuid.UUID) (generated.EventIssuer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(generated.EventIssuer), args.Error(1)
}

func (m *MockEventIssuerDataGateway) GetEventIssuersByEventID(ctx context.Context, eventID uuid.UUID) ([]generated.EventIssuer, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]generated.EventIssuer), args.Error(1)
}

func (m *MockEventIssuerDataGateway) UpdateEventIssuer(ctx context.Context, params generated.UpdateEventIssuerParams) (*generated.EventIssuer, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*generated.EventIssuer), args.Error(1)
}

func (m *MockEventIssuerDataGateway) DeleteEventIssuer(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *MockEventIssuerDataGateway) GetEventIssuerByEventIDAndIssuerCredentialID(ctx context.Context, eventID uuid.UUID, issuerCredentialID uuid.UUID) (generated.EventIssuer, error) {
	args := m.Called(ctx, eventID, issuerCredentialID)
	return args.Get(0).(generated.EventIssuer), args.Error(1)
}

func (m *MockEventIssuerDataGateway) UpdateEventIssuerSigningStatus(ctx context.Context, eventID uuid.UUID, issuerCredentialID uuid.UUID, isSigned int32) error {
	args := m.Called(ctx, eventID, issuerCredentialID, isSigned)
	return args.Error(0)
}

func (m *MockEventIssuerDataGateway) ResetAllEventIssuersSigningStatus(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *MockEventIssuerDataGateway) AllIssuersHaveSigned(ctx context.Context, eventID uuid.UUID) (bool, error) {
	args := m.Called(ctx, eventID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEventIssuerDataGateway) GetSignedIssuersCount(ctx context.Context, eventID uuid.UUID) (int64, error) {
	args := m.Called(ctx, eventID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockEventIssuerDataGateway) GetTotalIssuersCount(ctx context.Context, eventID uuid.UUID) (int64, error) {
	args := m.Called(ctx, eventID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockEventIssuerDataGateway) HasSignedIssuers(ctx context.Context, eventID uuid.UUID) (bool, error) {
	args := m.Called(ctx, eventID)
	return args.Bool(0), args.Error(1)
}

// MockEventCertificateConfigDataGateway implements event_datagateway.EventCertificateConfigDataGateway
type MockEventCertificateConfigDataGateway struct {
	mock.Mock
}

var _ event_datagateway.EventCertificateConfigDataGateway = (*MockEventCertificateConfigDataGateway)(nil)

func (m *MockEventCertificateConfigDataGateway) CreateEventCertificateConfig(ctx context.Context, params generated.CreateEventCertificateConfigParams) (*entity.EventCertificateConfig, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateConfig), args.Error(1)
}

func (m *MockEventCertificateConfigDataGateway) GetEventCertificateConfigByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificateConfig, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateConfig), args.Error(1)
}

func (m *MockEventCertificateConfigDataGateway) GetEventCertificateConfigByEventID(ctx context.Context, eventId uuid.UUID) (*entity.EventCertificateConfig, error) {
	args := m.Called(ctx, eventId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateConfig), args.Error(1)
}

func (m *MockEventCertificateConfigDataGateway) UpdateEventCertificateConfig(ctx context.Context, params generated.UpdateEventCertificateConfigParams) (*entity.EventCertificateConfig, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateConfig), args.Error(1)
}

func (m *MockEventCertificateConfigDataGateway) UpdateEventCertificateTextConfig(ctx context.Context, params generated.UpdateEventCertificateTextConfigParams) (*entity.EventCertificateConfig, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateConfig), args.Error(1)
}

func (m *MockEventCertificateConfigDataGateway) ToggleEventCertificateConfigPublished(ctx context.Context, params generated.ToggleEventCertificateConfigPublishedParams) (*entity.EventCertificateConfig, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateConfig), args.Error(1)
}

func (m *MockEventCertificateConfigDataGateway) DeleteEventCertificateConfig(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

// MockEventContractDataGateway implements event_datagateway.EventContractDataGateway
type MockEventContractDataGateway struct {
	mock.Mock
}

var _ event_datagateway.EventContractDataGateway = (*MockEventContractDataGateway)(nil)

func (m *MockEventContractDataGateway) CreateEventContract(ctx context.Context, params generated.CreateEventContractParams) (*entity.EventContract, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventContract), args.Error(1)
}

func (m *MockEventContractDataGateway) GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*entity.EventContract, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventContract), args.Error(1)
}

func (m *MockEventContractDataGateway) UpdateEventContract(ctx context.Context, params generated.UpdateEventContractParams) (*entity.EventContract, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventContract), args.Error(1)
}

func (m *MockEventContractDataGateway) DeleteEventContract(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

// MockUserSignatureDataGateway implements offchain_datagateway.UserSignatureDataGateway
type MockUserSignatureDataGateway struct {
	mock.Mock
}

var _ offchain_datagateway.UserSignatureDataGateway = (*MockUserSignatureDataGateway)(nil)

func (m *MockUserSignatureDataGateway) CreateUserSignature(ctx context.Context, params offchain_datagateway.CreateUserSignatureParameters) (*entity.UserSignature, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetUserSignatureByID(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) UpdateUserSignatureBroadcastedAt(ctx context.Context, id uuid.UUID, broadcastedAt *time.Time) (*entity.UserSignature, error) {
	args := m.Called(ctx, id, broadcastedAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) UpdateUserSignatureMarkAsExpiredAt(ctx context.Context, id uuid.UUID, markAsExpiredAt *time.Time) (*entity.UserSignature, error) {
	args := m.Called(ctx, id, markAsExpiredAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) DeleteUserSignature(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserSignatureDataGateway) GetUserSignaturesByCredentialID(ctx context.Context, authenticationCredentialId uuid.UUID) ([]entity.UserSignature, error) {
	args := m.Called(ctx, authenticationCredentialId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetPendingUserSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetStaleUserSignatures(ctx context.Context, createdBefore time.Time) ([]entity.UserSignature, error) {
	args := m.Called(ctx, createdBefore)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetBroadcastedUserSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetUserSignaturesByDeadlineBlockRange(ctx context.Context, minBlock int32, maxBlock int32) ([]entity.UserSignature, error) {
	args := m.Called(ctx, minBlock, maxBlock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetUserSignaturesExpiringBefore(ctx context.Context, deadline time.Time) ([]entity.UserSignature, error) {
	args := m.Called(ctx, deadline)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetEventJoinSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetPendingEventJoinSignatures(ctx context.Context) ([]entity.PendingEventJoin, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.PendingEventJoin), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetPendingCertificateClaimSignatures(ctx context.Context) ([]entity.PendingCertificateClaim, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.PendingCertificateClaim), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetCertificateClaimSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetOrphanedUserSignatures(ctx context.Context, createdBefore time.Time) ([]entity.UserSignature, error) {
	args := m.Called(ctx, createdBefore)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetUserSignatureWithUsageDetails(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) GetRecentUserSignatureActivity(ctx context.Context, since time.Time) ([]entity.UserSignature, error) {
	args := m.Called(ctx, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserSignature), args.Error(1)
}

func (m *MockUserSignatureDataGateway) CountUserSignaturesByCredentialID(ctx context.Context, authenticationCredentialId uuid.UUID) (int64, error) {
	args := m.Called(ctx, authenticationCredentialId)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockUserSignatureDataGateway) CountPendingUserSignatures(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockUserSignatureDataGateway) CountBroadcastedUserSignatures(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockUserSignatureDataGateway) UpdateUserSignatureAbortedAt(ctx context.Context, id uuid.UUID, abortedAt time.Time, reason entity.UserSignatureAbortReason) (*entity.UserSignature, error) {
	args := m.Called(ctx, id, abortedAt, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

// --- Shared Test Helpers ---

func strPtr(s string) *string {
	return &s
}

func stringPtr(s string) *string {
	return &s
}

func createMockConfig() *config.Config {
	return &config.Config{
		Blockchain: blockchain.BlockchainConfig{
			ChainID:                  1,
			DecmAccessManagerAddress: "0x1234567890123456789012345678901234567890",
		},
	}
}
