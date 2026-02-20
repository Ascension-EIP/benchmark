package config

import (
	"fmt"
	"net/url"
	"strconv"
)

type Config struct {
	DB DBConfig
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

func Load() (*Config, error) {
	cfg := &Config{
		DB: DBConfig{
			User:     getEnvWithFallback("DB_USER", nil, "root"),
			Password: getEnvWithFallback("DB_PASSWORD", nil, ""),
			Host:     getEnvWithFallback("DB_HOST", nil, "127.0.0.1"),
			Port:     getEnvWithFallback("DB_PORT", strconv.Atoi, 3306),
			Name:     getEnvWithFallback("DB_NAME", nil, "app"),
			Params:   getEnvWithFallback("DB_PARAMS", nil, "charset=utf8mb4&parseTime=True&loc=Local"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
