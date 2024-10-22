package orm

import (
	"github.com/google/uuid"
)

type Trip struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	StartTime string    `gorm:"type:DATE"`
	EndTime   string    `gorm:"type:DATE"`
	AreaID    string
	Users     []*User  `gorm:"many2many:trip_user;constraint:OnDelete:CASCADE;"`
	Places    []*Place `gorm:"many2many:trip_place;constraint:OnDelete:CASCADE;"`
}
