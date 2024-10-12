package service

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type ITripService interface {
	GetTripByID(ctx context.Context, id int) (model.Trip, error)
	CreateTrip(ctx context.Context, trip model.Trip) error
	UpdateTrip(ctx context.Context, trip model.Trip) error
	DeleteTrip(ctx context.Context, id int) error
}
