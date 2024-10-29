package model

import "github.com/google/uuid"

type Trip struct {
	//todo: add trip name
	ID        uuid.UUID `json:"id"`
	Users     []*User   `json:"users"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	AreaID    string    `json:"area_id"`
	Area      *Place    `json:"area"`
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
