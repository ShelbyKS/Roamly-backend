package dto

import (
	"github.com/google/uuid"
)

type TripResponse struct {
	ID                uuid.UUID     `json:"id"`
	Name              string        `json:"name"`
	Users             []GetUser     `json:"users"`
	StartTime         string        `json:"start_time"`
	EndTime           string        `json:"end_time"`
	AreaID            string        `json:"area_id"`
	Area              PlaceGoogle   `json:"area"`
	Places            []PlaceGoogle `json:"places"`
	RecommendedPlaces []PlaceGoogle `json:"recommended_places"`
	Events            []GetEvent    `json:"events"`
}
