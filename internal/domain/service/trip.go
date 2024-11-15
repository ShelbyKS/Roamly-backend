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
	GetUserRole(ctx context.Context, userID int, tripID uuid.UUID) (model.UserTripRole, error)
	GetTripByEventID(ctx context.Context, eventID uuid.UUID) (model.Trip, error)
	DetermineRecommendedPlaces(ctx context.Context, tripID uuid.UUID) error
}
