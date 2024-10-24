package clients

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type IGoogleApiClient interface {
	FindPlace(ctx context.Context, input string, fields []string) ([]model.GooglePlace, error)
	GetPlaceByID(ctx context.Context, id string, fields []string) (model.GooglePlace, error)
}
