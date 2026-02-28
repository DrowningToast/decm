package event_registration

import (
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	eventcontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/event_contract"
	authUc "apps/backend/core-api/internal/usecase/auth"
	eventUc "apps/backend/core-api/internal/usecase/event"

	"github.com/google/uuid"
)

type Participant struct {
	Name                string  `json:"name"`
	FirstName           string  `json:"first_name"`
	LastName            string  `json:"last_name"`
	Email               string  `json:"email"`
	PhoneNumber         *string `json:"phone_number"`
	AcademicInstitution *string `json:"academic_institution"`
}

type ImportEventParticipantsParameters struct {
	EventID      uuid.UUID
	Participants []Participant
}

type EventRegistrationUsecase struct {
	InboxMessageDg                   offchain_datagateway.InboxMessageDataGateway
	EventRegistrationInvitationDg    eventdatagateway.EventRegistrationInvitationDataGateway
	EventDg                          eventdatagateway.EventDataGateway
	EventContractDg                  eventdatagateway.EventContractDataGateway
	EventAttendeeDg                  eventdatagateway.EventAttendeeDataGateway
	EventRegistrationConfigurationDg eventdatagateway.EventRegistrationConfigDataGateway
	UserSignatureDg                  offchain_datagateway.UserSignatureDataGateway
	AuthenticationCredentialDg       offchain_datagateway.AuthenticationCredentialDataGateway

	EventContractFactoryDg eventcontract_datagateway.EventContractFactoryDataGateway
	BlockchainClientDg     blockchainclient_datagateway.BlockchainClientDataGateway

	AuthUsecase  authUc.AuthUsecase
	EventUsecase eventUc.EventUsecase
}

func NewEventRegistrationUsecase(
	inboxMessageDg offchain_datagateway.InboxMessageDataGateway,
	eventRegistrationInvitationDg eventdatagateway.EventRegistrationInvitationDataGateway,
	eventDg eventdatagateway.EventDataGateway,
	eventContractDg eventdatagateway.EventContractDataGateway,
	eventAttendeeDg eventdatagateway.EventAttendeeDataGateway,
	eventRegistrationConfigurationDg eventdatagateway.EventRegistrationConfigDataGateway,
	authenticationCredentialDg offchain_datagateway.AuthenticationCredentialDataGateway,
	authUsecase authUc.AuthUsecase,
	eventUsecase eventUc.EventUsecase,
	userSignatureDg offchain_datagateway.UserSignatureDataGateway,
	eventContractFactoryDg eventcontract_datagateway.EventContractFactoryDataGateway,
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway,
) *EventRegistrationUsecase {
	return &EventRegistrationUsecase{
		InboxMessageDg:                   inboxMessageDg,
		EventRegistrationInvitationDg:    eventRegistrationInvitationDg,
		EventDg:                          eventDg,
		EventContractDg:                  eventContractDg,
		EventAttendeeDg:                  eventAttendeeDg,
		EventRegistrationConfigurationDg: eventRegistrationConfigurationDg,
		AuthenticationCredentialDg:       authenticationCredentialDg,
		AuthUsecase:                      authUsecase,
		EventUsecase:                     eventUsecase,
		UserSignatureDg:                  userSignatureDg,
		EventContractFactoryDg:           eventContractFactoryDg,
		BlockchainClientDg:               blockchainClientDg,
	}
}
