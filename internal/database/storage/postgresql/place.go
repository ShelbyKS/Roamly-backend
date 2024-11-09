package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/google/uuid"

	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
	"gorm.io/gorm"
)

type PlaceStorage struct {
	db *gorm.DB
}

func NewPlaceStorage(db *gorm.DB) storage.IPlaceStorage {
	return &PlaceStorage{
		db: db,
	}
}

func (storage *PlaceStorage) GetPlaceByID(ctx context.Context, placeID string) (model.Place, error) {
	placeModel := &orm.Place{
		ID: placeID,
	}
	res := storage.db.WithContext(ctx).
		First(placeModel)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.Place{}, domain.ErrPlaceNotFound
	}

	if res.Error != nil {
		return model.Place{}, res.Error
	}

	place := PlaceConverter{}.ToDomain(*placeModel)

	return place, nil
}

func (storage *PlaceStorage) CreatePlace(ctx context.Context, place *model.Place) (model.Place, error) {
	placeModel := PlaceConverter{}.ToDb(*place)
	// log.Println("in storage", placeModel)

	res := storage.db.WithContext(ctx).Create(&placeModel)
	if res.Error != nil {
		return model.Place{}, fmt.Errorf("create place in db: %w", res.Error)
	}

	return *place, nil
}

func (storage *PlaceStorage) AppendPlaceToTrip(ctx context.Context, placeID string, tripID uuid.UUID) error {
	err := storage.db.
		Model(&orm.Trip{
			ID: tripID,
		}).
		Association("Places").
		Append(&orm.Place{
			ID: placeID,
		})
	if err != nil {
		return fmt.Errorf("failed to add place to trip: %v", err)
	}

	return nil
}

func (storage *PlaceStorage) DeletePlace(ctx context.Context, tripID uuid.UUID, placeID string) error {
	trip := &orm.Trip{
		ID: tripID,
	}

	place := &orm.Place{
		ID: placeID,
	}

	if err := storage.db.WithContext(ctx).Model(&trip).Association("Places").Delete(place); err != nil {
		return err
	}

	return nil
}

func (storage *PlaceStorage) AppendPlace(ctx context.Context, tripID uuid.UUID, placeID string) error {
	trip := &orm.Trip{
		ID: tripID,
	}

	place := &orm.Place{
		ID: placeID,
	}

	if err := storage.db.WithContext(ctx).Model(&trip).Association("Places").Delete(place); err != nil {
		return err
	}

	return nil
}

func (storage *PlaceStorage) UpdatePlace(ctx context.Context, place model.Place) error {
	placeDB := PlaceConverter{}.ToDb(place)

	tx := storage.db.WithContext(ctx).
		Model(&orm.Place{ID: place.ID}).
		Updates(&placeDB)

	return tx.Error
}
