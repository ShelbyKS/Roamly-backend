package clients

import (
	"context"
)

type IChatClient interface {
	PostPrompt(ctx context.Context, prompt string) (string, error)
}
