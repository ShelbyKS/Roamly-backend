package handler

import (
	"errors"
	"net/http"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/handler/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CreateEventRequest struct {
	PlaceID   string    `json:"place_id" binding:"required"`
	TripID    uuid.UUID `json:"trip_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
}

type UpdateEventRequest struct {
	PlaceID   string    `json:"place_id" binding:"required"`
	TripID    uuid.UUID `json:"trip_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
}

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
	{
		eventGroup.POST("/", handler.CreateEvent)
		eventGroup.GET("/", handler.GetEvent)
		eventGroup.PUT("/", handler.UpdateEvent)
		eventGroup.DELETE("/", handler.DeleteEvent)
	}
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
	var eventReq CreateEventRequest

	if err := c.BindJSON(&eventReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// todo: в конвертер
	event := model.Event{
		PlaceID:   eventReq.PlaceID,
		TripID:    eventReq.TripID,
		StartTime: eventReq.StartTime,
		EndTime:   eventReq.EndTime,
	}

	if err := h.eventService.CreateEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
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
	var eventReq UpdateEventRequest

	if err := c.BindJSON(&eventReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// todo: в конвертер
	event := model.Event{
		PlaceID:   eventReq.PlaceID,
		TripID:    eventReq.TripID,
		StartTime: eventReq.StartTime,
		EndTime:   eventReq.EndTime,
	}

	if err := h.eventService.UpdateEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Delete event
// @Description Delete an event by place ID and trip ID
// @Tags event
// @Accept json
// @Produce json
// @Param place_id query string true "Place ID"
// @Param trip_id query string true "Trip ID"
// @Success 204 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/event [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	placeID := c.Query("place_id")
	tripID, err := uuid.Parse(c.Query("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid trip ID"})
		return
	}

	if err := h.eventService.DeleteEvent(c.Request.Context(), placeID, tripID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid trip ID"})
		return
	}

	event, err := h.eventService.GetEventByID(c.Request.Context(), placeID, tripID)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"err": "Event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.EventConverter{}.ToDto(event))
}
