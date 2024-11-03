package model

import (
	"github.com/google/uuid"
)

type Event struct {
	ID      uuid.UUID
	PlaceID string
	TripID  uuid.UUID
	Place   Place
	Trip    Trip

	StartTime string
	EndTime   string
}

type Schedule struct {
	Events  []Event
	Payload map[string]any
}

type DistanceMatrix map[string]map[string]map[string]float64
