package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShelbyKS/Roamly-backend/internal/domain"

	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
	"gorm.io/datatypes"

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
	placeModel := &orm.Place{}
	res := storage.db.WithContext(ctx).
		Where("payload ->> 'place_id' = ?", placeID).
		First(placeModel)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.Place{}, domain.ErrPlaceNotFound
	}

	if res.Error != nil {
		return model.Place{}, res.Error
	}

	return model.Place{
		ID:     placeModel.Payload.Data().PlaceID,
		Name:   placeModel.Payload.Data().Name,
		Rating: placeModel.Payload.Data().Rating,
	}, nil
}

func (storage *PlaceStorage) CreatePlace(ctx context.Context, place model.Place) error {
	var trips []orm.Trip
	for _, i := range place.Trips {
		trips = append(trips, TripConverter{}.ToDb(*i))
	}

	placeModel := orm.Place{
		Payload: datatypes.NewJSONType(
			orm.PlacePayload{
				PlaceID: place.ID,
				Name:    place.Name,
				Rating:  place.Rating,
			},
		),
		Trips: &trips,
	}

	res := storage.db.WithContext(ctx).Create(&placeModel)
	if res.Error != nil {
		return fmt.Errorf("create phantom in db: %w", res.Error)
	}

	return nil
}
