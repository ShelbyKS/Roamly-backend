package chatgpt

import (
	"context"

	"github.com/ayush6624/go-chatgpt"
)

type ChatGPTClient struct {
	// apiKey string
	client *chatgpt.Client
}

func NewChatGPTClient(apiKey string) (*ChatGPTClient, error) {
	client, err := chatgpt.NewClient(apiKey)
	if err != nil {
		return nil, err
	}
	
	return &ChatGPTClient{
		client: client,
	}, nil
}

func (c *ChatGPTClient) PostPrompt(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.SimpleSend(ctx, prompt)
	
	return resp.Object, err
}
