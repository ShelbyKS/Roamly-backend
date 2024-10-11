package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(router *gin.Engine, userService service.IUserService) {
	handler := &UserHandler{
		userService: userService,
	}

	userGroup := router.Group("/user")
	{
		userGroup.POST("/register", handler.Register)
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	c.JSON(http.StatusOK, "you registered")
}
