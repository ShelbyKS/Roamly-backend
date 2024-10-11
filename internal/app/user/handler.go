package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(router *gin.Engine, service Service) {
	handler := &Handler{
		service: service,
	}

	userGroup := router.Group("/user")
	{
		userGroup.POST("/register", handler.Register)
	}
}

func (h *Handler) Register(c *gin.Context) {
	c.JSON(http.StatusOK, "you registered")
}
