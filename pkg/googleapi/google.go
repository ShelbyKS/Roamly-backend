package googleapi

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/pkg/googleapi/dto"
)

const (
	methodFindPlace       = "https://maps.googleapis.com/maps/api/place/findplacefromtext/json"
	methodGetPlaceData    = "https://maps.googleapis.com/maps/api/place/details/json"
	methodGetPlace        = "https://maps.googleapis.com/maps/api/place/textsearch/json"
	methodGetPlacePhoto   = "https://maps.googleapis.com/maps/api/place/photo"
	methodGetTimeMatrix   = "https://maps.googleapis.com/maps/api/distancematrix/json"
	methodGetPlacesNearby = "https://places.googleapis.com/v1/places:searchNearby"
	// methodGetPlacesNearby = "https://maps.googleapis.com/maps/api/place/nearbysearch/json"

	fieldMask = "places.id,places.formattedAddress,places.displayName,places.rating,places.location,places.photos,places.editorialSummary"
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

type EditorialSummary struct {
	Overview string `json:"overview"`
	Language string `json:"language"`
}

type Place struct {
	FormattedAddress string           `json:"formatted_address"`
	Geometry         Geometry         `json:"geometry"`
	Name             string           `json:"name"`
	Photos           []Photo          `json:"photos"`
	PlaceID          string           `json:"place_id"`
	Rating           float64          `json:"rating"`
	Types            []string         `json:"types"`
	EditorialSummary EditorialSummary `json:"editorial_summary"`
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

	params["language"] = "ru"

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

	params["language"] = "ru"

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

//func (c *GoogleApiClient) GetPlaces(ctx context.Context, query map[string]string) ([]Place, error) {
//	var result GeocodeResponse
//
//	query["key"] = c.apiKey
//	query["language"] = "ru"
//
//	_, err := c.client.R().
//		SetContext(ctx).
//		SetQueryParams(query).
//		SetResult(&result).
//		Get(methodGetPlace)
//
//	if err != nil {
//		return []Place{}, err
//	}
//
//	return result.Results, nil
//}

func (c *GoogleApiClient) GetPlaces(ctx context.Context, query map[string]string) ([]Place, error) {
	var result GeocodeResponse

	query["key"] = c.apiKey
	query["language"] = "ru"

	_, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(query).
		SetResult(&result).
		Get(methodGetPlace)

	if err != nil {
		return nil, err
	}

	if _, ok := query["type"]; ok {
		return result.Results, nil
	}

	lat, err := strconv.ParseFloat(strings.Split(query["location"], ",")[0], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude: %v", err)
	}

	lng, err := strconv.ParseFloat(strings.Split(query["location"], ",")[1], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude: %v", err)
	}

	radius, err := strconv.ParseFloat(query["radius"], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid radius: %v", err)
	}

	var filteredPlaces []Place
	for _, place := range result.Results {
		placeLat := place.Geometry.Location.Lat
		placeLng := place.Geometry.Location.Lng

		if isWithinRadius(lat, lng, placeLat, placeLng, radius) {
			filteredPlaces = append(filteredPlaces, place)
		}
	}

	return filteredPlaces, nil
}

func isWithinRadius(lat1, lng1, lat2, lng2, radius float64) bool {
	const EarthRadius = 6371e3

	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180

	deltaLat := lat2Rad - lat1Rad
	deltaLng := lng2Rad - lng1Rad

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := EarthRadius * c

	return distance <= radius
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
	Status               string   `json:"status"`
}

func (c *GoogleApiClient) GetTimeDistanceMatrix(ctx context.Context, placeIDs []string) (model.DistanceMatrix, error) {
	// Добавляем префикс "place_id:" к каждому идентификатору
	var placeIDsQuery []string
	for _, placeID := range placeIDs {
		placeIDsQuery = append(placeIDsQuery, "place_id:"+placeID)
	}

	fmt.Println("LEN: ", len(placeIDsQuery))

	placesParams := strings.Join(placeIDsQuery, "|")

	params := map[string]string{
		"origins":      placesParams,
		"destinations": placesParams,
		"key":          c.apiKey,
	}

	params["language"] = "ru"

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

	fmt.Println("GOOGLE result: ", result.Status)

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

func (c *GoogleApiClient) GetPlacesNearby(ctx context.Context,
	includedTypes []string,
	maxPlaces int,
	rankPrefernce string,
	lat float64,
	lng float64,
	radius float64,
	languageCode string) ([]model.GooglePlace, error) {

	var result dto.PlacesNearbyResult

	request := dto.PlacesNearbyRequest{
		IncludeTypes:   includedTypes,
		MaxResultCount: maxPlaces,
		RankPreference: rankPrefernce,
		LocationRestriction: dto.LocationRestriction{
			Circle: dto.Circle{
				Center: dto.Center{
					Latitude:  lat,
					Longitude: lng,
				},
				Radius: radius,
			},
		},
		LanguageCode: languageCode,
	}

	// query := map[string]string{}
	// // 	"fields": c.joinFields([]string{
	// // 		"formatted_address",
	// // 		"name",
	// // 		"rating",
	// // 		"geometry",
	// // 		"photo"}),
	// // }

	// query["language"] = "ru"
	// query["location"] = fmt.Sprintf("%f,%f", lat, lng)
	// query["radius"] = fmt.Sprintf("%f", radius)
	// query["type"] = includedTypes[0]
	// query["key"] = c.apiKey

	_, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Goog-Api-Key", c.apiKey).
		// SetQueryParams(query).
		SetHeader("X-Goog-FieldMask", fieldMask).
		SetBody(request).
		SetResult(&result).
		Post(methodGetPlacesNearby)

	log.Println(result.Results[0].ID)
	// log.Println(string(resp.Body()), "key:", c.apiKey)

	if err != nil {
		return []model.GooglePlace{}, err
	}

	places := make([]model.GooglePlace, len(result.Results))
	for i, place := range result.Results {
		places[i] = dto.PlaceConverter{}.ToDomain(place)
	}

	return places, nil
}
