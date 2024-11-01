package model

import "github.com/google/uuid"

type Trip struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	AreaID    string    `json:"area_id"`
	Area      *Place    `json:"area"`
	Users     []*User   `json:"users"`
	Places    []*Place  `json:"places"`
	Events    []Event   `json:"events"`
}

func (trip *Trip) GetTripPlaceIDs() []string {
	var placeIDs []string
	for _, place := range trip.Places {
		placeIDs = append(placeIDs, place.ID)
	}
	return placeIDs
}
