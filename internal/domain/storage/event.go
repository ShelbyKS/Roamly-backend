package storage

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/google/uuid"
)

type IEventStorage interface {
	GetEventByID(ctx context.Context, placeID string, tripID uuid.UUID) (model.Event, error)
	CreateEvent(ctx context.Context, event *model.Event) error
	UpdateEvent(ctx context.Context, event model.Event) error
	DeleteEvent(ctx context.Context, placeID string, tripID uuid.UUID) error
}
