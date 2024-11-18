package dto

import "github.com/google/uuid"

type InviteResponse struct {
	Token  string    `json:"token"`
	TripID uuid.UUID `json:"trip_id"`
	Access string    `json:"access"`
	Enable bool      `json:"enable"`
}
