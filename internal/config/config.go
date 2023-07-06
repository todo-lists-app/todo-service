package config

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type Config struct {
	Local
	Vault
	Mongo
}

func Build() (*Config, error) {
	cfg := &Config{}

	if err := BuildLocal(cfg); err != nil {
		return nil, logs.Errorf("build local: %w", err)
	}

	if err := BuildVault(cfg); err != nil {
		return nil, logs.Errorf("build vault: %w", err)
	}

	if err := BuildMongo(cfg); err != nil {
		return nil, logs.Errorf("build mongo: %w", err)
	}

	if err := env.Parse(cfg); err != nil {
		return nil, logs.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
