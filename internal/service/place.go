package service

import (
	"context"
	"fmt"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type PlaceService struct {
	placeStorage storage.IPlaceStorage
	tripStorage  storage.ITripStorage
}

func NewPlaceService(placeStorage storage.IPlaceStorage, tripStorage storage.ITripStorage) service.IPlaceService {
	return &PlaceService{
		placeStorage: placeStorage,
		tripStorage:  tripStorage,
	}
}

func (service *PlaceService) AddPlaceToTrip(ctx context.Context, tripID int, placeID string) error {
	//todo: check for this placeID in our bd
	// if exists: create relation for tripID

	//todo: go to google place api and take info about

	trip, err := service.tripStorage.GetTripByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("fail to get trip from storage: %w", err)
	}

	newPlace := model.Place{
		ID:    placeID,
		Trips: []*model.Trip{&trip},
	}

	err = service.placeStorage.CreatePlace(ctx, newPlace)
	if err != nil {
		return fmt.Errorf("fail to add new place: %w", err)
	}

	return nil
}
