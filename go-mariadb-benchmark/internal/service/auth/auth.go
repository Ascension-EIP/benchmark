package auth

import (
	"context"
	"time"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/config"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/dto/request"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/model"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/repo"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	r   repo.Auth
	cfg config.AuthConfig
}

func New(r repo.Auth, cfg config.AuthConfig) *Service {
	return &Service{
		r:   r,
		cfg: cfg,
	}
}

func (s *Service) Signup(c context.Context, req request.Signup) error {
	if err := ValidateUsername(req.Username); err != nil {
		return err
	}
	if err := ValidatePassword(req.Password); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.r.CreateUser(&model.User{
		Username: req.Username,
		Password: string(hash),
	}); err != nil {
		return err
	}
	return nil
}

func (s *Service) Login(c context.Context, req request.Login) (string, error) {
	user, err := s.r.GetUserByUsername(req.Username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(s.cfg.JWTExp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString(s.cfg.JWTKey)
	if err != nil {
		return "", err
	}
	return tokenSigned, nil
}
