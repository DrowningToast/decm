package event

import (
	"apps/backend/core-api/internal/usecase/event"
)

type Handler struct {
	EventUc *event.EventUsecase
}

func NewHandler(eventUc *event.EventUsecase) *Handler {
	return &Handler{
		EventUc: eventUc,
	}
}
