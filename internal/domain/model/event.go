package model

import (
	"github.com/google/uuid"
)

type Event struct {
	ID      uuid.UUID
	Name    string
	PlaceID string
	TripID  uuid.UUID

	StartTime string
	EndTime   string
}

type DistanceMatrix map[string]map[string]map[string]float64
