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

type InviteHandler struct {
	inviteService service.IInviteService
	tripService   service.ITripService
	lg            *logrus.Logger
}

func NewInviteHandler(
	router *gin.Engine,
	lg *logrus.Logger,
	inviteService service.IInviteService,
	tripService service.ITripService,
) {
	handler := &InviteHandler{
		lg:            lg,
		inviteService: inviteService,
		tripService:   tripService,
	}

	tripInviteGroup := router.Group("/api/v1/trip")
	tripInviteGroup.Use(middleware.Mw.AuthMiddleware())
	{
		tripInviteGroup.POST(
			"/invite/",
			middleware.AccessTripByTripIdFromBodyMiddleware(tripService, middleware.ForOwner),
			handler.EnableInvitation,
		)
		tripInviteGroup.DELETE(
			"/invite/",
			middleware.AccessTripByTripIdFromBodyMiddleware(tripService, middleware.ForOwner),
			handler.DisableInvitation,
		)

		tripInviteGroup.GET(
			"/:trip_id/invite",
			middleware.AccessTripMiddleware(tripService, middleware.ForOwnerAndEditor),
			handler.GetTripInvitations,
		)

		tripInviteGroup.POST("/join/:invite_token", handler.JoinTrip)
	}
}

type EnableInvitationRequest struct {
	TripID uuid.UUID `json:"trip_id" binding:"required"`
	Access string    `json:"access" binding:"required"`
}

// @Summary Enable trip invitation
// @Description Enable trip invitation by access
// @Tags invite
// @Accept json
// @Produce json
// @Param event body EnableInvitationRequest true "Invitation data"
// @Success 200 {object} map[string]string "invite_token: bla_bla"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/invite [post]
func (h *InviteHandler) EnableInvitation(c *gin.Context) {
	var req EnableInvitationRequest

	if err := c.BindJSON(&req); err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invitation, err := h.inviteService.EnableInvitation(c.Request.Context(), model.Invite{
		TripID: req.TripID,
		Access: req.Access,
	})

	if err != nil {
		h.lg.WithError(err).Errorf("failed to enable invite for trip %s with access %s", req.TripID, req.Access)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invite_token": invitation.Token})
}

type DisableInvitationRequest struct {
	TripID uuid.UUID `json:"trip_id" binding:"required"`
	Access string    `json:"access" binding:"required"`
}

// @Summary Disable trip invitation
// @Description Disable trip invitation by access
// @Tags invite
// @Accept json
// @Produce json
// @Param event body DisableInvitationRequest true "Invitation data"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/invite [delete]
func (h *InviteHandler) DisableInvitation(c *gin.Context) {
	var req DisableInvitationRequest

	if err := c.BindJSON(&req); err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.inviteService.DisableInvitation(c.Request.Context(), model.Invite{
		TripID: req.TripID,
		Access: req.Access,
	})

	if err != nil {
		h.lg.WithError(err).Errorf("failed to disable invite for trip %s with access %s", req.TripID, req.Access)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Trip invitations
// @Description Get trip invite tokens
// @Tags invite
// @Produce json
// @Success 200 {object} map[string][]dto.InviteResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/{trip_id}/invite [get]
func (h *InviteHandler) GetTripInvitations(c *gin.Context) {
	tripID, err := uuid.Parse(c.Param("trip_id"))
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID"})
		return
	}

	invitations, err := h.inviteService.GetTripInvitations(c.Request.Context(), tripID)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to get invite tokens for trip %s", tripID)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	invitationsDto := make([]dto.InviteResponse, len(invitations))
	for i, invitation := range invitations {
		invitationsDto[i] = dto.InviteConverter{}.ToDto(invitation)
	}

	c.JSON(http.StatusOK, gin.H{
		"invitations": invitationsDto,
	})
}

// @Summary Join trip
// @Description Join trip via invite_token
// @Tags invite
// @Accept json
// @Produce json
// @Param event path string true "Invite token"
// @Success 200 {object} map[string]string "trip_id: bla_bla"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/trip/join/{invite_token} [post]
func (h *InviteHandler) JoinTrip(c *gin.Context) {
	inviteToken := c.Param("invite_token")

	userID, ok := c.Get("user_id")
	if !ok {
		h.lg.Errorf("Fail to get user_id from context")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fail to get user_id from context"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		h.lg.Errorf("User ID is not an integer")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not an integer"})
		return
	}

	tripID, err := h.inviteService.JoinTrip(c.Request.Context(), inviteToken, userIDInt)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to join trip via invite token: %s", inviteToken)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip_id": tripID})
}
