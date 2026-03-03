package certificate_share

import (
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	certificatecontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/certificate_contract"
	"apps/backend/core-api/internal/entity"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockCertificateShareDataGateway implements event_datagateway.CertificateShareDataGateway
type MockCertificateShareDataGateway struct {
	mock.Mock
}

var _ event_datagateway.CertificateShareDataGateway = (*MockCertificateShareDataGateway)(nil)

func (m *MockCertificateShareDataGateway) CreateCertificateShare(ctx context.Context, params event_datagateway.CreateCertificateShareParameters) (*entity.CertificateShare, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CertificateShare), args.Error(1)
}

func (m *MockCertificateShareDataGateway) GetCertificateShareByHandle(ctx context.Context, handle string) (*entity.CertificateShare, error) {
	args := m.Called(ctx, handle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CertificateShare), args.Error(1)
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

func (m *MockEventCertificateDataGateway) GetEventCertificateByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificate, error) {
	args := m.Called(ctx, id)
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

// MockCertificateContractFactoryDg implements certificatecontract_datagateway.CertificateContractFactoryDataGateway
type MockCertificateContractFactoryDg struct {
	mock.Mock
}

var _ certificatecontract_datagateway.CertificateContractFactoryDataGateway = (*MockCertificateContractFactoryDg)(nil)

func (m *MockCertificateContractFactoryDg) GetContract(address common.Address) (certificatecontract_datagateway.CertificateContractDataGateway, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(certificatecontract_datagateway.CertificateContractDataGateway), args.Error(1)
}

// MockCertificateContractDg implements certificatecontract_datagateway.CertificateContractDataGateway
type MockCertificateContractDg struct {
	mock.Mock
}

var _ certificatecontract_datagateway.CertificateContractDataGateway = (*MockCertificateContractDg)(nil)

func (m *MockCertificateContractDg) GetTokenData(ctx context.Context, tokenId *big.Int) (*entity.CertificatePayload, error) {
	args := m.Called(ctx, tokenId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CertificatePayload), args.Error(1)
}

func (m *MockCertificateContractDg) UsedSignatures(ctx context.Context, signature []byte) (bool, error) {
	args := m.Called(ctx, signature)
	return args.Bool(0), args.Error(1)
}

func (m *MockCertificateContractDg) MintNft(ctx context.Context, params certificatecontract_datagateway.MintNftParams) (*big.Int, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockCertificateContractDg) RevokeCertificate(ctx context.Context, tokenId *big.Int, signedMessageDigest string, signature []byte) error {
	args := m.Called(ctx, tokenId, signedMessageDigest, signature)
	return args.Error(0)
}

func (m *MockCertificateContractDg) FilterCertificateMinted(ctx context.Context, fromBlock uint64, toBlock uint64, receiverAddresses []common.Address) ([]*certificatecontract_datagateway.CertificateMintedEvent, error) {
	args := m.Called(ctx, fromBlock, toBlock, receiverAddresses)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*certificatecontract_datagateway.CertificateMintedEvent), args.Error(1)
}

func (m *MockCertificateContractDg) ParseCertificateMinted(log types.Log) (*certificatecontract_datagateway.CertificateMintedEvent, error) {
	args := m.Called(log)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*certificatecontract_datagateway.CertificateMintedEvent), args.Error(1)
}

// --- helpers ---

func strPtr(s string) *string { return &s }
