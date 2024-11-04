package dto

import "github.com/google/uuid"

type GetEvent struct {
	ID        uuid.UUID `json:"id"`
	PlaceID   string    `json:"place_id"`
	TripID    uuid.UUID `json:"trip_id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
}
