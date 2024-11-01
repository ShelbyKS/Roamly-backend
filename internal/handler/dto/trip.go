package dto

import (
	"github.com/google/uuid"
)

type TripResponse struct {
	ID        uuid.UUID     `json:"id"`
	Users     []GetUser     `json:"users"`
	StartTime string        `json:"start_time"`
	EndTime   string        `json:"end_time"`
	AreaID    string        `json:"area_id"`
	Area      GooglePlace   `json:"area"`
	Places    []GooglePlace `json:"places"`
	Events    []GetEvent    `json:"events"`
}
