package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/config"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenSigned := c.GetHeader("Authorization")
		if tokenSigned == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		parts := strings.SplitN(tokenSigned, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims := model.JWTClaims{}
		token, err := jwt.ParseWithClaims(parts[1], &claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(cfg.JWTKey), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
