package handler

import (
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
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

	tripPlaceGroup := router.Group("/trip/place")
	{
		tripPlaceGroup.POST("/", handler.AddPlaceToTrip)
	}
}

type AddPlaceToTripRequest struct {
	TripID  int    `json:"trip_id"`
	PlaceID string `json:"place_id"`
}

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
