package googleapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

const (
	methodFindPlace     = "https://maps.googleapis.com/maps/api/place/findplacefromtext/json"
	methodGetPlaceData  = "https://maps.googleapis.com/maps/api/place/details/json"
	methodGetPlace      = "https://maps.googleapis.com/maps/api/place/textsearch/json"
	methodGetPlacePhoto = "https://maps.googleapis.com/maps/api/place/photo"
	methodGetTimeMatrix = "https://maps.googleapis.com/maps/api/distancematrix/json"
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

	_, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&result).
		Get(methodFindPlace)

	//log.Println("resp", resp, err)

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
	Result model.GooglePlace
	Status string `json:"status"`
	// ErrorMsg string `json:"error_message"`
}

func (c *GoogleApiClient) GetPlaceByID(ctx context.Context, placeID string, fields []string) (model.GooglePlace, error) {
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

	//log.Println("body", string(resp.Body()))
	if err != nil {
		//log.Println(string(resp.Body()), err)
		return model.GooglePlace{}, err
	}

	if result.Status != "OK" {
		// fmt.Println("Err:", result.ErrorMsg)
		return model.GooglePlace{}, fmt.Errorf("error: received status '%s'", result.Status)
		// log.Println()
	}

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

type Element struct {
	Status   string `json:"status"`
	Duration struct {
		Text  string `json:"text"`
		Value int    `json:"value"` // в секундах
	} `json:"duration"`
	Distance struct {
		Text  string `json:"text"`
		Value int    `json:"value"` // в метрах
	} `json:"distance"`
}

type Row struct {
	Elements []Element `json:"elements"`
}

type DistanceMatrixResponse struct {
	DestinationAddresses []string `json:"destination_addresses"`
	OriginAddresses      []string `json:"origin_addresses"`
	Rows                 []Row    `json:"rows"`
}

func (c *GoogleApiClient) GetTimeDistanceMatrix(ctx context.Context, placeIDs []string) (model.DistanceMatrix, error) {
	// Добавляем префикс "place_id:" к каждому идентификатору
	var placeIDsQuery []string
	for _, placeID := range placeIDs {
		placeIDsQuery = append(placeIDsQuery, "place_id:"+placeID)
	}

	placesParams := strings.Join(placeIDsQuery, "|")

	params := map[string]string{
		"origins":      placesParams,
		"destinations": placesParams,
		"key":          c.apiKey,
	}

	var result DistanceMatrixResponse

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&result).
		Get(methodGetTimeMatrix)

	if err != nil {
		return nil, fmt.Errorf("Failed to get call google api method %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error get time distance matrix: received status '%s'", resp.Status())
	}

	parsedMatrix := c.getParsedTimeDistance(placeIDs, result)

	return parsedMatrix, nil
}

func (c *GoogleApiClient) getParsedTimeDistance(placeIDs []string, response DistanceMatrixResponse) model.DistanceMatrix {
	result := make(map[string]map[string]map[string]float64)

	for i, origin := range placeIDs {
		result[origin] = make(map[string]map[string]float64)

		for j, destination := range placeIDs {
			// Извлекаем значение расстояния в километрах и времени в минутах
			distanceKm := float64(response.Rows[i].Elements[j].Distance.Value) / 1000.0
			durationMin := float64(response.Rows[i].Elements[j].Duration.Value) / 60.0

			result[origin][destination] = map[string]float64{
				"distance": distanceKm,
				"duration": durationMin,
			}
		}
	}

	//Example:
	//result[ChIJP-7oyAP00S0RtHCGoa9FgNM][ChIJP-7oyAP00S0RtHCGoa9FgNM] = { distance: 0.00 км, duration: 0.00 мин }
	//result[ChIJP-7oyAP00S0RtHCGoa9FgNM][ChIJq95xT4I30i0RaU3j93Diq8o] = { distance: 73.32 км, duration: 135.45 мин }
	//result[ChIJP-7oyAP00S0RtHCGoa9FgNM][ChIJ0Zo39DZH0i0R0acI0bADr2U] = { distance: 57.52 км, duration: 108.72 мин }
	//result[ChIJq95xT4I30i0RaU3j93Diq8o][ChIJq95xT4I30i0RaU3j93Diq8o] = { distance: 0.00 км, duration: 0.00 мин }
	//result[ChIJq95xT4I30i0RaU3j93Diq8o][ChIJ0Zo39DZH0i0R0acI0bADr2U] = { distance: 86.98 км, duration: 148.28 мин }
	//result[ChIJq95xT4I30i0RaU3j93Diq8o][ChIJP-7oyAP00S0RtHCGoa9FgNM] = { distance: 72.75 км, duration: 134.07 мин }
	//result[ChIJ0Zo39DZH0i0R0acI0bADr2U][ChIJP-7oyAP00S0RtHCGoa9FgNM] = { distance: 57.22 км, duration: 113.58 мин }
	//result[ChIJ0Zo39DZH0i0R0acI0bADr2U][ChIJq95xT4I30i0RaU3j93Diq8o] = { distance: 84.93 км, duration: 158.73 мин }
	//result[ChIJ0Zo39DZH0i0R0acI0bADr2U][ChIJ0Zo39DZH0i0R0acI0bADr2U] = { distance: 0.00 км, duration: 0.00 мин }

	for origin, destinations := range result {
		for destination, metrics := range destinations {
			if origin == destination {
				continue
			}
			fmt.Printf("result[%s][%s] = { distance: %.2f км, duration: %.2f мин }\n",
				origin, destination, metrics["distance"], metrics["duration"])
		}
	}

	return result
}
