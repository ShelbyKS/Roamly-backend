package orm

import (
	"github.com/google/uuid"
)

type Trip struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	StartTime string    `gorm:"type:TIME"`
	EndTime   string    `gorm:"type:TIME"`
	AreaID    string
	Events    []Event
	Users     []*User  `gorm:"many2many:trip_user;constraint:OnDelete:CASCADE;"`
	Places    []*Place `gorm:"many2many:trip_place;constraint:OnDelete:CASCADE;"`
}
