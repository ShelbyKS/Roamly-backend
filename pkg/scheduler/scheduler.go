package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	URL = "http://localhost:11434/api/generate"
)

type SchedulerClient struct {
	client *resty.Client
	url    string
}

func NewClient(url string) *SchedulerClient {
	return &SchedulerClient{
		client: resty.New(),
		url:    url,
	}
}

type Request struct {
	Model  string
	Prompt string
	Stream bool
}

type Response struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

func (c *SchedulerClient) PostPrompt(ctx context.Context, prompt string) (string, error) {
	req := Request{
        Model:  "llama3.2",
        Prompt: prompt,
        Stream: false,
    }
	
	var resp Response
	_, err := c.client.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&resp).
		Post(c.url)

	log.Println(resp, err)

	if err != nil {
		return "", err
	}

	return resp.Response, nil
}
