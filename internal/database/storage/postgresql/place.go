package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"

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
	//placeModel := orm.Place{}
	//
	//res := storage.db.WithContext(ctx).
	//	Where(&orm.Place{ID: placeID}).
	//	First(&placeModel)

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

	return model.Place{
		ID:     placeModel.ID,
		Name:   placeModel.Name,
		Photo:  placeModel.Photo,
		Rating: placeModel.Rating,
	}, nil
}

func (storage *PlaceStorage) CreatePlace(ctx context.Context, place *model.Place) (model.Place, error) {
	var trips []*orm.Trip
	for _, i := range place.Trips {
		trip := TripConverter{}.ToDb(*i)
		trips = append(trips, &trip)
	}

	placeModel := orm.Place{

		ID:     place.ID,
		Name:   place.Name,
		Rating: place.Rating,
		Photo:  place.Photo,
		Trips:  trips,
	}

	res := storage.db.WithContext(ctx).Create(&placeModel)
	if res.Error != nil {
		return model.Place{}, fmt.Errorf("create place in db: %w", res.Error)
	}

	return model.Place{
		ID:     placeModel.ID,
		Name:   placeModel.Name,
		Photo:  placeModel.Photo,
		Rating: placeModel.Rating,
	}, nil
}
