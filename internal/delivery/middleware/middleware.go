package middleware

import (
	"net/http"
	"strings"

	"github.com/SemenShakhray/doccash/internal/cache"
	"github.com/SemenShakhray/doccash/internal/config"
	"github.com/SemenShakhray/doccash/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AdminTokenRequired(adminToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(401, models.APIResponse{Error: &models.APIError{Code: 401, Text: "invalid admin token"}})
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token != adminToken {
			c.AbortWithStatusJSON(403, models.APIResponse{Error: &models.APIError{Code: 403, Text: "invalid admin token"}})
			return
		}

		c.Next()
	}
}

func AuthMiddleware(cfg *config.Config, cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.Auth.TokenSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, models.APIResponse{Error: &models.APIError{Code: 401, Text: "invalid token"}})
			return
		}

		isValid, err := cache.IsTokenValid(c, tokenStr)
		if err != nil || !isValid {
			c.AbortWithStatusJSON(401, models.APIResponse{Error: &models.APIError{Code: 401, Text: "invalid token"}})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, models.APIResponse{Error: &models.APIError{Code: 401, Text: "invalid claims into token"}})
			return
		}

		loginRaw, ok := claims["login"]
		if !ok {
			c.AbortWithStatusJSON(401, models.APIResponse{Error: &models.APIError{Code: 401, Text: "login not found in claims"}})
			return
		}

		login, ok := loginRaw.(string)
		if !ok {
			c.AbortWithStatusJSON(401, models.APIResponse{Error: &models.APIError{Code: 401, Text: "login not string"}})
			return
		}

		c.Set("login", login)
		c.Next()
	}
}
