package model

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID        int
	TripID    uuid.UUID
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt time.Time
}

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)
