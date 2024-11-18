package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// проверяет доступ к поездке по event_id из квери
func AccessTripByEventIdFromQueryMiddleware(tripService service.ITripService, userRoles []model.UserTripRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No user id"})
			c.Abort()
		}
		userIDInt, ok := userID.(int)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse user id: "})
			c.Abort()
		}

		eventID := c.Query("event_id")
		eventIDuuid, err := uuid.Parse(eventID)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse trip id: " + err.Error()})
			c.Abort()
			return
		}

		trip, err := tripService.GetTripByEventID(c.Request.Context(), eventIDuuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get trip by event id: " + err.Error()})
			c.Abort()
		}

		role, err := tripService.GetUserRole(c.Request.Context(), userIDInt, trip.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "can't get role: " + err.Error()})
			c.Abort()
			return
		}
		if !slices.Contains(userRoles, role) {
			log.Println("roles:", userRoles, role)
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			c.Abort()
		}

		c.Next()
	}
}

// проверяет доступ к поездке по trip_id из body
func AccessTripByTripIdFromBodyMiddleware(tripService service.ITripService, userRoles []model.UserTripRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No user id"})
			c.Abort()
			return
		}
		userIDInt, ok := userID.(int)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse user id"})
			c.Abort()
			return
		}

		// Клонируем тело запроса (дважды бади читать нельзя)
		var bodyBuffer bytes.Buffer
		_, err := io.Copy(&bodyBuffer, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body: " + err.Error()})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBuffer.Bytes()))

		var body struct {
			TripID string `json:"trip_id"`
		}
		if err := json.NewDecoder(&bodyBuffer).Decode(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body: " + err.Error()})
			c.Abort()
			return
		}

		tripIDuuid, err := uuid.Parse(body.TripID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse trip id:" + err.Error()})
			c.Abort()
			return
		}

		role, err := tripService.GetUserRole(c.Request.Context(), userIDInt, tripIDuuid)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "can't get role: " + err.Error()})
			c.Abort()
			return
		}

		if !slices.Contains(userRoles, role) {
			log.Println("roles:", userRoles, role)
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// проверяет доступ к поездке по id (является id event) из body
func AccessTripByIdOfEventFromBody(tripService service.ITripService, userRoles []model.UserTripRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No user id"})
			c.Abort()
			return
		}
		userIDInt, ok := userID.(int)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse user id"})
			c.Abort()
			return
		}

		// Клонируем тело запроса (дважды бади читать нельзя)
		var bodyBuffer bytes.Buffer
		_, err := io.Copy(&bodyBuffer, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body:" + err.Error()})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBuffer.Bytes()))

		var body struct {
			EventID string `json:"id"`
		}
		if err := json.NewDecoder(&bodyBuffer).Decode(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body: " + err.Error()})
			c.Abort()
			return
		}

		eventIDuuid, err := uuid.Parse(body.EventID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse event id: " + err.Error()})
			log.Println("eventID: ", body.EventID, err)
			c.Abort()
			return
		}

		trip, err := tripService.GetTripByEventID(c.Request.Context(), eventIDuuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get trip by event id: " + err.Error()})
			c.Abort()
			return
		}

		role, err := tripService.GetUserRole(c.Request.Context(), userIDInt, trip.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "can't get role: " + err.Error()})
			c.Abort()
			return
		}

		if !slices.Contains(userRoles, role) {
			log.Println("roles:", userRoles, role)
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
