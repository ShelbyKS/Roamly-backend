package storage

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/google/uuid"
)

type IAIChatStorage interface {
	SaveAIChatMessage(ctx context.Context, message model.ChatMessage) error
	GetMessagesByTripID(ctx context.Context, tripID uuid.UUID) ([]model.ChatMessage, error)
}
