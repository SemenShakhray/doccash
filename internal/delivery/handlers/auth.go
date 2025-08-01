package handlers

import (
	"log/slog"
	"regexp"

	"github.com/SemenShakhray/doccash/internal/models"
	"github.com/SemenShakhray/doccash/utils/logger/sl"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUser(c *gin.Context) {
	var req models.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid request", sl.Err(err))

		c.JSON(400, models.APIResponse{Error: &models.APIError{Code: 400, Text: "invalid request"}})
		return
	}

	if !validateLogin(req.Login) {
		h.log.Warn("invalid login", slog.String("login", req.Login))

		c.JSON(400, models.APIResponse{Error: &models.APIError{Code: 400, Text: "invalid login"}})
		return
	}

	if !validatePassword(req.Password) {
		h.log.Warn("invalid password", slog.String("password", req.Password))

		c.JSON(400, models.APIResponse{Error: &models.APIError{Code: 400, Text: "invalid password"}})
		return
	}

	err := h.service.RegisterUser(c, req)
	if err != nil {
		c.JSON(500, models.APIResponse{Error: &models.APIError{Code: 500, Text: err.Error()}})
		return
	}

	c.JSON(200, models.APIResponse{Response: map[string]interface{}{"login": req.Login}})
}

func (h *Handler) LoginUser(c *gin.Context) {
	var req models.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid request", sl.Err(err))

		c.JSON(400, models.APIResponse{Error: &models.APIError{Code: 400, Text: "invalid request"}})
		return
	}

	token, err := h.service.LoginUser(c, req)
	if err != nil {
		c.JSON(500, models.APIResponse{Error: &models.APIError{Code: 500, Text: err.Error()}})
		return
	}

	c.JSON(200, models.APIResponse{Response: map[string]interface{}{"token": token}})
}

func (h *Handler) LogoutUser(c *gin.Context) {
	token := c.Param("token")
	h.log.Info("logout", slog.String("token", token))

	err := h.service.LogoutUser(c, token)
	if err != nil {
		h.log.Error("logout error", sl.Err(err))

		c.JSON(500, models.APIResponse{Error: &models.APIError{Code: 500, Text: err.Error()}})
		return
	}
	h.log.Info("logout success", slog.String("token", token))

	c.JSON(200, models.APIResponse{Response: map[string]interface{}{token: true}})
}

func validateLogin(login string) bool {
	if len(login) < 8 {
		return false
	}
	match, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, login)
	return match
}

func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	upper := regexp.MustCompile(`[A-Z]`)
	lower := regexp.MustCompile(`[a-z]`)
	digits := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return upper.MatchString(password) &&
		lower.MatchString(password) &&
		digits.MatchString(password) &&
		special.MatchString(password)
}
