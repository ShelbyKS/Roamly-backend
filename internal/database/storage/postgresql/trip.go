package postgresql

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
	"github.com/google/uuid"

	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
)

type TripStorage struct {
	db *gorm.DB
}

func NewTripStorage(db *gorm.DB) storage.ITripStorage {
	return &TripStorage{
		db: db,
	}
}

func (storage *TripStorage) GetTripByID(ctx context.Context, id uuid.UUID) (model.Trip, error) {
	trip := orm.Trip{
		ID: id,
	}

	tx := storage.db.WithContext(ctx).
		Model(&orm.Trip{}).
		Preload("Area").
		Preload("Users").
		Preload("Places").
		Preload("Events").
		First(&trip)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return model.Trip{}, domain.ErrTripNotFound
	}

	if tx.Error != nil {
		return model.Trip{}, tx.Error
	}

	return TripConverter{}.ToDomain(trip), tx.Error
}

func (storage *TripStorage) GetTrips(ctx context.Context) ([]model.Trip, error) {
	var tripsOrm []orm.Trip
	tx := storage.db.WithContext(ctx).Find(&tripsOrm)

	if tx.Error != nil {
		return []model.Trip{}, tx.Error
	}

	trips := make([]model.Trip, len(tripsOrm))
	for i, trip := range tripsOrm {
		trips[i] = TripConverter{}.ToDomain(trip)
	}

	return trips, tx.Error
}

func (storage *TripStorage) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	trip := orm.Trip{
		ID: id,
	}

	tx := storage.db.WithContext(ctx).Delete(&trip)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(domain.ErrUserNotFound, tx.Error)
	}

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (storage *TripStorage) CreateTrip(ctx context.Context, trip model.Trip) error {
	tripDb := TripConverter{}.ToDb(trip)
	tx := storage.db.WithContext(ctx).Create(&tripDb)

	return tx.Error
}

func (storage *TripStorage) UpdateTrip(ctx context.Context, trip model.Trip) error {
	tripDb := TripConverter{}.ToDb(trip)

	tx := storage.db.WithContext(ctx).
		Model(&orm.Trip{ID: trip.ID}).
		Updates(&tripDb)

	return tx.Error
}
