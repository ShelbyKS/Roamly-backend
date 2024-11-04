package storage

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/google/uuid"
)

type IEventStorage interface {
	GetEventByID(ctx context.Context, eventID uuid.UUID) (model.Event, error)
	CreateEvent(ctx context.Context, event model.Event) error
	UpdateEvent(ctx context.Context, event model.Event) (model.Event, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	CreateBatchEvents(ctx context.Context, events *[]model.Event) error
	DeleteEventsByTrip(ctx context.Context, tripID uuid.UUID) error
}
