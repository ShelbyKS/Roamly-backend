package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
)

type TripHandler struct {
	lg               *logrus.Logger
	tripService      service.ITripService
	schedulerService service.ISchedulerService
	placesService    service.IPlaceService
}

func NewTripHandler(router *gin.Engine,
	lg *logrus.Logger,
	tripService service.ITripService,
	schedulerService service.ISchedulerService) {

	handler := &TripHandler{
		lg:               lg,
		tripService:      tripService,
		schedulerService: schedulerService,
	}

	tripGroup := router.Group("/api/v1/trip")
	{
		tripGroup.GET("/:trip_id", handler.GetTripByID)
		tripGroup.POST("/", handler.CreateTrip)
		tripGroup.PUT("/", handler.UpdateTrip)
		tripGroup.DELETE("/:trip_id", handler.DeleteTrip)

		tripGroup.GET("/:trip_id/schedule", handler.ScheduleTrip)
	}
}

// @Summary Get trip by ID
// @Description Get data of a specific trip by its ID
// @Tags trip
// @Produce json
// @Param trip_id path string true "Trip ID"
// @Success 200 {object} model.Trip
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/{trip_id} [get]
func (h *TripHandler) GetTripByID(c *gin.Context) {
	idString := c.Param("trip_id")
	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	trip, err := h.tripService.GetTripByID(c.Request.Context(), id)
	if errors.Is(err, domain.ErrTripNotFound) {
		h.lg.Warnf("Trip with id=%d not found", id)
		c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}
	if err != nil {
		h.lg.WithError(err).Errorf("Fail to get trip with id=%d", id)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip": trip})
}

// @Summary Delete a trip
// @Description Delete a trip by its ID
// @Tags trip
// @Produce json
// @Param trip_id path string true "Trip ID"
// @Success 200 {null} string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/{trip_id} [delete]
func (h *TripHandler) DeleteTrip(c *gin.Context) {
	idString := c.Param("trip_id")
	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.tripService.DeleteTrip(c.Request.Context(), id)
	if errors.Is(err, domain.ErrTripNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type CreateTripRequest struct {
	UserID    int    `json:"user_id" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	AreaID    string `json:"area_id" binding:"required"`
}

// @Summary Create a new trip
// @Description Create a new trip for the user
// @Tags trip
// @Accept json
// @Produce json
// @Param trip body CreateTripRequest true "Trip data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip [post]
func (h *TripHandler) CreateTrip(c *gin.Context) {
	var tripReq CreateTripRequest

	err := c.BindJSON(&tripReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id, err := h.tripService.CreateTrip(c.Request.Context(), model.Trip{
		StartTime: tripReq.StartTime,
		EndTime:   tripReq.EndTime,
		AreaID:    tripReq.AreaID,
		// TODO: ХААААААААРДКОД!!!!!!!!!!!!!!!!!!!!!!
		Users: []*model.User{
			{
				ID: tripReq.UserID,
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

type UpdateTripRequest struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	AreaID    string    `json:"area_id" binding:"required"`
}

// @Summary Update trip
// @Description Update trip data
// @Tags trip
// @Accept  json
// @Produce  json
// @Param trip body UpdateTripRequest true "Trip data"
// @Success 200 {object} model.Trip
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip [put]
func (h *TripHandler) UpdateTrip(c *gin.Context) {
	var tripReq UpdateTripRequest

	err := c.BindJSON(&tripReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.tripService.UpdateTrip(c.Request.Context(), model.Trip{
		ID:        tripReq.ID,
		StartTime: tripReq.StartTime,
		EndTime:   tripReq.EndTime,
		AreaID:    tripReq.AreaID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Schedule trip
// @Description Schedule places  in trip
// @Tags trip
// @Produce json
// @Param trip_id path string true "Trip ID"
// @Success 200 {null} string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/{trip_id}/schedule [get]
func (h *TripHandler) ScheduleTrip(c *gin.Context) {
	idString := c.Param("trip_id")
	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	trip, err := h.tripService.GetTripByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	// h.lg.Info("trip", trip)
	// h.lg.Info("service", h.schedulerService)
	// h.lg.Info("service", trip.Places)
	matrix := h.placesService.GetTimeMatrix(c.Request.Context(), trip.Places)

	schedule, err := h.schedulerService.GetSchedule(c.Request.Context(), trip, trip.Places, matrix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedule": schedule})
}
