package chatgpt

import (
	"context"
	"fmt"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/go-resty/resty/v2"
)

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
	Model       string              `json:"model"`
	Messages    []model.ChatMessage `json:"messages"`
	Temperature float64             `json:"temperature"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *ChatGPTClient) PostPrompt(ctx context.Context, messages []model.ChatMessage, model string) (string, error) {
	req := Request{
		Model:       model,
		Messages:    messages,
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
		//log.Println("respnose from api: ", string(res.Body()))
		return resp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("invalid response from API: %s", res.Body())
}
