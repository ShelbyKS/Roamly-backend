package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/clients"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
	"github.com/google/uuid"
)

type TripService struct {
	tripStorage     storage.ITripStorage
	placeStorage    storage.IPlaceStorage
	googleApiClient clients.IGoogleApiClient
}

func NewTripService(tripStorage storage.ITripStorage, placeStorage storage.IPlaceStorage, googleApiClient clients.IGoogleApiClient) service.ITripService {
	return &TripService{
		tripStorage:     tripStorage,
		placeStorage:    placeStorage,
		googleApiClient: googleApiClient,
	}
}

func (service *TripService) GetTripByID(ctx context.Context, id uuid.UUID) (model.Trip, error) {
	trip, err := service.tripStorage.GetTripByID(ctx, id)
	if errors.Is(err, domain.ErrTripNotFound) {
		return model.Trip{}, err
	}
	if err != nil {
		return model.Trip{}, fmt.Errorf("fail to get trip from storage: %w", err)
	}

	return trip, nil
}

func (service *TripService) GetTrips(ctx context.Context) ([]model.Trip, error) {
	trips, err := service.tripStorage.GetTrips(ctx)
	if err != nil {
		return []model.Trip{}, fmt.Errorf("fail to get trip from storage: %w", err)
	}

	return trips, nil
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
		areaGoogle, err := service.googleApiClient.GetPlaceByID(ctx, trip.AreaID, []string{
			"formatted_address",
			"name",
			"rating",
			"geometry",
			"photo",
		})

		log.Println("in service:", areaGoogle)

		area, err = service.placeStorage.CreatePlace(ctx, &model.Place{
			ID:          trip.AreaID,
			GooglePlace: areaGoogle,
		})
		if err != nil {
			return uuid.Nil, fmt.Errorf("fail to create area from storage: %w", err)
		}
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
