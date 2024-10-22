package googleapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/go-resty/resty/v2"
)

const (
	methodFindPlace    = "https://maps.googleapis.com/maps/api/place/findplacefromtext/json"
	methodGetPlaceData = "https://places.googleapis.com/v1/places"
)

type GoogleApiClient struct {
	client *resty.Client
	apiKey string
}

func NewClient(apiKey string) *GoogleApiClient {
	// init settings from cfg

	return &GoogleApiClient{
		client: resty.New(),
		apiKey: apiKey,
	}
}

type FindPlaceResponse struct {
	Candidates []model.Place
	Status     string `json:"status"`
	ErrorMsg   string `json:"error_message"`
}

func (c *GoogleApiClient) FindPlace(ctx context.Context, input string, fields []string) ([]model.Place, error) {
	params := map[string]string{
		"input":     input,
		"inputtype": "textquery",
		"fields":    strings.Join(fields, ","),
		"key":       c.apiKey,
	}

	var result FindPlaceResponse

	_, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&result).
		Get(methodFindPlace)

	if err != nil {
		return nil, err
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("error: received status '%s'", result.Status)
	}

	return result.Candidates, nil
}

func (c *GoogleApiClient) joinFields(fields []string) string {
	return fmt.Sprintf("%s", fields)
}

type GetPlaceDataResponse struct {
	Place    model.Place
	Status   string `json:"status"`
	ErrorMsg string `json:"error_message"`
}

func (c *Client) GetPlaceByID(ctx context.Context, placeID string, fields []string) (model.Place, error) {
	params := map[string]string{
		"place_id": placeID,
		"fields":   strings.Join(fields, ","),
		"key":      c.apiKey,
	}

	var result GetPlaceDataResponse

	_, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&result).
		Get(methodGetPlaceData)

	if err != nil {
		return model.Place{}, err
	}

	if result.Status != "OK" {
		fmt.Println("Err:", result.ErrorMsg)
		return model.Place{}, fmt.Errorf("error: received status '%s'", result.Status)
	}

	return result.Place, nil
}
