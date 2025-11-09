package event

import (
	"context"

	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
)

type ListEventsParameters struct {
	OnlyUserJoinedEvents  *auth.JwtClaims `query:"only_user_joined_events,omitempty" validate:"omitempty,boolean"`
	IncludeActiveEvents   bool            `query:"include_accepted_events,omitempty" validate:"omitempty,boolean"`
	IncludeInactiveEvents bool            `query:"include_inactive_events,omitempty" validate:"omitempty,boolean"`
	IncludeClosedEvents   bool            `query:"include_closed_events,omitempty" validate:"omitempty,boolean"`
}

func (uc *EventUsecase) ListEvents(ctx context.Context, params ListEventsParameters) ([]*entity.Event, error) {
	if params.OnlyUserJoinedEvents != nil {
		events, err := uc.EventDataGateway.ListEventsByEventAttendeeCredentialID(ctx, params.OnlyUserJoinedEvents.UserId, nil, nil)
		if err != nil {
			return nil, err
		}
		filteredEvents := make([]*entity.Event, 0)
		for _, event := range events {
			if event.EventStatus == entity.EventStatusActive && params.IncludeActiveEvents {
				filteredEvents = append(filteredEvents, event)
			}
			if event.EventStatus == entity.EventStatusInactive && params.IncludeInactiveEvents {
				filteredEvents = append(filteredEvents, event)
			}
			if event.EventStatus == entity.EventStatusClosed && params.IncludeClosedEvents {
				filteredEvents = append(filteredEvents, event)
			}
		}
		return filteredEvents, nil
	}
	events, err := uc.EventDataGateway.ListEvents(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	filteredEvents := make([]*entity.Event, 0)
	for _, event := range events {
		if event.EventStatus == entity.EventStatusActive && params.IncludeActiveEvents {
			filteredEvents = append(filteredEvents, event)
		}
		if event.EventStatus == entity.EventStatusInactive && params.IncludeInactiveEvents {
			filteredEvents = append(filteredEvents, event)
		}
		if event.EventStatus == entity.EventStatusClosed && params.IncludeClosedEvents {
			filteredEvents = append(filteredEvents, event)
		}
	}
	return events, nil
}
