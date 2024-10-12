package orm

import "gorm.io/datatypes"

type Place struct {
	ID      int     `gorm:"primary_key"`
	Trips   *[]Trip `gorm:"many2many:trip_place"`
	Payload datatypes.JSONType[PlacePayload]
}

type PlacePayload struct {
	PlaceID string  `json:"place_id"`
	Name    string  `json:"name"`
	Photo   string  `json:"photo"`
	Rating  float32 `json:"rating"`
}
