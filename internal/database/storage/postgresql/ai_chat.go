package postgresql

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type AIChatStorage struct {
	db *gorm.DB
}

func NewAIChatStorage(db *gorm.DB) storage.IAIChatStorage {
	return &AIChatStorage{
		db: db,
	}
}

func (storage *AIChatStorage) SaveAIChatMessage(ctx context.Context, message model.ChatMessage) error {
	messageDB := orm.AIChatMessage{
		TripID:    message.TripID,
		Role:      message.Role,
		Content:   message.Content,
		CreatedAt: time.Now(),
	}

	err := storage.db.WithContext(ctx).Create(&messageDB).Error
	if err != nil {
		log.Fatalf("failed to save ai chat message to trip: %v", err)
	}

	return nil
}

func (storage *AIChatStorage) GetMessagesByTripID(ctx context.Context, tripID uuid.UUID) ([]model.ChatMessage, error) {
	var messagesDB []orm.AIChatMessage

	err := storage.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("created_at ASC").
		Find(&messagesDB).Error
	if err != nil {
		return nil, err
	}

	var messagesDomain []model.ChatMessage
	for _, msg := range messagesDB {
		messagesDomain = append(messagesDomain, ChatMessageConverter{}.ToDomain(msg))
	}

	return messagesDomain, nil
}
