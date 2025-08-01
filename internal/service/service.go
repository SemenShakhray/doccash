package service

import (
	"log/slog"

	"github.com/SemenShakhray/doccash/internal/cache"
	"github.com/SemenShakhray/doccash/internal/config"
	"github.com/SemenShakhray/doccash/internal/delivery/handlers"
)

type Service struct {
	repo  Repositorer
	cache *cache.Cache
	cfg   *config.Config
	log   *slog.Logger
}

func NewService(repo Repositorer, cache *cache.Cache, cfg *config.Config, log *slog.Logger) handlers.Servicer {
	return &Service{
		repo:  repo,
		cache: cache,
		cfg:   cfg,
		log:   log,
	}
}

type Repositorer interface {
	Close() error
	RegisterUser(login string, passHash []byte) error
	GetUserPassword(login string) ([]byte, error)
}
