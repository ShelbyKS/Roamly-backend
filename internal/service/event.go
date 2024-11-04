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

type EventService struct {
	eventStorage storage.IEventStorage
	tripStorage  storage.ITripStorage
	placeStorage storage.IPlaceStorage
}

func NewEventService(eventStorage storage.IEventStorage, tripStorage storage.ITripStorage, placeStorage storage.IPlaceStorage) service.IEventService {
	return &EventService{
		eventStorage: eventStorage,
		tripStorage:  tripStorage,
		placeStorage: placeStorage,
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

	return event, nil
}

func (service *EventService) UpdateEvent(ctx context.Context, event model.Event) (model.Event, error) {
	updatedEvent, err := service.eventStorage.UpdateEvent(ctx, event)
	if errors.Is(err, domain.ErrEventNotFound) {
		return model.Event{}, err
	}
	if err != nil {
		return model.Event{}, fmt.Errorf("fail to update event in storage: %w", err)
	}

	updatedEvent, err = service.eventStorage.GetEventByID(ctx, updatedEvent.ID)
	if errors.Is(err, domain.ErrEventNotFound) {
		return model.Event{}, err
	}
	if err != nil {
		return model.Event{}, fmt.Errorf("fail to get event from storage: %w", err)
	}

	return updatedEvent, nil
}

func (service *EventService) DeleteEventsByTrip(ctx context.Context, tripID uuid.UUID) error {
	err := service.eventStorage.DeleteEventsByTrip(ctx, tripID)
	if err != nil {
		return fmt.Errorf("fail to delete events by trip ID: %w", err)
	}

	return nil
}
