package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/SemenShakhray/doccash/utils/logger/sl"
	"github.com/google/uuid"
)

func (r *Repository) RegisterUser(login string, passHash []byte) error {
	tx, err := r.DB.Begin()
	if err != nil {
		r.log.Error("begin tx", sl.Err(err))

		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM users WHERE username = $1)`,
		login).Scan(&exists)

	if err != nil {
		r.log.Error("check exists", sl.Err(err))

		return fmt.Errorf("check exists: %w", err)
	}
	if exists {
		r.log.Warn("user with username already exists", slog.String("username", login))

		return fmt.Errorf("user with %s already exists", login)
	}

	var UserID uuid.UUID
	err = tx.QueryRow(`
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		login, passHash).Scan(&UserID)

	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	r.log.Info("user registered", slog.String("username", login), slog.Any("user_id", UserID))

	return nil
}

func (r *Repository) GetUserPassword(login string) ([]byte, error) {
	var passHash []byte
	var user_id uuid.UUID

	err := r.DB.QueryRow(`
		SELECT user_id, password_hash 
		FROM users 
		WHERE username = $1`,
		login,
	).Scan(&user_id, &passHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Warn("user not found", slog.String("username", login))
			return nil, fmt.Errorf("user not found")
		}

		r.log.Error("failed to get user password", sl.Err(err))
		return nil, fmt.Errorf("get user password: %w", err)
	}

	return passHash, nil
}
