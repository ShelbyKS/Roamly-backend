package googleapi

import (
	"context"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type GoogleClient interface {
	FindPlace(ctx context.Context, input string, fields []string) ([]model.GooglePlace, error)
	GetPlaceByID(ctx context.Context, id string, fields []string) (model.GooglePlace, error)
}

var DefaultClient GoogleClient

func Init(apiKey string) {
	DefaultClient = NewClient(apiKey)
}
