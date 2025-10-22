package event

import (
	"apps/backend/common/log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EventResponse struct {
	ID                       uuid.UUID `json:"id"`
	ChainID                  int32     `json:"chain_id"`
	ContactNumber            string    `json:"contact_number"`
	OwnerCredentialID        uuid.UUID `json:"owner_credential_id"`
	BannerStorageKey         string    `json:"banner_storage_key"`
	IconStorageKey           string    `json:"icon_storage_key"`
	BannerPresignedURL       string    `json:"banner_presigned_url"`
	IconPresignedURL         string    `json:"icon_presigned_url"`
	Title                    string    `json:"title"`
	ShortDescription         string    `json:"short_description"`
	LongDescription          string    `json:"long_description"`
	StartDate                time.Time `json:"start_date"`
	EndDate                  time.Time `json:"end_date"`
	Location                 string    `json:"location"`
	GoogleMapQuery           string    `json:"google_map_query"`
	MaxAttendees             int32     `json:"max_attendees"`
	IsPublic                 bool      `json:"is_public"`
	IsBookingRequestRequired bool      `json:"is_booking_request_required"`
	IsVerified               bool      `json:"is_verified"`
	IsTicketTransferable     bool      `json:"is_ticket_transferable"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted event routes")

	eventGroup := r.Group("/events").Use(
		h.AuthenticationGuardMiddleware.Middleware,
	)

	eventGroup.Post("/", h.CreateEvent)
	eventGroup.Post("/:event_id/contracts", h.CreateEventContract)
	eventGroup.Post("/:event_id/issuers", h.CreateEventIssuer)

	eventGroup.Get("/:event_id", h.GetEventById)
	eventGroup.Get("/:event_id/contracts", h.GetEventContractByEventID)
	eventGroup.Get("/:event_id/issuers", h.GetEventIssuersByEventID)
	eventGroup.Get("/:event_id/issuers/:issuer_id", h.GetEventIssuerByID)

	eventGroup.Get("/owner-credentials/:owner_credential_id", h.GetEventsByOwnerCredentialsId)

	eventGroup.Put("/:event_id", h.UpdateEvent)
	eventGroup.Put("/:event_id/contracts", h.UpdateEventContract)
	eventGroup.Put("/:event_id/issuers", h.UpdateEventIssuer)

	// eventGroup.Delete("/:event_id/contracts", h.DeleteEventContract)
	eventGroup.Delete("/:event_id", h.DeleteEvent)
	eventGroup.Delete("/:event_id/issuers/:issuer_id", h.DeleteEventIssuer)
	// eventGroup.Delete("/:event_id/issuers/:issuer_id", h.DeleteEventContract)
}
