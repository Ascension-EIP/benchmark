package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		DB   DBConfig `envPrefix:"DB_"`
		Auth AuthConfig
	}

	DBConfig struct {
		Host     string `env:"HOST" envDefault:"localhost"`
		Port     int    `env:"PORT" envDefault:"3306"`
		Name     string `env:"NAME,unset" envDefault:"db"`
		User     string `env:"USER,unset" envDefault:"user"`
		Password string `env:"PASSWORD,unset" envDefault:"pass"`
		Params   string `env:"PARAMS" envDefault:"charset=utf8mb4&parseTime=True&loc=Local"`
	}

	AuthConfig struct {
		JWTKey string        `env:"JWT_KEY,unset" envDefault:"THIS_IS_NOT_REALLY_SECURED"`
		JWTExp time.Duration `env:"JWT_EXP" envDefault:"24h"`
	}
)

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
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
