package chatgpt

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// type ChatGPTClient struct {
// 	client *chatgpt.Client
// }

// func NewChatGPTClient(apiKey string) (*ChatGPTClient, error) {
// 	client, err := chatgpt.NewClient(apiKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ChatGPTClient{
// 		client: client,
// 	}, nil
// }

// func (c *ChatGPTClient) PostPrompt(ctx context.Context, prompt string) (string, error) {
// 	resp, err := c.client.Send(ctx, &chatgpt.ChatCompletionRequest{
// 		Model: chatgpt.ChatGPTModel("gpt-4o-mini"),
// 		Messages: []chatgpt.ChatMessage{
// 			{
// 				Role:    chatgpt.ChatGPTModelRoleSystem,
// 				Content: prompt,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("can't send message to chatgpt: %w", err)
// 	}

// 	return resp.Object, nil
// }

// package main

const (
	url = "https://api.openai.com/v1/chat/completions"
)

type ChatGPTClient struct {
	client *resty.Client
	apiKey string
}

func NewChatGPTClient(apiKey string) *ChatGPTClient {
	return &ChatGPTClient{
		client: resty.New(),
		apiKey: apiKey,
	}
}

type Request struct {
	Model       string           `json:"model"`
	Messages    []MessageRequest `json:"messages"`
	Temperature float64          `json:"temperature"`
}

type MessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *ChatGPTClient) PostPrompt(ctx context.Context, prompt string) (string, error) {
	req := Request{
		Model:       "gpt-4o-mini",
		Messages:    []MessageRequest{{Role: "user", Content: prompt}},
		Temperature: 0.7,
	}

	var resp Response
	res, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(c.apiKey).
		SetBody(req).
		SetResult(&resp).
		Post(url)

	if err != nil {
		return "", fmt.Errorf("%s:%w", res.Body(), err)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("invalid response from API: %s", res.Body())
}
