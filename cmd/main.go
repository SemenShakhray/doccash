package main

import (
	"log/slog"

	"github.com/SemenShakhray/doccash/internal/config"
	"github.com/SemenShakhray/doccash/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger()
	log.Info("load config", slog.Any("config", cfg))
}
