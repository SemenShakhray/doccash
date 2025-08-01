package token

import (
	"time"

	"github.com/SemenShakhray/doccash/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(login string, cfg *config.Config) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = login
	claims["exp"] = time.Now().Add(cfg.Auth.TokenTTL).Unix()

	tokenString, err := token.SignedString([]byte(cfg.Auth.TokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
