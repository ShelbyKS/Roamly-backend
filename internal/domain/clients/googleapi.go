package clients

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type IGoogleApiClient interface {
	FindPlace(ctx context.Context, input string, fields []string) ([]model.GooglePlace, error)
	GetPlaceByID(ctx context.Context, id string, fields []string) (model.GooglePlace, error)
	GetTimeDistanceMatrix(ctx context.Context, placeIDs []string) (model.DistanceMatrix, error)
	GetNewPlaces(ctx context.Context,
		textQuery string,
		includedType string,
		pageSize int,
		lat float64,
		lng float64,
		radius float64,
		languageCode string,
		pageToken string) ([]model.GooglePlace, error)
	GetPlacesNearby(ctx context.Context,
		includedTypes []string,
		maxPlaces int,
		rankPrefernce string,
		lat float64,
		lng float64,
		radius float64,
		languageCode string) ([]model.GooglePlace, error)
}
