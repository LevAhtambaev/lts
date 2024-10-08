package config

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config - структура конфигурации.
// Содержит все конфигурационные данные о сервисе
type Config struct {
	PostgresConfig PostgresConfig `yaml:"postgres" mapstructure:"postgres"`
}

// PostgresConfig - конфигурация для клиента PostgreSQL
type PostgresConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"pass" mapstructure:"pass"`
	Name     string `yaml:"name" mapstructure:"name"`
}

func Read(ctx context.Context, path string) (Config, error) {
	v := viper.New()

	if path == "" {
		return Config{}, fmt.Errorf("read: empty configuration path")
	}

	v.SetConfigType("yaml")
	v.SetConfigFile(path)
	v.WatchConfig()

	err := v.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("read: ReadInConfig error: %w", err)
	}

	cfg := &Config{}
	err = v.Unmarshal(cfg)
	if err != nil {
		return Config{}, fmt.Errorf("read: Unmarshal error: %w", err)
	}

	log.Ctx(ctx).Debug().Interface("cfg", cfg).Msg("config parsed")

	return *cfg, nil
}
