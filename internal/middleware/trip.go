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

type TripMiddleware struct {
	tripService    service.ITripService
	validUserRoles []model.UserTripRole
}

var (
	ForOwner          = []model.UserTripRole{model.Owner}
	ForOwnerAndEditor = []model.UserTripRole{model.Owner, model.Editor}
	ForAll            = []model.UserTripRole{model.Owner, model.Editor, model.Reader}
)


func AccessTripMiddleware(tripService service.ITripService, userRoles []model.UserTripRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No user id"})
			c.Abort()
		}
		userIDInt, ok := userID.(int)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse user id"})
			c.Abort()
		}

		tripID := c.Param("trip_id")
		tripIDuuid, err := uuid.Parse(tripID)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse trip id"})
			c.Abort()
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
		}

		c.Next()
	}
}

func AccessTripFromBodyMiddleware(tripService service.ITripService, userRoles []model.UserTripRole) gin.HandlerFunc {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBuffer.Bytes()))

		var body struct {
			TripID string `json:"id"`
		}
		if err := json.NewDecoder(&bodyBuffer).Decode(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
			c.Abort()
			return
		}

		tripIDuuid, err := uuid.Parse(body.TripID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse trip id"})
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

