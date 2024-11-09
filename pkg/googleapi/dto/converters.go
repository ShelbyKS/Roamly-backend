package dto

import "github.com/ShelbyKS/Roamly-backend/internal/domain/model"

type PlaceConverter struct{}

func (PlaceConverter) ToDomain(place GooglePlace) model.GooglePlace {
	return model.GooglePlace{
		FormattedAddress: place.FormattedAddress,
		Name:             place.DisplayName.Text,
		Rating:           place.Rating,
		Geometry: model.Geometry{
			Location: model.Location{
				Lat: place.Location.Latitude,
				Lng: place.Location.Longitude,
			},
		},
		Photos: PhotoConverter{}.ConvertPhotos(place.Photos),
	}
}

type PhotoConverter struct{}

func (PhotoConverter) ConvertPhotos(oldPhotos []Photo) []model.Photo {
	newPhotos := make([]model.Photo, len(oldPhotos))
	for i, oldPhoto := range oldPhotos {
		newPhotos[i] = model.Photo{PhotoReference: oldPhoto.Name}
	}
	return newPhotos
}
