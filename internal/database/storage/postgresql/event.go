package postgresql

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

func (storage *EventStorage) GetEventByID(ctx context.Context, eventID uuid.UUID) (model.Event, error) {
	event := orm.Event{
		ID: eventID,
	}

	tx := storage.db.WithContext(ctx).
		//Preload("Place").
		//Preload("Trip").
		First(&event)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return model.Event{}, domain.ErrEventNotFound
	}

	if tx.Error != nil {
		return model.Event{}, tx.Error
	}

	return EventConverter{}.ToDomain(event), nil
}

func (storage *EventStorage) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	event := orm.Event{
		ID: eventID,
	}

	tx := storage.db.WithContext(ctx).Delete(&event)

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (storage *EventStorage) CreateEvent(ctx context.Context, event model.Event) error {
	eventDb := EventConverter{}.ToDb(event)

	tx := storage.db.WithContext(ctx).Create(&eventDb)

	return tx.Error
}

func (storage *EventStorage) UpdateEvent(ctx context.Context, event model.Event) (model.Event, error) {
	eventDb := EventConverter{}.ToDb(event)
	tx := storage.db.WithContext(ctx).
		Model(&orm.Event{ID: eventDb.ID}).
		Updates(eventDb)

	if tx.Error != nil {
		return model.Event{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return model.Event{}, domain.ErrEventNotFound
	}

	return EventConverter{}.ToDomain(eventDb), nil
}

func (storage *EventStorage) CreateBatchEvents(ctx context.Context, events *[]model.Event) error {
	var eventsDb []orm.Event
	for i, event := range *events {
		id := uuid.New()
		event.ID = id
		(*events)[i].ID = id
		eventsDb = append(eventsDb, EventConverter{}.ToDb(event))
	}

	tx := storage.db.WithContext(ctx).Create(&eventsDb)

	return tx.Error
}

func (storage *EventStorage) DeleteEventsByTrip(ctx context.Context, tripID uuid.UUID) error {
	tx := storage.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Delete(&orm.Event{})

	return tx.Error
}
