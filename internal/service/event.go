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

func (service *EventService) GetEventByID(ctx context.Context, placeID string, tripID uuid.UUID) (model.Event, error) {
	event, err := service.eventStorage.GetEventByID(ctx, placeID, tripID)
	if err != nil {
		return model.Event{}, fmt.Errorf("fail to get event from storage: %w", err)
	}

	return event, nil
}

func (service *EventService) DeleteEvent(ctx context.Context, placeID string, tripID uuid.UUID) error {
	err := service.eventStorage.DeleteEvent(ctx, placeID, tripID)
	if err != nil {
		return fmt.Errorf("fail to delete event from storage: %w", err)
	}

	return nil
}

func (service *EventService) CreateEvent(ctx context.Context, event model.Event) error {
	_, err := service.placeStorage.GetPlaceByID(ctx, event.PlaceID)
	if err != nil && !errors.Is(err, domain.ErrPlaceNotFound) {
		return fmt.Errorf("fail to get place from storage: %w", err)
	}

	_, err = service.tripStorage.GetTripByID(ctx, event.TripID)
	if err != nil && !errors.Is(err, domain.ErrTripNotFound) {
		return fmt.Errorf("fail to get trip from storage: %w", err)
	}

	// event.Place = place
	// event.Trip = trip

	err = service.eventStorage.CreateEvent(ctx, &event)
	if err != nil {
		return fmt.Errorf("fail to create event in storage: %w", err)
	}

	return nil
}

func (service *EventService) UpdateEvent(ctx context.Context, event model.Event) error {
	err := service.eventStorage.UpdateEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("fail to update event in storage: %w", err)
	}

	return nil
}
