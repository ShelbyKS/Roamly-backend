package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/middleware"
	"github.com/google/uuid"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/handler/dto"
	"github.com/ShelbyKS/Roamly-backend/pkg/googleapi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PlaceHandler struct {
	lg           *logrus.Logger
	placeService service.IPlaceService
	client       googleapi.GoogleApiClient
}

func NewPlaceHandler(router *gin.Engine, lg *logrus.Logger, placeService service.IPlaceService, client googleapi.GoogleApiClient) {
	handler := &PlaceHandler{
		lg:           lg,
		placeService: placeService,
		client:       client,
	}

	router.GET("/api/v1/place", middleware.Mw.AuthMiddleware(), handler.GetPlaces)
	router.GET("/api/v1/place/find", middleware.Mw.AuthMiddleware(), handler.FindPlaces)
	router.GET("/api/v1/place/photo", middleware.Mw.AuthMiddleware(), handler.GetPhoto)
	router.DELETE("/api/v1/trip/:trip_id/place/:place_id", middleware.Mw.AuthMiddleware(), handler.DeletePlaceFromTrip)
	router.GET("api/v1/place/recomendations", handler.GetPlacesNearby)

	tripPlaceGroup := router.Group("/api/v1/trip/place")
	tripPlaceGroup.Use(middleware.Mw.AuthMiddleware())
	{
		tripPlaceGroup.POST("/", handler.AddPlaceToTrip)
	}
}

type AddPlaceToTripRequest struct {
	TripID  string `json:"trip_id" form:"trip_id" binding:"required"`
	PlaceID string `json:"place_id" form:"place_id" binding:"required"`
}

// @Summary Add place to trip
// @Description Add a place to a specific trip by their IDs
// @Tags place
// @Accept json
// @Produce json
// @Param trip-place body AddPlaceToTripRequest true "JSON containing trip and place IDs"
// @Success 200 {object} dto.TripResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/trip/place [post]
func (h *PlaceHandler) AddPlaceToTrip(c *gin.Context) {
	var req AddPlaceToTripRequest

	err := c.Bind(&req)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tripUUID, err := uuid.Parse(req.TripID)
	if err != nil {
		h.lg.WithError(err).Errorf("invalid trip_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip_id format"})
		return
	}

	trip, err := h.placeService.AddPlaceToTrip(c.Request.Context(), tripUUID, req.PlaceID)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to add place to trip")
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip": dto.TripConverter{}.ToDto(trip)})
}

// @Summary Delete place from trip
// @Description Delete place from a specific trip by their IDs
// @Tags place
// @Accept json
// @Produce json
// @Success 200 {object} dto.TripResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/trip/{trip_id}/place/{place_id} [delete]
func (h *PlaceHandler) DeletePlaceFromTrip(c *gin.Context) {
	tripID := c.Param("trip_id")
	placeID := c.Param("place_id")

	tripUUID, err := uuid.Parse(tripID)
	if err != nil {
		h.lg.WithError(err).Errorf("invalid trip_id format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip_id format"})
		return
	}

	trip, err := h.placeService.DeletePlace(c.Request.Context(), tripUUID, placeID)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to remove place from trip")
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip": dto.TripConverter{}.ToDto(trip)})
}

// @Summary Find places
// @Description Find places by searchString
// @Tags place
// @Produce json
// @Param searchString query string true "Search string to search places"
// @Success 200 {object} map[string][]dto.GooglePlace "List of found places"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/place/find [get]
func (h *PlaceHandler) FindPlaces(c *gin.Context) {
	searchString := c.Query("searchString")

	places, err := h.placeService.FindPlace(c.Request.Context(), searchString)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to find place by %s", searchString)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	placesDto := make([]dto.GooglePlace, len(places))
	for i, place := range places {
		placesDto[i] = dto.GooglePlaceConverter{}.ToDto(place.GooglePlace)
	}

	c.JSON(http.StatusOK, gin.H{
		"places": placesDto,
	})
}

// @Summary Get places
// @Description Find places by name or location
// @Tags place
// @Produce json
// @Param name query string true "Name to search places by"
// @Param type query string false "Type of place"
// @Param lat query string false "Latitude for location-based search"
// @Param lng query string false "Longitude for location-based search"
// @Success 200 {object} map[string][]model.Place "List of found places"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/place [get]
func (h *PlaceHandler) GetPlaces(c *gin.Context) {
	// todo: логику унести в сервис, клиента обернуть в сервис
	qeuryMap := map[string]string{}

	name := c.Query("name")
	qeuryMap["query"] = name

	qeuryMap["fields"] = "formatted_address,name,rating,geometry"

	typeQuery, hasType := c.GetQuery("type")
	if hasType {
		qeuryMap["type"] = typeQuery
	} else {
		lat := c.Query("lat")
		lng := c.Query("lng")
		qeuryMap["location"] = fmt.Sprintf("%s,%s", lat, lng)
		qeuryMap["radius"] = "20000"
	}

	places, err := h.client.GetPlaces(c.Request.Context(), qeuryMap)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to get places from google")
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": places})
}

// @Summary Get place photo
// @Description Get a photo of a place by photo reference
// @Tags place
// @Produce image/jpeg
// @Param reference query string true "Photo reference ID"
// @Success 200 {file} jpeg "Binary image data of the place photo"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/place/photo [get]
func (h *PlaceHandler) GetPhoto(c *gin.Context) {
	reference := c.Query("reference")

	file, err := h.client.GetPlacePhoto(c.Request.Context(), reference)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to get place %s photo", reference)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	fileName := "photo.jpg"

	c.Header("Content-Disposition", `inline; filename="`+fileName+`"`)
	c.Data(http.StatusOK, "text/plain", file)
}

func (h *PlaceHandler) GetPlacesNearby(c *gin.Context) {
	lat, ok := c.GetQuery("lat")
	latFloat, err := strconv.ParseFloat(lat, 64)
	if !ok || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't get lat:" + err.Error()})
		return
	}

	lng, ok := c.GetQuery("lng")
	lngFloat, err := strconv.ParseFloat(lng, 64)
	if !ok || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't get lng:" + err.Error()})
		return
	}
	placesTypesString := c.Query("types")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't get types"})
		return
	}

	placesTypes := strings.Split(placesTypesString, ",")

	places, err := h.placeService.GetPlacesNearby(c.Request.Context(), latFloat, lngFloat, placesTypes)

	c.JSON(http.StatusOK, places)
}
