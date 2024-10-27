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

func (storage *TripStorage) GetTrips(ctx context.Context, userId int) ([]model.Trip, error) {
	var user orm.User

    err := storage.db.
		Preload("Trips").
		Preload("Trips.Area").
		Preload("Trips.Users").
		First(&user, userId).Error
	if err != nil {
		return []model.Trip{}, err
	}

	trips := make([]model.Trip, len(user.Trips))
	for i, trip := range user.Trips {
		trips[i] = TripConverter{}.ToDomain(*trip)
	}

	return trips, err
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
