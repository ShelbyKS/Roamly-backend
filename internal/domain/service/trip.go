package service

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/google/uuid"
)

type ITripService interface {
	GetTripByID(ctx context.Context, id uuid.UUID) (model.Trip, error)
	CreateTrip(ctx context.Context, trip model.Trip) (uuid.UUID, error)
	UpdateTrip(ctx context.Context, trip model.Trip) error
	GetTrips(ctx context.Context, userId int) ([]model.Trip, error)
	DeleteTrip(ctx context.Context, id uuid.UUID) error
}
