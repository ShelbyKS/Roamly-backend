package googleapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/go-resty/resty/v2"
)

const (
	methodFindPlace     = "https://maps.googleapis.com/maps/api/place/findplacefromtext/json"
	methodGetPlaceData  = "https://maps.googleapis.com/maps/api/place/details/json"
	methodGetPlace      = "https://maps.googleapis.com/maps/api/place/textsearch/json"
	methodGetPlacePhoto = "https://maps.googleapis.com/maps/api/place/photo"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Geometry struct {
	Location Location `json:"location"`
}

type Photo struct {
	PhotoReference string `json:"photo_reference"`
}

type Place struct {
	FormattedAddress string   `json:"formatted_address"`
	Geometry         Geometry `json:"geometry"`
	Name             string   `json:"name"`
	Photos           []Photo  `json:"photos"`
	PlaceID          string   `json:"place_id"`
	Rating           float64  `json:"rating"`
	Types            []string `json:"types"`
}

type Result struct {
	results []Place
}

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
	Candidates []model.GooglePlace
	Status     string `json:"status"`
	// ErrorMsg   string `json:"error_message"`
}

func (c *GoogleApiClient) FindPlace(ctx context.Context, input string, fields []string) ([]model.GooglePlace, error) {
	params := map[string]string{
		"input":     input,
		"inputtype": "textquery",
		"fields":    strings.Join(fields, ","),
		"key":       c.apiKey,
	}

	var result FindPlaceResponse

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&result).
		Get(methodFindPlace)

	log.Println("resp", resp, err)	

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
	Result   model.GooglePlace
	Status   string `json:"status"`
	// ErrorMsg string `json:"error_message"`
}

func (c *GoogleApiClient) GetPlaceByID(ctx context.Context, placeID string, fields []string) (model.GooglePlace, error) {
	params := map[string]string{
		"place_id": placeID,
		"fields":   strings.Join(fields, ","),
		"key":      c.apiKey,
	}

	var result GetPlaceDataResponse

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&result).
		Get(methodGetPlaceData)

	log.Println("body", string(resp.Body()))
	if err != nil {
		log.Println(string(resp.Body()), err)
		return model.GooglePlace{}, err
	}

	if result.Status != "OK" {
		// fmt.Println("Err:", result.ErrorMsg)
		return model.GooglePlace{}, fmt.Errorf("error: received status '%s'", result.Status)
		// log.Println()
	}

	log.Println("googleapi", result.Result)
	result.Result.PlaceID = placeID

	return result.Result, nil
}

type GeocodeResponse struct {
	HtmlAttributions []string `json:"html_attributions"`
	Results          []Place  `json:"results"`
	Status           string   `json:"status"`
}

func (c *GoogleApiClient) GetPlaces(ctx context.Context, query map[string]string) ([]Place, error) {
	var result GeocodeResponse

	query["key"] = c.apiKey

	_, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(query).
		SetResult(&result).
		Get(methodGetPlace)

	if err != nil {
		return []Place{}, err
	}

	return result.Results, nil
}

func (c *GoogleApiClient) GetPlacePhoto(ctx context.Context, reference string) ([]byte, error) {
	params := map[string]string{
		"maxwidth":        "200",
		"photo_reference": reference,
		"key":             c.apiKey,
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		Get(methodGetPlacePhoto)

	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return []byte{}, fmt.Errorf("error get photo: received status '%s'", resp.Status())
	}

	return resp.Body(), nil
}
