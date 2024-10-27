package handler

import (
	"fmt"
	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/middleware"
	"net/http"

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

	router.GET("/api/v1/place", handler.GetPlaces)
	router.GET("/api/v1/place/find", handler.FindPlaces)
	router.GET("/api/v1/place/photo", handler.GetPhoto)

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
// @Description Add place to trip by id
// @Tags place
// @Accept json
// @Produce json
// @Param trip-place body AddPlaceToTripRequest true "Place and trip IDs"
// @Success 200 {object} model.Trip
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/place [post]
func (h *PlaceHandler) AddPlaceToTrip(c *gin.Context) {
	var req AddPlaceToTripRequest

	err := c.Bind(&req)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	tripUUID, err := uuid.Parse(req.TripID)
	if err != nil {
		h.lg.WithError(err).Errorf("invalid trip_id")
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid trip_id format"})
		return
	}

	err = h.placeService.AddPlaceToTrip(c.Request.Context(), tripUUID, req.PlaceID)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to add place to trip")
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Find place
// @Description Find places by searchString
// @Tags place
// @Produce json
// @Param searchString query true "SearchString"
// @Success 200 {object} model.Trip
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/place [get]
func (h *PlaceHandler) FindPlaces(c *gin.Context) {
	searchString := c.Query("searchString")

	places, err := h.placeService.FindPlace(c.Request.Context(), searchString)
	if err != nil {
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()})
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
// @Description Find places by searchString
// @Tags place
// @Produce json
// @Param searchString query true "SearchString"
// @Success 200 {object} model.Trip
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
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
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": places})
}

// @Summary Get place photo
// @Description Find places by searchString
// @Tags place
// @Produce json
// @Param searchString query true "SearchString"
// @Success 200 {object} model.Trip
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/place [get]
func (h *PlaceHandler) GetPhoto(c *gin.Context) {
	reference := c.Query("reference")

	file, err := h.client.GetPlacePhoto(c.Request.Context(), reference)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to get place %s photo", reference)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()})
		return
	}

	fileName := "photo.jpg"

	c.Header("Content-Disposition", `inline; filename="`+fileName+`"`)
	c.Data(http.StatusOK, "text/plain", file)
}
