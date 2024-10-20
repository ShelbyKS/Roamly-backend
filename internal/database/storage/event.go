package storage

import (
	"context"
	"errors"

	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type EventStorage struct {
	db *gorm.DB
}

func NewEventStorage(db *gorm.DB) storage.IEventStorage {
	return &EventStorage{
		db: db,
	}
}

func (storage *EventStorage) GetEventByID(ctx context.Context, placeID string, tripID uuid.UUID) (model.Event, error) {
	event := orm.Event{
		PlaceID: placeID,
		TripID:  tripID,
	}

	tx := storage.db.WithContext(ctx).
		Preload("Place").
		Preload("Trip").
		First(&event)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(domain.ErrEventNotFound, tx.Error)
	}

	if tx.Error != nil {
		return model.Event{}, tx.Error
	}

	return EventConverter{}.ToDomain(event), nil
}

func (storage *EventStorage) DeleteEvent(ctx context.Context, placeID string, tripID uuid.UUID) error {
	event := orm.Event{
		PlaceID: placeID,
		TripID:  tripID,
	}

	tx := storage.db.WithContext(ctx).Delete(&event)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(domain.ErrEventNotFound, tx.Error)
	}

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (storage *EventStorage) CreateEvent(ctx context.Context, event model.Event) error {
	eventDb := EventConverter{}.ToDb(event)
	eventDb.Trip = orm.Trip{}
	tx := storage.db.WithContext(ctx).Create(&eventDb)
	return tx.Error
}

func (storage *EventStorage) UpdateEvent(ctx context.Context, event model.Event) error {
	eventDb := EventConverter{}.ToDb(event)
	tx := storage.db.WithContext(ctx).
		Model(&orm.Event{PlaceID: event.PlaceID, TripID: event.TripID}).
		Updates(eventDb)

	return tx.Error
}
