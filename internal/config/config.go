package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Env  string `env-required:"true" yaml:"env" env:"ENV"`
		HTTP `yaml:"http"`
	}

	HTTP struct {
		Address     string        `env-required:"true" yaml:"address" env:"HTTP_ADDRESS"`
		Timeout     time.Duration `env-required:"true" yaml:"timeout" env:"HTTP_TIMEOUT"`
		IdleTimeout time.Duration `env-required:"true" yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
