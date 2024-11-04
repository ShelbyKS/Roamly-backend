package handler

import (
	"github.com/ShelbyKS/Roamly-backend/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/handler/dto"
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

	userGroup := router.Group("/api/v1/user")
	userGroup.Use(middleware.Mw.AuthMiddleware())
	{
		userGroup.GET("/:user_id", handler.GetUserByID)
		userGroup.PUT("/", handler.UpdateUser)
	}
}

// @Summary Get user by ID
// @Description Retrieves a user by their ID.
// @Tags user
// @Produce  json
// @Param user_id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/user/{user_id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idString := c.Param("user_id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse query")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		h.lg.WithError(err).Errorf("Fail to get user with id=%d", id)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": dto.UserConverter{}.ToDto(user)})
}

type UpdateUserRequest struct {
	ID       int    `json:"id" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Update user details
// @Description Updates the details of an existing user.
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body UpdateUserRequest true "User data"
// @Success 200 {string} map[string]string
// @Failure 400 {object} object{err=string}
// @Failure 401 {object} object{err=string}
// @Failure 404 {object} object{err=string}
// @Failure 500 {object} object{err=string}
// @Router /api/v1/user [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var userReq UpdateUserRequest

	err := c.BindJSON(&userReq)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.userService.UpdateUser(c, model.User{
		ID:       userReq.ID,
		Login:    userReq.Login,
		Password: userReq.Password,
	})
	if err != nil {
		h.lg.WithError(err).Errorf("Fail to update user with id=%d", userReq.ID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
