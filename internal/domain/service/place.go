package service

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/google/uuid"
)

type IPlaceService interface {
	AddPlaceToTrip(ctx context.Context, tripID uuid.UUID, placeID string) error
	GetTimeMatrix(ctx context.Context, places []*model.Place) [][]int
	FindPlace(ctx context.Context, searchString string) ([]model.Place, error)
}
