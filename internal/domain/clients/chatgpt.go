package clients

import (
	"context"
)

type IChatClient interface {
	PostPrompt(ctx context.Context, prompt, model string) (string, error)
}

const (
	ModelChatGPT4o     = "gpt-4o"
	ModelChatGPT4oMini = "gpt-4o-mini"
)
