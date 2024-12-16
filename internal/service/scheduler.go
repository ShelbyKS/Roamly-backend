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
	openAIClient    clients.IChatClient
	googleApi       clients.IGoogleApiClient
	tripStorage     storage.ITripStorage
	eventStorage    storage.IEventStorage
	placeStorage    storage.IPlaceStorage
	sessionStorage  storage.ISessionStorage
	messageProducer clients.IMessageProdcuer
}

func NewShedulerService(
	openAIClient clients.IChatClient,
	googleApi clients.IGoogleApiClient,
	tripStorage storage.ITripStorage,
	eventStorage storage.IEventStorage,
	placeStorage storage.IPlaceStorage,
	sessionStorage storage.ISessionStorage,
	messageProducer clients.IMessageProdcuer,
) service.ISchedulerService {
	return &SchedulerService{
		openAIClient:    openAIClient,
		googleApi:       googleApi,
		tripStorage:     tripStorage,
		eventStorage:    eventStorage,
		placeStorage:    placeStorage,
		sessionStorage:  sessionStorage,
		messageProducer: messageProducer,
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

	resp, err := s.openAIClient.PostPrompt(ctx, []model.ChatMessage{{
		Role:    model.RoleUser,
		Content: prompt,
	}}, clients.ModelChatGPT4o)
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

	users := trip.Users
	// var cookies []string
	for _, user := range users {
		cooks, _ := s.sessionStorage.GetTokensByUserID(ctx, user.ID)
		// cookies = append(cookies, cooks...)
		var message model.NotifyMessage
		message.Payload.Action = "trip_events_update"
		message.Payload.TripID = trip.ID
		message.Payload.Author = fmt.Sprintf("%d", user.ID)
		message.Payload.Message = "Поездка спланирована"
		message.Clients = cooks
		s.messageProducer.SendMessage(message)
	}

	return trip, nil
}

func (s *SchedulerService) AutoScheduleTrip(ctx context.Context, tripID uuid.UUID) (model.Trip, error) {
	trip, err := s.tripStorage.GetTripByID(ctx, tripID)
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to get trip for schedule: %w", err)
	}

	timeDistMatrix, err := s.googleApi.GetTimeDistanceMatrix(ctx, trip.GetTripRecommendedPlaceIDs())
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to get time distance matrix: %w", err)
	}

	prompt := s.generateRequestString(trip, trip.RecommendedPlaces, timeDistMatrix)

	resp, err := s.openAIClient.PostPrompt(ctx, []model.ChatMessage{{
		Role:    model.RoleUser,
		Content: prompt,
	}}, clients.ModelChatGPT4o)
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to get openai response: %w", err)
	}

	events, err := ParseSchedule(resp)
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to parse schedule: %w", err)
	}

	//todo: batch
	for _, place := range trip.RecommendedPlaces {
		err = s.placeStorage.AppendPlaceToTrip(ctx, place.ID, trip.ID)
		if err != nil {
			return model.Trip{}, fmt.Errorf("failed to append place: %w", err)
		}
	}
	trip.Places = trip.RecommendedPlaces

	err = s.eventStorage.DeleteEventsByTrip(ctx, trip.ID)
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to delete current events: %w", err)
	}

	err = s.eventStorage.CreateBatchEvents(ctx, &events)
	if err != nil {
		return model.Trip{}, fmt.Errorf("failed to save events: %w", err)
	}
	trip.Events = events

	users := trip.Users
	for _, user := range users {
		cooks, _ := s.sessionStorage.GetTokensByUserID(ctx, user.ID)
		var message model.NotifyMessage
		message.Payload.Action = "trip_events_update"
		message.Payload.TripID = trip.ID
		message.Payload.Author = fmt.Sprintf("%d", user.ID)
		message.Payload.Message = "Поездка спланирована"
		message.Clients = cooks
		s.messageProducer.SendMessage(message)
	}

	return trip, nil
}

func (s *SchedulerService) generateRequestString(trip model.Trip, places []*model.Place, timeMatrix model.DistanceMatrix) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("TripID: %s\n", trip.ID.String()))

	sb.WriteString("\nДата поездки:\n")
	sb.WriteString(fmt.Sprintf("С: %s\n", trip.StartTime))
	sb.WriteString(fmt.Sprintf("По: %s\n", trip.EndTime))

	sb.WriteString("\nМеста поездки - Название:PlaceID:Время на посещение в минутах\n")
	for i, place := range places {
		sb.WriteString(fmt.Sprintf("%d. %s:%s:%d\n", i+1,
			place.GooglePlace.Name,
			place.GooglePlace.PlaceID,
			place.RecommendedVisitingDuration,
		))
	}
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
		"НЕ ПЛАНИРУЙ ПОСЕЩЕНИЕ МЕСТ РАНЕЕ 9 УТРА И НЕ СТАВЬ БОЛЬШЕ ТРЁХ МЕСТ В ДЕНЬ\n" +
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
