package orm

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Event struct {
	PlaceID string    `gorm:"primaryKey"`
	TripID  uuid.UUID `gorm:"primaryKey"`
	Place   Place
	Trip    Trip

	StartTime string `gorm:"type:TIME"`
	EndTime   string `gorm:"type:TIME"`
	Payload   datatypes.JSON
}
