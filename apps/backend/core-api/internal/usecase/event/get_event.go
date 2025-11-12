package event

import (
	"context"
	"errors"
	"time"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
)

type EventResponse struct {
	ID                       uuid.UUID          `json:"id"`
	ChainID                  int32              `json:"chain_id"`
	ContactNumber            string             `json:"contact_number"`
	OwnerCredentialID        uuid.UUID          `json:"owner_credential_id"`
	BannerStorageKey         string             `json:"banner_storage_key"`
	IconStorageKey           string             `json:"icon_storage_key"`
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
	IsPublic                 bool               `json:"is_public"`
	IsBookingRequestRequired bool               `json:"is_booking_request_required"`
	IsVerified               bool               `json:"is_verified"`
	IsTicketTransferable     bool               `json:"is_ticket_transferable"`
	CreatedAt                time.Time          `json:"created_at"`
	UpdatedAt                time.Time          `json:"updated_at"`
	EventStatus              entity.EventStatus `json:"event_status"`
}

type EventViewModel struct {
	EventResponse

	IsInvited bool `json:"is_invited,omitempty"`
	IsJoined  bool `json:"is_joined,omitempty"`
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
	event, err := u.GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	email := currentUser.Email
	walletAddress := currentUser.WalletAddress

	isInvited := false
	isJoined := false

	invitation, inbox, err := u.EventRegistrationInvitationDg.GetEventRegistrationInvitationByEventIDAndCredential(ctx, eventId, currentUser.UserId, email, &walletAddress)
	if err != nil {
		return nil, err
	}
	if invitation != nil {
		isInvited = true
	}
	if inbox != nil {
		isJoined = true
	}
	// check if user is invited to the event, checks if the user is joined or not
	if isInvited {
		attendee, err := u.EventAttendeeDg.GetEventAttendeeByEventIDAndCredentialID(ctx, eventId, currentUser.UserId)
		if err != nil {
			return nil, err
		}
		if attendee != nil {
			isJoined = true
		}
	}

	eventResponse, err := u.ToEventResponse(ctx, event)
	if err != nil {
		return nil, err
	}
	if eventResponse == nil {
		return nil, customerror.NewWithPreset(&customerror.ErrInternalServer, errors.New("event response is nil"))
	}

	return &EventViewModel{
		EventResponse: *eventResponse,
		IsInvited:     isInvited,
		IsJoined:      isJoined,
	}, nil
}

func (u *EventUsecase) GetEventCertificatesByEventID(ctx context.Context, eventID uuid.UUID, currentUser *auth.JwtClaims) ([]*entity.EventCertificate, error) {
	// Get certificates for the event
	certificates, err := u.EventCertificateDataGateway.GetEventCertificatesByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}
