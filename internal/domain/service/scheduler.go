package service

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type ISchedulerService interface {
	GetSchedule(ctx context.Context, trip model.Trip, places []*model.Place, timeMatrix [][]int) (model.Schedule, error)
}
