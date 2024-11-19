package handler

import (
	"fmt"
	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/middleware"
	"net/http"

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

	// api := router.Group("/api/v1")
	router.GET("/api/v1/place", middleware.Mw.AuthMiddleware(), handler.GetPlaces)
	router.GET("/api/v1/place/find", middleware.Mw.AuthMiddleware(), handler.FindPlaces)
	router.GET("/api/v1/place/photo", middleware.Mw.AuthMiddleware(), handler.GetPhoto)
	router.GET("api/v1/place/recomendations", handler.GetPlacesNearby)
}

type AddPlaceToTripRequest struct {
	TripID  string `json:"trip_id" form:"trip_id" binding:"required"`
	PlaceID string `json:"place_id" form:"place_id" binding:"required"`
}

// @Summary Find places
// @Description Find places by searchString
// @Tags place
// @Produce json
// @Param searchString query string true "Search string to search places"
// @Success 200 {object} map[string][]dto.PlaceGoogle "List of found places"
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

	placesDto := make([]dto.PlaceGoogle, len(places))
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

	qeuryMap["fields"] = "formatted_address,name,rating,geometry,photos,types,editorial_summary"

	typeQuery, hasType := c.GetQuery("type")
	if hasType {
		qeuryMap["type"] = typeQuery
	} else {
		lat := c.Query("lat")
		lng := c.Query("lng")
		radius := c.Query("radius")
		qeuryMap["location"] = fmt.Sprintf("%s,%s", lat, lng)
		qeuryMap["radius"] = radius
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
	//lat, ok := c.GetQuery("lat")
	//latFloat, err := strconv.ParseFloat(lat, 64)
	//if !ok || err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "can't get lat:" + err.Error()})
	//	return
	//}
	//
	//lng, ok := c.GetQuery("lng")
	//lngFloat, err := strconv.ParseFloat(lng, 64)
	//if !ok || err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "can't get lng:" + err.Error()})
	//	return
	//}
	//
	//radius, ok := c.GetQuery("radius")
	//radiusFloat, err := strconv.ParseFloat(radius, 64)
	//if !ok || err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "can't get radius:" + err.Error()})
	//	return
	//}
	//
	//maxPlaces, ok := c.GetQuery("max_places")
	//maxPlacesInt, err := strconv.Atoi(maxPlaces)
	//if !ok || err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "can't get max_places:" + err.Error()})
	//	return
	//}
	//
	//placesTypesString := c.Query("types")
	//if !ok {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "can't get types"})
	//	return
	//}
	//
	//placesTypes := strings.Split(placesTypesString, ",")
	//
	//places, err := h.placeService.GetPlacesNearby(c.Request.Context(),
	//	radiusFloat,
	//	latFloat,
	//	lngFloat,
	//	placesTypes,
	//	maxPlacesInt)
	//if err != nil {
	//	c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": "can't get places nerby: " + err.Error()})
	//}
	//
	//c.JSON(http.StatusOK, places)
	c.JSON(http.StatusOK, gin.H{})
}
