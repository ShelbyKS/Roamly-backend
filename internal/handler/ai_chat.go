package handler

import (
	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/handler/dto"
	"github.com/ShelbyKS/Roamly-backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AIChatHandler struct {
	lg            *logrus.Logger
	aiChatService service.IAIChatService
	tripService   service.ITripService
}

func NewAIChatHandler(
	router *gin.Engine,
	lg *logrus.Logger,
	aiChatService service.IAIChatService,
	tripService service.ITripService,
) {
	handler := &AIChatHandler{
		lg:            lg,
		aiChatService: aiChatService,
		tripService:   tripService,
	}

	chatGroup := router.Group("/api/v1/chat")
	chatGroup.Use(middleware.Mw.AuthMiddleware())
	{
		chatGroup.GET("/:trip_id",
			middleware.AccessTripMiddleware(tripService, middleware.ForOwnerAndEditor),
			handler.GetChatHistory)

		chatGroup.POST("/:trip_id",
			middleware.AccessTripMiddleware(tripService, middleware.ForOwnerAndEditor),
			handler.SentMessage)
	}
}

// @Summary Get chat history
// @Description Get chat messages by trip
// @Tags chat
// @Accept json
// @Produce json
// @Param trip_id query string true "Trip ID"
// @Success 200 {object} map[string][]dto.ChatMessageResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/chat/{trip_id} [get]
func (h *AIChatHandler) GetChatHistory(c *gin.Context) {
	tripIDString := c.Param("trip_id")
	tripID, err := uuid.Parse(tripIDString)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID"})
		return
	}

	messages, err := h.aiChatService.GetAIChatMessages(c.Request.Context(), tripID)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to get chat history for trip %s", tripID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	messagesDto := make([]dto.ChatMessageResponse, len(messages))
	for i, msg := range messages {
		messagesDto[i] = dto.AIChatConverter{}.ToDto(msg)
	}

	c.JSON(http.StatusOK, gin.H{"messages": messagesDto})
}

type SentMessageRequest struct {
	Message string `json:"message"`
}

// @Summary Sent message
// @Description Sent message to ai chat
// @Tags chat
// @Produce json
// @Param trip_id path string true "Trip ID"
// @Param message body SentMessageRequest true "Message"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/chat/{trip_id} [post]
func (h *AIChatHandler) SentMessage(c *gin.Context) {
	userIDany, ok := c.Get("user_id")
	if !ok {
		h.lg.Warningln("No user_id in context")
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user_id in context"})
		return
	}
	userID, ok := userIDany.(int)
	if !ok {
		h.lg.Warningln("failed to parse user_id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse user_id to int"})
		return
	}

	tripIdString := c.Param("trip_id")
	tripID, err := uuid.Parse(tripIdString)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID"})
		return
	}

	var messageReq SentMessageRequest

	err = c.Bind(&messageReq)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.aiChatService.SentMessage(c.Request.Context(), model.ChatMessage{
		TripID:  tripID,
		Role:    model.RoleUser,
		Content: messageReq.Message,
	}, userID)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to sent message to ai chat in trip: %s", tripID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
