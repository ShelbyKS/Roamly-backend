package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
	"github.com/google/uuid"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/clients"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
)

type SchedulerService struct {
	openAIClient clients.IChatClient
	googleApi    clients.IGoogleApiClient
	tripStorage  storage.ITripStorage
}

func NewShedulerService(
	openAIClient clients.IChatClient,
	googleApi clients.IGoogleApiClient,
	tripStorage storage.ITripStorage,
) service.ISchedulerService {
	return &SchedulerService{
		openAIClient: openAIClient,
		googleApi:    googleApi,
		tripStorage:  tripStorage,
	}
}

func (s *SchedulerService) GetSchedule(ctx context.Context, tripID uuid.UUID) (model.Schedule, error) {
	trip, err := s.tripStorage.GetTripByID(ctx, tripID)
	if err != nil {
		return model.Schedule{}, fmt.Errorf("failed to get trip for schedule: %w", err)
	}

	timeDistMatrix, err := s.googleApi.GetTimeDistanceMatrix(ctx, trip.GetTripPlaceIDs())
	if err != nil {
		return model.Schedule{}, fmt.Errorf("failed to get time distance matrix: %w", err)
	}

	prompt, err := s.generateRequestString(trip, trip.Places, timeDistMatrix)
	if err != nil {
		return model.Schedule{}, err
	}

	fmt.Println("PROMT: ", prompt)

	resp, err := s.openAIClient.PostPrompt(ctx, prompt)
	if err != nil {
		return model.Schedule{}, err
	}

	return ParseResponse(resp)
}

func (s *SchedulerService) generateRequestString(trip model.Trip, places []*model.Place, timeMatrix model.DistanceMatrix) (string, error) {
	var sb strings.Builder

	sb.WriteString(`СПЛАНИРУЙ ПОЕЗДКУ ТОЛЬКО ПО ДАННЫМ МЕСТАМ,
ИСПОЛЬЗУЯ ИНФОРМАЦИЮ, ПРИВЕДЕННУЮ НИЖЕ, ВЕРНИ PlaceID И ВРЕМЯ ПОСЕЩЕНИЯ,
В ОТВЕТЕ ОПИШИ ИМЕННО ТОЛЬКО JSON  объект, КОТОРЫЙ БУДЕТ ОПИСЫВАТЬ СПЛАНИРОВАННОЕ РАСПИСАНИЕ,
КРОМЕ JSON В ОТЕТЕ НИЧЕГО НЕ ДОЛЖНО БЫТЬ, МАРШУРТ ДОЛЖЕН БЫТЬ ОПТИМАЛЬНЫМ И УЧИТЫВАТЬ ВРЕМЯ НА ДОРОГУ МЕЖДУ МЕСТАМИ,
КОТОРЫЕ УКАЗАНЫ В МАТРИЦЕ ВРЕМЕНИ И РАССТОЯНИЙ:`)
	sb.WriteString(fmt.Sprintf("TripID: %s\n", trip.ID.String()))

	sb.WriteString("Дата поездки:\n")
	sb.WriteString(fmt.Sprintf("С: %s\n", trip.StartTime))
	sb.WriteString(fmt.Sprintf("По: %s\n", trip.EndTime))

	sb.WriteString("PlaceID поездки:\n")
	for i, place := range places {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, place.GooglePlace.PlaceID))
	}

	//sb.WriteString("Время работы:\n")
	//for i, place := range places {
	//	sb.WriteString(fmt.Sprintf("%d. %s - %s\n", i+1, place.Opening, place.Closing))
	//}

	sb.WriteString("Матрица времени и расстояния между местами:\n")
	//for i := range timeMatrix {
	//	for j := range timeMatrix[i] {
	//		sb.WriteString(fmt.Sprintf("%d ", timeMatrix[i][j]))
	//	}
	//	sb.WriteString("\n")
	//}

	for origin, destinations := range timeMatrix {
		for destination, metrics := range destinations {
			if origin == destination {
				continue
			}
			sb.WriteString(fmt.Sprintf("%s : %s = { distance: %.2f км, duration: %.2f мин }\n",
				origin, destination, metrics["distance"], metrics["duration"]))
		}
	}

	sb.WriteString(`
ФОРМАТ JSON ДОЛЖЕН СООТВЕТСТВОВАТЬ СЛЕДУЮЩЕЙ СТРУКТУРЕ: 
type Event struct {
	PlaceID string    
	TripID  uuid.UUID
	StartTime string 
	EndTime   string 
}

type Schedule struct {
	Events  []Event
}

В ОТВЕТЕ ВЕРНИ ТОЛЬКО JSON  объект, КОТОРЫЙ БУДЕТ ОПИСЫВАТЬ СПЛАНИРОВАННОЕ РАСПИСАНИЕ. БЕЗ ЛИШНИХ КОММЕНТАРИЕВ И БЕЗ ФОРМАТИРОВАНИЯ.
`)
	return sb.String(), nil
}

func ParseResponse(response string) (model.Schedule, error) {
	var schedule model.Schedule

	err := json.Unmarshal([]byte(response), &schedule)
	if err != nil {
		return model.Schedule{}, err
	}
	return schedule, nil
}
