package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"log/slog"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type EventResponse struct {
	Id                       uuid.UUID          `json:"id"`
	ChainID                  int32              `json:"chain_id"`
	ContactNumber            string             `json:"contact_number"`
	ContactAddress           string             `json:"contact_address"`
	OwnerCredentialID        uuid.UUID          `json:"owner_credential_id"`
	BannerPresignedURL       string             `json:"banner_presigned_url"`
	IconPresignedURL         string             `json:"icon_presigned_url"`
	Title                    string             `json:"title"`
	ShortDescription         string             `json:"short_description"`
	LongDescription          string             `json:"long_description"`
	StartDate                time.Time          `json:"start_date"`
	EndDate                  time.Time          `json:"end_date"`
	Location                 string             `json:"location"`
	GoogleMapQuery           string             `json:"google_map_query"`
	MaxAttendees             int32              `json:"max_attendees"`
	AttendeesCount           int32              `json:"attendees_count"`
	IsPublic                 bool               `json:"is_public"`
	IsBookingRequestRequired bool               `json:"is_booking_request_required"`
	IsVerified               bool               `json:"is_verified"`
	IsTicketTransferable     bool               `json:"is_ticket_transferable"`
	CreatedAt                time.Time          `json:"created_at"`
	UpdatedAt                time.Time          `json:"updated_at"`
	EventStatus              entity.EventStatus `json:"event_status"`
	EventType                entity.EventType   `json:"event_type"`
}

type RegistrationConfigResponse struct {
	Id                                   uuid.UUID  `json:"id"`
	EventId                              uuid.UUID  `json:"event_id"`
	FinalCallForRegistration             *time.Time `json:"final_call_for_registration,omitempty"`
	IsIdentityVerificationRequired       bool       `json:"is_identity_verification_required"`
	FirstNameRequirementStatus           int        `json:"first_name_requirement_status"`
	LastNameRequirementStatus            int        `json:"last_name_requirement_status"`
	EmailRequirementStatus               int        `json:"email_requirement_status"`
	BioRequirementStatus                 int        `json:"bio_requirement_status"`
	PhoneNumberRequirementStatus         int        `json:"phone_number_requirement_status"`
	AddressRequirementStatus             int        `json:"address_requirement_status"`
	AcademicInstitutionRequirementStatus int        `json:"academic_institution_requirement_status"`
	AcademicEmailRequirementStatus       int        `json:"academic_email_requirement_status"`

	IsRegistrationPasswordRequired bool `json:"is_registration_password_required"`
}

type EventContractResponse struct {
	ID                           uuid.UUID `json:"id"`
	EventID                      uuid.UUID `json:"event_id"`
	AccessManagerContractAddress string    `json:"access_manager_contract_address"`
	EventContractAddress         string    `json:"event_contract_address"`
	TicketContractAddress        *string   `json:"ticket_contract_address,omitempty"`
	CertificateContractAddress   *string   `json:"certificate_contract_address,omitempty"`
}

type EventViewModel struct {
	EventResponse
	RegistrationConfig RegistrationConfigResponse `json:"registration_config"`
	EventContract      EventContractResponse      `json:"event_contract"`

	IsInvited bool `json:"is_invited"`
	IsJoined  bool `json:"is_joined,omitempty"`
	IsFull    bool `json:"is_full"`
}

func (u *EventUsecase) ListEventsByOwnerCredentialID(ctx context.Context, ownerCredentialID uuid.UUID, limitCount int32, offsetCount int32) ([]*entity.Event, error) {
	events, err := u.EventDataGateway.ListEventsByOwnerCredentialID(ctx, ownerCredentialID, limitCount, offsetCount)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (u *EventUsecase) GetEventById(ctx context.Context, eventId uuid.UUID) (*entity.Event, error) {
	event, err := u.EventDataGateway.GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (u *EventUsecase) GetEventViewModelByEventId(ctx context.Context, eventId uuid.UUID, currentUser *auth.JwtClaims) (*EventViewModel, error) {
	// Get event with registration config and contracts in a single query
	event, registrationConfig, eventContract, err := u.EventDataGateway.GetViewModelById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	email := currentUser.Email
	walletAddress := currentUser.WalletAddress

	isInvited := false
	isJoined := false

	invitation, _, err := u.EventRegistrationInvitationDg.GetEventRegistrationInvitationByEventIDAndCredential(ctx, eventId, currentUser.UserId, email, &walletAddress)
	if err != nil {
		var customError *customerror.Err
		if errors.As(err, &customError) {
			if *customError.Code != customerror.ErrNotFound.Code {
				return nil, errors.Wrap(err, "failed to get event registration invitation by event id and credential")
			}
			slog.InfoContext(ctx, "No row found in event_registration_invitations table (handled as Not Found)", "event_id", eventId, "user_id", currentUser.UserId)
		}
	}
	if invitation != nil {
		isInvited = true
	}
	// check if user is invited to the event, checks if the user is joined or not
	if isInvited {
		attendee, err := u.EventAttendeeDg.GetEventAttendeeByEventIdAndCredentialId(ctx, eventId, currentUser.UserId)
		if err != nil {
			var customError *customerror.Err
			if errors.As(err, &customError) {
				if *customError.Code != customerror.ErrNotFound.Code {
					return nil, errors.Wrap(err, "failed to get event attendee by event id and credential id")
				}
			} else {
				return nil, err
			}
		}
		if attendee != nil {
			isJoined = true
		}
	}

	// Convert to EventResponse with presigned URLs
	eventResponse, err := u.ToEventResponse(ctx, event)
	if err != nil {
		return nil, err
	}
	if eventResponse == nil {
		return nil, customerror.NewWithPreset(&customerror.ErrInternalServer, errors.New("event response is nil"))
	}

	// Convert RegistrationConfig to response format
	registrationConfigResponse := RegistrationConfigResponse{
		Id:                                   registrationConfig.ID,
		EventId:                              registrationConfig.EventID,
		FinalCallForRegistration:             registrationConfig.FinalCallForRegistration,
		IsIdentityVerificationRequired:       registrationConfig.IsIdentityVerificationRequired,
		FirstNameRequirementStatus:           registrationConfig.FirstNameRequirementStatus,
		LastNameRequirementStatus:            registrationConfig.LastNameRequirementStatus,
		EmailRequirementStatus:               registrationConfig.EmailRequirementStatus,
		BioRequirementStatus:                 registrationConfig.BioRequirementStatus,
		PhoneNumberRequirementStatus:         registrationConfig.PhoneNumberRequirementStatus,
		AddressRequirementStatus:             registrationConfig.AddressRequirementStatus,
		AcademicInstitutionRequirementStatus: registrationConfig.AcademicInstitutionRequirementStatus,
		AcademicEmailRequirementStatus:       registrationConfig.AcademicEmailRequirementStatus,
	}

	// Convert EventContract to response format
	eventContractResponse := EventContractResponse{
		ID:                           eventContract.ID,
		EventID:                      eventContract.EventID,
		AccessManagerContractAddress: eventContract.AccessManagerContractAddress,
		EventContractAddress:         eventContract.EventContractAddress,
		TicketContractAddress:        eventContract.TicketContractAddress,
		CertificateContractAddress:   eventContract.CertificateContractAddress,
	}

	return &EventViewModel{
		EventResponse:      *eventResponse,
		RegistrationConfig: registrationConfigResponse,
		EventContract:      eventContractResponse,
		IsInvited:          isInvited,
		IsJoined:           isJoined,
		IsFull:             event.MaxAttendees > 0 && event.AttendeesCount >= event.MaxAttendees,
	}, nil
}

func (u *EventUsecase) GetEventCertificatesByEventId(ctx context.Context, eventID uuid.UUID, currentUser *auth.JwtClaims) ([]*entity.EventCertificate, error) {
	// Get certificates for the event
	certificates, err := u.EventCertificateDataGateway.GetEventCertificatesByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}
