package event

import (
	"apps/backend/core-api/config"
	"apps/backend/services/auth"
	"context"
	"log/slog"
	"mime/multipart"

	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	eventcontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/event_contract"
	storage_datagateway "apps/backend/core-api/internal/datagateway/storage"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

type EventUsecase struct {
	EventDataGateway                     eventdatagateway.EventDataGateway
	EventContractDataGateway             eventdatagateway.EventContractDataGateway
	EventContractFactoryDg               eventcontract_datagateway.EventContractFactoryDataGateway
	EventIssuerDataGateway               eventdatagateway.EventIssuerDataGateway
	EventCertificateDataGateway          eventdatagateway.EventCertificateDataGateway
	CertificateShareDataGateway          eventdatagateway.CertificateShareDataGateway
	EventCertificateSignatureDataGateway eventdatagateway.EventCertificateSignatureDataGateway
	EventCertificateConfigDg             eventdatagateway.EventCertificateConfigDataGateway
	EventCertificateFontFamilyDg         eventdatagateway.EventCertificateFontFamilyDataGateway
	AuthenticationCredentialDg           offchain_datagateway.AuthenticationCredentialDataGateway
	EventRegistrationInvitationDg        eventdatagateway.EventRegistrationInvitationDataGateway
	EventAttendeeDg                      eventdatagateway.EventAttendeeDataGateway
	UserSignatureDg                      offchain_datagateway.UserSignatureDataGateway
	InboxMessageDg                       offchain_datagateway.InboxMessageDataGateway
	BlockchainClientDg                   blockchainclient_datagateway.BlockchainClientDataGateway
	S3DataGateway                        storage_datagateway.S3DataGateway
	logger                               *slog.Logger
	authService                          *auth.AuthService
	cfg                                  *config.Config
	ethClient                            *ethclient.Client
	UploadEventBanner                    func(ctx context.Context, entityID uuid.UUID, bannerFile *multipart.FileHeader) (string, error)
	UploadEventIcon                      func(ctx context.Context, entityID uuid.UUID, iconFile *multipart.FileHeader) (string, error)
	deployCertificateContract            func(ctx context.Context, transactor *bind.TransactOpts, accessManagerAddr, eventAddr common.Address) (common.Address, error)
}

func NewEventUsecase(eventDataGateway eventdatagateway.EventDataGateway, eventContractDataGateway eventdatagateway.EventContractDataGateway, eventContractFactoryDg eventcontract_datagateway.EventContractFactoryDataGateway, eventIssuerDataGateway eventdatagateway.EventIssuerDataGateway, eventCertificateDataGateway eventdatagateway.EventCertificateDataGateway, certificateShareDataGateway eventdatagateway.CertificateShareDataGateway, eventCertificateSignatureDataGateway eventdatagateway.EventCertificateSignatureDataGateway, eventCertificateConfigDg eventdatagateway.EventCertificateConfigDataGateway, eventCertificateFontFamilyDg eventdatagateway.EventCertificateFontFamilyDataGateway, authenticationCredentialDg offchain_datagateway.AuthenticationCredentialDataGateway, eventRegistrationInvitationDg eventdatagateway.EventRegistrationInvitationDataGateway, eventAttendeeDg eventdatagateway.EventAttendeeDataGateway, userSignatureDg offchain_datagateway.UserSignatureDataGateway, inboxMessageDg offchain_datagateway.InboxMessageDataGateway, blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway, ethClient *ethclient.Client, s3DataGateway storage_datagateway.S3DataGateway, logger *slog.Logger, authService *auth.AuthService, cfg *config.Config) *EventUsecase {
	uc := &EventUsecase{
		EventDataGateway:                     eventDataGateway,
		EventContractDataGateway:             eventContractDataGateway,
		EventContractFactoryDg:               eventContractFactoryDg,
		EventIssuerDataGateway:               eventIssuerDataGateway,
		EventAttendeeDg:                      eventAttendeeDg,
		UserSignatureDg:                      userSignatureDg,
		EventCertificateDataGateway:          eventCertificateDataGateway,
		CertificateShareDataGateway:          certificateShareDataGateway,
		EventCertificateSignatureDataGateway: eventCertificateSignatureDataGateway,
		EventCertificateConfigDg:             eventCertificateConfigDg,
		EventCertificateFontFamilyDg:         eventCertificateFontFamilyDg,
		AuthenticationCredentialDg:           authenticationCredentialDg,
		EventRegistrationInvitationDg:        eventRegistrationInvitationDg,
		InboxMessageDg:                       inboxMessageDg,
		BlockchainClientDg:                   blockchainClientDg,
		S3DataGateway:                        s3DataGateway,
		logger:                               logger,
		authService:                          authService,
		cfg:                                  cfg,
		ethClient:                            ethClient,
	}
	// Initialize upload function fields with their default implementations
	uc.UploadEventBanner = uc.uploadEventBannerImpl
	uc.UploadEventIcon = uc.uploadEventIconImpl
	uc.deployCertificateContract = uc.deployCertificateContractImpl
	return uc
}
