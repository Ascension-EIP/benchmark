package config

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	DB   DBConfig
	Auth AuthConfig
}

func (cfg *Config) validate() error {
	return nil
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
	Params   string
}

func (c DBConfig) DSN() string {
	escapedPassword := url.QueryEscape(c.Password)

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		c.User,
		escapedPassword,
		c.Host,
		c.Port,
		c.Name,
		c.Params,
	)
}

type AuthConfig struct {
	JWTKey []byte
	JWTExp time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		DB: DBConfig{
			User:     getEnvWithFallback("DB_USER", nil, "user"),
			Password: getEnvWithFallback("DB_PASSWORD", nil, "pass"),
			Host:     getEnvWithFallback("DB_HOST", nil, "127.0.0.1"),
			Port:     getEnvWithFallback("DB_PORT", strconv.Atoi, 3306),
			Name:     getEnvWithFallback("DB_NAME", nil, "db"),
			Params:   getEnvWithFallback("DB_PARAMS", nil, "charset=utf8mb4&parseTime=True&loc=Local"),
		},
		Auth: AuthConfig{
			JWTKey: getEnvWithFallback("JWT_KEY", func(s string) ([]byte, error) { return []byte(s), nil }, []byte("THIS_IS_NOT_VERY_SECURED")),
			JWTExp: getEnvWithFallback("JWT_EXP", time.ParseDuration, 24*time.Hour),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
