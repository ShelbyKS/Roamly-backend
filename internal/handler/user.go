package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
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
		// userGroup.GET("/", handler.GetUserByLogin)
		userGroup.GET("/:user_id", handler.GetUserByID)
		userGroup.POST("/", handler.CreateUser)
		userGroup.PUT("/", handler.UpdateUser)
	}
}

func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	idString := ctx.Param("user_id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := h.userService.GetUserByID(ctx.Request.Context(), id)
	if errors.Is(err, domain.ErrUserNotFound) {
		h.lg.Warnf("User with id=%d not found", id)
		ctx.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}
	if err != nil {
		h.lg.WithError(err).Errorf("Fail to get user with id=%d", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) GetUserByLogin(ctx *gin.Context) {
	login := ctx.Param("login")

	user, err := h.userService.GetUserByLogin(ctx.Request.Context(), login)
	if errors.Is(err, domain.ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

type CreateUserRequest struct {
	ID       int    `json:"id" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var userReq CreateUserRequest

	err := ctx.BindJSON(&userReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.userService.CreateUser(ctx, model.User{
		ID:       userReq.ID,
		Login:    userReq.Login,
		Password: userReq.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

type UpdateUserRequest struct {
	ID       int    `json:"id" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var userReq UpdateUserRequest

	err := ctx.BindJSON(&userReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.userService.UpdateUser(ctx, model.User{
		ID:       userReq.ID,
		Login:    userReq.Login,
		Password: userReq.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
