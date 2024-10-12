package service

import (
	"context"
	"fmt"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type TripService struct {
	tripStorage storage.ITripStorage
}

func NewTripService(tripStorage storage.ITripStorage) service.ITripService {
	return &TripService{
		tripStorage: tripStorage,
	}
}

func (service *TripService) GetTripByID(ctx context.Context, id int) (model.Trip, error) {
	trip, err := service.tripStorage.GetTripByID(ctx, id)
	if err != nil {
		return model.Trip{}, fmt.Errorf("fail to get trip from storage: %w", err)
	}

	return trip, nil
}

func (service *TripService) DeleteTrip(ctx context.Context, id int) error {
	err := service.tripStorage.DeleteTrip(ctx, id)
	if err != nil {
		return fmt.Errorf("fail to delete trip from storage: %w", err)
	}

	return nil
}

func (service *TripService) CreateTrip(ctx context.Context, trip model.Trip) error {
	err := service.tripStorage.CreateTrip(ctx, trip)
	if err != nil {
		return fmt.Errorf("fail to create trip from storage: %w", err)
	}

	return nil
}

func (service *TripService) UpdateTrip(ctx context.Context, trip model.Trip) error {
	err := service.tripStorage.UpdateTrip(ctx, trip)
	if err != nil {
		return fmt.Errorf("fail to update trip from storage: %w", err)
	}

	return nil
}
