package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/SemenShakhray/doccash/internal/cache"
	"github.com/SemenShakhray/doccash/internal/config"
	"github.com/SemenShakhray/doccash/internal/delivery/handlers"
	"github.com/SemenShakhray/doccash/internal/delivery/router"
	"github.com/SemenShakhray/doccash/internal/repository/postgres"
	"github.com/SemenShakhray/doccash/internal/service"
	"github.com/SemenShakhray/doccash/utils/logger"
	"github.com/SemenShakhray/doccash/utils/logger/sl"
)

func RunServer() {
	cfg := config.MustLoad()

	log := logger.SetupLogger()
	log.Info("load config", slog.Any("config", cfg))

	repo, err := postgres.NewRepository(cfg, log)
	if err != nil {
		log.Error("failed to create repository", sl.Err(err))

		os.Exit(1)
	}
	defer repo.Close()

	cache := cache.NewCache(cfg)
	defer cache.Client.Close()

	service := service.NewService(repo, cache, cfg, log)

	handler := handlers.NewHandler(service, log, cfg)

	router := router.NewRouter(&handler, cfg, cache)

	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port)),
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", sl.Err(err))
			os.Exit(1)
		}
	}()

	sig := <-sigint
	log.Info("received signal", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Info("failed to stop server", sl.Err(err))
	}
}
