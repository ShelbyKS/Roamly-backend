package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ShelbyKS/Roamly-backend/internal/middleware"
	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/handler/dto"
)

type EventHandler struct {
	eventService service.IEventService
	tripService  service.ITripService
	lg           *logrus.Logger
}

func NewEventHandler(
	router *gin.Engine,
	lg *logrus.Logger,
	eventService service.IEventService,
	tripService service.ITripService,
) {
	handler := &EventHandler{
		lg:           lg,
		eventService: eventService,
		tripService:  tripService,
	}

	tripEventGroup := router.Group("/api/v1/trip/event")
	tripEventGroup.Use(middleware.Mw.AuthMiddleware())
	{
		tripEventGroup.POST("/",
			middleware.AccessTripByTripIdFromBodyMiddleware(tripService, middleware.ForOwnerAndEditor),
			handler.CreateEvent)
		tripEventGroup.GET("/",
			middleware.AccessTripByEventIdFromQueryMiddleware(tripService, middleware.ForAll),
			handler.GetEvent)
		tripEventGroup.PUT("/",
			middleware.AccessTripByIdOfEventFromBody(tripService, middleware.ForOwnerAndEditor),
			handler.UpdateEvent)
		tripEventGroup.DELETE("/",
			middleware.AccessTripByEventIdFromQueryMiddleware(tripService, middleware.ForOwnerAndEditor),
			handler.DeleteEvent)
	}

	router.DELETE("/api/v1/trip/:trip_id/event",
		middleware.Mw.AuthMiddleware(),
		middleware.AccessTripMiddleware(tripService, middleware.ForOwnerAndEditor),
		handler.DeleteAllEvents,
	)
}

type CreateEventRequest struct {
	Name      string    `json:"name"`
	PlaceID   string    `json:"place_id"`
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
// @Success 201 {object} dto.GetEvent
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/event [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest

	if err := c.BindJSON(&req); err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// todo: в конвертер
	event := model.Event{
		Name:      req.Name,
		PlaceID:   req.PlaceID,
		TripID:    req.TripID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	event, err := h.eventService.CreateEvent(c.Request.Context(), event)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to create event %s for trip %s", req.PlaceID, req.TripID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"event": dto.EventConverter{}.ToDto(event)})
}

type UpdateEventRequest struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	Name      string    `json:"name"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
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
// @Router /api/v1/trip/event [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	var req UpdateEventRequest

	if err := c.BindJSON(&req); err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent, err := h.eventService.UpdateEvent(c.Request.Context(), model.Event{
		ID:        req.ID,
		Name:      req.Name,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		h.lg.WithError(err).Errorf("failed to update event %s", req.ID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": dto.EventConverter{}.ToDto(updatedEvent)})
}

// @Summary Delete event
// @Description Delete an event by ID
// @Tags event
// @Accept json
// @Produce json
// @Param event_id query string true "Event ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/event [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Query("event_id"))
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := h.eventService.DeleteEvent(c.Request.Context(), eventID); err != nil {
		h.lg.WithError(err).Errorf("failed to delete event %s", eventID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get event
// @Description Get event by ID
// @Tags event
// @Accept json
// @Produce json
// @Param event_id query string true "Event ID"
// @Success 200 {object} model.Event
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/event [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Query("event_id"))
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to get event %s from trip %s", eventID, event.TripID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": dto.EventConverter{}.ToDto(event)})
}

// @Summary Delete trip events
// @Description Delete all events by trip ID
// @Tags event
// @Accept json
// @Produce json
// @Param trip_id query string true "Trip ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/{trip_id}/event [delete]
func (h *EventHandler) DeleteAllEvents(c *gin.Context) {
	tripID, err := uuid.Parse(c.Param("trip_id"))
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID"})
		return
	}

	if err := h.eventService.DeleteEventsByTrip(c.Request.Context(), tripID); err != nil {
		h.lg.WithError(err).Errorf("failed to delete events for trip %s", tripID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
