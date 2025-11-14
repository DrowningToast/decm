package event_registration

import (
	datagateway "apps/backend/core-api/internal/datagateway"
	eventdatagateway "apps/backend/core-api/internal/datagateway/event"
	authUc "apps/backend/core-api/internal/usecase/auth"
	eventUc "apps/backend/core-api/internal/usecase/event"

	"github.com/google/uuid"
)

type Participant struct {
	Name                string `json:"name"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Email               string `json:"email"`
	PhoneNumber         string `json:"phone_number"`
	AcademicInstitution string `json:"academic_institution"`
}

type ImportEventParticipantsParameters struct {
	EventID      uuid.UUID
	Participants []Participant
}

type EventRegistrationUsecase struct {
	InboxMessageDg                   datagateway.InboxMessageDataGateway
	EventRegistrationInvitationDg    datagateway.EventRegistrationInvitationDataGateway
	EventDg                          eventdatagateway.EventDataGateway
	EventContractDg                  eventdatagateway.EventContractDataGateway
	EventAttendeeDg                  datagateway.EventAttendeeDataGateway
	EventRegistrationConfigurationDg eventdatagateway.EventRegistrationConfigDataGateway
	AuthenticationCredentialDg       datagateway.AuthenticationCredentialDataGateway

	AuthUsecase  authUc.AuthUsecase
	EventUsecase eventUc.EventUsecase
}

func NewEventRegistrationInvitationUsecase(
	inboxMessageDg datagateway.InboxMessageDataGateway,
	eventRegistrationInvitationDg datagateway.EventRegistrationInvitationDataGateway,
	eventDg eventdatagateway.EventDataGateway,
	eventContractDg eventdatagateway.EventContractDataGateway,
	eventAttendeeDg datagateway.EventAttendeeDataGateway,
	eventRegistrationConfigurationDg eventdatagateway.EventRegistrationConfigDataGateway,
	authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway,
	authUsecase authUc.AuthUsecase,
	eventUsecase eventUc.EventUsecase,
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
	}
}
