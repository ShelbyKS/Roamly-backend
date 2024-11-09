package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/clients"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type PlaceService struct {
	placeStorage storage.IPlaceStorage
	tripStorage  storage.ITripStorage
	eventStorage storage.IEventStorage
	googleApi    clients.IGoogleApiClient
	openAIClient clients.IChatClient
}

func NewPlaceService(
	placeStorage storage.IPlaceStorage,
	tripStorage storage.ITripStorage,
	googleApi clients.IGoogleApiClient,
	eventStorage storage.IEventStorage,
	openAIClient clients.IChatClient,
) service.IPlaceService {

	return &PlaceService{
		placeStorage: placeStorage,
		tripStorage:  tripStorage,
		googleApi:    googleApi,
		eventStorage: eventStorage,
		openAIClient: openAIClient,
	}
}

func (service *PlaceService) FindPlace(ctx context.Context, searchString string) ([]model.Place, error) {
	places, err := service.googleApi.FindPlace(ctx, searchString, []string{
		"formatted_address",
		"name",
		"rating",
		"geometry",
		"photo",
	})
	if err != nil {
		return []model.Place{}, fmt.Errorf("fail to find place: %w", err)
	}

	placesDomain := make([]model.Place, len(places))
	for i, place := range places {
		placesDomain[i] = model.Place{
			ID:          place.PlaceID,
			GooglePlace: place,
		}
	}

	return placesDomain, nil
}

func (service *PlaceService) DeletePlace(ctx context.Context, tripID uuid.UUID, placeID string) (model.Trip, error) {
	trip, err := service.tripStorage.GetTripByID(ctx, tripID)
	if err != nil {
		return model.Trip{}, fmt.Errorf("trip not found: %w", err)
	}

	err = service.placeStorage.DeletePlace(ctx, tripID, placeID)
	if err != nil {
		return model.Trip{}, fmt.Errorf("can't delete place: %w", err)
	}

	err = service.eventStorage.DeleteEventsByPlace(ctx, tripID, placeID)
	if err != nil {
		return model.Trip{}, fmt.Errorf("can't delete related events: %w", err)
	}

	trip, err = service.tripStorage.GetTripByID(ctx, tripID)
	if err != nil {
		return model.Trip{}, fmt.Errorf("trip after deleting not found: %w", err)
	}

	return trip, nil
}

func (service *PlaceService) AddPlaceToTrip(ctx context.Context, tripID uuid.UUID, placeID string) (model.Trip, error) {
	trip, err := service.tripStorage.GetTripByID(ctx, tripID)
	if err != nil {
		return model.Trip{}, fmt.Errorf("fail to get trip from storage: %w", err)
	}

	place, err := service.placeStorage.GetPlaceByID(ctx, placeID)
	if err != nil && !errors.Is(err, domain.ErrPlaceNotFound) {
		return model.Trip{}, fmt.Errorf("can't get place by id %w", err)
	}

	if !errors.Is(err, domain.ErrPlaceNotFound) {
		err := service.placeStorage.AppendPlaceToTrip(ctx, place.ID, trip.ID)
		if err != nil {
			return model.Trip{}, fmt.Errorf("can't append place to trip: %w", err)
		}

		trip.Places = append(trip.Places, &place)
		return trip, nil
	}

	googlePlace, err := service.googleApi.GetPlaceByID(ctx, placeID, []string{
		"formatted_address",
		"name",
		"rating",
		"geometry",
		"photo",
	})
	if err != nil {
		return model.Trip{}, fmt.Errorf("can't get place from api: %w", err)
	}
	place = model.Place{
		ID:          placeID,
		GooglePlace: googlePlace,
		// Opening:     time.Now(),
		// Closing:     time.Now(),
	}

	place.Trips = []*model.Trip{&trip}

	newPlace, err := service.placeStorage.CreatePlace(ctx, &place)
	if err != nil {
		return model.Trip{}, fmt.Errorf("fail to add new place: %w", err)
	}

	trip.Places = append(trip.Places, &newPlace)

	return trip, nil

}

func (service *PlaceService) GetPlacesNearby(ctx context.Context, lat float64, lng float64, placesTypes []string) ([]model.GooglePlace, error) {
	maxPlaces := 10
	rankPreference := "DISTANCE"
	radius := 20000.

	places, err := service.googleApi.GetPlacesNearby(ctx,
		placesTypes,
		maxPlaces,
		rankPreference,
		lat,
		lng,
		radius)

	if err != nil {
		return []model.GooglePlace{}, fmt.Errorf("error get places nerby: %w", err)
	}

	return places, nil
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

func (service *PlaceService) DetermineRecommendedDuration(ctx context.Context, placeID string) error {
	place, err := service.placeStorage.GetPlaceByID(ctx, placeID)
	if err != nil {
		return fmt.Errorf("can't get place by id %w", err)
	}

	//todo: сделать какой-то отдельный файл для промптов
	var prompt strings.Builder
	prompt.WriteString(fmt.Sprintf("Определи оптимальное время для посещения %s\n", place.GooglePlace.Name))
	prompt.WriteString("Напиши только  число - время в минутах")

	recommendedDurationStr, err := service.openAIClient.PostPrompt(ctx, prompt.String(), clients.ModelChatGPT4oMini)
	if err != nil {
		return fmt.Errorf("can't get recommended duration: %w", err)
	}
	recommendedDurationInt, err := strconv.Atoi(recommendedDurationStr)
	if err != nil {
		return fmt.Errorf("recommended duration has wrong format: %w", err)
	}

	place.RecommendedVisitingDuration = recommendedDurationInt

	err = service.placeStorage.UpdatePlace(ctx, place)
	if err != nil {
		return fmt.Errorf("can't update place recommended duration: %w", err)
	}

	return nil
}
