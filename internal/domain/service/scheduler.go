package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type ISchedulerService interface {
	GetSchedule(ctx context.Context, tripID uuid.UUID) (model.Schedule, error)
}
