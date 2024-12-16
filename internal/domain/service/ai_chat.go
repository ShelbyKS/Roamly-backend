package service

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/google/uuid"
)

type IAIChatService interface {
	SentMessage(ctx context.Context, message model.ChatMessage, userID int) error
	GetAIChatMessages(ctx context.Context, tripID uuid.UUID) ([]model.ChatMessage, error)
}
