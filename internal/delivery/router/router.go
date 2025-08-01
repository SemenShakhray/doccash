package router

import (
	"github.com/SemenShakhray/doccash/internal/cache"
	"github.com/SemenShakhray/doccash/internal/config"
	"github.com/SemenShakhray/doccash/internal/delivery/handlers"
	"github.com/SemenShakhray/doccash/internal/delivery/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(h *handlers.Handler, cfg *config.Config, cache *cache.Cache) *gin.Engine {
	r := gin.Default()
	group := r.Group("/api")

	group.POST("/register", middleware.AdminTokenRequired(cfg.Auth.AdminToken), h.RegisterUser)
	group.POST("/auth", h.LoginUser)
	group.DELETE("/auth/:token", h.LogoutUser)

	return r
}
