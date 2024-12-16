package utils

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/clients"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type NotifyUtils struct {
	tripStorage     storage.ITripStorage
	sessionStorage  storage.ISessionStorage
	messageProducer clients.IMessageProdcuer
}

func NewNotifyUtils(
	tripStorage storage.ITripStorage,
	sessionStorage storage.ISessionStorage,
	messageProducer clients.IMessageProdcuer,
) NotifyUtils {
	return NotifyUtils{
		tripStorage:     tripStorage,
		sessionStorage:  sessionStorage,
		messageProducer: messageProducer,
	}
}

func (utils *NotifyUtils) FormAndSendNotifyMessage(
	ctx context.Context,
	tripID uuid.UUID,
	action string,
	message string,
	authorID int,
) error {
	trip, err := utils.tripStorage.GetTripByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to get trip: %w", err)
	}

	var notifyMessage model.NotifyMessage
	for _, user := range trip.Users {
		sessionTokens, err := utils.sessionStorage.GetTokensByUserID(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to get user session tokens: %w", err)
		}

		notifyMessage.Clients = append(notifyMessage.Clients, sessionTokens...)
	}

	notifyMessage.Payload.Action = action
	notifyMessage.Payload.TripID = trip.ID
	notifyMessage.Payload.Author = fmt.Sprintf("%d", authorID)
	notifyMessage.Payload.Message = message
	err = utils.messageProducer.SendMessage(notifyMessage)
	if err != nil {
		return fmt.Errorf("failed to send action %s: %w", action, err)
	}

	return nil
}
