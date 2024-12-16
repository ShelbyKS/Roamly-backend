package orm

import (
	"github.com/google/uuid"
	"time"
)

type AIChatMessage struct {
	ID        int       `gorm:"primaryKey"`
	TripID    uuid.UUID `gorm:"not null;index"`
	Role      string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}
