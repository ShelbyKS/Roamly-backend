package model

import "time"

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Geometry struct {
	Location Location `json:"location"`
}

type Photo struct {
	PhotoReference string `json:"photo_reference"`
}

type GooglePlace struct {
	FormattedAddress string   `json:"formatted_address"`
	Geometry         Geometry `json:"geometry"`
	Name             string   `json:"name"`
	Photos           []Photo  `json:"photos"`
	PlaceID          string   `json:"place_id"`
	Rating           float64  `json:"rating"`
	Types            []string `json:"types"`
	EditorialSummary string   `json:"editorial_summary"`
	Vicinity         string   `json:"vicinity"`
}

type Place struct {
	ID                          string `json:"id"`
	Trips                       []*Trip
	Closing                     time.Time
	Opening                     time.Time
	GooglePlace                 GooglePlace
	RecommendedVisitingDuration int
}
