package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
	"github.com/google/uuid"
)

type TripService struct {
	tripStorage  storage.ITripStorage
	placeStorage storage.IPlaceStorage
}

func NewTripService(tripStorage storage.ITripStorage, placeStorage storage.IPlaceStorage) service.ITripService {
	return &TripService{
		tripStorage:  tripStorage,
		placeStorage: placeStorage,
	}
}

func (service *TripService) GetTripByID(ctx context.Context, id uuid.UUID) (model.Trip, error) {
	trip, err := service.tripStorage.GetTripByID(ctx, id)
	if err != nil {
		return model.Trip{}, fmt.Errorf("fail to get trip from storage: %w", err)
	}

	return trip, nil
}

func (service *TripService) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	err := service.tripStorage.DeleteTrip(ctx, id)
	if err != nil {
		return fmt.Errorf("fail to delete trip from storage: %w", err)
	}

	return nil
}

func (service *TripService) CreateTrip(ctx context.Context, trip model.Trip) (uuid.UUID, error) {
	area, err := service.placeStorage.GetPlaceByID(ctx, trip.AreaID)
	if err != nil && !errors.Is(err, domain.ErrPlaceNotFound) {
		return uuid.Nil, fmt.Errorf("fail to get area from storage: %w", err)
	}

	if errors.Is(err, domain.ErrPlaceNotFound) {
		//todo: go to google place api and take info about
		// create area in db
	}

	trip.Area = &area
	trip.ID = uuid.New()

	err = service.tripStorage.CreateTrip(ctx, trip)
	if err != nil {
		return uuid.Nil, fmt.Errorf("fail to create trip from storage: %w", err)
	}

	return trip.ID, nil
}

func (service *TripService) UpdateTrip(ctx context.Context, trip model.Trip) error {
	err := service.tripStorage.UpdateTrip(ctx, trip)
	if err != nil {
		return fmt.Errorf("fail to update trip from storage: %w", err)
	}

	return nil
}
