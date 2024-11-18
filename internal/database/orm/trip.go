package orm

import (
	"github.com/google/uuid"
)

type Trip struct {
	ID                uuid.UUID `gorm:"primaryKey"`
	Name              string
	StartTime         string `gorm:"type:TIME"`
	EndTime           string `gorm:"type:TIME"`
	AreaID            string
	Area              Place
	Users             []*User  `gorm:"many2many:trip_users;constraint:OnDelete:CASCADE;"`
	Places            []*Place `gorm:"many2many:trip_place;constraint:OnDelete:CASCADE;"`
	RecommendedPlaces []*Place `gorm:"many2many:trip_recommended_place;constraint:OnDelete:CASCADE;"`
	Events            []Event  `gorm:"constraint:OnDelete:CASCADE;"`
}

type TripUser struct {
	UserID   int       `gorm:"primaryKey"`
	TripID   uuid.UUID `gorm:"primaryKey"`
	UserRole int
}
