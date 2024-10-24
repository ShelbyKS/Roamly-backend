package model

import "github.com/google/uuid"

type Trip struct {
	ID        uuid.UUID `json:"id"`
	Users     []*User   `json:"users"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	AreaID    string    `json:"area_id"`
	Area      *Place    `json:"area"`
	Places    []*Place  `json:"places"`
	Events    []Event   `json:"events"`
}
