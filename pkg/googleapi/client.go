package googleapi

import (
	"context"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type GoogleClient interface {
	FindPlace(ctx context.Context, input string, fields []string) ([]model.Place, error)
}

var DefaultClient GoogleClient

func Init(apiKey string) {
	DefaultClient = NewClient(apiKey)
}
