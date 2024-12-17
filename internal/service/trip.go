package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"sync"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/clients"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type TripService struct {
	tripStorage     storage.ITripStorage
	placeStorage    storage.IPlaceStorage
	googleApiClient clients.IGoogleApiClient
	openAIClient    clients.IChatClient
	sessionStorage  storage.ISessionStorage
	messageProducer clients.IMessageProdcuer
	aiChatStorage   storage.IAIChatStorage
}

func NewTripService(
	tripStorage storage.ITripStorage,
	placeStorage storage.IPlaceStorage,
	googleApiClient clients.IGoogleApiClient,
	openAIClient clients.IChatClient,
	sessionStorage storage.ISessionStorage,
	messageProducer clients.IMessageProdcuer,
	aiChatStorage storage.IAIChatStorage,
) service.ITripService {
	return &TripService{
		tripStorage:     tripStorage,
		placeStorage:    placeStorage,
		googleApiClient: googleApiClient,
		openAIClient:    openAIClient,
		sessionStorage:  sessionStorage,
		messageProducer: messageProducer,
		aiChatStorage:   aiChatStorage,
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

func (service *TripService) GetTrips(ctx context.Context, userId int) ([]model.Trip, error) {
	trips, err := service.tripStorage.GetTrips(ctx, userId)
	if err != nil {
		return []model.Trip{}, fmt.Errorf("fail to get trip from storage: %w", err)
	}

	return trips, nil
}

func (service *TripService) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	err := service.tripStorage.DeleteTrip(ctx, id)
	if errors.Is(err, domain.ErrTripNotFound) {
		return err
	}
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
		areaGoogle, err := service.googleApiClient.GetPlaceByID(ctx, trip.AreaID, []string{
			"formatted_address",
			"name",
			"rating",
			"geometry",
			"photo",
		})

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

	err = service.tripStorage.CreateTrip(ctx, trip, model.Owner)
	if err != nil {
		return uuid.Nil, fmt.Errorf("fail to create trip from storage: %w", err)
	}

	aiChatMsg := model.ChatMessage{
		TripID:  trip.ID,
		Role:    model.RoleSystem,
		Content: "Ты помощник для планирования путешествия",
	}
	_, err = service.openAIClient.PostPrompt(ctx, []model.ChatMessage{aiChatMsg}, clients.ModelChatGPT4o)

	err = service.aiChatStorage.SaveAIChatMessage(ctx, aiChatMsg)
	if err != nil {
		return uuid.Nil, fmt.Errorf("fail to save AI chat message: %w", err)
	}

	return trip.ID, nil
}

func (service *TripService) DetermineRecommendedPlaces(ctx context.Context, tripID uuid.UUID) error {
	trip, err := service.tripStorage.GetTripByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("fail to get trip from storage: %w", err)
	}

	recommendedPlacesNames, err := service.getRecommendedPlacesNames(ctx, trip.Area.GooglePlace.Name)
	if err != nil {
		return fmt.Errorf("fail to get recommended places names from openai: %w", err)
	}

	recommendedPlacesDomain, err := service.GetRecommendedPlacesDomain(ctx, recommendedPlacesNames, trip.Area.GooglePlace.Name)
	if err != nil {
		return fmt.Errorf("fail to get recommended places domains from google: %w", err)
	}

	trip.RecommendedPlaces = recommendedPlacesDomain

	err = service.tripStorage.UpdateTrip(ctx, trip)
	if err != nil {
		return fmt.Errorf("fail to update trip from storage: %w", err)
	}

	users := trip.Users
	for _, user := range users {
		cooks, _ := service.sessionStorage.GetTokensByUserID(ctx, user.ID)
		var message model.NotifyMessage
		message.Payload.Action = "trip_auto_planning_enable"
		message.Payload.TripID = trip.ID
		message.Payload.Author = fmt.Sprintf("%d", user.ID)
		message.Payload.Message = "Вам доступно экспресс планирование!"
		message.Clients = cooks
		service.messageProducer.SendMessage(message)
	}

	return nil
}

func (service *TripService) getRecommendedPlacesNames(ctx context.Context, area string) ([]string, error) {
	//todo: сделать какой-то отдельный файл для промптов
	var prompt strings.Builder
	prompt.WriteString(fmt.Sprintf("Какие главные достопримечательности нужно посетить в %s\n", area))
	prompt.WriteString("Без описания, через запятую, 9 штук")

	recommendedPlacesStr, err := service.openAIClient.PostPrompt(ctx, []model.ChatMessage{{
		Role:    model.RoleUser,
		Content: prompt.String(),
	}}, clients.ModelChatGPT4oMini)

	if err != nil {
		return nil, fmt.Errorf("can't get recommended duration: %w", err)
	}

	places := strings.Split(recommendedPlacesStr, ", ")

	return places, nil
}

func (service *TripService) GetRecommendedPlacesDomain(ctx context.Context, recommendedPlacesNames []string, area string) ([]*model.Place, error) {
	var recommendedPlacesDomain []*model.Place

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for _, recommendedPlace := range recommendedPlacesNames {
		wg.Add(1)
		go func() {
			defer wg.Done()

			searchStr := fmt.Sprintf("%s %s", area, recommendedPlace)
			places, err := service.googleApiClient.FindPlace(ctx, searchStr, []string{
				"formatted_address",
				"name",
				"rating",
				"geometry",
				"photo",
				"place_id",
			})
			if err != nil {
				fmt.Println("fail to find place %s: %w", recommendedPlace, err)
				return
			}

			//todo: сделать какой-то отдельный файл для промптов
			var prompt strings.Builder
			prompt.WriteString(fmt.Sprintf("Определи оптимальное время для посещения %s\n", places[0].Name))
			prompt.WriteString("Напиши только  число - время в минутах")

			recommendedDurationStr, err := service.openAIClient.PostPrompt(ctx, []model.ChatMessage{{
				Role:    model.RoleUser,
				Content: prompt.String(),
			}}, clients.ModelChatGPT4oMini)

			if err != nil {
				fmt.Println("can't get recommended duration: %w", err)
				return
			}
			recommendedDurationInt, err := strconv.Atoi(recommendedDurationStr)
			if err != nil {
				fmt.Println("recommended duration has wrong format: %w", err)
				return
			}

			placeDomain := model.Place{
				ID:                          places[0].PlaceID,
				GooglePlace:                 places[0],
				RecommendedVisitingDuration: recommendedDurationInt,
			}

			_, err = service.placeStorage.CreatePlace(ctx, &placeDomain)
			if err != nil && !errors.Is(err, domain.ErrPlaceAlreadyExists) {
				fmt.Println("fail to create place: %s: %w", recommendedPlace, err)
				return
			}

			mu.Lock()
			recommendedPlacesDomain = append(recommendedPlacesDomain, &placeDomain)
			mu.Unlock()
		}()

	}
	wg.Wait()

	return recommendedPlacesDomain, nil
}

func (service *TripService) UpdateTrip(ctx context.Context, trip model.Trip) error {
	err := service.tripStorage.UpdateTrip(ctx, trip)
	if errors.Is(err, domain.ErrTripNotFound) {
		return err
	}
	if err != nil {
		return fmt.Errorf("fail to update trip from storage: %w", err)
	}

	tripFound, _ := service.tripStorage.GetTripByID(ctx, trip.ID)
	users := tripFound.Users
	for _, user := range users {
		cooks, _ := service.sessionStorage.GetTokensByUserID(ctx, user.ID)
		// cookies = append(cookies, cooks...)
		var message model.NotifyMessage
		message.Payload.Action = "trip_update"
		message.Payload.TripID = trip.ID
		message.Payload.Author = fmt.Sprintf("%d", user.ID)
		message.Payload.Message = "Поездка обновилась"
		message.Clients = cooks
		service.messageProducer.SendMessage(message)
	}

	return nil
}

func (service *TripService) GetUserRole(ctx context.Context, userID int, tripID uuid.UUID) (model.UserTripRole, error) {
	role, err := service.tripStorage.GetUserRole(ctx, userID, tripID)
	if err != nil {
		return -1, fmt.Errorf("can't get user role: %w", err)
	}

	return role, nil
}

func (service *TripService) GetTripByEventID(ctx context.Context, eventID uuid.UUID) (model.Trip, error) {
	trip, err := service.tripStorage.GetTripByEventID(ctx, eventID)
	if err != nil {
		return model.Trip{}, err
	}

	return trip, nil
}

func (service *TripService) RemoveUserFromTrip(ctx context.Context, userID int, tripID uuid.UUID) error {
	err := service.tripStorage.RemoveUserFromTrip(ctx, userID, tripID)
	if err != nil {
		return err
	}

	return nil
}
