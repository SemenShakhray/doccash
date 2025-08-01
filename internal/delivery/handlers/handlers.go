package handlers

import (
	"context"
	"log/slog"

	"github.com/SemenShakhray/doccash/internal/config"
	"github.com/SemenShakhray/doccash/internal/models"
)

type Handler struct {
	service Servicer
	cfg     *config.Config
	log     *slog.Logger
}

func NewHandler(serv Servicer, log *slog.Logger, cfg *config.Config) Handler {
	return Handler{
		service: serv,
		log:     log,
		cfg:     cfg,
	}
}

type Servicer interface {
	RegisterUser(ctx context.Context, req models.AuthRequest) error
	LoginUser(ctx context.Context, req models.AuthRequest) (string, error)
	LogoutUser(ctx context.Context, token string) error
}
