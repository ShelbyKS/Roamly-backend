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

type EventService struct {
	eventStorage    storage.IEventStorage
	tripStorage     storage.ITripStorage
	placeStorage    storage.IPlaceStorage
	sessionStorage  storage.ISessionStorage
	messageProducer clients.IMessageProdcuer
}

func NewEventService(eventStorage storage.IEventStorage,
	tripStorage storage.ITripStorage,
	placeStorage storage.IPlaceStorage,
	sessionStorage storage.ISessionStorage,
	messageProducer clients.IMessageProdcuer) service.IEventService {
	return &EventService{
		eventStorage:    eventStorage,
		tripStorage:     tripStorage,
		placeStorage:    placeStorage,
		sessionStorage:  sessionStorage,
		messageProducer: messageProducer,
	}
}

func (service *EventService) GetEventByID(ctx context.Context, eventID uuid.UUID) (model.Event, error) {
	event, err := service.eventStorage.GetEventByID(ctx, eventID)
	if errors.Is(err, domain.ErrEventNotFound) {
		return model.Event{}, err
	}
	if err != nil {
		return model.Event{}, fmt.Errorf("fail to get event from storage: %w", err)
	}

	return event, nil
}

func (service *EventService) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	err := service.eventStorage.DeleteEvent(ctx, eventID)
	if errors.Is(err, domain.ErrEventNotFound) {
		return err
	}
	if err != nil {
		return fmt.Errorf("fail to delete event from storage: %w", err)
	}

	tripFound, err := service.tripStorage.GetTripByEventID(ctx, eventID)

	users := tripFound.Users
	for _, user := range users {
		cooks, _ := service.sessionStorage.GetTokensByUserID(ctx, user.ID)
		// cookies = append(cookies, cooks...)
		var message model.Message
		message.Payload.Action = "trip_events_update"
		message.Payload.TripID = tripFound.ID
		message.Payload.Author = fmt.Sprintf("%d", user.ID)
		message.Payload.Message = "Из поездки удалено событие"
		message.Clients = cooks
		service.messageProducer.SendMessage(message)
	}

	return nil
}

func (service *EventService) CreateEvent(ctx context.Context, event model.Event) (model.Event, error) {
	//_, err := service.placeStorage.GetPlaceByID(ctx, event.PlaceID)
	//if err != nil && !errors.Is(err, domain.ErrPlaceNotFound) {
	//	return fmt.Errorf("fail to get place from storage: %w", err)
	//}
	//
	//_, err = service.tripStorage.GetTripByID(ctx, event.TripID)
	//if err != nil && !errors.Is(err, domain.ErrTripNotFound) {
	//	return fmt.Errorf("fail to get trip from storage: %w", err)
	//}
	//todo: check that event time is between trip date

	event.ID = uuid.New()

	err := service.eventStorage.CreateEvent(ctx, event)
	if err != nil {
		return model.Event{}, fmt.Errorf("fail to create event in storage: %w", err)
	}

	tripFound, _ := service.tripStorage.GetTripByID(ctx, event.TripID)
	users := tripFound.Users
	// var cookies []string
	for _, user := range users {
		cooks, _ := service.sessionStorage.GetTokensByUserID(ctx, user.ID)
		// cookies = append(cookies, cooks...)
		var message model.Message
		message.Payload.Action = "trip_events_update"
		message.Payload.TripID = event.TripID
		message.Payload.Author = fmt.Sprintf("%d", user.ID)
		message.Payload.Message = "В поездке создано новое событие"
		message.Clients = cooks
		service.messageProducer.SendMessage(message)
	}

	return event, nil
}

func (service *EventService) UpdateEvent(ctx context.Context, event model.Event) (model.Event, error) {
	log.Println("START_UPDATING_EVENT: ")
	updatedEvent, err := service.eventStorage.UpdateEvent(ctx, event)
	if errors.Is(err, domain.ErrEventNotFound) {
		log.Println("START_UPDATING_EVENT: NOT FOUND")
		return model.Event{}, err
	}
	if err != nil {
		log.Println("START_UPDATING_EVENT: ERR:", err)
		return model.Event{}, fmt.Errorf("fail to update event in storage: %w", err)
	}
	

	updatedEvent, err = service.eventStorage.GetEventByID(ctx, updatedEvent.ID)
	if errors.Is(err, domain.ErrEventNotFound) {
		log.Println("START_UPDATING_EVENT: NOT FOUND 2x")
		return model.Event{}, err
	}
	if err != nil {
		log.Println("START_UPDATING_EVENT: ERR 2x", err)
		return model.Event{}, fmt.Errorf("fail to get event from storage: %w", err)
	}

	tripFound, err := service.tripStorage.GetTripByEventID(ctx, event.ID)
	log.Println("TRIP_FOUND: ", tripFound, "error:", err, "TRIP_ID: ", tripFound.ID)
	users := tripFound.Users
	log.Println("USERS: ", tripFound.Users)
	// var cookies []string
	for _, user := range users {
		cooks, _ := service.sessionStorage.GetTokensByUserID(ctx, user.ID)
		// cookies = append(cookies, cooks...)
		var message model.Message
		message.Payload.Action = "trip_events_update"
		message.Payload.TripID = tripFound.ID
		message.Payload.Author = fmt.Sprintf("%d", user.ID)
		message.Payload.Message = "Событие поездки обновлено"
		message.Clients = cooks
		service.messageProducer.SendMessage(message)
	}

	return updatedEvent, nil
}

func (service *EventService) DeleteEventsByTrip(ctx context.Context, tripID uuid.UUID) error {
	err := service.eventStorage.DeleteEventsByTrip(ctx, tripID)
	if err != nil {
		return fmt.Errorf("fail to delete events by trip ID: %w", err)
	}

	tripFound, _ := service.tripStorage.GetTripByID(ctx, tripID)
	users := tripFound.Users
	// var cookies []string
	for _, user := range users {
		cooks, _ := service.sessionStorage.GetTokensByUserID(ctx, user.ID)
		// cookies = append(cookies, cooks...)
		var message model.Message
		message.Payload.Action = "trip_events_update"
		message.Payload.TripID = tripID
		message.Payload.Author = fmt.Sprintf("%d", user.ID)
		message.Payload.Message = "Из поездки удалено событие"
		message.Clients = cooks
		service.messageProducer.SendMessage(message)
	}

	return nil
}
