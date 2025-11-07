package entity

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

type EventStatus string

const (
	EventStatusActive   EventStatus = "active"
	EventStatusInactive EventStatus = "inactive"
	EventStatusClosed   EventStatus = "closed"
)

const (
	EventTypePublic  EventType = "public"
	EventTypePrivate EventType = "private"
	EventTypeInvite  EventType = "invite"
)

type Event struct {
	ID                       uuid.UUID   `json:"id"`
	EventType                EventType   `json:"event_type"`
	ChainID                  int         `json:"chain_id"`
	ContactNumber            string      `json:"contact_number"`
	ContactAddress           string      `json:"contact_address"`
	OwnerCredentialID        uuid.UUID   `json:"owner_credential_id"`
	BannerStorageKey         string      `json:"banner_storage_key"`
	IconStorageKey           string      `json:"icon_storage_key"`
	Title                    string      `json:"title"`
	ShortDescription         string      `json:"short_description"`
	LongDescription          string      `json:"long_description"`
	StartDate                time.Time   `json:"start_date"`
	EndDate                  time.Time   `json:"end_date"`
	Location                 string      `json:"location"`
	GoogleMapQuery           string      `json:"google_map_query"`
	MaxAttendees             int         `json:"max_attendees"`
	IsPublic                 bool        `json:"is_public"`
	IsBookingRequestRequired bool        `json:"is_booking_request_required"`
	IsVerified               bool        `json:"is_verified"`
	IsTicketTransferable     bool        `json:"is_ticket_transferable"`
	CreatedAt                time.Time   `json:"created_at"`
	UpdatedAt                time.Time   `json:"updated_at"`
	EventStatus              EventStatus `json:"event_status"`
}
