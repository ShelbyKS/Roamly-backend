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
	eventStorage storage.IEventStorage
}

func NewShedulerService(
	openAIClient clients.IChatClient,
	googleApi clients.IGoogleApiClient,
	tripStorage storage.ITripStorage,
	eventStorage storage.IEventStorage,
) service.ISchedulerService {
	return &SchedulerService{
		openAIClient: openAIClient,
		googleApi:    googleApi,
		tripStorage:  tripStorage,
		eventStorage: eventStorage,
	}
}

func (s *SchedulerService) ScheduleTrip(ctx context.Context, tripID uuid.UUID) (model.Trip, error) {
	trip, err := s.tripStorage.GetTripByID(ctx, tripID)
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to get trip for schedule: %w", err)
	}

	timeDistMatrix, err := s.googleApi.GetTimeDistanceMatrix(ctx, trip.GetTripPlaceIDs())
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to get time distance matrix: %w", err)
	}

	prompt := s.generateRequestString(trip, trip.Places, timeDistMatrix)

	//fmt.Println("PROMT: ", prompt)

	resp, err := s.openAIClient.PostPrompt(ctx, prompt)
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to get openai response: %w", err)
	}

	events, err := ParseSchedule(resp)
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to parse schedule: %w", err)
	}

	err = s.eventStorage.DeleteEventsByTrip(ctx, trip.ID)
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to delete current events: %w", err)
	}

	err = s.eventStorage.CreateBatchEvents(ctx, &events)
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to save events: %w", err)
	}
	trip.Events = events

	return trip, nil
}

func (s *SchedulerService) generateRequestString(trip model.Trip, places []*model.Place, timeMatrix model.DistanceMatrix) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("TripID: %s\n", trip.ID.String()))

	sb.WriteString("\nДата поездки:\n")
	sb.WriteString(fmt.Sprintf("С: %s\n", trip.StartTime))
	sb.WriteString(fmt.Sprintf("По: %s\n", trip.EndTime))

	sb.WriteString("\nМеста поездки - Название (PlaceID):\n")
	for i, place := range places {
		sb.WriteString(fmt.Sprintf("%d. %s (%s)\n", i+1, place.GooglePlace.Name, place.GooglePlace.PlaceID))
	}
	sb.WriteString("\nСКОЛЬКО ВРЕМЕНИ ЗАЙМЁТ ПОСЕЩЕНИЕ ЭТИХ МЕСТ? ВЫДЕЛИ СРЕДНЕЕ ВРЕМЯ\n")
	//sb.WriteString("Время работы:\n")
	//for i, place := range places {
	//	sb.WriteString(fmt.Sprintf("%d. %s - %s\n", i+1, place.Opening, place.Closing))
	//}

	sb.WriteString("\nМатрица времени и расстояния между местами:\n")
	for origin, destinations := range timeMatrix {
		for destination, metrics := range destinations {
			if origin == destination {
				continue
			}
			sb.WriteString(fmt.Sprintf("%s : %s = { distance: %.2f км, duration: %.2f мин }\n",
				origin, destination, metrics["distance"], metrics["duration"]))
		}
	}

	sb.WriteString("\nФОРМАТ JSON ДОЛЖЕН СООТВЕТСТВОВАТЬ СЛЕДУЮЩЕЙ СТРУКТУРЕ: \n" +
		"type Event struct {\n" +
		"    PlaceID string\n" +
		"    TripID uuid.UUID\n" +
		"    StartTime string\n" +
		"    EndTime string\n" +
		"}\n\n" +
		"НУЖНО ВЕРНУТЬ []Event (массив Event)\n" +
		//"В ОТВЕТЕ ВЕРНИ ТОЛЬКО СПЛАНИРОВАННОЕ РАСПИСАНИЕ. БЕЗ ЛИШНИХ КОММЕНТАРИЕВ И БЕЗ ФОРМАТИРОВАНИЯ.\n" +
		"ОКРУГЛЯЙ ВРЕМЯ НАЧАЛА СОБЫТИЯ И КОНЦА ДО ЦЕЛЫХ ЧАСА ИЛИ ПОЛОВИНЫ, ДАВАЯ ЗАПАС НА ПЕРЕМЕЩЕНИЕ МЕЖДУ ОБЪЕКТАМИ.\n" +
		"НЕ ПЛАНИРУЙ ПОСЕЩЕНИЕ МЕСТ РАНЕЕ 9 УТРА\n" +
		"ТАКЖЕ В РАСПИСАНИИ НУЖНО УЧИТЫВАТЬ ВРЕМЯ НА ПРИЁМЫ ПИЩИ И ПОХОДЫ В ТУАЛЕТ.\n" +
		"НУЖНО РАСПРЕДЕЛЯТЬ РАВНОМЕРНО ПОСЕЩЕНИЕ МЕСТ ПО ДАТАМ ПОЕЗДКИ. " +
		"ОДНОМ МЕСТО МОЖНО ПОСЕТИТЬ ТОЛЬКО 1 РАЗ ЗА ПОЕЗДКУ.\n" +
		"БЕЗ ЛИШНИХ КОММЕНТАРИЕВ И БЕЗ ФОРМАТИРОВАНИЯ ПО ТИПУ \\`\\`\\`json\\`\\`\\`.\n")
	return sb.String()
}

func ParseSchedule(response string) ([]model.Event, error) {
	var events []model.Event

	err := json.Unmarshal([]byte(response), &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}
