package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
)

type UserHandler struct {
	lg          *logrus.Logger
	userService service.IUserService
}

func NewUserHandler(router *gin.Engine, lg *logrus.Logger, userService service.IUserService) {
	handler := &UserHandler{
		lg:          lg,
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
