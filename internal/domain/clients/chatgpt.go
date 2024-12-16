package clients

import (
	"context"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type IChatClient interface {
	PostPrompt(ctx context.Context, messages []model.ChatMessage, model string) (string, error)
}

const (
	ModelChatGPT4o     = "gpt-4o"
	ModelChatGPT4oMini = "gpt-4o-mini"
)
