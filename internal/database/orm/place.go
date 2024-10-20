package orm

import "time"

type Place struct {
	ID      string  `gorm:"primary_key"`
	Trips   []*Trip `gorm:"many2many:trip_place;constraint:OnDelete:CASCADE;"`
	Closing time.Time
	Opening time.Time
	Name    string
	Photo   string
	Rating  float32
}
