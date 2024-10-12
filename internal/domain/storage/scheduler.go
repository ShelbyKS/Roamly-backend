package storage

import (
	"context"
)

type ISchedulerClient interface {
	PostPrompt(ctx context.Context, prompt string) (string, error)
}
