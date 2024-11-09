package orm

import "time"

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Geometry struct {
	Location Location `gorm:"embedded"`
}

type Photo struct {
	PhotoReference string `json:"photo_reference"`
}

type GooglePlace struct {
	FormattedAddress string   `json:"formatted_address"`
	Geometry         Geometry `gorm:"embedded"`
	Name             string   `json:"name"`
	Photos           []Photo  `gorm:"json"`
	PlaceID          string   `json:"place_id"`
	Rating           float64  `json:"rating"`
	Types            []string `gorm:"json"`
}

type Place struct {
	ID    string  `gorm:"primary_key"`
	Trips []*Trip `gorm:"many2many:trip_place;constraint:OnDelete:CASCADE;"`
	// Closing time.Time
	// Opening time.Time
	// Name        string
	// Photo string
	// Rating      float32
	GooglePlace                 GooglePlace `gorm:"embedded"`
	RecommendedVisitingDuration time.Duration
}
