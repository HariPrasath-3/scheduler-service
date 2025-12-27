package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(ctx context.Context) (*AppConfig, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	applyDefaults(&cfg)
	validate(&cfg)

	return &cfg, nil
}

func validate(cfg *AppConfig) {
	if len(cfg.Redis.Hosts) == 0 {
		log.Fatal("redis.hosts must not be empty")
	}
}

func applyDefaults(cfg *AppConfig) {
	// Redis defaults
	if cfg.Redis.PoolSize == 0 {
		cfg.Redis.PoolSize = 10
	}
	if cfg.Redis.MinIdleConns == 0 {
		cfg.Redis.MinIdleConns = 5
	}
}
