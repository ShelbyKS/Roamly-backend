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

type TripHandler struct {
	lg               *logrus.Logger
	tripService      service.ITripService
	schedulerService service.ISchedulerService
}

func NewTripHandler(router *gin.Engine,
	lg *logrus.Logger,
	tripService service.ITripService,
	schedulerService service.ISchedulerService) {


	handler := &TripHandler{
		lg:          lg,
		tripService: tripService,
		schedulerService: schedulerService,
	}

	userGroup := router.Group("/trip")
	{
		userGroup.GET("/:trip_id", handler.GetTripByID)
		userGroup.POST("/", handler.CreateTrip)
		userGroup.PUT("/", handler.UpdateTrip)
		userGroup.DELETE("/:trip_id", handler.DeleteTrip)

		userGroup.GET("/:trip_id/schedule", handler.ScheduleTrip)
	}
}

func (h *TripHandler) GetTripByID(ctx *gin.Context) {
	idString := ctx.Param("trip_id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	trip, err := h.tripService.GetTripByID(ctx.Request.Context(), id)
	if errors.Is(err, domain.ErrTripNotFound) {
		h.lg.Warnf("Trip with id=%d not found", id)
		ctx.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}
	if err != nil {
		h.lg.WithError(err).Errorf("Fail to get trip with id=%d", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"trip": trip})
}

func (h *TripHandler) DeleteTrip(ctx *gin.Context) {
	idString := ctx.Param("trip_id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.tripService.DeleteTrip(ctx.Request.Context(), id)
	if errors.Is(err, domain.ErrTripNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

type CreateTripRequest struct {
	ID        int    `json:"id" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Region    string `json:"region" binding:"required"`
	UserId    int    `json:"user_id" binding:"required"`
}

func (h *TripHandler) CreateTrip(ctx *gin.Context) {
	var tripReq CreateTripRequest

	err := ctx.BindJSON(&tripReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.tripService.CreateTrip(ctx, model.Trip{
		ID:        tripReq.ID,
		StartTime: tripReq.StartTime,
		EndTime:   tripReq.EndTime,
		Region:    tripReq.Region,
		
		// TODO: ХААААААААРДКОД!!!!!!!!!!!!!!!!!!!!!!
		Users: []*model.User{
			&model.User{
				ID: 1,
			},
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

type UpdateTripRequest struct {
	ID        int    `json:"id" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Region    string `json:"region" binding:"required"`
}

func (h *TripHandler) UpdateTrip(ctx *gin.Context) {
	var tripReq UpdateTripRequest

	err := ctx.BindJSON(&tripReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.tripService.UpdateTrip(ctx, model.Trip{
		ID:        tripReq.ID,
		StartTime: tripReq.StartTime,
		EndTime:   tripReq.EndTime,
		Region:    tripReq.Region,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *TripHandler) ScheduleTrip(ctx *gin.Context) {
	idString := ctx.Param("trip_id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	trip, err := h.tripService.GetTripByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	// h.lg.Info("trip", trip)
	// h.lg.Info("service", h.schedulerService)
	// h.lg.Info("service", trip.Places)

	schedule, err := h.schedulerService.GetSchedule(ctx.Request.Context(), trip.Places)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"schedule" : schedule})
}
