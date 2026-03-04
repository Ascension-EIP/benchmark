package model

import "../../../../go-mariadb-benchmark/internal/model/github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}

