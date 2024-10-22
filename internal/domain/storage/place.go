package storage

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type IPlaceStorage interface {
	GetPlaceByID(ctx context.Context, placeID string) (model.Place, error)
	CreatePlace(ctx context.Context, place *model.Place) (model.Place, error)
}
