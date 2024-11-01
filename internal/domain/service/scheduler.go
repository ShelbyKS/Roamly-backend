package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type ISchedulerService interface {
	ScheduleTrip(ctx context.Context, tripID uuid.UUID) (model.Trip, error)
}
