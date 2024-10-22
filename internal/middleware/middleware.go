package middleware

import (
	"errors"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Middleware struct {
	sessionStorage storage.ISessionStorage
}

func InitMiddleware(sessionStorage storage.ISessionStorage) *Middleware {
	return &Middleware{sessionStorage}
}

var Mw *Middleware

func (mw *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token"})
			c.Abort()
			return
		}

		session, err := mw.sessionStorage.SessionExists(c.Request.Context(), sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
			c.Abort()
			return
		}

		c.Set("user_id", session.UserID)
		c.Next()
	}
}

func (mw *Middleware) UnauthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.Next()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		_, err = mw.sessionStorage.SessionExists(c.Request.Context(), sessionToken)
		if err == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is already logged in"})
			c.Abort()
			return
		}

		c.Next()
	}
}
