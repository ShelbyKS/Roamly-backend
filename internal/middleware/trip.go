package middleware

import (
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

func InitTripMiddleware(tripService service.ITripService, userRoles []model.UserTripRole) *TripMiddleware {
	return &TripMiddleware{
		tripService:    tripService,
		validUserRoles: userRoles,
	}
}

func (mw *TripMiddleware) AccessTripMiddleware() gin.HandlerFunc {
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

		role, err := mw.tripService.GetUserRole(c.Request.Context(), userIDInt, tripIDuuid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "can't get role: " + err.Error()})
			c.Abort()
			return
		}
		if !slices.Contains(mw.validUserRoles, role) {
			log.Println("roles:",mw.validUserRoles, role)
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			c.Abort() 
		}

		// c.Set("user_role", role)
		c.Next()
	}
}

// var Mw *Middleware
