package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/middleware"
)

type AuthHandler struct {
	lg          *logrus.Logger
	authService service.IAuthService
}

func NewAuthHandler(router *gin.Engine, lg *logrus.Logger, authService service.IAuthService) {
	handler := &AuthHandler{
		lg:          lg,
		authService: authService,
	}

	userGroup := router.Group("/api/v1/auth")
	{
		userGroup.POST("/register", middleware.Mw.UnauthMiddleware(), handler.Register)
		userGroup.POST("/login", middleware.Mw.UnauthMiddleware(), handler.Login)
		userGroup.POST("/logout", middleware.Mw.AuthMiddleware(), handler.Logout)
		userGroup.GET("/check", middleware.Mw.AuthMiddleware(), handler.CheckAuth)
	}
}

type RegisterRequest struct {
	Login    string `json:"login" form:"login" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// @Summary Register a new user
// @Description Register a new user with the provided details.
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body RegisterRequest true "User data"
// @Success 200 {object} object{body=object{user_id=int}}
// @Failure 400 {object} object{err=string}
// @Failure 500 {object} object{err=string}
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	err := c.Bind(&req)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to parse request body"})
		return
	}

	newUser, err := h.authService.Register(c.Request.Context(), model.User{
		Login:    req.Login,
		Email:    req.Email,
		Password: []byte(req.Password),
	})
	if err != nil {
		h.lg.WithError(err).Errorf("failed to register new user with email=%s", req.Email)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()})
		return
	}

	session, err := h.authService.Login(c.Request.Context(), newUser)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to login after registration with email=%s", newUser.Email)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()})
		return
	}

	expiresIn := int(session.ExpiresAt.Sub(time.Now()).Seconds())
	c.SetCookie("session_token", session.Token, expiresIn, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"user_id": newUser.ID})
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// @Summary Login a user
// @Description Authenticate a user with the provided credentials.
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body LoginRequest true "User credentials"
// @Success 200 {object} object{body=object{user_id=int}}
// @Failure 400 {object} object{err=string}
// @Failure 404 {object} object{err=string}
// @Failure 500 {object} object{err=string}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	err := c.Bind(&req)
	if err != nil {
		h.lg.WithError(err).Errorf("failed to parse body")
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to parse request body"})
		return
	}

	session, err := h.authService.Login(c.Request.Context(), model.User{
		Email:    req.Email,
		Password: []byte(req.Password),
	})
	if err != nil {
		h.lg.WithError(err).Errorf("failed to login  with email=%s", req.Email)
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()}) //todo: mb not 500
		return
	}

	expiresIn := int(session.ExpiresAt.Sub(time.Now()).Seconds())
	c.SetCookie("session_token", session.Token, expiresIn, "/", "roamly.ru", true, true)
	c.SetSameSite(http.SameSiteNoneMode) //todo: delete for prod

	c.JSON(http.StatusOK, gin.H{"user_id": session.UserID})
}

// @Summary Logout a user
// @Description Logout a user. Delete session.
// @Tags user
// @Accept  json
// @Produce  json
// @Success 204
// @Failure 400 {object} object{err=string}
// @Failure 401 {object} object{err=string}
// @Failure 404 {object} object{err=string}
// @Failure 500 {object} object{err=string}
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		h.lg.WithError(err).Errorf("failed to get cookie")
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	err = h.authService.Logout(c.Request.Context(), model.Session{
		Token: sessionToken,
	})
	if err != nil {
		h.lg.WithError(err).Errorf("failed to logout")
		c.JSON(domain.GetStatusCodeByError(err), gin.H{"err": err.Error()}) //todo: mb not 500
		return
	}

	c.SetCookie("session_token", "", -1, "/", "", false, true)

	c.Status(http.StatusNoContent)
}

// @Summary Check auth
// @Description Check if user is authenticated.
// @Tags user
// @Produce  json
// @Success 200 {object} object{body=object{user_id=int}}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/check [get]
func (h *AuthHandler) CheckAuth(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		h.lg.Errorf("Fail to get user_id from context")
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Fail to get user_id from context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}
