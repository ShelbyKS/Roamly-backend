package handler

import (
	"net/http"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PlaceHandler struct {
	lg           *logrus.Logger
	placeService service.IPlaceService
}

func NewPlaceHandler(router *gin.Engine, lg *logrus.Logger, placeService service.IPlaceService) {
	handler := &PlaceHandler{
		lg:           lg,
		placeService: placeService,
	}

	router.GET("/api/v1/place", handler.FindPlaces)

	tripPlaceGroup := router.Group("/api/v1/trip/place")
	{
		tripPlaceGroup.POST("/", handler.AddPlaceToTrip)
	}
}

type AddPlaceToTripRequest struct {
	TripID  uuid.UUID `json:"trip_id"`
	PlaceID string    `json:"place_id"`
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

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.placeService.AddPlaceToTrip(c.Request.Context(), req.TripID, req.PlaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"places": places})
}
