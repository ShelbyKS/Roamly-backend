package clients

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type IGoogleApiClient interface {
	FindPlace(ctx context.Context, input string, fields []string) ([]model.Place, error)
}
