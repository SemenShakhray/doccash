package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"

	"github.com/SemenShakhray/doccash/internal/config"
	"github.com/SemenShakhray/doccash/internal/service"
	"github.com/SemenShakhray/doccash/utils/logger/sl"
)

type Repository struct {
	DB  *sql.DB
	log *slog.Logger
}

func NewRepository(cfg *config.Config, log *slog.Logger) (service.Repositorer, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Error("failed to create database connection pool", sl.Err(err))

		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	if err := db.Ping(); err != nil {
		log.Error("failed to ping database", sl.Err(err))

		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Repository{
		DB:  db,
		log: log,
	}, nil
}

func (r *Repository) Close() error {
	return r.DB.Close()
}
