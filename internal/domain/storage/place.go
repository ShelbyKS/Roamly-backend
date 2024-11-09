package storage

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/google/uuid"
)

type IPlaceStorage interface {
	GetPlaceByID(ctx context.Context, placeID string) (model.Place, error)
	DeletePlace(ctx context.Context, tripID uuid.UUID, placeID string) error
	CreatePlace(ctx context.Context, place *model.Place) (model.Place, error)
	AppendPlaceToTrip(ctx context.Context, placeID string, tripID uuid.UUID) error
	UpdatePlace(ctx context.Context, place model.Place) error
}
