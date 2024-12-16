package model

import "github.com/google/uuid"

type Trip struct {
	ID                uuid.UUID     `json:"id"`
	Name              string        `json:"name"`
	StartTime         string        `json:"start_time"`
	EndTime           string        `json:"end_time"`
	AreaID            string        `json:"area_id"`
	Area              *Place        `json:"area"`
	Users             []*User       `json:"users"`
	Places            []*Place      `json:"places"`
	RecommendedPlaces []*Place      `json:"recommended_places"`
	Events            []Event       `json:"events"`
	AIChat            []ChatMessage `json:"ai_chat"`
}

func (trip *Trip) GetTripPlaceIDs() []string {
	var placeIDs []string
	for _, place := range trip.Places {
		placeIDs = append(placeIDs, place.ID)
	}
	return placeIDs
}

func (trip *Trip) GetTripRecommendedPlaceIDs() []string {
	var placeIDs []string
	for _, place := range trip.RecommendedPlaces {
		placeIDs = append(placeIDs, place.ID)
	}
	return placeIDs
}

func (trip *Trip) GetTopRecommendations() ([]string, []*Place) {
	size := 10
	if len(trip.RecommendedPlaces) < size {
		size = len(trip.RecommendedPlaces)
	}

	placeIDs := make([]string, size)
	places := make([]*Place, size)

	for i := 0; i < size; i++ {
		placeIDs[i] = trip.RecommendedPlaces[i].ID
		places[i] = trip.RecommendedPlaces[i]
	}

	return placeIDs, places
}
