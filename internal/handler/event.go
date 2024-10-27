package handler

import (
	"github.com/ShelbyKS/Roamly-backend/internal/middleware"
	"net/http"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/handler/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type EventHandler struct {
	eventService service.IEventService
	lg           *logrus.Logger
}

func NewEventHandler(router *gin.Engine,
	lg *logrus.Logger, eventService service.IEventService) {

	handler := &EventHandler{
		lg:           lg,
		eventService: eventService,
	}

	eventGroup := router.Group("/api/v1/event")
	eventGroup.Use(middleware.Mw.AuthMiddleware())
	{
		eventGroup.POST("/", handler.CreateEvent)
		eventGroup.GET("/", handler.GetEvent)
		eventGroup.PUT("/", handler.UpdateEvent)
		eventGroup.DELETE("/", handler.DeleteEvent)
	}
}

type CreateEventRequest struct {
	PlaceID   string    `json:"place_id" binding:"required"`
	TripID    uuid.UUID `json:"trip_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
}

// @Summary Create event
// @Description Create a new event
// @Tags event
// @Accept json
// @Produce json
// @Param event body CreateEventRequest true "Event data"
// @Success 201 {object} model.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/event [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest

	if err := c.BindJSON(&req); err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// todo: в конвертер
	event := model.Event{
		PlaceID:   req.PlaceID,
		TripID:    req.TripID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if err := h.eventService.CreateEvent(c.Request.Context(), event); err != nil {
		h.lg.WithError(err).Errorf("failed to create event %s for trip %s", req.PlaceID, req.TripID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

type UpdateEventRequest struct {
	PlaceID   string    `json:"place_id" binding:"required"`
	TripID    uuid.UUID `json:"trip_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
}

// @Summary Update event
// @Description Update event data
// @Tags event
// @Accept json
// @Produce json
// @Param event body UpdateEventRequest true "Event data"
// @Success 200 {object} model.Event
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/event [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	var req UpdateEventRequest

	if err := c.BindJSON(&req); err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// todo: в конвертер
	event := model.Event{
		PlaceID:   req.PlaceID,
		TripID:    req.TripID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if err := h.eventService.UpdateEvent(c.Request.Context(), event); err != nil {
		h.lg.WithError(err).Errorf("failed to update event %s from trip %s ", req.PlaceID, req.TripID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

type DeleteEventRequest struct {
	PlaceID string    `json:"place_id" binding:"required"`
	TripID  uuid.UUID `json:"trip_id" binding:"required"`
}

// @Summary Delete event
// @Description Delete an event by place ID and trip ID
// @Tags event
// @Accept json
// @Produce json
// @Param event body DeleteEventRequest true "Event data"
// @Success 204 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/event [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	var req DeleteEventRequest

	if err := c.BindJSON(&req); err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := h.eventService.DeleteEvent(c.Request.Context(), req.PlaceID, req.TripID); err != nil {
		h.lg.WithError(err).Errorf("failed to delete event %s from trip %s", req.PlaceID, req.TripID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// @Summary Get event
// @Description Get event by place ID and trip ID
// @Tags event
// @Accept json
// @Produce json
// @Param place_id query string true "Place ID"
// @Param trip_id query string true "Trip ID"
// @Success 200 {object} model.Event
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/event [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	placeID := c.Query("place_id")
	tripID, err := uuid.Parse(c.Query("trip_id"))
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse query")
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid trip ID"})
		return
	}

	event, err := h.eventService.GetEventByID(c.Request.Context(), placeID, tripID)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to get event %s from trip %s", placeID, tripID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.EventConverter{}.ToDto(event))
}
