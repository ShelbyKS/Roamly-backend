package orm

import (
	"github.com/google/uuid"
)

type Event struct {
	ID      uuid.UUID `gorm:"primary_key"`
	PlaceID string
	TripID  uuid.UUID `gorm:"index"`
	Place   Place
	Trip    Trip

	StartTime string `gorm:"type:TIME"`
	EndTime   string `gorm:"type:TIME"`
}
