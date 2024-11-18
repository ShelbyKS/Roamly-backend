package orm

import (
	"database/sql"
	"github.com/google/uuid"
)

type Invite struct {
	Token  string    `gorm:"primaryKey"`
	TripID uuid.UUID `gorm:"index:idx_trip_id;not null"`
	Trip   Trip      `gorm:"constraint:OnDelete:CASCADE;"`
	Access string    `gorm:"not null;check:access IN ('reader', 'editor');index:idx_trip_id_access"`
	Enable sql.NullBool
}
