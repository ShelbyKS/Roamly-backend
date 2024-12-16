package model

import "github.com/google/uuid"

type NotifyMessage struct {
	Payload struct {
		Action  string    `json:"action"`
		Author  string    `json:"author"`
		TripID  uuid.UUID `json:"trip_id"`
		Message string    `json:"message"`
	} `json:"payload"`
	Clients []string `json:"clients"`
}
