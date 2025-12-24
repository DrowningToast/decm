package event

import (
	"apps/backend/core-api/internal/entity"
	"context"

	"github.com/google/uuid"
)

func (uc *EventUsecase) ListEventAttendees(ctx context.Context, eventId uuid.UUID) ([]entity.EventAttendee, error) {
	// Verify event exists
	_, err := uc.EventDataGateway.GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	// Get all attendees for the event
	attendees, err := uc.EventAttendeeDg.ListEventAttendeesByEventID(ctx, eventId)
	if err != nil {
		return nil, err
	}

	return attendees, nil
}
