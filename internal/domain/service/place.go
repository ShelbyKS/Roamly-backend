package service

import (
	"context"
)

type IPlaceService interface {
	AddPlaceToTrip(ctx context.Context, tripID int, placeID string) error
}
