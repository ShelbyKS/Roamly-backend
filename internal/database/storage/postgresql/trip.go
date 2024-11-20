package postgresql

import (
	"context"
	"errors"
	"log"

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
		Preload("TripUsers").
		Preload("Places").
		Preload("RecommendedPlaces").
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

	err := storage.db.WithContext(ctx).
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

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return domain.ErrTripNotFound
	}

	return nil
}

func (storage *TripStorage) CreateTrip(ctx context.Context, trip model.Trip, userRole model.UserTripRole) error {
	tripDb := TripConverter{}.ToDb(trip)

	tripUser := orm.TripUsers{
		UserID:   trip.Users[0].ID,
		TripID:   tripDb.ID,
		UserRole: int(userRole),
	}

	tripDb.Users = []*orm.User{}

	tx := storage.db.WithContext(ctx).Begin()
	if err := tx.Create(&tripDb).Error; err != nil {
		tx.Rollback()
		log.Println("create trip err:", err)
		return err
	}

	if err := tx.Create(&tripUser).Error; err != nil {
		tx.Rollback()
		log.Println("create tripUser err:", err)
		return err
	}

	return tx.Commit().Error
}

func (storage *TripStorage) UpdateTrip(ctx context.Context, trip model.Trip) error {
	tripDb := TripConverter{}.ToDb(trip)

	tx := storage.db.WithContext(ctx).
		Model(&orm.Trip{ID: trip.ID}).
		Updates(&tripDb)

	if tx.RowsAffected == 0 {
		return domain.ErrTripNotFound
	}

	return tx.Error
}

func (storage *TripStorage) GetUserRole(ctx context.Context, userID int, tripID uuid.UUID) (model.UserTripRole, error) {
	tripUser := orm.TripUsers{
		UserID: userID,
		TripID: tripID,
	}

	err := storage.db.WithContext(ctx).
		First(&tripUser).Error

	if err != nil {
		return 0, err
	}

	return model.UserTripRole(tripUser.UserRole), nil
}

func (storage *TripStorage) GetTripByEventID(ctx context.Context, eventID uuid.UUID) (model.Trip, error) {
	event := orm.Event{
		ID: eventID,
	}

	err := storage.db.
		WithContext(ctx).
		Preload("Trip").
		Preload("Trip.Users").
		First(&event).Error

	if err != nil {
		return model.Trip{}, err
	}

	return TripConverter{}.ToDomain(event.Trip), nil
}
