package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/clients"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/google/uuid"
	"golang.org/x/exp/rand"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type PlaceService struct {
	placeStorage storage.IPlaceStorage
	tripStorage  storage.ITripStorage
	googleApi    clients.IGoogleApiClient
}

func NewPlaceService(placeStorage storage.IPlaceStorage,
		tripStorage storage.ITripStorage,
		googleApi clients.IGoogleApiClient) service.IPlaceService {
	
		return &PlaceService{
		placeStorage: placeStorage,
		tripStorage:  tripStorage,
		googleApi: googleApi,
	}
}

func (service *PlaceService) FindPlace(ctx context.Context, searchString string) ([]model.Place, error) {
	places, err := service.googleApi.FindPlace(ctx,
		searchString,
		[]string{
			"formatted_address",
			"name",
			"rating",
			"geometry"})
	if err != nil {
		return []model.Place{}, fmt.Errorf("fail to find place: %w", err)
	}

	return places, nil
}

func (service *PlaceService) AddPlaceToTrip(ctx context.Context, tripID uuid.UUID, placeID string) error {
	//todo: check for this placeID in our bd
	// if exists: create relation for tripID

	//todo: go to google place api and take info about

	// мне кажется вот это не особо нужно(можно добавлять просто по айди поездки, а не всей поездке)
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

func (service *PlaceService) GetTimeMatrix(ctx context.Context, places []*model.Place) [][]int {
	if len(places) == 0 {
		return [][]int{}
	}

	matrixSize := len(places)
	timeMatrix := make([][]int, matrixSize)

	rand.Seed(uint64(time.Now().UnixNano()))

	for i := 0; i < matrixSize; i++ {
		timeMatrix[i] = make([]int, matrixSize)
		for j := 0; j < matrixSize; j++ {
			if i == j {
				timeMatrix[i][j] = 0
			} else {
				timeMatrix[i][j] = rand.Intn(60) + 1
			}
		}
	}

	return timeMatrix
}
