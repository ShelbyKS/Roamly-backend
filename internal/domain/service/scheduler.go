package service

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type ISchedulerService interface {
	GetSchedule(ctx context.Context, places []model.Place) (model.Schedule, error)
}
