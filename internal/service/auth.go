package service

import (
	"context"
	"fmt"

	"github.com/SemenShakhray/doccash/internal/models"
	"github.com/SemenShakhray/doccash/utils/logger/sl"
	"github.com/SemenShakhray/doccash/utils/token"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterUser(ctx context.Context, req models.AuthRequest) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to generate password hash", sl.Err(err))

		return fmt.Errorf("failed to generate password hash: %w", err)
	}

	err = s.repo.RegisterUser(req.Login, passHash)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}
	return nil
}

func (s *Service) LoginUser(ctx context.Context, req models.AuthRequest) (string, error) {
	hashPasword, err := s.repo.GetUserPassword(req.Login)
	if err != nil {
		s.log.Error("failed to get user password", sl.Err(err))

		return "", fmt.Errorf("failed to get user password: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(hashPasword, []byte(req.Password))
	if err != nil {
		s.log.Error("failed to compare password", sl.Err(err))

		return "", fmt.Errorf("failed to compare password: %w", err)
	}

	token, err := token.GenerateToken(req.Login, s.cfg)
	if err != nil {
		s.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	err = s.cache.SetToken(ctx, token, s.cfg.Auth.TokenTTL)
	if err != nil {
		s.log.Error("failed to set token", sl.Err(err))

		return "", fmt.Errorf("failed to set token: %w", err)
	}

	return token, nil
}

func (s *Service) LogoutUser(ctx context.Context, token string) error {
	return s.cache.DeleteToken(ctx, token)
}
